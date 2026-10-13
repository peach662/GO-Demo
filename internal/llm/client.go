package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client 调用 OpenAI 兼容的 Chat Completions API（官方或中转站均可）。
type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewClient 用密钥、服务根地址和模型名构造客户端。
// baseURL 一般是 https://xxx.com 或 https://xxx.com/v1，不要在这里写完整 path。
func NewClient(apiKey, baseURL, model string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		// 中转/模型有时较慢，给足超时，避免请求一直挂起。
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// chatRequest / chatResponse 只建模我们用到的字段，忽略其余 JSON。
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// Chat 发送一条 user 消息，返回 assistant 文本。
// 这是最小闭环：还不支持多轮历史、tools、流式输出。
func (c *Client) Chat(ctx context.Context, userMessage string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("llm client is nil")
	}
	if userMessage == "" {
		return "", fmt.Errorf("user message is empty")
	}

	body, err := json.Marshal(chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "user", Content: userMessage},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal chat request: %w", err)
	}

	url := c.chatCompletionsURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	// 非 2xx 时把状态码和一小段 body 带上，排错中转/鉴权/模型名时很有用。
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("chat api status %d: %s", resp.StatusCode, truncate(string(respBody), 512))
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("unmarshal chat response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("chat api returned no choices")
	}
	content := parsed.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("chat api returned empty content")
	}
	return content, nil
}

// chatCompletionsURL 兼容两种 BaseURL：
//   - https://host          → https://host/v1/chat/completions
//   - https://host/v1       → https://host/v1/chat/completions
// 避免拼成 /v1/v1/...。
func (c *Client) chatCompletionsURL() string {
	base := c.baseURL
	if strings.HasSuffix(base, "/v1") {
		return base + "/chat/completions"
	}
	return base + "/v1/chat/completions"
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
