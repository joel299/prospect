package composio

import "context"

type ToolExecutor interface {
	Execute(context.Context, string, map[string]any) (map[string]any, error)
}
type Client struct{ BaseURL, APIKey string }

func (c Client) Execute(_ context.Context, _ string, _ map[string]any) (map[string]any, error) {
	return nil, ErrNotConfigured
}

var ErrNotConfigured = errorString("composio tool executor is not configured")

type errorString string

func (e errorString) Error() string { return string(e) }
