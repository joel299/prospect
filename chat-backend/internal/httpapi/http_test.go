package httpapi

import (
	"context"
	"github.com/iainfinito/chat-backend/internal/domain"
	"github.com/iainfinito/chat-backend/internal/providers/composio"
	"github.com/iainfinito/chat-backend/internal/providers/omniroute"
	"github.com/iainfinito/chat-backend/internal/providers/ryze"
	"github.com/iainfinito/chat-backend/internal/realtime"
	"github.com/iainfinito/chat-backend/internal/service"
	"net/http/httptest"
	"strings"
	"testing"
)

type testGenerator struct{}

func (testGenerator) Generate(_ context.Context, req domain.AgentRequest) (domain.AgentResponse, error) {
	return domain.AgentResponse{Content: "API_OK: " + req.User}, nil
}

func TestAgentEndpointDoesNotSendOutbound(t *testing.T) {
	h := realtime.New()
	s := service.New(ryze.HTTPClient{}, testGenerator{}, composio.Client{}, h)
	srv := httptest.NewServer((API{S: s, H: h}).Handler())
	defer srv.Close()
	r, e := srv.Client().Post(srv.URL+"/api/v1/test/agent", "application/json", strings.NewReader(`{"message":"teste"}`))
	if e != nil || r.StatusCode != 200 {
		t.Fatalf("agent test %v %v", r.StatusCode, e)
	}
}

func TestHealthAndInbound(t *testing.T) {
	h := realtime.New()
	s := service.New(ryze.HTTPClient{}, omniroute.Client{}, composio.Client{}, h)
	srv := httptest.NewServer((API{S: s, H: h}).Handler())
	defer srv.Close()
	r, e := srv.Client().Get(srv.URL + "/health")
	if e != nil || r.StatusCode != 200 {
		t.Fatalf("health %v %v", r.StatusCode, e)
	}
	r, e = srv.Client().Post(srv.URL+"/api/v1/webhooks/ryze", "application/json", strings.NewReader(`{"conversation_id":"c","external_id":"x","content":"oi"}`))
	if e != nil || r.StatusCode != 202 {
		t.Fatalf("inbound %v %v", r.StatusCode, e)
	}
}
