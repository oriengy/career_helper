package domain

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"app_server/model"
)

var (
	//go:embed demo/male.json
	demoDataMaleJSON []byte
	//go:embed demo/female.json
	demoDataFemaleJSON []byte
	//go:embed demo/feature_guide.json
	demoDataFeatureGuideJSON []byte

	demoJON = map[string][]byte{
		"male":          demoDataMaleJSON,
		"female":        demoDataFemaleJSON,
		"feature_guide": demoDataFeatureGuideJSON,
		"":              demoDataFeatureGuideJSON,
	}

	demoCases = map[string]DemoData{}
)

func init() {
	for demoType, demoData := range demoJON {
		var data DemoData
		if err := json.Unmarshal(demoData, &data); err != nil {
			slog.Error("failed to unmarshal demo data", "error", err, "demo_type", demoType)
			continue
		}
		demoCases[demoType] = data
	}
}

// DemoData 演示数据结构
type DemoData struct {
	Name      string     `json:"name"`
	UserID    uint       `json:"user_id"`
	DemoCases []DemoCase `json:"demo_cases"`
}

// DemoChatSession 演示会话数据
type DemoCase struct {
	ChatSession model.ChatSession   `json:"chat_session"`
	Profile     model.Profile       `json:"profile"`
	Messages    []model.ChatMessage `json:"messages"`
}

// LoadDemoData 加载演示数据
func LoadDemoData(demoType string) DemoData {
	// 构建文件路径
	filename := fmt.Sprintf("demo/%s.json", demoType)
	// 尝试从文件系统读取
	if data, err := os.ReadFile(filename); err == nil {
		var demoData DemoData
		if err := json.Unmarshal(data, &demoData); err == nil {
			slog.Info("loaded demo data from file", "type", demoType, "file", filename)
			return demoData
		}
	}

	// 如果文件读取失败，返回内存中的数据（嵌入的或init中加载的）
	return demoCases[demoType]
}

func LoadDemoDataFromJsonString(jsonString string) DemoData {
	var demoData DemoData
	if err := json.Unmarshal([]byte(jsonString), &demoData); err != nil {
		slog.Error("failed to unmarshal demo data", "error", err)
		return DemoData{}
	}
	return demoData
}
