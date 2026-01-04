package user

import (
	"context"
	"log/slog"

	"app_server/domain"
	"app_server/model"
	"app_server/pkg/db"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

func WxLoginFirstOrCreate(ctx context.Context, openID, unionID string) (user *model.User, err error) {
	externalID := lo.Ternary(unionID != "", unionID, openID)

	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// 查询用户
		if err := tx.Where("external_id = ?", externalID).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// 如果用户不存在，则创建用户
		if user.ID == 0 {
			user = &model.User{
				ExternalId: externalID,
			}
			err = domain.FindOrRegisterUser(ctx, user)
		}
		return nil
	})
	slog.Info("wx login first or create", "external_id", externalID, "user_id", user.ID, "err", err)

	return user, err
}

// PhoneLoginFirstOrCreate 根据手机号查找或创建用户
func PhoneLoginFirstOrCreate(ctx context.Context, phone string) (user *model.User, err error) {
	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		// 查询用户
		if err := tx.Where("phone = ?", phone).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		// 如果用户不存在，则创建用户
		if user.ID == 0 {
			user = &model.User{
				Phone: phone,
			}
			err = domain.FindOrRegisterUser(ctx, user)
		}
		return nil
	})
	slog.Info("phone login first or create", "phone", phone, "user_id", user.ID, "err", err)

	return user, err
}
