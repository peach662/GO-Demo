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

// chatRequest 是发往 /v1/chat/completions 的请求体。
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	// Tools 可选：有 Tool Calling 时带上；普通 Chat 不填（omitempty 不会出现在 JSON 里）。
	Tools []Tool `json:"tools,omitempty"`
}

// chatMessage 是 messages 数组里的一条消息。
// 三种常见形态共用这一个结构体，靠 Role + 可选字段区分：
//   - user/assistant 普通文本：Role + Content
//   - assistant 要调工具：Role=assistant + ToolCalls
//   - 本地把工具结果喂回模型：Role=tool + ToolCallID + Content
type chatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// Chat 发送一条 user 消息，返回 assistant 文本（不含 tools）。
func (c *Client) Chat(ctx context.Context, userMessage string) (string, error) {
	if userMessage == "" {
		return "", fmt.Errorf("user message is empty")
	}
	msg, err := c.ChatWithTools(ctx, []chatMessage{
		{Role: "user", Content: userMessage},
	}, nil)
	if err != nil {
		return "", err
	}
	if msg.Content == "" {
		return "", fmt.Errorf("chat api returned empty content")
	}
	return msg.Content, nil
}

// ChatWithTools 发送多轮 messages，并可附带工具定义。
// 返回的是「助手那一条消息」：可能只有 Content，也可能带 ToolCalls（此时 Content 常为空）。
func (c *Client) ChatWithTools(ctx context.Context, messages []chatMessage, tools []Tool) (chatMessage, error) {
	if c == nil {
		return chatMessage{}, fmt.Errorf("llm client is nil")
	}
	if len(messages) == 0 {
		return chatMessage{}, fmt.Errorf("messages is empty")
	}

	body, err := json.Marshal(chatRequest{
		Model:    c.model,
		Messages: messages,
		Tools:    tools,
	})
	if err != nil {
		return chatMessage{}, fmt.Errorf("marshal chat request: %w", err)
	}

	respBody, err := c.postChat(ctx, body)
	if err != nil {
		return chatMessage{}, err
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return chatMessage{}, fmt.Errorf("unmarshal chat response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return chatMessage{}, fmt.Errorf("chat api returned no choices")
	}
	return parsed.Choices[0].Message, nil
}

func (c *Client) postChat(ctx context.Context, body []byte) ([]byte, error) {
	url := c.chatCompletionsURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("chat api status %d: %s", resp.StatusCode, truncate(string(respBody), 512))
	}
	return respBody, nil
}

// chatCompletionsURL 兼容两种 BaseURL：
//   - https://host          → https://host/v1/chat/completions
//   - https://host/v1       → https://host/v1/chat/completions
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
