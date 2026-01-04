package model

import (
	"strconv"
	"time"

	"app_server/pkg/fn"
	"app_server/proto/message"
)

// FriendMessageToConsultMessage 将历史消息转换为ConsultMessage
func FriendMessageToConsultMessage(profileID uint, role string, content string) *ChatMessage {
	now := time.Now()
	return &ChatMessage{
		UserID:    0, // 这个字段需要在服务层中设置
		SessionID: 0, // 历史消息没有会话ID
		ParentID:  0, // 历史消息没有父消息ID
		Role:      role,
		MsgAt:     now,
		MsgType:   "HISTORY", // 历史消息类型固定为HISTORY
		Content:   content,
		Tags:      []string{"friend_message", "profile_" + strconv.Itoa(int(profileID))},
	}
}

// FriendMessageFromProto 将proto模型转换为数据库模型
func FriendMessageFromProto(protoMessage *message.ConsultMessage) *ChatMessage {
	msg := ConsultMessageFromProto(protoMessage)
	msg.MsgType = "HISTORY" // 强制设置类型为HISTORY
	return msg
}

// ProtoToConsultMessage 将proto消息转换为ConsultMessage
func ProtoToConsultMessage(protoMessage *message.ConsultMessage) *ChatMessage {
	msgAt := time.Now()
	if protoMessage.MsgAt != nil && protoMessage.MsgAt.IsValid() {
		msgAt = protoMessage.MsgAt.AsTime()
	}

	return &ChatMessage{
		UserID:    fn.Atoi[uint](protoMessage.UserId),
		SessionID: fn.Atoi[uint](protoMessage.SessionId),
		ParentID:  fn.Atoi[uint](protoMessage.ParentId),
		Role:      protoMessage.Role,
		MsgAt:     msgAt,
		MsgType:   "HISTORY", // 固定设置为HISTORY
		Content:   protoMessage.Content,
		Tags:      protoMessage.Tags,
	}
}
