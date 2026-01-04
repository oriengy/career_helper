package domain

import (
	"app_server/model"
	"app_server/pkg/db"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

func CreateChatSession(ctx context.Context, userID uint, profile *model.Profile) (chatSession *model.ChatSession, err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// 创建新的聊天会话
		chatSession = &model.ChatSession{
			Name:   profile.Name,
			UserID: userID,
			Avatar: profile.Avatar,
		}
		// 快速创建个人资料
		if profile != nil {
			if err := tx.Create(profile).Error; err != nil {
				slog.Error("create profile error", "error", err)
				return err
			}
			chatSession.ProfileID = profile.ID
		}
		// 保存到数据库
		if err := tx.Create(chatSession).Error; err != nil {
			slog.Error("create chat session error", "error", err)
			return err
		}
		return nil
	})
	return chatSession, err
}
