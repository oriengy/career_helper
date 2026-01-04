package chat

import (
	"context"
	"log/slog"
	"strconv"

	"app_server/domain"
	"app_server/model"
	"app_server/pkg/db"
	"app_server/pkg/fn"
	"app_server/proto/chat"
	"app_server/service/auth"

	connect "connectrpc.com/connect"
	"github.com/samber/lo"
)

type ChatService struct{}

func (s *ChatService) ListChatSessions(ctx context.Context, req *connect.Request[chat.ListChatSessionsRequest]) (*connect.Response[chat.ListChatSessionsResponse], error) {
	// 构建查询条件
	query := db.GetDB().Model(&model.ChatSession{})

	// 添加用户ID条件
	query = query.Where("user_id = ?", auth.GetUserID(ctx))

	pageToken := req.Msg.PageToken
	pageSize := req.Msg.PageSize
	if pageSize == 0 {
		pageSize = 20
	}

	if pageToken != "" {
		query = query.Where("id > ?", pageToken)
	}

	// 查询数据
	var chatSessions []model.ChatSession
	if err := query.Order("id DESC").Find(&chatSessions).Error; err != nil {
		slog.Error("list chat sessions error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 查询好友信息
	profileIds := make([]uint, 0)
	for _, chatSession := range chatSessions {
		if chatSession.ProfileID != 0 {
			profileIds = append(profileIds, chatSession.ProfileID)
		}
	}
	if len(profileIds) > 0 {
		profiles := make([]model.Profile, 0)
		db.GetDB().Where("id IN (?)", profileIds).Find(&profiles)
		profilesMap := lo.SliceToMap(profiles, func(profile model.Profile) (uint, model.Profile) {
			return profile.ID, profile
		})
		for i, chatSession := range chatSessions {
			if profile, ok := profilesMap[chatSession.ProfileID]; ok {
				chatSessions[i].Name = profile.Name
				chatSessions[i].Avatar = profile.Avatar
			}
		}
	}

	return connect.NewResponse(&chat.ListChatSessionsResponse{
		Data: fn.Map(chatSessions, model.ChatSession.ToProto),
	}), nil
}

func (s *ChatService) CreateChatSession(ctx context.Context, req *connect.Request[chat.CreateChatSessionRequest]) (*connect.Response[chat.CreateChatSessionResponse], error) {
	userID := auth.GetUserID(ctx)
	msg := req.Msg

	chatSession, err := domain.CreateChatSession(ctx, userID, &model.Profile{
		UserID: userID,
		Name:   msg.Profile.GetName(),
		ImName: msg.Profile.GetImName(),
		Avatar: msg.Profile.GetAvatar(),
		Gender: msg.Profile.GetGender(),
	})

	if err != nil {
		slog.Error("create chat session transaction error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&chat.CreateChatSessionResponse{
		ChatSession: chatSession.ToProto(),
	}), nil
}

func (s *ChatService) DeleteChatSession(ctx context.Context, req *connect.Request[chat.DeleteChatSessionRequest]) (*connect.Response[chat.DeleteChatSessionResponse], error) {
	msg := req.Msg

	// 验证会话ID
	id, err := strconv.Atoi(msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// 删除会话
	if err := db.GetDB().Where("id = ? AND user_id = ?", id, auth.GetUserID(ctx)).Delete(&model.ChatSession{}).Error; err != nil {
		slog.Error("delete chat session error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&chat.DeleteChatSessionResponse{}), nil
}

func (s *ChatService) UpdateChatSession(ctx context.Context, req *connect.Request[chat.UpdateChatSessionRequest]) (*connect.Response[chat.UpdateChatSessionResponse], error) {
	msg := req.Msg

	// 验证会话ID
	id, err := strconv.Atoi(msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// 查询原会话
	var chatSession model.ChatSession
	if err := db.GetDB().Where("id = ? AND user_id = ?", id, auth.GetUserID(ctx)).First(&chatSession).Error; err != nil {
		slog.Error("get chat session error", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// 更新字段
	updates := map[string]any{}
	if msg.Name != "" {
		updates["name"] = msg.Name
	}
	if msg.Avatar != "" {
		updates["avatar"] = msg.Avatar
	}

	// 更新数据库
	if err := db.GetDB().Model(&chatSession).Updates(updates).Error; err != nil {
		slog.Error("update chat session error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的数据
	if err := db.GetDB().First(&chatSession, chatSession.ID).Error; err != nil {
		slog.Error("get updated chat session error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&chat.UpdateChatSessionResponse{
		ChatSession: chatSession.ToProto(),
	}), nil
}
