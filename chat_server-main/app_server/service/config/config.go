package config

import (
	"context"
	"log/slog"

	"app_server/domain/appconfig"
	"app_server/model"
	"app_server/pkg/db"
	"app_server/proto/config"

	"connectrpc.com/connect"
)

type ConfigService struct{}

func (s *ConfigService) GetConfig(ctx context.Context, req *connect.Request[config.GetConfigRequest]) (*connect.Response[config.GetConfigResponse], error) {
	var configs []model.Config
	query := db.GetDB().Model(&model.Config{})

	// 构建查询条件
	if len(req.Msg.Keys) > 0 {
		query = query.Where("k IN ?", req.Msg.Keys)
	}
	if req.Msg.App != "" {
		query = query.Where("app = ? OR app = ''", req.Msg.App)
	}
	if req.Msg.Platform != "" {
		query = query.Where("platform = ? OR platform = ''", req.Msg.Platform)
	}
	if req.Msg.Env != "" {
		query = query.Where("env = ? OR env = ''", req.Msg.Env)
	}

	// 执行查询
	if err := query.Find(&configs).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 按key分组，找出符合版本要求的配置
	configMap := appconfig.FilterConfigByVersion(configs, req.Msg.Version)
	// 构建响应
	resp := &config.GetConfigResponse{
		Configs: make([]*config.Config, 0, len(configMap)),
	}

	for _, cfg := range configMap {
		resp.Configs = append(resp.Configs, &config.Config{
			Key:   cfg.Key,
			Value: cfg.Value,
		})
	}
	slog.Info("get config", "configs", configs, "respConfigs", resp.Configs, "version", req.Msg.Version, "env", req.Msg.Env)

	return connect.NewResponse(resp), nil
}
