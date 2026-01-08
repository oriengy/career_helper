package message

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"app_server/model"
	"app_server/pkg/aiapi"
	"app_server/pkg/db"
	"app_server/pkg/fn"
	"app_server/pkg/idgen"
	"app_server/pkg/oai"
	"app_server/pkg/ossc"
	"app_server/proto/message"
	"app_server/service/auth"

	connect "connectrpc.com/connect"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// ChatMessageService 统一的消息服务实现
type ChatMessageService struct{}

// ListChatMessages 查询消息列表 - 合并原来的 ListConsultMessages 和 ListFriendMessages
func (s *ChatMessageService) ListChatMessages(ctx context.Context, connectReq *connect.Request[message.ListChatMessagesRequest]) (*connect.Response[message.ListChatMessagesResponse], error) {
	req := connectReq.Msg

	// session_id 是必填的
	if req.SessionId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}
	sessionID := fn.Atoi[uint](req.SessionId)
	userID := auth.GetUserID(ctx)

	// 构建查询条件
	query := db.GetDB().Model(&model.ChatMessage{})

	// 添加用户ID和会话ID
	query = query.Where("user_id = ? AND session_id = ?", userID, sessionID)

	// 处理消息类型过滤
	if req.MsgType != "" {
		query = query.Where("msg_type = ?", req.MsgType)
	}

	// 处理角色过滤
	if len(req.Roles) > 0 {
		query = query.Where("role IN ?", req.Roles)
	}

	// 根据ID列表过滤
	if len(req.Ids) > 0 {
		query = query.Where("id IN ?", req.Ids)
	}

	// 根据父消息ID列表过滤
	if len(req.ParentIds) > 0 {
		query = query.Where("parent_id IN ?", req.ParentIds)
	}

	// 处理分页
	var lastID int
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 1000 // 默认每页1000条
	}

	if req.PageToken != "" {
		lastID, _ = strconv.Atoi(req.PageToken)
		if lastID > 0 {
			query = query.Where("id < ?", lastID)
		}
	}

	// 查询数据并排序
	var dbMessages []model.ChatMessage
	if err := query.Order("id DESC").Limit(pageSize).Find(&dbMessages).Error; err != nil {
		slog.Error("list chat messages error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if len(dbMessages) == 0 {
		// 从配置表加载新建会话引导消息
		var guideConfig model.Config
		err := db.GetDB().Where("k = ?", "guide_msg:on_new_chat").First(&guideConfig).Error
		if err != nil || guideConfig.Value == "" {
			// 配置不存在或为空，返回空消息列表
			return connect.NewResponse(&message.ListChatMessagesResponse{}), nil
		}

		// 尝试将配置值反序列化为 ChatMessage 数组
		jsoniter.Unmarshal([]byte(guideConfig.Value), &dbMessages)

		// 转换为 proto 消息数组
		baseTime := time.Now()

		mutable.Reverse(dbMessages)
		for i := range dbMessages {
			dbMsg := &dbMessages[i]
			dbMsg.UserID = userID
			dbMsg.SessionID = sessionID
			dbMsg.Tags = append(dbMsg.Tags, "disable_interact")
			dbMsg.MsgAt = baseTime
		}
	}

	// 反转消息顺序，让较早的消息排在前面
	mutable.Reverse(dbMessages)

	// 从 dbMessages 中分离翻译消息和其他消息
	translationMap := make(map[uint]model.ChatMessage) // parentID -> 最新翻译消息
	aiConsultMap := make(map[uint]uint)                // parentID -> 最新 AI CONSULT 消息
	var filteredMessages []*message.ChatMessage

	// 1. 第一遍遍历：收集所有翻译消息和 AI CONSULT 消息
	for _, msg := range dbMessages {
		if msg.MsgType == model.MessageTypeTranslate && msg.ParentID > 0 {
			// 保留最新的翻译（ID 最大的）
			if existing, ok := translationMap[msg.ParentID]; !ok || msg.ID > existing.ID {
				translationMap[msg.ParentID] = msg
			}
		} else if msg.MsgType == model.MessageTypeConsult && msg.ParentID > 0 {
			// 保留最新的 AI CONSULT 消息（ID 最大的）
			if existingID, ok := aiConsultMap[msg.ParentID]; !ok || msg.ID > existingID {
				aiConsultMap[msg.ParentID] = msg.ID
			}
		}
	}

	// 2. 第二遍遍历：处理消息，过滤旧的 AI CONSULT 消息
	for _, msg := range dbMessages {
		// 跳过翻译类型的消息
		if msg.MsgType == model.MessageTypeTranslate {
			continue
		}

		// 如果是 AI CONSULT 消息，检查是否是最新的
		if msg.MsgType == model.MessageTypeConsult && msg.ParentID > 0 {
			// 如果不是最新的 AI CONSULT 消息，跳过
			if latestID := aiConsultMap[msg.ParentID]; msg.ID != latestID {
				continue
			}
		}

		protoMsg := msg.ToProto()

		// 如果有对应的翻译，添加到 translate_content
		if translation, ok := translationMap[msg.ID]; ok {
			protoMsg.TranslateContent = &translation.Content
		}

		filteredMessages = append(filteredMessages, protoMsg)
	}

	// 处理结果
	var nextPageToken string
	if len(dbMessages) == pageSize {
		nextPageToken = strconv.Itoa(int(dbMessages[0].ID))
	}

	return connect.NewResponse(&message.ListChatMessagesResponse{
		Messages:      filteredMessages,
		NextPageToken: nextPageToken,
	}), nil
}

