package omniroute

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/iainfinito/chat-backend/internal/domain"
)

type Client struct {
	BaseURL, APIKey, Model string
	HTTP                   *http.Client
}

func (c Client) Generate(ctx context.Context, r domain.AgentRequest) (domain.AgentResponse, error) {
	if strings.TrimSpace(c.BaseURL) == "" {
		return domain.AgentResponse{}, fmt.Errorf("omniroute not configured")
	}
	b, _ := json.Marshal(map[string]any{"model": c.Model, "system": r.System, "user": r.User, "history": r.History, "tools": r.Tools})
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, endpoint(c.BaseURL), bytes.NewReader(b))
	if e != nil {
		return domain.AgentResponse{}, e
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", bearer(c.APIKey))
	}
	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	resp, e := client.Do(req)
	if e != nil {
		return domain.AgentResponse{}, e
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return domain.AgentResponse{}, fmt.Errorf("omniroute status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return parseResponse(resp)
}

func endpoint(base string) string {
	base = strings.TrimRight(base, "/")
	if strings.HasSuffix(base, "/v1/chat/completions") {
		return base
	}
	if strings.HasSuffix(base, "/v1") {
		return base + "/chat/completions"
	}
	return base + "/v1/chat/completions"
}

func bearer(key string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(key)), "bearer ") {
		return key
	}
	return "Bearer " + key
}

func parseResponse(resp *http.Response) (domain.AgentResponse, error) {
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "text/event-stream") {
		return parseSSE(resp.Body)
	}
	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return domain.AgentResponse{}, err
	}
	return parseJSON(raw)
}

func parseJSON(raw []byte) (domain.AgentResponse, error) {
	var direct domain.AgentResponse
	if json.Unmarshal(raw, &direct) == nil && (direct.Content != "" || len(direct.ToolCalls) > 0) {
		return direct, nil
	}
	var envelope struct {
		Choices []struct {
			Message struct {
				Content   string            `json:"content"`
				ToolCalls []domain.ToolCall `json:"tool_calls"`
			} `json:"message"`
			Text string `json:"text"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return domain.AgentResponse{}, err
	}
	if len(envelope.Choices) == 0 {
		return domain.AgentResponse{}, fmt.Errorf("omniroute response has no choices")
	}
	content := envelope.Choices[0].Message.Content
	if content == "" {
		content = envelope.Choices[0].Text
	}
	return domain.AgentResponse{Content: content, ToolCalls: envelope.Choices[0].Message.ToolCalls}, nil
}

func parseSSE(r io.Reader) (domain.AgentResponse, error) {
	var content strings.Builder
	var calls []domain.ToolCall
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" || data == "[DONE]" {
			continue
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content   string            `json:"content"`
					ToolCalls []domain.ToolCall `json:"tool_calls"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if json.Unmarshal([]byte(data), &chunk) != nil || len(chunk.Choices) == 0 {
			continue
		}
		content.WriteString(chunk.Choices[0].Delta.Content)
		calls = append(calls, chunk.Choices[0].Delta.ToolCalls...)
	}
	if err := s.Err(); err != nil {
		return domain.AgentResponse{}, err
	}
	if content.Len() == 0 && len(calls) == 0 {
		return domain.AgentResponse{}, fmt.Errorf("omniroute SSE contained no content")
	}
	return domain.AgentResponse{Content: content.String(), ToolCalls: calls}, nil
}
