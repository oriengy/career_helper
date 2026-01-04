package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// UserFile 用户文件表
type UserFile struct {
	gorm.Model
	UserID       uint         `json:"user_id" gorm:"index;not null;comment:用户ID"`
	OriginalName string       `json:"original_name" gorm:"not null;comment:原始文件名"`
	FileSize     int64        `json:"file_size" gorm:"not null;comment:文件大小(字节)"`
	FileType     string       `json:"file_type" gorm:"comment:文件MIME类型"`
	FileExt      string       `json:"file_ext" gorm:"comment:文件扩展名"`
	OssKey       string       `json:"oss_key" gorm:"not null;comment:OSS存储路径"`
	PublicURL    string       `json:"public_url" gorm:"type:text;comment:OSS公网访问URL"`
	PublicExpire sql.NullTime `json:"public_expire" gorm:"comment:公网访问有效期"`
	FileHash     string       `json:"file_hash" gorm:"index;comment:文件哈希值(用于去重)"`
	Status       int          `json:"status" gorm:"default:1;comment:状态:1正常,2已删除"`
	UsageType    string       `json:"usage_type" gorm:"comment:用途:avatar,chat_image等"`
	ExpiresAt    sql.NullTime `json:"expires_at" gorm:"index;comment:过期时间"`
}

func (UserFile) TableName() string {
	return "user_file"
}

// 文件状态常量
const (
	FileStatusNormal  = 1
	FileStatusDeleted = 2
)

// 文件用途类型常量
const (
	UsageTypeAvatar     = "avatar"      // 头像
	UsageTypeChatImage  = "chat_image"  // 聊天图片
	UsageTypeTempUpload = "temp_upload" // 临时上传
)


// GetExpirationTime 根据用途类型获取过期时间
func GetExpirationTime(usageType string) *time.Time {
	now := time.Now()
	switch usageType {
	case UsageTypeTempUpload:
		// 临时文件：24小时后过期
		expTime := now.Add(24 * time.Hour)
		return &expTime
	case UsageTypeChatImage:
		// 聊天图片：30天后过期
		expTime := now.Add(30 * 24 * time.Hour)
		return &expTime
	case UsageTypeAvatar:
		// 头像：永不过期
		return nil
	default:
		// 默认：7天后过期
		expTime := now.Add(7 * 24 * time.Hour)
		return &expTime
	}
}