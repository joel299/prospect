package ryze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/iainfinito/chat-backend/internal/domain"
)

type Client interface {
	SendText(context.Context, domain.OutboundMessage) (string, error)
	History(context.Context, string) ([]domain.Message, error)
}

type HTTPClient struct {
	BaseURL  string
	APIKey   string
	Instance string
	HTTP     *http.Client
}

func (c HTTPClient) client() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

func (c HTTPClient) endpoint(path string) (string, error) {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.APIKey) == "" || strings.TrimSpace(c.Instance) == "" {
		return "", fmt.Errorf("ryze not configured")
	}
	return strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(fmt.Sprintf(path, c.Instance), "/"), nil
}

func (c HTTPClient) do(ctx context.Context, method, url string, body any, out any) error {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return err
	}
	req.Header.Set("token", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("ryze status %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c HTTPClient) SendText(ctx context.Context, m domain.OutboundMessage) (string, error) {
	url, err := c.endpoint("api/message/text/%s")
	if err != nil {
		return "", err
	}
	var out struct {
		MessageID string `json:"messageId"`
		ID        string `json:"id"`
	}
	err = c.do(ctx, http.MethodPost, url, map[string]string{"number": m.LeadID, "message": m.Content}, &out)
	if err != nil {
		return "", err
	}
	if out.MessageID != "" {
		return out.MessageID, nil
	}
	return out.ID, nil
}

func (c HTTPClient) History(ctx context.Context, number string) ([]domain.Message, error) {
	url, err := c.endpoint("api/chat/history/%s")
	if err != nil {
		return nil, err
	}
	var raw any
	if err = c.do(ctx, http.MethodPost, url, map[string]string{"number": number}, &raw); err != nil {
		return nil, err
	}
	return parseHistory(raw), nil
}

func parseHistory(raw any) []domain.Message {
	var items []any
	switch v := raw.(type) {
	case []any:
		items = v
	case map[string]any:
		for _, k := range []string{"messages", "data", "history"} {
			if x, ok := v[k].([]any); ok {
				items = x
				break
			}
		}
	}
	out := make([]domain.Message, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		content, _ := m["message"].(string)
		if content == "" {
			content, _ = m["body"].(string)
		}
		id, _ := m["id"].(string)
		if id == "" {
			id, _ = m["messageId"].(string)
		}
		direction := "inbound"
		if x, ok := m["fromMe"].(bool); ok && x {
			direction = "outbound"
		}
		out = append(out, domain.Message{Provider: "ryze", ExternalID: id, Direction: direction, Content: content, Status: "received"})
	}
	return out
}
