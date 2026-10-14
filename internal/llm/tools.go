package llm

import "encoding/json"

// Tool 是发给模型的「工具说明书」（请求 JSON 里 tools 数组的每一项）。
// 注意：这不是执行结果，只是告诉模型「你可以调用哪些函数」。
type Tool struct {
	Type     string       `json:"type"` // 固定填 "function"
	Function ToolFunction `json:"function"`
}

// ToolFunction 描述一个可调用函数。
type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	// Parameters 一般是 JSON Schema 对象；用 RawMessage 避免再套一层 map。
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

// ToolCall 是模型在回复里说的「我要调用这个工具」。
// 出现在 assistant 消息的 tool_calls 数组里。
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"` // 一般为 "function"
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction 里的 Arguments 在 API 里是 JSON 字符串，不是对象。
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
