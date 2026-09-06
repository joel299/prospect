package domain

import "time"

type Conversation struct {
	ID, LeadID, Channel, Provider  string
	AgentEnabled, MeetingScheduled bool
	CreatedAt, UpdatedAt           time.Time
}
type Message struct {
	ID, ConversationID, LeadID, Provider, ExternalID, Direction, Content, Status string
	CreatedAt                                                                    time.Time
}
type InboundMessage struct {
	ConversationID string `json:"conversation_id"`
	LeadID         string `json:"lead_id"`
	ExternalID     string `json:"external_id"`
	Content        string `json:"content"`
}
type OutboundMessage struct{ ConversationID, LeadID, Content string }
type ToolCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}
type AgentRequest struct {
	System, User string
	History      []Message
	Tools        []string
}
type AgentResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}
