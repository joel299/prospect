package buffer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/iainfinito/chat-backend/internal/domain"
)

// HistoryClient reads messages captured by the external message-buffer API.
// HistoryURL must be provided by deployment config and may contain {lead_id}.
type HistoryClient struct {
	HistoryURL, AccessToken string
	HTTP                    *http.Client
}

func (c HistoryClient) Recent(ctx context.Context, leadID string, limit int) ([]domain.Message, error) {
	if strings.TrimSpace(c.HistoryURL) == "" {
		return nil, fmt.Errorf("buffer history URL is not configured")
	}
	u := strings.ReplaceAll(c.HistoryURL, "{lead_id}", url.PathEscape(leadID))
	if !strings.Contains(c.HistoryURL, "{lead_id}") {
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u += fmt.Sprintf("%slimit=%d", sep, limit)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}
	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("buffer history status %d", resp.StatusCode)
	}
	var raw any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return parseMessages(raw), nil
}

func parseMessages(raw any) []domain.Message {
	var items []any
	if m, ok := raw.(map[string]any); ok {
		for _, key := range []string{"messages", "data", "history", "items"} {
			if v, ok := m[key].([]any); ok {
				items = v
				break
			}
		}
	} else if v, ok := raw.([]any); ok {
		items = v
	}
	out := make([]domain.Message, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		content := firstString(m, "content", "message", "text", "body")
		if strings.TrimSpace(content) == "" {
			continue
		}
		direction := firstString(m, "direction", "role")
		if direction == "user" {
			direction = "inbound"
		}
		if direction == "assistant" || direction == "agent" {
			direction = "outbound"
		}
		if direction == "" {
			if fromMe, ok := m["fromMe"].(bool); ok && fromMe {
				direction = "outbound"
			} else {
				direction = "inbound"
			}
		}
		out = append(out, domain.Message{Provider: "buffer", ExternalID: firstString(m, "id", "message_id", "messageId", "external_id"), Direction: direction, Content: content, Status: "captured"})
	}
	return out
}
func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}
