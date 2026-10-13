package llm

import(
	"context"
	"errors"
)
type Service struct {
    client *Client          // todo 里是 repo / redis / msg
}

func NewService(client *Client) *Service{
	return &Service{
		client: client,
	}
}

func (s *Service) Chat(ctx context.Context, message string) (string, error) {
    // 1. message 为空 → 返回 error
	if message == ""{
		return "", errors.New("message is empty")
	}
    // 2. 否则 return s.client.Chat(ctx, message)
	return s.client.Chat(ctx,message)
}