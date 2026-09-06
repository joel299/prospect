package omniroute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/iainfinito/chat-backend/internal/domain"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL, APIKey, Model string
	HTTP                   *http.Client
}

func (c Client) Generate(ctx context.Context, r domain.AgentRequest) (domain.AgentResponse, error) {
	if c.BaseURL == "" {
		return domain.AgentResponse{}, fmt.Errorf("omniroute not configured")
	}
	b, _ := json.Marshal(map[string]any{"model": c.Model, "system": r.System, "user": r.User, "history": r.History, "tools": r.Tools})
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.BaseURL, "/")+"/v1/chat/completions", bytes.NewReader(b))
	if e != nil {
		return domain.AgentResponse{}, e
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	resp, e := c.HTTP.Do(req)
	if e != nil {
		return domain.AgentResponse{}, e
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return domain.AgentResponse{}, fmt.Errorf("omniroute status %d", resp.StatusCode)
	}
	var out domain.AgentResponse
	e = json.NewDecoder(resp.Body).Decode(&out)
	return out, e
}
