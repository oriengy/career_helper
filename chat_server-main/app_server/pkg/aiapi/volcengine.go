package aiapi

import (
	"app_server/pkg/openaic"
	"context"
	"errors"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ParseImageChat 解析图片中的聊天内容
func ParseImageChat(imageUrl string) ([]string, error) {
	// 构建请求
	resp, err := openaic.Get().CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openaic.Model.Ocr,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: "提取图中的聊天记录。忽略引用的内容。每行一句，每行下列格式开始\n【朋友】xxx\n【自己】yyy",
						},
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL: imageUrl,
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	// 检查是否有结果
	if len(resp.Choices) == 0 {
		return nil, errors.New("未获取到解析结果")
	}

	// 解析聊天记录，按行分割
	content := resp.Choices[0].Message.Content
	lines := strings.Split(content, "\n")

	var chatLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && (strings.HasPrefix(line, "【朋友】") || strings.HasPrefix(line, "【自己】")) {
			chatLines = append(chatLines, line)
		}
	}

	return chatLines, nil
}

// 解析单行聊天记录，返回角色和内容
func ParseChatLine(line string) (role string, content string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", false
	}

	// 检查【朋友】或【自己】格式
	if after, ok0 :=strings.CutPrefix(line, "【朋友】"); ok0  {
		return "FRIEND", after, true
	} else if after, ok0 :=strings.CutPrefix(line, "【自己】"); ok0  {
		return "SELF", after, true
	}

	return "", "", false
}