// CreateChatMessage 创建消息 - 合并原来的 SendConsultMessage 和 CreateFriendMessage
func (s *ChatMessageService) CreateChatMessage(ctx context.Context, connectReq *connect.Request[message.CreateChatMessageRequest]) (*connect.Response[message.CreateChatMessageResponse], error) {
	req := connectReq.Msg
	userID := auth.GetUserID(ctx)

	// 如果没有消息需要创建，直接返回
	if len(req.Messages) == 0 {
		return connect.NewResponse(&message.CreateChatMessageResponse{}), nil
	}

	// 开启事务
	tx := db.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 将proto对象转换为数据库对象
	dbMessages := make([]model.ChatMessage, len(req.Messages))
	for i, protoMsg := range req.Messages {
		// 创建一个新的 ChatMessage
		dbMsg := &model.ChatMessage{
			UserID:    userID,
			SessionID: fn.Atoi[uint](protoMsg.SessionId),
			ParentID:  fn.Atoi[uint](protoMsg.ParentId),
			Role:      protoMsg.Role,
			MsgType:   protoMsg.MsgType,
			Content:   protoMsg.Content,
			Tags:      protoMsg.Tags,
			MsgAt:     time.Now(),
		}

		// 使用提供的消息时间（如果有）
		if protoMsg.MsgAt != nil && protoMsg.MsgAt.IsValid() {
			dbMsg.MsgAt = protoMsg.MsgAt.AsTime()
		}

		dbMessages[i] = *dbMsg
	}

	// 批量创建消息
	if err := tx.Create(&dbMessages).Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 转换为proto消息
	protoMessages := fn.Map(dbMessages, model.ChatMessage.ToProto)

	return connect.NewResponse(&message.CreateChatMessageResponse{
		Messages: protoMessages,
	}), nil
}

