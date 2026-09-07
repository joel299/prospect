package omniroute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iainfinito/chat-backend/internal/domain"
)

func TestGenerateParsesOpenAIResponseAndFullEndpoint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer demo-key" {
			t.Fatalf("auth=%q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"resposta OmniRoute"}}]}`))
	}))
	defer srv.Close()
	out, err := (Client{BaseURL: srv.URL + "/v1/chat/completions", APIKey: "demo-key", Model: "auto", HTTP: srv.Client()}).Generate(context.Background(), domain.AgentRequest{User: "oi"})
	if err != nil || out.Content != "resposta OmniRoute" {
		t.Fatalf("out=%#v err=%v", out, err)
	}
}

func TestParseSSE(t *testing.T) {
	resp := httptest.NewRecorder()
	resp.Header().Set("Content-Type", "text/event-stream")
	_ = resp
	out, err := parseSSE(strings.NewReader(`data: {"choices":[{"delta":{"content":"Oi"}}]}

data: [DONE]
`))
	if err != nil || out.Content != "Oi" {
		t.Fatalf("out=%#v err=%v", out, err)
	}
}
