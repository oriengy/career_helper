package model

import (
	"app_server/pkg/fn"
	"app_server/proto/chat"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ChatSession struct {
	gorm.Model
	Name      string `json:"name"`
	UserID    uint   `json:"user_id"`
	ProfileID uint   `json:"profile_id"`
	Avatar    string `json:"avatar"`
}

func (c ChatSession) ToProto() *chat.ChatSession {
	return &chat.ChatSession{
		Id:        fn.Itoa(c.ID),
		Name:      c.Name,
		ProfileId: lo.Ternary(c.ProfileID == 0, "", fn.Itoa(c.ProfileID)),
		UserId:    fn.Itoa(c.UserID),
		Avatar:    c.Avatar,
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}
}

func ChatFromProto(protoChat *chat.ChatSession) *ChatSession {
	return &ChatSession{
		Model: gorm.Model{
			ID:        fn.Atoi[uint](protoChat.Id),
			CreatedAt: protoChat.CreatedAt.AsTime(),
			UpdatedAt: protoChat.UpdatedAt.AsTime(),
		},
		Name:      protoChat.Name,
		UserID:    fn.Atoi[uint](protoChat.UserId),
		ProfileID: fn.Atoi[uint](protoChat.ProfileId),
		Avatar:    protoChat.Avatar,
	}
}