// UpdateChatMessage 更新消息
func (s *ChatMessageService) UpdateChatMessage(ctx context.Context, connectReq *connect.Request[message.UpdateChatMessageRequest]) (*connect.Response[message.UpdateChatMessageResponse], error) {
	req := connectReq.Msg
	userID := auth.GetUserID(ctx)

	// 如果没有消息需要更新，直接返回
	if len(req.Messages) == 0 {
		return connect.NewResponse(&message.UpdateChatMessageResponse{}), nil
	}

	// 开启事务
	tx := db.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var updatedMessages []model.ChatMessage

	// 逐个更新消息
	for _, protoMsg := range req.Messages {
		if protoMsg.Id == "" {
			continue // 跳过没有ID的消息
		}

		id, _ := strconv.Atoi(protoMsg.Id)
		if id <= 0 {
			continue
		}

		// 构建更新字段
		updates := map[string]interface{}{}

		// 只更新非空字段
		if protoMsg.Content != "" {
			updates["content"] = protoMsg.Content
		}
		if protoMsg.Role != "" {
			updates["role"] = protoMsg.Role
		}
		if protoMsg.MsgType != "" {
			updates["msg_type"] = protoMsg.MsgType
		}
		if len(protoMsg.Tags) > 0 {
			updates["tags"] = protoMsg.Tags
		}

		// 更新消息
		var dbMessage model.ChatMessage
		if err := tx.Model(&model.ChatMessage{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(updates).
			First(&dbMessage, id).Error; err != nil {
			tx.Rollback()
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		updatedMessages = append(updatedMessages, dbMessage)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 转换为proto消息
	protoMessages := fn.Map(updatedMessages, model.ChatMessage.ToProto)

	return connect.NewResponse(&message.UpdateChatMessageResponse{
		Messages: protoMessages,
	}), nil
}

// DeleteChatMessage 删除消息 - 合并原来的 RecallConsultMessage 和 DeleteFriendMessage
func (s *ChatMessageService) DeleteChatMessage(ctx context.Context, connectReq *connect.Request[message.DeleteChatMessageRequest]) (*connect.Response[message.DeleteChatMessageResponse], error) {
	req := connectReq.Msg
	userID := auth.GetUserID(ctx)

	// 如果没有消息需要删除，直接返回
	if len(req.Ids) == 0 {
		return connect.NewResponse(&message.DeleteChatMessageResponse{
			DeletedCount: 0,
		}), nil
	}

	// 把字符串ID转换为整数ID
	var ids []int
	for _, idStr := range req.Ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return connect.NewResponse(&message.DeleteChatMessageResponse{
			DeletedCount: 0,
		}), nil
	}

	// 删除消息
	result := db.GetDB().Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.ChatMessage{})
	if result.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, result.Error)
	}

	return connect.NewResponse(&message.DeleteChatMessageResponse{
		DeletedCount: int32(result.RowsAffected),
	}), nil
}

