package buffer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/iainfinito/chat-backend/internal/domain"
)

type Consumed struct {
	Success bool `json:"success"`
	Found   bool `json:"found"`
	Buffer  struct {
		ID       string `json:"id"`
		TenantID string `json:"tenant_id"`
		Phone    string `json:"phone"`
		Messages []struct {
			MessageID         string `json:"message_id"`
			Content           string `json:"content"`
			NormalizedContent string `json:"normalized_content"`
		} `json:"messages"`
	} `json:"buffer"`
}

type Consumer struct {
	BaseURL, AccessToken, TenantID, ConsumerID string
	HTTP                                       *http.Client
}

func (c Consumer) Publish(ctx context.Context, in domain.InboundMessage) error {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.AccessToken) == "" || strings.TrimSpace(c.TenantID) == "" {
		return fmt.Errorf("buffer publisher not configured")
	}
	body, err := json.Marshal(map[string]any{
		"tenant_id":  c.TenantID,
		"channel":    "whatsapp",
		"phone":      in.LeadID,
		"message_id": in.ExternalID,
		"content":    in.Content,
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.BaseURL, "/")+"/api/message-buffer/messages", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	h := c.HTTP
	if h == nil {
		h = http.DefaultClient
	}
	resp, err := h.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("buffer publish status %d", resp.StatusCode)
	}
	return nil
}

func (c Consumer) ConsumeNext(ctx context.Context) (Consumed, error) {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.AccessToken) == "" || strings.TrimSpace(c.TenantID) == "" || strings.TrimSpace(c.ConsumerID) == "" {
		return Consumed{}, fmt.Errorf("buffer consumer not configured")
	}
	body, _ := json.Marshal(map[string]string{"tenant_id": c.TenantID, "channel": "whatsapp", "consumer_id": c.ConsumerID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.BaseURL, "/")+"/api/message-buffer/consume/next", bytes.NewReader(body))
	if err != nil {
		return Consumed{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	h := c.HTTP
	if h == nil {
		h = http.DefaultClient
	}
	resp, err := h.Do(req)
	if err != nil {
		return Consumed{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Consumed{}, fmt.Errorf("buffer consume status %d", resp.StatusCode)
	}
	var out Consumed
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Consumed{}, err
	}
	return out, nil
}
func (c Consumer) Recent(ctx context.Context, _ string, _ int) ([]domain.Message, error) {
	out, err := c.ConsumeNext(ctx)
	if err != nil || !out.Found {
		return nil, err
	}
	msgs := make([]domain.Message, 0, len(out.Buffer.Messages))
	for _, m := range out.Buffer.Messages {
		content := m.Content
		if content == "" {
			content = m.NormalizedContent
		}
		if content != "" {
			msgs = append(msgs, domain.Message{Provider: "buffer", ExternalID: m.MessageID, Direction: "inbound", Content: content, Status: "locked"})
		}
	}
	return msgs, nil
}
