package model

import (
	"gorm.io/gorm"
)

// Config 配置表
type Config struct {
	gorm.Model
	Key      string `gorm:"column:k"`
	Value    string
	App      string
	Version  string
	Platform string
	Env      string
}

// TableName 指定表名
func (Config) TableName() string {
	return "config"
}
