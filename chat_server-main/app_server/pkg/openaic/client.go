package openaic

import "github.com/sashabaranov/go-openai"

var (
	client *openai.Client
	Model  Models
)

type Config struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Models  Models `mapstructure:"models"`
}

type Models struct {
	Chat string `mapstructure:"chat"`
	Ocr  string `mapstructure:"ocr"`
}

func Init(cfg Config) {
	// 初始化OpenAI客户端
	config := openai.DefaultConfig(cfg.ApiKey)
	config.BaseURL = cfg.BaseURL
	Model = cfg.Models
	client = openai.NewClientWithConfig(config)
}

func Get() *openai.Client {
	return client
}
