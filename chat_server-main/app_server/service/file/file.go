package file

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"app_server/model"
	"app_server/pkg/db"
	"app_server/pkg/ossc"

	"gorm.io/gorm"
)

// Service 文件服务
type Service struct {
	db *gorm.DB
}

// NewService 创建文件服务实例
func NewService() *Service {
	return &Service{
		db: db.GetDB(),
	}
}

// CreateFileRecord 创建文件记录
func (s *Service) CreateFileRecord(userID uint, originalName string, fileSize int64, fileType, fileExt, ossKey, fileHash, usageType string) (*model.UserFile, error) {
	// 根据用途类型获取过期时间
	expirationTime := model.GetExpirationTime(usageType)

	file := &model.UserFile{
		UserID:       userID,
		OriginalName: originalName,
		FileSize:     fileSize,
		FileType:     fileType,
		FileExt:      fileExt,
		OssKey:       ossKey,
		FileHash:     fileHash,
		Status:       model.FileStatusNormal,
		UsageType:    usageType,
	}

	if expirationTime != nil {
		file.ExpiresAt = sql.NullTime{Time: *expirationTime, Valid: true}
	}

	// 生成公网访问URL
	publicURL, err := s.generatePublicURL(ossKey)
	if err != nil {
		return nil, fmt.Errorf("生成公网访问URL失败: %w", err)
	}
	file.PublicURL = publicURL

	// 设置公网URL过期时间（10年）
	publicExpire := time.Now().Add(10 * 365 * 24 * time.Hour)
	file.PublicExpire = sql.NullTime{Time: publicExpire, Valid: true}

	// 保存到数据库
	if err := s.db.Create(file).Error; err != nil {
		return nil, fmt.Errorf("创建文件记录失败: %w", err)
	}

	return file, nil
}

// GetFileByID 根据ID获取文件
func (s *Service) GetFileByID(userID uint, fileID uint) (*model.UserFile, error) {
	var file model.UserFile
	err := s.db.Where("id = ? AND user_id = ? AND status = ?", fileID, userID, model.FileStatusNormal).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文件不存在或无权访问")
		}
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}

	// 检查公网URL是否过期，如果过期则重新生成
	if file.PublicExpire.Valid && file.PublicExpire.Time.Before(time.Now()) {
		publicURL, err := s.generatePublicURL(file.OssKey)
		if err == nil {
			file.PublicURL = publicURL
			file.PublicExpire = sql.NullTime{Time: time.Now().Add(10 * 365 * 24 * time.Hour), Valid: true}
			s.db.Model(&file).Updates(map[string]interface{}{
				"public_url":    file.PublicURL,
				"public_expire": file.PublicExpire,
			})
		}
	}

	return &file, nil
}

// GetFileByHash 根据文件哈希查找文件（用于去重）
func (s *Service) GetFileByHash(userID uint, fileHash string) (*model.UserFile, error) {
	var file model.UserFile
	err := s.db.Where("user_id = ? AND file_hash = ? AND status = ?", userID, fileHash, model.FileStatusNormal).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}
	return &file, nil
}

// DeleteFile 删除文件（软删除）
func (s *Service) DeleteFile(userID uint, fileID uint) error {
	result := s.db.Model(&model.UserFile{}).
		Where("id = ? AND user_id = ?", fileID, userID).
		Update("status", model.FileStatusDeleted)

	if result.Error != nil {
		return fmt.Errorf("删除文件失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("文件不存在或无权删除")
	}

	return nil
}

// ListFiles 查询文件列表
func (s *Service) ListFiles(userID uint, usageType string, status int, limit int) ([]*model.UserFile, error) {
	var files []*model.UserFile
	query := s.db.Where("user_id = ?", userID)

	if usageType != "" {
		query = query.Where("usage_type = ?", usageType)
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", model.FileStatusNormal)
	}

	err := query.Order("id DESC").Limit(limit).Find(&files).Error
	if err != nil {
		return nil, fmt.Errorf("查询文件列表失败: %w", err)
	}

	return files, nil
}

// CalculateFileHash 计算文件哈希值
func CalculateFileHash(reader io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", fmt.Errorf("计算文件哈希失败: %w", err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// generatePublicURL 生成公网访问URL
func (s *Service) generatePublicURL(ossKey string) (string, error) {
	// 生成10年有效期的签名URL
	publicURL, err := ossc.GetPublic().UserFileBucket().SignURL(ossKey, "GET", 10*365*24*60*60)
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %w", err)
	}
	return publicURL, nil
}

// UpdateFileUsageType 更新文件用途类型
func (s *Service) UpdateFileUsageType(userID uint, fileID uint, usageType string) error {
	// 根据新的用途类型获取过期时间
	expirationTime := model.GetExpirationTime(usageType)

	updates := map[string]interface{}{
		"usage_type": usageType,
	}

	if expirationTime != nil {
		updates["expires_at"] = sql.NullTime{Time: *expirationTime, Valid: true}
	} else {
		updates["expires_at"] = sql.NullTime{Valid: false}
	}

	result := s.db.Model(&model.UserFile{}).
		Where("id = ? AND user_id = ? AND status = ?", fileID, userID, model.FileStatusNormal).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("更新文件用途类型失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("文件不存在或无权更新")
	}

	return nil
}
