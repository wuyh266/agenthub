package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const endpoint = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
const model = "glm-5.3-flash"

type glmClient struct {
	apikey     string
	httpClient *http.Client
}

// 真实返回的JSON
// json
// {
// 	"choices": [
// 	{
// 		"index": 0,
// 		"message": {
// 		"role": "assistant",
// 		"content": "模型回复的那段话"
// 		},
// 		"finish_reason": "stop"
// 	}
// 	],
// 	"usage": { "prompt_tokens": 10, "completion_tokens": 5 },
// 	"id": "xxx",
// 	"created": 1234567890
// }

// GLMResponse 是大模型响应客户端结构体
type glmResponse struct {
	Choices []chatChoice `json:"choices"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

// GLMRequest 是客户端请求大模型结构体
type glmRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewGLMClient(apikey string) *glmClient {
	return &glmClient{
		apikey: apikey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
func (c *glmClient) Chat(ctx context.Context, messages []Message) (string, error) {
	// 这里实现调用 GLM 接口的逻辑，使用 c.apikey 和 c.http.Client 进行请求
	// 例如，构建请求体，发送 HTTP 请求，处理响应等
	chatmessage := make([]chatMessage, len(messages))
	for i, m := range messages {
		chatmessage[i] = chatMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	glmReq := glmRequest{
		Model:    model,
		Messages: chatmessage,
	}
	data, err := json.Marshal(glmReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apikey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GLM API 返回状态码 %d: %s", resp.StatusCode, body)
	}

	var glmResp glmResponse
	err = json.Unmarshal(body, &glmResp) // 这玩意儿只返回一个err值
	if err != nil {
		return "", err
	}
	if len(glmResp.Choices) == 0 {
		return "", errors.New("GLM API 返回的 choices 为空")
	}
	return glmResp.Choices[0].Message.Content, nil // 返回模型的回复内容
}