// ParseImageMessages 解析图片中的消息
func (s *ChatMessageService) ParseImageMessages(ctx context.Context, connectReq *connect.Request[message.ParseImageMessagesRequest]) (*connect.Response[message.ParseImageMessagesResponse], error) {
	req := connectReq.Msg
	userID := auth.GetUserID(ctx)

	// 验证参数
	sessionID := fn.Atoi[uint](req.SessionId)
	if sessionID == 0 || req.ImageUrl == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	// 获取图片公开访问链接
	imageUrl := req.ImageUrl
	publicURL, err := ossc.GetPublic().UserFileBucket().SignURL(imageUrl, "GET", 3600)
	slog.Info("publicURL", "publicURL", publicURL)
	if err != nil {
		slog.Error("sign image url error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 调用火山API解析图片中的聊天记录
	chatLines, err := aiapi.ParseImageChat(publicURL)
	if err != nil {
		slog.Error("parse image chat error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 开始事务
	tx := db.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 解析聊天记录并创建消息
	now := time.Now()
	var parsedMessages []model.ChatMessage

	for _, line := range chatLines {
		role, content, ok := aiapi.ParseChatLine(line)
		if !ok {
			continue
		}

		// 创建消息对象
		msg := model.ChatMessage{
			UserID:    userID,
			SessionID: sessionID,
			Role:      role,
			MsgAt:     now,
			MsgType:   "HISTORY",
			Content:   content,
			Tags:      []string{"parsed_from_image"},
		}

		parsedMessages = append(parsedMessages, msg)
	}

	// 如果没有解析出有效消息
	if len(parsedMessages) == 0 {
		return connect.NewResponse(&message.ParseImageMessagesResponse{
			Message: "未能从图片中解析出有效聊天记录",
		}), nil
	}

	// 从数据库查询最近的消息以避免重复
	var lastMessages []model.ChatMessage
	tx.Where("user_id = ? AND session_id = ? AND msg_type = ?", userID, sessionID, "HISTORY").
		Order("id DESC").
		Limit(20).
		Find(&lastMessages)

	// 检查是否存在重复消息
	dbMessages := fn.Filter(parsedMessages, func(msg model.ChatMessage) bool {
		return !lo.ContainsBy(lastMessages, func(lastMsg model.ChatMessage) bool {
			return lastMsg.Content == msg.Content && lastMsg.Role == msg.Role
		})
	})

	if len(dbMessages) == 0 {
		return connect.NewResponse(&message.ParseImageMessagesResponse{
			Message: "解析消息全部重复，没有新增消息",
		}), nil
	}

	// 批量创建消息
	if err := tx.Create(&dbMessages).Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&message.ParseImageMessagesResponse{
		Success:  true,
		Message:  "解析成功",
		Messages: fn.Map(dbMessages, model.ChatMessage.ToProto),
	}), nil
}

// buildChatHistoryWithExclude 构建聊天历史记录，支持排除某个消息之后的内容（用于 regenerate）
func (s *ChatMessageService) buildChatHistoryWithExclude(ctx context.Context, allMessages []model.ChatMessage, userProfile, friendProfile *model.Profile) []openai.ChatCompletionMessage {
	// 处理翻译消息去重 - 保留最新的翻译
	translationMap := make(map[uint]model.ChatMessage) // parentID -> 最新翻译
	var filteredMessages []model.ChatMessage

	for _, msg := range allMessages {
		if msg.MsgType == model.MessageTypeTranslate {
			// 翻译消息，保存或更新
			if msg.ID > translationMap[msg.ParentID].ID {
				translationMap[msg.ParentID] = msg
			}
		} else {
			// 非翻译消息，直接加入
			filteredMessages = append(filteredMessages, msg)
		}
	}

	// 构建 OpenAI 消息列表
	var openaiMessages []openai.ChatCompletionMessage

	// 1. 获取系统提示词
	systemPrompt := s.getSystemPrompt(ctx, friendProfile)
	openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: systemPrompt})

	// 2. 添加用户和朋友的 profile 信息作为上下文
	if userProfile != nil || friendProfile != nil {
		var profileContext strings.Builder
		const prefix = "对话双方信息：\n"
		profileContext.WriteString(prefix)

		if userProfile != nil && userProfile.FormatPropertyLinesString() != "" {
			profileContext.WriteString("\n**用户资料**\n")
			profileContext.WriteString(userProfile.FormatPropertyLinesString())
			profileContext.WriteString("\n")
		}

		if friendProfile != nil && friendProfile.FormatPropertyLinesString() != "" {
			// 动态获取对方名称
			friendName := "对方"
			if friendProfile.Name != "" {
				friendName = friendProfile.Name
			}
			profileContext.WriteString(fmt.Sprintf("\n**%s资料**\n", friendName))
			profileContext.WriteString(friendProfile.FormatPropertyLinesString())
			profileContext.WriteString("\n")
		}

		if profileContext.Len() > len(prefix) {
			openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: profileContext.String(),
			})
		}
	}

	// 3. 处理消息历史
	var currentHistoryBatch []string

	for _, msg := range filteredMessages {
		switch msg.MsgType {
		case model.MessageTypeHistory:
			historyLine := msg.HistoryCnString()
			currentHistoryBatch = append(currentHistoryBatch, historyLine)
			// 获取对应的翻译（如果有）
			if trans, ok := translationMap[msg.ID]; ok {
				currentHistoryBatch = append(currentHistoryBatch, trans.HistoryCnString())
			}

		case model.MessageTypeConsult:
			// 先打包之前的历史消息
			if len(currentHistoryBatch) > 0 {
				openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: "二人聊天记录：\n" + strings.Join(currentHistoryBatch, "\n"),
				})
				currentHistoryBatch = []string{}
			}

			// 根据角色处理 CONSULT 消息
			oaiMsg := openai.ChatCompletionMessage{Content: msg.Content}
			switch msg.Role {
			case model.MessageRoleUser:
				oaiMsg.Role = openai.ChatMessageRoleUser // 用户咨询使用 User 角色
			case model.MessageRoleAI:
				oaiMsg.Role = openai.ChatMessageRoleAssistant // AI 回复使用 Assistant 角色
			default:
				oaiMsg.Role = openai.ChatMessageRoleUser // 默认使用 User 角色
			}
			openaiMessages = append(openaiMessages, oaiMsg)
		}
	}

	return openaiMessages
}

