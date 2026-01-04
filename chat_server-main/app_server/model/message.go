package model

import (
	"fmt"
	"strconv"
	"time"

	"app_server/pkg/fn"
	"app_server/proto/message"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	SessionID uint      `json:"session_id"`
	ParentID  uint      `json:"parent_id"`
	ProfileID uint      `json:"profile_id"`
	Role      string    `json:"role"`
	MsgType   string    `json:"msg_type"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags" gorm:"serializer:json"`
	MsgAt     time.Time `json:"msg_at"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}

func (m ChatMessage) ToProto() *message.ChatMessage {
	return &message.ChatMessage{
		Id:        strconv.Itoa(int(m.ID)),
		Content:   m.Content,
		UserId:    fn.Itoa(m.UserID),
		Role:      m.Role,
		MsgAt:     timestamppb.New(m.MsgAt),
		MsgType:   m.MsgType,
		SessionId: fn.Itoa(m.SessionID),
		ParentId:  fn.Itoa(m.ParentID),
		Tags:      m.Tags,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

// ToConsultProto 转换为旧的 ConsultMessage 格式（向后兼容）
func (m ChatMessage) ToConsultProto() *message.ConsultMessage {
	return &message.ConsultMessage{
		Id:        strconv.Itoa(int(m.ID)),
		Content:   m.Content,
		UserId:    fn.Itoa(m.UserID),
		Role:      m.Role,
		MsgAt:     timestamppb.New(m.MsgAt),
		MsgType:   m.MsgType,
		SessionId: fn.Itoa(m.SessionID),
		ParentId:  fn.Itoa(m.ParentID),
		ProfileId: fn.Itoa(m.ProfileID),
		Tags:      m.Tags,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

// ChatMessageFromProto 从新的 ChatMessage proto 转换
func ChatMessageFromProto(protoMessage *message.ChatMessage) *ChatMessage {
	m := &ChatMessage{
		Content:   protoMessage.Content,
		UserID:    fn.Atoi[uint](protoMessage.UserId),
		Role:      protoMessage.Role,
		MsgType:   protoMessage.MsgType,
		ParentID:  fn.Atoi[uint](protoMessage.ParentId),
		SessionID: fn.Atoi[uint](protoMessage.SessionId),
		Tags:      protoMessage.Tags,
	}

	if protoMessage.Id != "" {
		m.ID = uint(fn.Atoi[int](protoMessage.Id))
	}

	if protoMessage.MsgAt != nil {
		m.MsgAt = protoMessage.MsgAt.AsTime()
	}

	if protoMessage.CreatedAt != nil {
		m.CreatedAt = protoMessage.CreatedAt.AsTime()
	}

	if protoMessage.UpdatedAt != nil {
		m.UpdatedAt = protoMessage.UpdatedAt.AsTime()
	}

	return m
}

// ConsultMessageFromProto 从旧的 ConsultMessage proto 转换（向后兼容）
func ConsultMessageFromProto(protoMessage *message.ConsultMessage) *ChatMessage {
	return &ChatMessage{
		Model: gorm.Model{
			ID: uint(fn.Atoi[int](protoMessage.Id)),
		},
		Content:   protoMessage.Content,
		UserID:    fn.Atoi[uint](protoMessage.UserId),
		Role:      protoMessage.Role,
		MsgAt:     protoMessage.MsgAt.AsTime(),
		MsgType:   protoMessage.MsgType,
		ParentID:  fn.Atoi[uint](protoMessage.ParentId),
		SessionID: fn.Atoi[uint](protoMessage.SessionId),
		ProfileID: fn.Atoi[uint](protoMessage.ProfileId),
		Tags:      protoMessage.Tags,
	}
}

func (m ChatMessage) RoleCnString() string {
	switch m.Role {
	case MessageRoleAI:
		return "AI"
	case MessageRoleSelf:
		return "用户"
	case MessageRoleFriend:
		return "朋友"
	case MessageRoleUser:
		return "用户"
	default:
		return m.Role
	}
}

func (m ChatMessage) TypeCnString() string {
	switch m.MsgType {
	case MessageTypeHistory:
		return "聊天历史"
	case MessageTypeTranslate:
		return "AI解读"
	case MessageTypeConsult:
		return "咨询"
	default:
		return m.MsgType
	}
}

func (m ChatMessage) HistoryCnString() string {
	switch m.MsgType {
	case MessageTypeHistory:
		return fmt.Sprintf("%s:%s", m.RoleCnString(), m.Content)
	case MessageTypeTranslate:
		return fmt.Sprintf("<%s>%s</%s>", m.TypeCnString(), m.Content, m.TypeCnString())
	default:
		return fmt.Sprintf("%s:%s", m.RoleCnString(), m.Content)
	}
}

const (
	MessageRoleSelf   = "SELF"
	MessageRoleFriend = "FRIEND"
	MessageRoleAI     = "AI"
	MessageRoleUser   = "USER"
)

const (
	MessageTypeHistory   = "HISTORY"
	MessageTypeTranslate = "TRANSLATE"
	MessageTypeConsult   = "CONSULT"
)
