package oai

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"app_server/pkg/openaic"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

// Client 封装了 OpenAI 客户端的调用逻辑
type Client struct {
	client *openai.Client
	debug  bool
}

// NewClient 创建一个新的 OpenAI 客户端包装器
func NewClient() *Client {
	// 检查是否在 debug 模式
	debug := os.Getenv("DEBUG") == "true" || os.Getenv("OAI_DEBUG") == "true" || viper.GetString("logLevel") == "debug"

	return &Client{
		client: openaic.Get(),
		debug:  debug,
	}
}

// ChatCompletionRequest 包装了聊天完成请求的参数
type ChatCompletionRequest struct {
	Messages []openai.ChatCompletionMessage
	Model    string
	Stream   bool
}

// CreateChatCompletion 调用 OpenAI API 并处理日志
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (content string, err error) {
	// 如果没有指定模型，使用默认模型
	if req.Model == "" {
		req.Model = openaic.Model.Chat
	}

	// 在 debug 模式下打印请求
	if c.debug {
		defer func() {
			c.logRequest(req, content)
		}()
	}

	// 构建 OpenAI 请求
	openaiReq := openai.ChatCompletionRequest{
		Model:    req.Model,
		Messages: req.Messages,
		Stream:   req.Stream,
	}

	// 调用 OpenAI API
	resp, err := c.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		slog.Error("OpenAI API 调用失败",
			"error", err,
			"model", req.Model,
		)
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	// 获取响应内容
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("OpenAI API 返回空响应")
	}

	content = resp.Choices[0].Message.Content

	return content, nil
}

// CreateChatCompletionSimple 简化版本的聊天完成请求，接受单个用户消息
func (c *Client) CreateChatCompletionSimple(ctx context.Context, userPrompt string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}

	return c.CreateChatCompletion(ctx, ChatCompletionRequest{
		Messages: messages,
	})
}

// CreateChatCompletionWithSystem 带系统提示词的聊天完成请求
func (c *Client) CreateChatCompletionWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}

	return c.CreateChatCompletion(ctx, ChatCompletionRequest{
		Messages: messages,
	})
}

// logRequest 美化打印请求日志
func (c *Client) logRequest(req ChatCompletionRequest, respStr string) {
	fmt.Println("\n" + strings.Repeat("=", 30))
	fmt.Println("🤖 OpenAI API 请求")
	fmt.Printf("💬 消息数量: %d\n", len(req.Messages))

	for _, msg := range req.Messages {
		fmt.Printf("\n[%s]:", msg.Role)
		content := msg.Content
		// 缩进消息内容
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			fmt.Printf("   %s\n", line)
		}
	}

	if respStr != "" {
		fmt.Printf("\n[AI RESPONSE]: %s\n", respStr)
	}

	fmt.Println(strings.Repeat("=", 30))
}

// FormatMessagesForLog 格式化消息列表用于日志记录
func FormatMessagesForLog(messages []openai.ChatCompletionMessage) string {
	var parts []string
	for _, msg := range messages {
		parts = append(parts, fmt.Sprintf("[%s]: %s", msg.Role, msg.Content))
	}
	return strings.Join(parts, "\n")
}

func Get() *Client {
	return NewClient()
}
