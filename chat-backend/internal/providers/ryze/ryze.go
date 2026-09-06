package ryze

import (
	"context"
	"fmt"
	"github.com/iainfinito/chat-backend/internal/domain"
)

type Client interface {
	SendText(context.Context, domain.OutboundMessage) (string, error)
	History(context.Context, string) ([]domain.Message, error)
}
type HTTPClient struct{ BaseURL, APIKey, Instance string }

func (c HTTPClient) SendText(_ context.Context, m domain.OutboundMessage) (string, error) {
	if c.BaseURL == "" {
		return "", fmt.Errorf("ryze not configured")
	}
	return "", fmt.Errorf("ryze adapter requires endpoint mapping")
}
func (c HTTPClient) History(_ context.Context, _ string) ([]domain.Message, error) {
	if c.BaseURL == "" {
		return nil, fmt.Errorf("ryze not configured")
	}
	return nil, fmt.Errorf("ryze adapter requires endpoint mapping")
}
