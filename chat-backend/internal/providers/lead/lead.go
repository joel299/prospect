package lead

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SupabaseClient reads the commercial lead row. The backend uses the server
// credential only; this client is never exposed to the frontend.
type SupabaseClient struct {
	BaseURL, APIKey, Table string
	HTTP                   *http.Client
}

func (c SupabaseClient) Get(ctx context.Context, leadID string) (map[string]any, string, error) {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.APIKey) == "" {
		return nil, "", fmt.Errorf("supabase lead client is not configured")
	}
	table := c.Table
	if table == "" {
		table = "prospect_leads_google"
	}
	endpoint := strings.TrimRight(c.BaseURL, "/") + "/rest/v1/" + url.PathEscape(table) + "?select=*&id=eq." + url.QueryEscape(leadID) + "&limit=1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("apikey", c.APIKey)
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("supabase leads status %d", resp.StatusCode)
	}
	var rows []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("lead %q not found", leadID)
	}
	row := rows[0]
	state := ""
	for _, key := range []string{"state", "lead_status", "pipeline_stage", "status"} {
		if value, ok := row[key].(string); ok && strings.TrimSpace(value) != "" {
			state = value
			break
		}
	}
	return row, state, nil
}
