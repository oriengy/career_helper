package model

import (
	"app_server/pkg/fn"
	"app_server/proto/user"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string
	ImName     string
	ExternalId string
	Phone      string
	Avatar     string
	ProfileID  uint
}

func (User) TableName() string {
	return "user"
}

func (u User) ToProto() *user.User {
	return &user.User{
		Id:         fn.Itoa(u.ID),
		Name:       u.Name,
		ImName:     u.ImName,
		ExternalId: u.ExternalId,
		Phone:      u.Phone,
		Avatar:     u.Avatar,
		CreatedAt:  timestamppb.New(u.CreatedAt),
		UpdatedAt:  timestamppb.New(u.UpdatedAt),
		ProfileId:  fn.Itoa(u.ProfileID),
	}
}

func FromProto(protoUser *user.User) *User {
	return &User{
		Model: gorm.Model{
			ID:        uint(fn.Atoi[int](protoUser.Id)),
			CreatedAt: protoUser.CreatedAt.AsTime(),
			UpdatedAt: protoUser.UpdatedAt.AsTime(),
		},
		Name:       protoUser.Name,
		ImName:     protoUser.ImName,
		ExternalId: protoUser.ExternalId,
		Phone:      protoUser.Phone,
		Avatar:     protoUser.Avatar,
		ProfileID:  fn.Atoi[uint](protoUser.ProfileId),
	}
}
