package message

import (
	"context"
	"log/slog"
	"strconv"

	"app_server/model"
	"app_server/pkg/db"
	"app_server/pkg/fn"
	"app_server/proto/message"
	"app_server/service/auth"

	connect "connectrpc.com/connect"
	"github.com/samber/lo/mutable"
)

type ConsultMessageService struct{}

func (s *ConsultMessageService) ListConsultMessages(ctx context.Context, connectReq *connect.Request[message.ListConsultMessagesRequest]) (*connect.Response[message.ListConsultMessagesResponse], error) {
	req := connectReq.Msg

	// 构建查询条件
	query := db.GetDB().Model(&model.ChatMessage{})

	// 添加用户ID
	query = query.Where("user_id = ?", auth.GetUserID(ctx))

	// 处理消息类型过滤
	if req.MsgType != "" {
		query = query.Where("msg_type = ?", req.MsgType)
	}
	// 根据ID列表过滤
	if len(req.Ids) > 0 {
		query = query.Where("id IN ?", req.Ids)
	}
	// 根据会话ID列表过滤
	if len(req.SessionIds) > 0 {
		query = query.Where("session_id IN ?", req.SessionIds)
	}
	// 根据父消息ID列表过滤
	if len(req.ParentIds) > 0 {
		query = query.Where("parent_id IN ?", req.ParentIds)
	}

	// 处理分页
	var lastID int
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 100 // 默认每页100条
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
		slog.Error("list consult messages error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 反转消息顺序，让较早的消息排在前面
	mutable.Reverse(dbMessages)

	// 处理结果
	var nextPageToken string
	if len(dbMessages) == pageSize {
		nextPageToken = strconv.Itoa(int(dbMessages[0].ID))
	}

	return connect.NewResponse(&message.ListConsultMessagesResponse{
		Messages:      fn.Map(dbMessages, model.ChatMessage.ToConsultProto),
		NextPageToken: nextPageToken,
	}), nil
}

func (s *ConsultMessageService) SendConsultMessage(ctx context.Context, req *connect.Request[message.SendConsultMessageRequest]) (*connect.Response[message.SendConsultMessageResponse], error) {
	return nil, nil
}

func (s *ConsultMessageService) RecallConsultMessage(ctx context.Context, req *connect.Request[message.RecallConsultMessageRequest]) (*connect.Response[message.RecallConsultMessageResponse], error) {
	return nil, nil
}
