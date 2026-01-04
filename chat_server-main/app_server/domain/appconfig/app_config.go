package appconfig

import (
	"context"
	"log/slog"

	"app_server/model"
	"app_server/pkg/db"
)

func LoadAppConfigByKeyVersion(ctx context.Context, key string, param Cond) string {
	query := db.GetDB().Model(&model.Config{})

	// 构建查询条件
	query = query.Where("k = ?", key)
	query = query.Where("env IN (?, '')", param.Env)

	// 执行查询
	var configs []model.Config
	if err := query.Find(&configs).Error; err != nil {
		slog.Error("failed to load app config", "error", err)
		return ""
	}

	// 按key分组，找出符合版本要求的配置
	config := FilterConfigByKeyVersion(configs, key, param.Version)
	if config == nil {
		return ""
	}

	return config.Value
}

type Cond struct {
	Version string
	Env     string
}
