package message

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"app_server/model"
	"app_server/pkg/aiapi"
	"app_server/pkg/db"
	"app_server/pkg/fn"
	"app_server/pkg/ossc"
	"app_server/proto/message"
	"app_server/service/auth"

	connect "connectrpc.com/connect"
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

type FriendMessageService struct{}

// 将ConsultMessage转换为适合展示的形式
func adaptConsultMessageForFriend(msg *model.ChatMessage) *message.ConsultMessage {
	return msg.ToConsultProto()
}

func (s *FriendMessageService) ListFriendMessages(ctx context.Context, connectReq *connect.Request[message.ListFriendMessagesRequest]) (*connect.Response[message.ListFriendMessagesResponse], error) {
	req := connectReq.Msg

	// 构建查询条件
	query := db.GetDB().Model(&model.ChatMessage{})

	// 添加用户ID
	query = query.Where("user_id = ?", auth.GetUserID(ctx))

	// 设置消息类型为HISTORY
	query = query.Where("msg_type = ?", "HISTORY")

	// 添加朋友ID
	if req.ProfileId != "" {
		query = query.Where("profile_id = ?", req.ProfileId)
	}

	// 处理分页
	var lastID int
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20 // 默认每页20条
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
		slog.Error("list friend messages error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 反转消息顺序
	mutable.Reverse(dbMessages)

	// 处理结果
	var nextPageToken string
	if len(dbMessages) == pageSize {
		nextPageToken = strconv.Itoa(int(dbMessages[0].ID))
	}

	// 使用map转换数据库模型到proto模型
	protoMessages := fn.Map(dbMessages, model.ChatMessage.ToConsultProto)

	return connect.NewResponse(&message.ListFriendMessagesResponse{
		Messages:      protoMessages,
		NextPageToken: nextPageToken,
	}), nil
}

func (s *FriendMessageService) CreateFriendMessage(ctx context.Context, connectReq *connect.Request[message.CreateFriendMessageRequest]) (*connect.Response[message.CreateFriendMessageResponse], error) {
	req := connectReq.Msg

	userID := auth.GetUserID(ctx)

	// 如果没有消息需要创建，直接返回
	if len(req.Messages) == 0 {
		return connect.NewResponse(&message.CreateFriendMessageResponse{}), nil
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
		// 创建一个新的ConsultMessage
		dbMsg := &model.ChatMessage{
			UserID:    userID,
			SessionID: fn.Atoi[uint](protoMsg.SessionId),
			Role:      protoMsg.Role,
			MsgAt:     time.Now(),
			MsgType:   "HISTORY",
			Content:   protoMsg.Content,
			ProfileID: fn.Atoi[uint](protoMsg.ProfileId),
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
	protoMessages := fn.Map(dbMessages, model.ChatMessage.ToConsultProto)

	return connect.NewResponse(&message.CreateFriendMessageResponse{
		Messages: protoMessages,
	}), nil
}

func (s *FriendMessageService) UpdateFriendMessage(ctx context.Context, connectReq *connect.Request[message.UpdateFriendMessageRequest]) (*connect.Response[message.UpdateFriendMessageResponse], error) {
	req := connectReq.Msg

	userID := auth.GetUserID(ctx)

	// 如果没有消息需要更新，直接返回
	if len(req.Messages) == 0 {
		return connect.NewResponse(&message.UpdateFriendMessageResponse{}), nil
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

		// 只更新非零值字段
		updates := map[string]interface{}{
			"content": protoMsg.Content,
		}

		if protoMsg.Role != "" {
			updates["role"] = protoMsg.Role
		}

		// 更新消息
		var dbMessage model.ChatMessage
		if err := tx.Model(&model.ChatMessage{}).Where("id = ? AND user_id = ? AND msg_type = ?", id, userID, "HISTORY").Updates(updates).First(&dbMessage, id).Error; err != nil {
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
	protoMessages := fn.Map(updatedMessages, model.ChatMessage.ToConsultProto)

	return connect.NewResponse(&message.UpdateFriendMessageResponse{
		Messages: protoMessages,
	}), nil
}

func (s *FriendMessageService) DeleteFriendMessage(ctx context.Context, connectReq *connect.Request[message.DeleteFriendMessageRequest]) (*connect.Response[message.DeleteFriendMessageResponse], error) {
	req := connectReq.Msg

	userID := auth.GetUserID(ctx)

	// 如果没有消息需要删除，直接返回
	if len(req.Ids) == 0 {
		return connect.NewResponse(&message.DeleteFriendMessageResponse{}), nil
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
		return connect.NewResponse(&message.DeleteFriendMessageResponse{}), nil
	}

	// 删除消息
	if err := db.GetDB().Where("id IN ? AND user_id = ? AND msg_type = ?", ids, userID, "HISTORY").Delete(&model.ChatMessage{}).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&message.DeleteFriendMessageResponse{}), nil
}

func (s *FriendMessageService) ParseImageMessages(ctx context.Context, connectReq *connect.Request[message.ParseImageMessagesRequest]) (*connect.Response[message.ParseImageMessagesResponse], error) {
	req := connectReq.Msg
	userID := auth.GetUserID(ctx)

	// 验证参数
	sessionID := fn.Atoi[uint](req.SessionId)
	if sessionID == 0 || req.ImageUrl == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	// 获取会话信息以确定profileID
	var chatSession model.ChatSession
	if err := db.GetDB().Where("id = ? AND user_id = ?", sessionID, userID).First(&chatSession).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	profileID := chatSession.ProfileID

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
			ProfileID: profileID,
			Role:      role,
			MsgAt:     now,
			MsgType:   "HISTORY",
			Content:   content,
			Tags:      []string{"friend_message"},
		}

		parsedMessages = append(parsedMessages, msg)
	}

	// 如果没有解析出有效消息
	if len(parsedMessages) == 0 {
		return connect.NewResponse(&message.ParseImageMessagesResponse{
			Message: "未能从图片中解析出有效聊天记录",
		}), nil
	}

	// 从数据库查询最近条消息 避免重复
	var lastMessages []model.ChatMessage
	tx.Where("user_id = ? AND profile_id = ? AND msg_type = ?", userID, profileID, "HISTORY").Order("id DESC").Limit(20).Find(&lastMessages)

	// 检查是否存在重复消息
	dbMessages := fn.Filter(parsedMessages, func(msg model.ChatMessage) bool {
		return !lo.ContainsBy(lastMessages, func(lastMsg model.ChatMessage) bool {
			return lastMsg.Content == msg.Content && lastMsg.Role == msg.Role
		})
	})

	if len(dbMessages) == 0 {
		return connect.NewResponse(&message.ParseImageMessagesResponse{
			Message: "解析消息全部重复 没有新增消息",
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
