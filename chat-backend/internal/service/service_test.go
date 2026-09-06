package service

import (
	"context"
	"github.com/iainfinito/chat-backend/internal/domain"
	"github.com/iainfinito/chat-backend/internal/providers/composio"
	"github.com/iainfinito/chat-backend/internal/providers/omniroute"
	"github.com/iainfinito/chat-backend/internal/providers/ryze"
	"github.com/iainfinito/chat-backend/internal/realtime"
	"testing"
)

type fakeGen struct{ n int }

func (f *fakeGen) Generate(_ context.Context, r domain.AgentRequest) (domain.AgentResponse, error) {
	f.n++
	return domain.AgentResponse{Content: "resposta curta"}, nil
}
func TestInboundIsIdempotentAndPublishes(t *testing.T) {
	h := realtime.New()
	s := New(ryze.HTTPClient{}, omniroute.Client{}, composio.Client{}, h)
	m, e := s.Inbound(context.Background(), domain.InboundMessage{ConversationID: "c1", LeadID: "l1", ExternalID: "x1", Content: "oi"})
	if e != nil || m.Direction != "inbound" {
		t.Fatalf("inbound: %#v %v", m, e)
	}
	if _, e = s.Inbound(context.Background(), domain.InboundMessage{ConversationID: "c1", LeadID: "l1", ExternalID: "x1", Content: "oi"}); e != ErrDuplicate {
		t.Fatalf("expected duplicate, got %v", e)
	}
}
func TestOutboundValidation(t *testing.T) {
	s := New(nil, nil, nil, nil)
	if _, e := s.Send(context.Background(), domain.OutboundMessage{ConversationID: "c", Content: " "}); e == nil {
		t.Fatal("expected validation error")
	}
}
