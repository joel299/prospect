package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/iainfinito/chat-backend/internal/domain"
	"github.com/iainfinito/chat-backend/internal/providers/composio"

	"github.com/iainfinito/chat-backend/internal/providers/ryze"
	"github.com/iainfinito/chat-backend/internal/realtime"
	"strings"
	"sync"
	"time"
)

var ErrDuplicate = errors.New("duplicate message")

type Store struct {
	mu   sync.RWMutex
	conv map[string]domain.Conversation
	msgs map[string][]domain.Message
	idem map[string]string
}
type Service struct {
	Store *Store
	Ryze  ryze.Client
	LLM   interface {
		Generate(context.Context, domain.AgentRequest) (domain.AgentResponse, error)
	}
	Tools composio.ToolExecutor
	Hub   *realtime.Hub
}

func New(r ryze.Client, llm Generator, t composio.ToolExecutor, h *realtime.Hub) *Service {
	return &Service{Store: &Store{conv: map[string]domain.Conversation{}, msgs: map[string][]domain.Message{}, idem: map[string]string{}}, Ryze: r, LLM: llm, Tools: t, Hub: h}
}
func (s *Service) Inbound(ctx context.Context, in domain.InboundMessage) (domain.Message, error) {
	if strings.TrimSpace(in.Content) == "" || in.ConversationID == "" || in.ExternalID == "" {
		return domain.Message{}, fmt.Errorf("invalid inbound message")
	}
	s.Store.mu.Lock()
	defer s.Store.mu.Unlock()
	if _, ok := s.Store.idem["inbound:"+in.ExternalID]; ok {
		return domain.Message{}, ErrDuplicate
	}
	now := time.Now().UTC()
	m := domain.Message{ID: hash(in.ExternalID), ConversationID: in.ConversationID, LeadID: in.LeadID, Provider: "ryze", ExternalID: in.ExternalID, Direction: "inbound", Content: in.Content, Status: "received", CreatedAt: now}
	s.Store.msgs[in.ConversationID] = append(s.Store.msgs[in.ConversationID], m)
	s.Store.idem["inbound:"+in.ExternalID] = m.ID
	if s.Hub != nil {
		s.Hub.Publish(realtime.Event{Type: "conversation.message.inbound", Data: m})
	}
	go s.respond(context.Background(), in.ConversationID, in.LeadID)
	return m, nil
}
func (s *Service) respond(ctx context.Context, cid, lid string) {
	s.Store.mu.RLock()
	history := append([]domain.Message(nil), s.Store.msgs[cid]...)
	s.Store.mu.RUnlock()
	if s.LLM == nil {
		return
	}
	r, e := s.LLM.Generate(ctx, domain.AgentRequest{System: "Você é um SDR. Responda em no máximo 3 linhas, sem inventar dados.", User: history[len(history)-1].Content, History: history})
	if e != nil || strings.TrimSpace(r.Content) == "" {
		return
	}
	if len(strings.Split(r.Content, "\n")) > 3 {
		return
	}
	_, _ = s.Send(ctx, domain.OutboundMessage{ConversationID: cid, LeadID: lid, Content: r.Content})
}
func (s *Service) Send(_ context.Context, out domain.OutboundMessage) (domain.Message, error) {
	if strings.TrimSpace(out.Content) == "" || out.ConversationID == "" {
		return domain.Message{}, fmt.Errorf("invalid outbound message")
	}
	m := domain.Message{ID: hash(out.ConversationID + out.Content + time.Now().String()), ConversationID: out.ConversationID, LeadID: out.LeadID, Provider: "ryze", Direction: "outbound", Content: out.Content, Status: "queued", CreatedAt: time.Now().UTC()}
	s.Store.mu.Lock()
	s.Store.msgs[out.ConversationID] = append(s.Store.msgs[out.ConversationID], m)
	s.Store.mu.Unlock()
	if s.Hub != nil {
		s.Hub.Publish(realtime.Event{Type: "conversation.message.outbound", Data: m})
	}
	return m, nil
}
func (s *Service) Messages(cid string) []domain.Message {
	s.Store.mu.RLock()
	defer s.Store.mu.RUnlock()
	return append([]domain.Message(nil), s.Store.msgs[cid]...)
}
func hash(v string) string { b := sha256.Sum256([]byte(v)); return hex.EncodeToString(b[:]) }
