package model

import (
	"strconv"

	"app_server/pkg/fn"
	"app_server/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string
	Content string
	UserID  string
}

func (Post) TableName() string {
	return "post"
}

func (p Post) ToProto() *proto.Post {
	return &proto.Post{
		Id:        strconv.Itoa(int(p.ID)),
		Title:     p.Title,
		Content:   p.Content,
		UserId:    p.UserID,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}

func PostFromProto(protoPost *proto.Post) *Post {
	return &Post{
		Model: gorm.Model{
			ID:        uint(fn.Atoi[int](protoPost.Id)),
			CreatedAt: protoPost.CreatedAt.AsTime(),
			UpdatedAt: protoPost.UpdatedAt.AsTime(),
		},
		Title:   protoPost.Title,
		Content: protoPost.Content,
		UserID:  protoPost.UserId,
	}
}
