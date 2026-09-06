package ryze

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iainfinito/chat-backend/internal/domain"
)

func TestHTTPClientSendTextUsesRyzeContract(t *testing.T) {
	var gotToken string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("token")
		if r.Method != http.MethodPost || r.URL.Path != "/api/message/text/instance" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["number"] != "5511999999999" || body["message"] != "oi" {
			t.Fatalf("body: %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"messageId": "m-1"})
	}))
	defer srv.Close()

	id, err := (HTTPClient{BaseURL: srv.URL, APIKey: "secret", Instance: "instance", HTTP: srv.Client()}).SendText(context.Background(), domain.OutboundMessage{LeadID: "5511999999999", Content: "oi"})
	if err != nil || id != "m-1" || gotToken != "secret" {
		t.Fatalf("id=%q err=%v token=%q", id, err, gotToken)
	}
}

func TestHTTPClientHistoryParsesMessages(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"messages": []any{map[string]any{"id": "m-1", "message": "oi", "fromMe": false}}})
	}))
	defer srv.Close()
	ms, err := (HTTPClient{BaseURL: srv.URL, APIKey: "secret", Instance: "instance", HTTP: srv.Client()}).History(context.Background(), "5511999999999")
	if err != nil || len(ms) != 1 || ms[0].Content != "oi" {
		t.Fatalf("messages=%#v err=%v", ms, err)
	}
}
