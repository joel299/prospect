package service

import (
	"context"
	"testing"
	"time"

	"github.com/iainfinito/chat-backend/internal/domain"
	"github.com/iainfinito/chat-backend/internal/realtime"
)

type e2eLLM struct{}

func (e2eLLM) Generate(context.Context, domain.AgentRequest) (domain.AgentResponse, error) {
	return domain.AgentResponse{Content: "resposta E2E"}, nil
}

type e2eRyze struct{}

func (e2eRyze) SendText(context.Context, domain.OutboundMessage) (string, error) {
	return "out-e2e", nil
}
func (e2eRyze) History(context.Context, string) ([]domain.Message, error) { return nil, nil }

func TestConversationVerticalSliceE2E(t *testing.T) {
	s := New(e2eRyze{}, e2eLLM{}, nil, realtime.New())
	_, err := s.Inbound(context.Background(), domain.InboundMessage{ConversationID: "conv-e2e", LeadID: "5511999999999", ExternalID: "in-e2e", Content: "Olá, tenho interesse"})
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if len(s.Messages("conv-e2e")) == 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	messages := s.Messages("conv-e2e")
	if len(messages) != 2 {
		t.Fatalf("expected inbound + outbound, got %d: %#v", len(messages), messages)
	}
	if messages[1].Direction != "outbound" || messages[1].Content != "resposta E2E" {
		t.Fatalf("unexpected outbound: %#v", messages[1])
	}
}
