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
	mu     sync.RWMutex
	conv   map[string]domain.Conversation
	msgs   map[string][]domain.Message
	idem   map[string]string
	memory map[string]AgentMemory
}
type AgentMemory struct {
	Summary string
	Facts   map[string]string
}
type LeadContextProvider interface {
	Get(context.Context, string) (map[string]any, string, error)
}
type BufferHistoryProvider interface {
	Recent(context.Context, string, int) ([]domain.Message, error)
}
type AgentMemoryProvider interface {
	Load(context.Context, string) (AgentMemory, error)
	Save(context.Context, string, AgentMemory) error
}
type Service struct {
	Store *Store
	Ryze  ryze.Client
	LLM   interface {
		Generate(context.Context, domain.AgentRequest) (domain.AgentResponse, error)
	}
	Tools  composio.ToolExecutor
	Hub    *realtime.Hub
	Leads  LeadContextProvider
	Buffer BufferHistoryProvider
	Memory AgentMemoryProvider
}

func New(r ryze.Client, llm Generator, t composio.ToolExecutor, h *realtime.Hub) *Service {
	return &Service{Store: &Store{conv: map[string]domain.Conversation{}, msgs: map[string][]domain.Message{}, idem: map[string]string{}, memory: map[string]AgentMemory{}}, Ryze: r, LLM: llm, Tools: t, Hub: h}
}

func (s *Service) SetContextProviders(leads LeadContextProvider, buffer BufferHistoryProvider) {
	s.Leads, s.Buffer = leads, buffer
}

func (s *Service) SetMemoryProvider(memory AgentMemoryProvider) { s.Memory = memory }

// TestAgent exercises the configured LLM without sending an outbound message.
func (s *Service) TestAgent(ctx context.Context, req domain.AgentRequest) (domain.AgentResponse, error) {
	if s.LLM == nil {
		return domain.AgentResponse{}, fmt.Errorf("omniroute is not configured")
	}
	if strings.TrimSpace(req.User) == "" {
		return domain.AgentResponse{}, fmt.Errorf("message is required")
	}
	if strings.TrimSpace(req.System) == "" {
		req.System = "Você é um agente de teste. Responda de forma curta e objetiva."
	}
	return s.LLM.Generate(ctx, req)
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
	if _, ok := s.Store.conv[in.ConversationID]; !ok {
		s.Store.conv[in.ConversationID] = domain.Conversation{ID: in.ConversationID, LeadID: in.LeadID, Channel: "ryze", Provider: "ryze", AgentEnabled: true, CreatedAt: now, UpdatedAt: now}
	}
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
	if len(history) == 0 {
		return
	}
	current := history[len(history)-1].Content
	lead := map[string]any{"id": lid}
	state := "(não carregado)"
	if s.Leads != nil {
		if fields, leadState, err := s.Leads.Get(ctx, lid); err == nil {
			lead = fields
			state = leadState
		}
	}
	recent := history
	if s.Buffer != nil {
		if buffered, err := s.Buffer.Recent(ctx, lid, 20); err == nil && len(buffered) > 0 {
			recent = buffered
			for i := len(buffered) - 1; i >= 0; i-- {
				if buffered[i].Direction == "inbound" && strings.TrimSpace(buffered[i].Content) != "" {
					current = buffered[i].Content
					break
				}
			}
		}
	}
	mem := AgentMemory{}
	if s.Memory != nil {
		mem, _ = s.Memory.Load(ctx, lid)
	} else {
		s.Store.mu.RLock()
		mem = s.Store.memory[lid]
		s.Store.mu.RUnlock()
	}
	prompt := BuildPromptUser(PromptContext{Lead: lead, State: state, ConversationSummary: mem.Summary, RecentMessages: recent, CurrentMessage: current})
	r, e := s.LLM.Generate(ctx, domain.AgentRequest{System: PromptSystem, User: prompt, PromptCache: PromptSystem})
	if e != nil || strings.TrimSpace(r.Content) == "" {
		return
	}
	if len(strings.Split(r.Content, "\n")) > 3 {
		return
	}
	_, _ = s.Send(ctx, domain.OutboundMessage{ConversationID: cid, LeadID: lid, Content: r.Content})
	// Keep a bounded, provider-independent memory snapshot. Production should inject
	// a PostgreSQL-backed AgentMemoryProvider; the map remains test/development only.
	mem.Summary = summarize(recent)
	if s.Memory != nil {
		_ = s.Memory.Save(ctx, lid, mem)
	} else {
		s.Store.mu.Lock()
		s.Store.memory[lid] = mem
		s.Store.mu.Unlock()
	}
}

func summarize(messages []domain.Message) string {
	if len(messages) == 0 {
		return ""
	}
	start := 0
	if len(messages) > 6 {
		start = len(messages) - 6
	}
	var b strings.Builder
	for _, m := range messages[start:] {
		fmt.Fprintf(&b, "%s: %s\n", m.Direction, strings.TrimSpace(m.Content))
	}
	return strings.TrimSpace(b.String())
}
func (s *Service) Send(ctx context.Context, out domain.OutboundMessage) (domain.Message, error) {
	if strings.TrimSpace(out.Content) == "" || out.ConversationID == "" || strings.TrimSpace(out.LeadID) == "" {
		return domain.Message{}, fmt.Errorf("invalid outbound message")
	}
	m := domain.Message{ID: hash(out.ConversationID + out.Content + time.Now().String()), ConversationID: out.ConversationID, LeadID: out.LeadID, Provider: "ryze", Direction: "outbound", Content: out.Content, Status: "queued", CreatedAt: time.Now().UTC()}
	if s.Ryze != nil {
		externalID, err := s.Ryze.SendText(ctx, out)
		if err != nil {
			return domain.Message{}, err
		}
		m.ExternalID = externalID
		m.Status = "accepted"
	}
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
func (s *Service) Conversations() []domain.Conversation {
	s.Store.mu.RLock()
	defer s.Store.mu.RUnlock()
	out := make([]domain.Conversation, 0, len(s.Store.conv))
	for _, c := range s.Store.conv {
		out = append(out, c)
	}
	return out
}
func hash(v string) string { b := sha256.Sum256([]byte(v)); return hex.EncodeToString(b[:]) }