// getSystemPrompt 获取系统提示词，按优先级从不同来源获取
func (s *ChatMessageService) getSystemPrompt(ctx context.Context, friendProfile *model.Profile) string {
	// 1. 优先从 friendProfile.Prompt 获取
	if friendProfile != nil && friendProfile.Prompt != "" {
		return friendProfile.Prompt
	}

	// 2. 从 config 表获取
	var config model.Config
	if err := db.GetDB().Model(&model.Config{}).
		Where("k = ?", "prompt:consult:default").
		First(&config).Error; err == nil && config.Value != "" {
		return config.Value
	}

	// 3. 使用默认值
	return `你是一个专业的职场沟通顾问，帮助用户更好地理解和回应领导或者上司的消息。
你的任务是基于对话历史，为用户提供专业、有帮助的回复建议。
回复要简洁明了，易于理解，并且具有高情商。`
}

// callAIForReply 调用 AI 生成回复
func (s *ChatMessageService) callAIForReply(ctx context.Context, openaiMessages []openai.ChatCompletionMessage) (string, error) {
	// 使用新的 OAI 包调用 OpenAI
	content, err := oai.Get().CreateChatCompletion(ctx, oai.ChatCompletionRequest{
		Messages: openaiMessages,
	})
	if err != nil {
		slog.Error("AI completion error", "error", err)
		return "", err
	}

	return content, nil
}

// formatPromptForLog 格式化提示词用于日志记录
func (s *ChatMessageService) formatPromptForLog(openaiMessages []openai.ChatCompletionMessage) string {
	return oai.FormatMessagesForLog(openaiMessages)
}

// SendConsultMessage 发送咨询消息 - 专门用于AI咨询回复
func (s *ChatMessageService) SendConsultMessage(ctx context.Context, connectReq *connect.Request[message.SendConsultMessageRequest]) (*connect.Response[message.SendConsultMessageResponse], error) {
	req := connectReq.Msg
	userID := auth.GetUserID(ctx)

	// 验证参数
	sessionID := fn.Atoi[uint](req.SessionId)
	if sessionID == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	database := db.GetDB()

	// 获取会话信息以获取 profile 信息
	var chatSession model.ChatSession
	if err := database.Model(&model.ChatSession{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		First(&chatSession).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("会话未找到"))
	}

	// 获取 user profile
	var user model.User
	var userProfile model.Profile
	if err := database.Model(&model.User{}).
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("用户未找到"))
	}
	if err := database.Model(&model.Profile{}).
		Where("user_id = ? AND id = ?", userID, user.ProfileID).
		First(&userProfile).Error; err != nil {
		slog.Warn("未找到用户Profile", "userID", userID)
	}

	// 获取 friend profile
	var friendProfile model.Profile
	if chatSession.ProfileID > 0 {
		if err := database.Model(&model.Profile{}).
			Where("id = ?", chatSession.ProfileID).
			First(&friendProfile).Error; err != nil {
			slog.Warn("未找到好友Profile", "profileID", chatSession.ProfileID)
		}
	}

	// 获取当前会话的所有消息
	var allMessages []model.ChatMessage
	if err := database.Model(&model.ChatMessage{}).
		Where("session_id = ? AND user_id = ?", sessionID, userID).
		Order("id ASC").
		Find(&allMessages).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var userConsultMsg model.ChatMessage

	// 判断是否是 regenerate 操作
	targetID := fn.Atoi[uint](req.GetTargetId())
	if targetID > 0 { // regenerate 模式
		// 验证 target_id 是最后一条AI CONSULT消息
		var wasLast bool
		var targetMsg model.ChatMessage
		for i := range allMessages {
			msg := allMessages[len(allMessages)-1-i]
			if msg.MsgType == model.MessageTypeConsult && msg.Role == model.MessageRoleAI {
				wasLast = msg.ID == targetID
				targetMsg = msg
				break
			}
		}
		if !wasLast {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("target_id must be the last message"))
		}

		// 查找lastUserConsultMsg
		userConsultMsg, _ = lo.Find(allMessages, func(msg model.ChatMessage) bool {
			return msg.ID == targetMsg.ParentID
		})
		if userConsultMsg.ID == 0 {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("target_id must be the last message"))
		}

		// 大于等于 target_id 的消息都不参与构建历史
		allMessages = fn.Filter(allMessages, func(msg model.ChatMessage) bool {
			return msg.ID < targetID
		})
	} else {
		// 正常模式
		if req.Content == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("content is required"))
		}

		userConsultMsg = model.ChatMessage{
			UserID:    userID,
			SessionID: sessionID,
			ParentID:  0, // 正常模式下没有 parent_id
			Role:      model.MessageRoleUser,
			MsgType:   model.MessageTypeConsult,
			Content:   req.Content,
			MsgAt:     time.Now(),
		}
		userConsultMsg.ID = idgen.Uint()
		allMessages = append(allMessages, userConsultMsg)
	}

	// 构建聊天历史
	openaiMessages := s.buildChatHistoryWithExclude(ctx, allMessages, &userProfile, &friendProfile)

	// 调用 AI 生成回复
	replyContent, err := s.callAIForReply(ctx, openaiMessages)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 3. 创建回复消息
	replyMsg := model.ChatMessage{
		UserID:    userID,
		SessionID: sessionID,
		ParentID:  userConsultMsg.ID,
		Role:      model.MessageRoleAI,
		MsgType:   model.MessageTypeConsult,
		Content:   replyContent,
		Tags:      []string{"ai_reply"},
		MsgAt:     time.Now(),
	}
	replyMsg.ID = idgen.Uint()
	createMsgs := []model.ChatMessage{replyMsg}
	if targetID == 0 {
		createMsgs = append([]model.ChatMessage{userConsultMsg}, createMsgs...)
	}
	if err := database.Create(&createMsgs).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 输出调试日志
	log.Printf("SendConsultMessage success\n%s\n [last_reply]:\n%s\n", s.formatPromptForLog(openaiMessages), replyContent)

	slog.Info("SendConsultMessage success",
		"consultId", userConsultMsg.ID,
		"replyId", replyMsg.ID,
		"content", req.Content,
		"reply", replyContent)

	return connect.NewResponse(&message.SendConsultMessageResponse{
		Consult: userConsultMsg.ToProto(),
		Reply:   replyMsg.ToProto(),
	}), nil
}

