package model

import (
	"fmt"
	"strconv"
	"strings"

	"app_server/pkg/fn"
	// 使用新生成的profile包
	"app_server/proto/profile"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ProfileData 存储个人资料的详细信息
type ProfileData struct {
	Description string `json:"description,omitempty"`
	// 其他需要的字段
}

// Profile 个人资料数据模型
type Profile struct {
	gorm.Model
	UserID       uint       `json:"user_id"`
	Name         string     `json:"name"`
	ImName       string     `json:"im_name"`
	Avatar       string     `json:"avatar"`
	AvatarFileID uint       `json:"avatar_file_id" gorm:"column:avatar_file_id;comment:头像文件ID"`
	Age          int        `json:"age"`
	Gender       string     `json:"gender"`
	Prompt       string     `json:"prompt"` // 提示词 不对外暴露
	Intro        string     `json:"desc"`
	Custom       []Property `json:"custom" gorm:"serializer:json"`
}

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (Profile) TableName() string {
	return "profile"
}

// ToProto 将数据库模型转换为proto模型
func (m Profile) ToProto() *profile.Profile {
	proto := &profile.Profile{
		Id:        strconv.Itoa(int(m.ID)),
		UserId:    fn.Itoa(m.UserID),
		Name:      m.Name,
		ImName:    m.ImName,
		Avatar:    m.Avatar,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
		Age:       int32(m.Age),
		Gender:    m.Gender,
		Intro:     m.Intro,
	}

	if m.AvatarFileID > 0 {
		proto.AvatarFileId = fn.Itoa(m.AvatarFileID)
	}

	// 转换自定义属性
	for _, prop := range m.Custom {
		proto.Custom = append(proto.Custom, &profile.Property{
			Name:  prop.Name,
			Value: prop.Value,
		})
	}

	return proto
}

// ProfileFromProto 将proto模型转换为数据库模型
func ProfileFromProto(protoProfile *profile.Profile) *Profile {
	p := &Profile{
		Model: gorm.Model{
			ID: uint(fn.Atoi[int](protoProfile.Id)),
		},
		UserID: fn.Atoi[uint](protoProfile.UserId),
		Name:   protoProfile.Name,
		ImName: protoProfile.ImName,
		Avatar: protoProfile.Avatar,
		Age:    int(protoProfile.Age),
		Gender: protoProfile.Gender,
		Intro:  protoProfile.Intro,
	}

	if protoProfile.AvatarFileId != "" {
		p.AvatarFileID = fn.Atoi[uint](protoProfile.AvatarFileId)
	}

	// 转换自定义属性
	for _, prop := range protoProfile.Custom {
		p.Custom = append(p.Custom, Property{
			Name:  prop.Name,
			Value: prop.Value,
		})
	}

	return p
}

func (p *Profile) GetGender() string {
	if p == nil {
		return ""
	}
	return p.Gender
}

func (p *Profile) GetGenderCn() string {
	if p == nil {
		return ""
	}
	switch p.Gender {
	case "male":
		return "男性"
	case "female":
		return "女性"
	}
	return ""
}

func (p *Profile) FormatPropertyLines() []string {
	if p == nil {
		return nil
	}
	var lines []string
	var properties []Property

	// 添加性别信息
	if p.Gender != "" {
		properties = append(properties, Property{Name: "性别", Value: p.GetGenderCn()})
	}

	// 添加年龄信息
	if p.Age > 0 {
		properties = append(properties, Property{Name: "年龄", Value: fmt.Sprintf("%d岁", p.Age)})
	}

	// 可以根据需要添加更多信息
	if p.Intro != "" {
		properties = append(properties, Property{Name: "简介", Value: p.Intro})
	}

	// 添加自定义属性
	for _, property := range p.Custom {
		properties = append(properties, Property{Name: property.Name, Value: property.Value})
	}
	for _, property := range properties {
		lines = append(lines, fmt.Sprintf("%s:%s", property.Name, property.Value))
	}
	return lines
}

func (p *Profile) FormatPropertyLinesString() string {
	return strings.Join(p.FormatPropertyLines(), "\n")
}
