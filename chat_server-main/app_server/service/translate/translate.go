package translate

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"app_server/model"
	"app_server/pkg/db"
	"app_server/pkg/idgen"
	"app_server/pkg/oai"
	"app_server/proto/translate"
	"app_server/service/auth"

	connect "connectrpc.com/connect"
)

type TranslateService struct{}

func (s *TranslateService) Translate(ctx context.Context, req *connect.Request[translate.TranslateRequest]) (*connect.Response[translate.TranslateResponse], error) {
	// 获取翻译请求参数
	content := req.Msg.Content
	from := TranslateTarget(req.Msg.From).CnString()
	to := TranslateTarget(req.Msg.To).CnString()
	history := req.Msg.History

	if from == "" || to == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("参数错误"))
	}

	// 获取提示模板
	var promptTemplate string //:= cfg.Viper().GetString("wechat.mp.translator.prompt")
	if promptTemplate == "" {
		promptTemplate = `你是一个帮助用户和领导沟通 揣摩领导潜台词 和领导话语含义的助手。
		翻译时请保持内容简短。
		
		请理解两人的对话上下文：
		%s
		
		将这句%s说的话翻译给%s听:
		%s`
	}

	// 构建翻译提示
	prompt := fmt.Sprintf(promptTemplate, history, from, to, content)

	slog.Info("translate", "prompt", prompt)

	// 使用新的 OAI 包调用 OpenAI
	translatedContent, err := oai.Get().CreateChatCompletionSimple(ctx, prompt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	slog.Info("translated", "plain", content, "translated", translatedContent)

	return connect.NewResponse(&translate.TranslateResponse{
		Content: translatedContent,
	}), nil
}

func (s *TranslateService) TranslateFriendMessage(ctx context.Context, req *connect.Request[translate.TranslateFriendMessageRequest]) (*connect.Response[translate.TranslateFriendMessageResponse], error) {
	userID := auth.GetUserID(ctx)
	friendMessageID := req.Msg.FriendMessageId
	from := req.Msg.From
	to := req.Msg.To
	if from == "" || to == "" || friendMessageID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("参数错误"))
	}

	// 先查一条消息
	var friendMessage model.ChatMessage
	if err := db.GetDB().Model(&model.ChatMessage{}).Where("id = ? AND user_id = ? AND msg_type = ?", friendMessageID, userID, "HISTORY").First(&friendMessage).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 再查消息前后一天时间的
	var friendMessages []model.ChatMessage
	if err := db.GetDB().Model(&model.ChatMessage{}).Where("session_id = ? AND user_id = ? AND msg_type = ? AND id BETWEEN ? AND ?",
		friendMessage.SessionID, userID, model.MessageTypeHistory,
		idgen.FromTime(friendMessage.CreatedAt.Add(-24*time.Hour)), idgen.FromTime(friendMessage.CreatedAt.Add(24*time.Hour))).
		Order("id ASC").
		Find(&friendMessages).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 将消息转换为上下文
	history := ""
	for _, msg := range friendMessages {
		history += fmt.Sprintf("%s: %s\n", msg.RoleCnString(), msg.Content)
	}

	translatedContent, err := s.Translate(ctx, &connect.Request[translate.TranslateRequest]{
		Msg: &translate.TranslateRequest{
			Content: friendMessage.Content,
			From:    from,
			To:      to,
			History: history,
		},
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 写入到 咨询消息
	consultMsg := model.ChatMessage{
		UserID:    friendMessage.UserID,
		SessionID: friendMessage.SessionID, // 使用原消息的 SessionID
		Content:   translatedContent.Msg.Content,
		ParentID:  friendMessage.ID,
		Role:      model.MessageRoleAI,
		MsgType:   model.MessageTypeTranslate,
		MsgAt:     time.Now(),
		Tags:      []string{to},
	}
	if err := db.GetDB().Create(&consultMsg).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&translate.TranslateFriendMessageResponse{
		Content: translatedContent.Msg.Content,
	}), nil
}

type TranslateTarget string

const (
	TranslateTargetMale   TranslateTarget = "MALE"
	TranslateTargetFemale TranslateTarget = "FEMALE"
)

func (t TranslateTarget) CnString() string {
	switch t {
	case TranslateTargetMale:
		return "男性"
	case TranslateTargetFemale:
		return "女性"
	}
	return ""
}

// buildProfileString 构建profile字符串，用于prompt模板
func buildProfileString(profile *model.Profile, who string) string {
	if profile == nil {
		return ""
	}

	profileStr := profile.FormatPropertyLinesString()
	if profileStr == "" {
		return ""
	}
	return fmt.Sprintf("\n**%s资料**\n%s\n", who, profileStr)
}

