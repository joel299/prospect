package domain

import "time"

type Conversation struct {
	ID               string    `json:"id"`
	LeadID           string    `json:"lead_id"`
	Channel          string    `json:"channel"`
	Provider         string    `json:"provider"`
	AgentEnabled     bool      `json:"agent_enabled"`
	MeetingScheduled bool      `json:"meeting_scheduled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	LeadID         string    `json:"lead_id"`
	Provider       string    `json:"provider"`
	ExternalID     string    `json:"external_id"`
	Direction      string    `json:"direction"`
	Content        string    `json:"content"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
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