// FeedbackToMessage 用户反馈 - 将 attitude 更新到消息的 tags 中
func (s *ChatMessageService) FeedbackToMessage(ctx context.Context, connectReq *connect.Request[message.FeedbackToMessageRequest]) (*connect.Response[message.FeedbackToMessageResponse], error) {
	req := connectReq.Msg
	userID := auth.GetUserID(ctx)

	// 验证参数
	if req.SessionId == "" || req.MessageId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("session_id, message_id and attitude are required"))
	}

	// 验证 attitude 值
	if !lo.Contains(validAttitudes, req.Attitude) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("attitude must be one of: %v", validAttitudes))
	}

	// 转换消息ID
	messageID, err := strconv.ParseUint(req.MessageId, 10, 64)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid message_id"))
	}

	// 查询消息，验证用户权限
	var dbMessage model.ChatMessage
	if err := db.GetDB().Where("id = ? AND user_id = ? AND session_id = ?", messageID, userID, req.SessionId).
		First(&dbMessage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("message not found or access denied"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 准备更新的 tags
	updatedTags := dbMessage.Tags

	// 移除已存在的 attitude 标签（up/down）
	filteredTags := lo.Filter(updatedTags, func(tag string, _ int) bool {
		return !lo.Contains(validAttitudes, tag)
	})

	// 添加新的 attitude 标签
	filteredTags = append(filteredTags, req.Attitude)
	dbMessage.Tags = filteredTags

	// 更新消息的 tags
	if err := db.GetDB().Model(&model.ChatMessage{}).
		Where("id = ?", messageID).
		Updates(&dbMessage).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	slog.Info("FeedbackToMessage success",
		"userId", userID,
		"messageId", messageID,
		"attitude", req.Attitude,
		"updatedTags", filteredTags)

	return connect.NewResponse(&message.FeedbackToMessageResponse{
		Success: true,
	}), nil
}

// 定义允许的 attitude 值
var validAttitudes = []string{"up", "down", ""}