// TranslateV2 新版翻译接口，支持从config表加载prompt模板
func (s *TranslateService) TranslateV2(ctx context.Context, req *connect.Request[translate.TranslateV2Request]) (*connect.Response[translate.TranslateV2Response], error) {
	userID := auth.GetUserID(ctx)
	chatSessionID := req.Msg.ChatSessionId
	targetMessageID := req.Msg.TargetMessageId

	if chatSessionID == "" || targetMessageID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("chat_session_id和target_message_id不能为空"))
	}

	// 1. 查询目标消息
	var targetMessage model.ChatMessage
	if err := db.GetDB().Model(&model.ChatMessage{}).
		Where("id = ? AND user_id = ? AND session_id = ?", targetMessageID, userID, chatSessionID).
		First(&targetMessage).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("消息未找到"))
	}

	// 2. 查询chat_session获取friend_profile_id
	var chatSession model.ChatSession
	if err := db.GetDB().Model(&model.ChatSession{}).
		Where("id = ? AND user_id = ?", chatSessionID, userID).
		First(&chatSession).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("会话未找到"))
	}

	// 3. 根据消息角色确定prompt key
	var promptKey string
	switch targetMessage.Role {
	case model.MessageRoleSelf, model.MessageRoleUser:
		promptKey = "prompt:translate:to_friend"
	case model.MessageRoleFriend:
		promptKey = "prompt:translate:to_user"
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("不支持的消息角色: %s", targetMessage.Role))
	}

	// 4. 从config表加载prompt模板
	var config model.Config
	if err := db.GetDB().Model(&model.Config{}).
		Where("k = ?", promptKey).
		First(&config).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("未找到配置"))
	}

	// 5. 查询user profile和friend profile
	var user model.User
	if err := db.GetDB().Model(&model.User{}).
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("用户未找到"))
	}
	var userProfile model.Profile
	if err := db.GetDB().Model(&model.Profile{}).
		Where("user_id = ? AND id = ?", userID, user.ProfileID).
		First(&userProfile).Error; err != nil {
		slog.Warn("未找到用户Profile", "userID", userID, "error", err)
	}

	var friendProfile model.Profile
	if chatSession.ProfileID > 0 {
		if err := db.GetDB().Model(&model.Profile{}).
			Where("id = ?", chatSession.ProfileID).
			First(&friendProfile).Error; err != nil {
			slog.Warn("未找到好友Profile", "profileID", chatSession.ProfileID, "error", err)
		}
	}

	// 6. 查询前后24小时的对话历史
	var chatMessages []model.ChatMessage
	if err := db.GetDB().Model(&model.ChatMessage{}).
		Where("session_id = ? AND user_id = ? AND msg_type = ? AND id BETWEEN ? AND ?",
			targetMessage.SessionID, userID, model.MessageTypeHistory,
			idgen.FromTime(targetMessage.CreatedAt.Add(-24*time.Hour)),
			idgen.FromTime(targetMessage.CreatedAt.Add(24*time.Hour))).
		Order("id ASC").
		Find(&chatMessages).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 7. 构建聊天上下文
	var chatContext strings.Builder
	for _, msg := range chatMessages {
		chatContext.WriteString(msg.HistoryCnString() + "\n")
	}

	// 8. 替换prompt模板中的变量
	prompt := config.Value

	// 构建user_profile和friend_profile字符串
	userProfileStr := buildProfileString(&userProfile, "用户")
	// 动态获取对方名称，优先使用 profile 名称，否则使用默认值"对方"
	friendName := "对方"
	if friendProfile.Name != "" {
		friendName = friendProfile.Name
	}
	friendProfileStr := buildProfileString(&friendProfile, friendName)

	// 替换模板变量
	prompt = strings.ReplaceAll(prompt, "{{user_profile}}", userProfileStr)
	prompt = strings.ReplaceAll(prompt, "{{friend_profile}}", friendProfileStr)
	prompt = strings.ReplaceAll(prompt, "{{chat_context}}", chatContext.String())
	prompt = strings.ReplaceAll(prompt, "{{src_message}}", targetMessage.HistoryCnString())

	// 9. 调用AI API进行翻译
	translatedContent, err := oai.Get().CreateChatCompletionSimple(ctx, prompt)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	slog.Info("TranslateV2 completed", "original", targetMessage.Content, "translated", translatedContent)

	// 11. 创建新的咨询消息
	consultMsg := model.ChatMessage{
		UserID:    targetMessage.UserID,
		SessionID: targetMessage.SessionID,
		Content:   translatedContent,
		ParentID:  targetMessage.ID,
		Role:      model.MessageRoleAI,
		MsgType:   model.MessageTypeTranslate,
		MsgAt:     time.Now(),
	}
	if err := db.GetDB().Create(&consultMsg).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&translate.TranslateV2Response{
		NewMessageId: fmt.Sprintf("%d", consultMsg.ID),
		Content:      translatedContent,
	}), nil
}
