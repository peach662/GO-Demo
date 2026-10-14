package llm

import(
	"context"
	"errors"
	"encoding/json"
	"awesomeProject/internal/todo"
)
type Service struct {
    client *Client    
	Todo *todo.Service      // todo 里是 repo / redis / msg
}

func NewService(client *Client,todo *todo.Service) *Service{
	return &Service{
		client: client,
		Todo: todo,
	}
}

// Chat 完成一轮带工具的对话。
// userID 来自登录态，只用于本地查待办，不会交给模型填写。
// 模型第一次若只回文本，直接返回；若点了 list_todos，则本地查库后再问第二次。
func (s *Service) Chat(ctx context.Context,userID int, message string) (string, error) {
	if message == ""{
		return "", errors.New("message is empty")
	}else if userID == 0{
		return "", errors.New("userID is empty")
	}
	// messages 是发给模型的对话记录。chatMessage 定义在同包 client.go，所以这里不用再声明。
	messages := []chatMessage{{Role: "user",Content:message}}
	// tools 是工具菜单。listTodoTool 只描述「可以点 list_todos」，本身不查库。
	tools:= []Tool{listTodoTool()}

	// ChatWithTools 发一次 HTTP，返回助手那一条消息（文本或 tool_calls）。
	assistant,err:= s.client.ChatWithTools(ctx,messages,tools)
	if err!=nil {
		return "", err
	}
	// 没有点单：Content 就是给用户的回答。
	if len(assistant.ToolCalls) == 0{
		return assistant.Content,nil
	}

	// 第二次请求必须带上模型这条点单，所以把 assistant 追加进记录。
	// append 返回新切片，必须赋回 messages，否则外面的变量还是旧的。
	messages = append(messages,assistant)
	if s.Todo == nil {
		return "",errors.New("todo service is nil")
	}
	for _,toolCall:= range assistant.ToolCalls{
		if toolCall.Function.Name == "list_todos"{
			// List 用 userID 查当前用户的待办（数据库或缓存），模型访问不到这里。
			todos,err:= s.Todo.List(ctx,userID)
			if err!=nil {
				return "", errors.New("failed to list todos")
			}
			// Marshal 把待办编成 JSON 字节；模型只能读字符串，所以下面再用 string。
			response,err:= json.Marshal(todos)
			if err!=nil {
				return "", errors.New("failed to marshal todos")
			}

			// 追加一条 role=tool 的消息。ToolCallID 对上这次点单的 ID，Content 是查库结果。
			messages = append(messages,chatMessage{Role: "tool",ToolCallID: toolCall.ID,Content: string(response)})
		} else {
			return "", errors.New("unknown tool")
		}
	}
	// 第二次 ChatWithTools：历史里已有用户原话、模型点单、工具结果，返回自然语言。
	second,err:=s.client.ChatWithTools(ctx,messages,tools)
	if err!=nil{
		return "",err
	}
	return second.Content,nil
}

// listTodoTool 返回 list_todos 的说明书，给模型看，不执行查询。
func listTodoTool() Tool{
	return Tool{
		Type: "function",
		Function:ToolFunction{
			Name: "list_todos",
			Description: "……列出当前用户待办……",
			Parameters:json.RawMessage(`{"type":"object","properties":{}}`),
			},
		}

}

