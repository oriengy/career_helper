package profile

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"app_server/domain"
	"app_server/model"
	"app_server/pkg/db"
	"app_server/pkg/fn"
	"app_server/pkg/ossc"
	"app_server/proto/profile"
	"app_server/service/auth"
	"app_server/service/file"

	"connectrpc.com/connect"
)

type ProfileService struct{}

// 生成头像的签名URL，有效期10年
func generateAvatarSignedURL(avatarPath string) string {
	if avatarPath == "" {
		return ""
	}

	// 10年的秒数
	publicURL, err := ossc.GetPublic().UserFileBucket().SignURL(avatarPath, "GET", 10*365*24*60*60)
	if err != nil {
		slog.Warn("generate public avatar URL error", "error", err, "avatar", avatarPath)
		return avatarPath // 如果生成失败，返回原始路径
	}

	return publicURL
}

func (s *ProfileService) ListProfiles(ctx context.Context, req *connect.Request[profile.ListProfilesRequest]) (*connect.Response[profile.ListProfilesResponse], error) {
	msg := req.Msg

	// 构建查询条件
	query := db.GetDB().Model(&model.Profile{})

	// 添加用户ID条件
	query = query.Where("user_id = ?", auth.GetUserID(ctx))

	// 添加搜索条件
	if msg.SearchName != "" {
		query = query.Where("name LIKE ? OR im_name LIKE ?", "%"+msg.SearchName+"%", "%"+msg.SearchName+"%")
	}

	// 处理分页
	var pageSize int
	if msg.PageSize != "" {
		pageSize = fn.Atoi[int](msg.PageSize)
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	if msg.PageToken != "" {
		query = query.Where("id > ?", msg.PageToken)
	}

	// 查询数据
	var profiles []model.Profile
	if err := query.Order("id ASC").Limit(pageSize).Find(&profiles).Error; err != nil {
		slog.Error("list profiles error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 转换为 proto 并处理头像
	protoProfiles := make([]*profile.Profile, len(profiles))
	fileService := file.NewService()
	userID := auth.GetUserID(ctx)

	for i, p := range profiles {
		protoProfile := p.ToProto()

		// 如果有 avatar_file_id，尝试获取最新的 public_url
		if p.AvatarFileID > 0 {
			userFile, err := fileService.GetFileByID(userID, p.AvatarFileID)
			if err == nil && userFile.PublicURL != "" {
				protoProfile.Avatar = userFile.PublicURL
			}
		}

		protoProfiles[i] = protoProfile
	}

	// 处理下一页令牌
	var nextPageToken string
	if len(profiles) == pageSize {
		nextPageToken = strconv.Itoa(int(profiles[len(profiles)-1].ID))
	}

	return connect.NewResponse(&profile.ListProfilesResponse{
		Profiles:      protoProfiles,
		NextPageToken: nextPageToken,
	}), nil
}

func (s *ProfileService) CreateProfile(ctx context.Context, req *connect.Request[profile.CreateProfileRequest]) (*connect.Response[profile.CreateProfileResponse], error) {
	msg := req.Msg

	// 验证性别字段
	if msg.Gender != "" && msg.Gender != "male" && msg.Gender != "female" {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("gender must be 'male' or 'female', got: %s", msg.Gender))
	}

	// 创建新的个人资料
	newProfile := &model.Profile{
		UserID: auth.GetUserID(ctx),
		Name:   msg.Name,
		ImName: msg.ImName,
		Avatar: msg.Avatar,
		Gender: msg.Gender,
		Age:    int(msg.Age),
		Intro:  msg.Intro,
	}

	// 转换自定义属性
	for _, prop := range msg.Custom {
		newProfile.Custom = append(newProfile.Custom, model.Property{
			Name:  prop.Name,
			Value: prop.Value,
		})
	}

	// 保存到数据库
	if err := db.GetDB().Create(newProfile).Error; err != nil {
		slog.Error("create profile error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 获取proto消息
	profileProto := newProfile.ToProto()

	return connect.NewResponse(&profile.CreateProfileResponse{
		Profile: profileProto,
	}), nil
}

// ProfileUpdate 用于更新的结构体，使用指针字段
type ProfileUpdate struct {
	Name         *string         `gorm:"column:name"`
	ImName       *string         `gorm:"column:im_name"`
	Avatar       *string         `gorm:"column:avatar"`
	AvatarFileID *uint           `gorm:"column:avatar_file_id"`
	Gender       *string         `gorm:"column:gender"`
	Age          *int            `gorm:"column:age"`
	Intro        *string         `gorm:"column:intro"`
	Custom       *[]model.Property `gorm:"column:custom;serializer:json"`
}

func (s *ProfileService) UpdateProfile(ctx context.Context, req *connect.Request[profile.UpdateProfileRequest]) (*connect.Response[profile.UpdateProfileResponse], error) {
	msg := req.Msg
	ctx = context.WithValue(ctx, "env_version", req.Header().Get("X-Env-Version"))

	// 验证个人资料ID
	id, err := strconv.Atoi(msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// 验证性别字段（提前验证，避免不必要的数据库查询）
	if msg.Gender != "" && msg.Gender != "male" && msg.Gender != "female" {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("gender must be 'male' or 'female', got: %s", msg.Gender))
	}

	// 查询原个人资料
	var profileModel model.Profile
	if err := db.GetDB().Where("id = ? AND user_id = ?", id, auth.GetUserID(ctx)).First(&profileModel).Error; err != nil {
		slog.Error("get profile error", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// 构建更新结构体
	updates := ProfileUpdate{}
	
	if msg.Name != "" {
		updates.Name = &msg.Name
	}
	if msg.ImName != "" {
		updates.ImName = &msg.ImName
	}
	if msg.Gender != "" {
		updates.Gender = &msg.Gender
	}
	if msg.Age > 0 {
		age := int(msg.Age)
		updates.Age = &age
	}
	if msg.Intro != "" {
		updates.Intro = &msg.Intro
	}
	
	// 处理自定义属性
	if len(msg.Custom) > 0 {
		var custom []model.Property
		for _, prop := range msg.Custom {
			custom = append(custom, model.Property{
				Name:  prop.Name,
				Value: prop.Value,
			})
		}
		updates.Custom = &custom
	}

	// 处理头像上传
	// 优先使用 avatar_file_id
	if msg.AvatarFileId != "" {
		fileID := fn.Atoi[uint](msg.AvatarFileId)
		userID := auth.GetUserID(ctx)

		// 获取文件信息
		fileService := file.NewService()
		userFile, err := fileService.GetFileByID(userID, fileID)
		if err != nil {
			slog.Error("get file by id error", "error", err, "file_id", fileID)
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("头像文件不存在或无权访问"))
		}

		// 更新文件用途类型为头像
		if userFile.UsageType != model.UsageTypeAvatar {
			if err := fileService.UpdateFileUsageType(userID, fileID, model.UsageTypeAvatar); err != nil {
				slog.Error("update file usage type error", "error", err, "file_id", fileID)
				// 不阻塞主流程
			}
		}

		// 更新 avatar_file_id 和 avatar（缓存公网URL）
		updates.AvatarFileID = &fileID
		updates.Avatar = &userFile.PublicURL

	} else if msg.Avatar != "" {
		// 兼容旧的头像上传方式
		// 如果头像是临时上传的文件，移动到用户专门的目录
		if strings.HasPrefix(msg.Avatar, "uploads/") {
			userId := auth.GetUserID(ctx)
			// 提取文件名
			_, fileName := filepath.Split(msg.Avatar)
			// 构建用户专属目录路径
			userAvatarPath := fmt.Sprintf("user/%d/%s", userId, fileName)

			// 使用OSS客户端复制文件到新位置
			bucket := ossc.Get().UserFileBucket()
			if _, err := bucket.CopyObjectFrom(bucket.BucketName, msg.Avatar, userAvatarPath); err != nil {
				slog.Error("copy avatar error", "error", err, "src", msg.Avatar, "dst", userAvatarPath)
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			// 删除原临时文件
			if err := bucket.DeleteObject(msg.Avatar); err != nil {
				slog.Warn("delete temp avatar error", "error", err, "path", msg.Avatar)
				// 继续执行，不返回错误
			}

			// 签名新头像路径
			publicURL, err := ossc.GetPublic().UserFileBucket().SignURL(userAvatarPath, "GET", 10*365*24*60*60)
			if err != nil {
				slog.Warn("generate public avatar URL error", "error", err, "avatar", userAvatarPath)
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			// 更新头像路径
			updates.Avatar = &publicURL
		} else {
			updates.Avatar = &msg.Avatar
		}
	}

	// 更新数据库
	if err := db.GetDB().Model(&profileModel).Updates(updates).Error; err != nil {
		slog.Error("update profile error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的数据
	if err := db.GetDB().First(&profileModel, profileModel.ID).Error; err != nil {
		slog.Error("get updated profile error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 检查是否需要创建 demo 数据
	userID := auth.GetUserID(ctx)
	// 检查用户是否已有标记为 'demo' 的消息
	var count int64
	if err := db.GetDB().Model(&model.ChatMessage{}).
		Where("user_id = ? AND JSON_CONTAINS(tags, '\"demo\"')", userID).
		Count(&count).Error; err != nil {
		slog.Error("check demo messages error", "error", err)
		// 不返回错误，继续执行
	} else if count == 0 {
		// 用户没有 demo 数据，创建 demo 数据
		user := &model.User{}
		user.ID = userID
		// 传递更新后的 profile
		if err := domain.CopyDemoDataForNewUser(ctx, user, &profileModel); err != nil {
			slog.Error("failed to copy demo data for user", "error", err, "user_id", userID, "gender", profileModel.Gender)
			// 不返回错误，让用户的 profile 更新成功
		}
	}

	// 获取proto消息
	profileProto := profileModel.ToProto()

	return connect.NewResponse(&profile.UpdateProfileResponse{
		Profile: profileProto,
	}), nil
}

func (s *ProfileService) DeleteProfile(ctx context.Context, req *connect.Request[profile.DeleteProfileRequest]) (*connect.Response[profile.DeleteProfileResponse], error) {
	msg := req.Msg

	// 验证个人资料ID
	id, err := strconv.Atoi(msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// 删除个人资料
	if err := db.GetDB().Where("id = ? AND user_id = ?", id, auth.GetUserID(ctx)).Delete(&model.Profile{}).Error; err != nil {
		slog.Error("delete profile error", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&profile.DeleteProfileResponse{}), nil
}

func (s *ProfileService) GetProfile(ctx context.Context, req *connect.Request[profile.GetProfileRequest]) (*connect.Response[profile.GetProfileResponse], error) {
	msg := req.Msg

	// 验证资料ID
	id, err := strconv.Atoi(msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// 查询资料
	var profileModel model.Profile
	query := db.GetDB().Where("id = ?", id)

	// 用户只能查看自己的资料
	// 这里获取当前登录用户ID进行校验
	userID := auth.GetUserID(ctx)

	// 用户只能查看自己的资料
	query = query.Where("user_id = ?", userID)

	if err := query.First(&profileModel).Error; err != nil {
		slog.Error("get profile error", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if profileModel.AvatarFileID > 0 {
		profileModel.Avatar = generateAvatarSignedURL(profileModel.Avatar)
	}

	// 获取proto消息
	profileProto := profileModel.ToProto()

	// 如果有 avatar_file_id，尝试获取最新的 public_url
	if profileModel.AvatarFileID > 0 {
		userFile, err := file.NewService().GetFileByID(userID, profileModel.AvatarFileID)
		if err == nil && userFile.PublicURL != "" {
			profileProto.Avatar = userFile.PublicURL
		}
	}

	// 返回资料信息
	return connect.NewResponse(&profile.GetProfileResponse{
		Profile: profileProto,
	}), nil
}
