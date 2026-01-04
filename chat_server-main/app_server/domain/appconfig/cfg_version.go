package appconfig

import (
	"strconv"
	"strings"

	"app_server/model"
)

func FilterConfigByKeyVersion(configs []model.Config, key string, clientVersion string) *model.Config {
	var retCfg *model.Config
	for i := range configs {
		cfg := &configs[i]
		if cfg.Key != key {
			continue
		}

		// 过滤掉 配置版本 > 客户端版本 的配置
		if compareVersion(cfg.Version, clientVersion) > 0 {
			continue
		}
		// 如果已存在，选择版本更大的
		if compareVersion(cfg.Version, getVersion(retCfg)) > 0 {
			retCfg = cfg
		}

	}

	return retCfg
}

func getVersion(cfg *model.Config) string {
	if cfg == nil || cfg.Version == "" {
		return "0.0.0"
	}
	return cfg.Version
}

func FilterConfigByVersion(configs []model.Config, clientVersion string) map[string]*model.Config {
	// 按key分组，找出符合版本要求的配置
	configMap := make(map[string]*model.Config)

	for i := range configs {
		cfg := &configs[i]
		// 如果客户端没有提供版本，或者配置版本为空，直接使用
		if clientVersion == "" || cfg.Version == "" {
			if existing, ok := configMap[cfg.Key]; !ok || cfg.ID > existing.ID {
				configMap[cfg.Key] = cfg
			}
			continue
		}

		// 比较版本
		if compareVersion(cfg.Version, clientVersion) <= 0 {
			// 配置版本 <= 客户端版本
			if existing, ok := configMap[cfg.Key]; ok {
				// 如果已存在，选择版本更大的
				if compareVersion(cfg.Version, existing.Version) > 0 {
					configMap[cfg.Key] = cfg
				}
			} else {
				configMap[cfg.Key] = cfg
			}
		}
	}

	return configMap
}

// compareVersion 比较两个版本号 (major.minor.patch格式)
// 返回值: -1 表示 v1 < v2, 0 表示 v1 = v2, 1 表示 v1 > v2
func compareVersion(v1, v2 string) int {
	if v1 == v2 {
		return 0
	}
	if v1 == "" {
		v1 = "0.0.0"
	}
	if v2 == "" {
		v2 = "0.0.0"
	}
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// 补齐版本号长度
	maxLen := max(len(parts2), len(parts1))

	for i := range maxLen {
		var num1, num2 int

		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}

		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}

	return 0
}
