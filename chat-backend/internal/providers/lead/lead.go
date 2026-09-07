package lead

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client reads the lead row from the configured lead API. The URL is explicit
// because the project may use Supabase or another CRM; no endpoint is guessed.
type Client struct {
	URL, Token string
	HTTP       *http.Client
}

func (c Client) Get(ctx context.Context, leadID string) (map[string]any, string, error) {
	if strings.TrimSpace(c.URL) == "" {
		return nil, "", fmt.Errorf("lead API URL is not configured")
	}
	u := strings.ReplaceAll(c.URL, "{lead_id}", url.PathEscape(leadID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Accept", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("apikey", c.Token)
	}
	h := c.HTTP
	if h == nil {
		h = http.DefaultClient
	}
	resp, err := h.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("lead API status %d", resp.StatusCode)
	}
	var raw any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, "", err
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return nil, "", fmt.Errorf("lead API returned invalid object")
	}
	if data, ok := m["data"].(map[string]any); ok {
		m = data
	}
	state := ""
	for _, k := range []string{"state", "lead_status", "pipeline_stage", "status"} {
		if v, ok := m[k].(string); ok && v != "" {
			state = v
			break
		}
	}
	return m, state, nil
}
