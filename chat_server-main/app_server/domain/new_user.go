package domain

import (
	"context"
	"log/slog"

	"app_server/domain/appconfig"
	"app_server/domain/exterr"
	"app_server/model"
	"app_server/pkg/ctxkv"
	"app_server/pkg/db"
	"app_server/pkg/idgen"

	"gorm.io/gorm"
)

func FindOrRegisterUser(ctx context.Context, user *model.User) (err error) {
	// 参数验证 - 微信登录需要external_id，手机登录需要phone
	if user.ExternalId == "" && user.Phone == "" {
		return exterr.ErrParamMissing
	}

	tx := db.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			slog.Error("panic occurred in FindOrRegisterUser", "panic", r)
			err = exterr.ErrInternal
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				slog.Error("commit transaction failed", "error", commitErr)
				err = exterr.ErrInternal
			}
		}
	}()

	// 检查用户是否已存在 - 根据external_id和phone选择查询条件
	var existingUser model.User
	var queryErr error
	if user.ExternalId != "" && user.Phone != "" {
		// 微信登录且有手机号的情况
		queryErr = tx.Where("external_id = ? OR phone = ?", user.ExternalId, user.Phone).First(&existingUser).Error
	} else if user.ExternalId != "" {
		// 只有微信登录的情况
		queryErr = tx.Where("external_id = ?", user.ExternalId).First(&existingUser).Error
	} else if user.Phone != "" {
		// 只有手机号登录的情况
		queryErr = tx.Where("phone = ?", user.Phone).First(&existingUser).Error
	}

	if queryErr != nil {
		if queryErr != gorm.ErrRecordNotFound {
			slog.Error("failed to query existing user", "error", queryErr, "external_id", user.ExternalId, "phone", user.Phone)
			return exterr.ErrInternal
		}
	}
	if existingUser.ID != 0 {
		// 用户已存在，更新传入的user对象的ID并返回
		*user = existingUser
		return nil
	}

	// 注册新用户
	err = RegisterNewUser(ctx, tx, user)
	if err != nil {
		slog.Error("failed to register new user", "error", err, "external_id", user.ExternalId)
		return err
	}

	return nil
}

func RegisterNewUser(ctx context.Context, tx *gorm.DB, user *model.User) (err error) {
	// 创建用户
	if err = tx.Create(&user).Error; err != nil {
		slog.Error("failed to create user", "error", err)
		return exterr.ErrInternal
	}

	// 创建用户基础Profile
	profile := &model.Profile{
		UserID: user.ID,
		Name:   user.Name,
		ImName: user.ImName,
	}
	if err = tx.Create(&profile).Error; err != nil {
		slog.Error("failed to create user profile", "error", err)
		return exterr.ErrInternal
	}

	// 更新用户的ProfileID
	user.ProfileID = profile.ID
	if err = tx.Save(&user).Error; err != nil {
		slog.Error("failed to update user profile_id", "error", err)
		return exterr.ErrInternal
	}

	return nil
}

// CopyDemoDataForNewUser 为新用户复制演示数据
func CopyDemoDataForNewUser(ctx context.Context, user *model.User, profile *model.Profile) (err error) {
	tx := db.GetDB().Begin()

	// 根据用户性别加载对应的演示数据
	demoJsonString := appconfig.LoadAppConfigByKeyVersion(ctx, "demo:"+profile.GetGender(), appconfig.Cond{
		Version: ctxkv.GetCtxKvString(ctx, "X-App-Version"),
		Env:     ctxkv.GetCtxKvString(ctx, "X-App-Env"),
	})

	var demoData DemoData
	if demoJsonString != "" {
		demoData = LoadDemoDataFromJsonString(demoJsonString)
	} else {
		demoData = LoadDemoData("feature_guide")
	}

	if len(demoData.DemoCases) == 0 {
		slog.Warn("no demo cases found in demo data")
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			slog.Error("panic occurred in FindOrRegisterUser", "panic", r)
			err = exterr.ErrInternal
		} else if err != nil {
			tx.Rollback()
			err = exterr.ErrInternal
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				slog.Error("commit transaction failed", "error", commitErr)
				err = exterr.ErrInternal
			}
		}
	}()

	profileOldToNew := make(map[uint]uint)
	chatOldToNew := make(map[uint]uint)
	messageOldToNew := make(map[uint]uint)

	// 收集所有要创建的数据
	var allProfiles []model.Profile
	var allChats []model.ChatSession
	var allMessages []model.ChatMessage

	// 第一步：为所有实体生成新ID并建立映射关系
	for _, demoCase := range demoData.DemoCases {
		profileOldToNew[demoCase.Profile.ID] = idgen.Uint()
		chatOldToNew[demoCase.ChatSession.ID] = idgen.Uint()
		for _, message := range demoCase.Messages {
			messageOldToNew[message.ID] = idgen.Uint()
		}
	}

	// 第二步：使用映射关系创建所有实体
	for _, demoCase := range demoData.DemoCases {
		// 创建Profile
		newProfile := demoCase.Profile
		newProfile.ID = profileOldToNew[demoCase.Profile.ID]
		newProfile.UserID = user.ID
		allProfiles = append(allProfiles, newProfile)

		// 创建ChatSession
		newChat := demoCase.ChatSession
		newChat.ID = chatOldToNew[demoCase.ChatSession.ID]
		newChat.UserID = user.ID
		newChat.ProfileID = profileOldToNew[demoCase.Profile.ID]
		allChats = append(allChats, newChat)

		// 创建消息
		for _, message := range demoCase.Messages {
			newMessage := message
			newMessage.ID = messageOldToNew[message.ID]
			newMessage.UserID = user.ID
			newMessage.SessionID = chatOldToNew[demoCase.ChatSession.ID]

			// 如果消息有ProfileID，映射到新的ProfileID
			newMessage.ProfileID = profileOldToNew[message.ProfileID]

			// 映射ParentID
			newMessage.ParentID = messageOldToNew[message.ParentID]
			// 为演示数据添加标识
			newMessage.Tags = append(message.Tags, "demo")

			allMessages = append(allMessages, newMessage)
		}
	}

	// 第三步：批量写入数据库
	if len(allProfiles) > 0 {
		if err = tx.Create(&allProfiles).Error; err != nil {
			slog.Error("failed to create profiles", "error", err)
			return
		}
	}

	if len(allChats) > 0 {
		if err = tx.Create(&allChats).Error; err != nil {
			slog.Error("failed to create chat sessions", "error", err)
			return
		}
	}

	if len(allMessages) > 0 {
		if err = tx.Create(&allMessages).Error; err != nil {
			slog.Error("failed to create messages", "error", err)
			return
		}
	}

	slog.Info("successfully copied demo data for new user",
		"user_id", user.ID,
		"cases_count", len(demoData.DemoCases),
		"profiles_count", len(allProfiles),
		"chats_count", len(allChats),
		"messages_count", len(allMessages))

	return nil
}
