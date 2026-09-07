package httpapi

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/iainfinito/chat-backend/internal/domain"
	"github.com/iainfinito/chat-backend/internal/providers/composio"
	"github.com/iainfinito/chat-backend/internal/realtime"
	"github.com/iainfinito/chat-backend/internal/service"
	"net/http"
	"strings"
)

type API struct {
	S     *service.Service
	H     *realtime.Hub
	Tools composio.ToolExecutor
}

//go:embed openapi.yaml
var openAPISpec []byte

func (a API) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /health", a.health)
	m.HandleFunc("GET /openapi.yaml", a.openapi)
	m.HandleFunc("GET /docs", a.scalar)
	m.HandleFunc("POST /api/v1/webhooks/ryze", a.inbound)
	m.HandleFunc("GET /api/v1/conversations", a.conversations)
	m.HandleFunc("GET /api/v1/conversations/{id}/messages", a.messages)
	m.HandleFunc("POST /api/v1/conversations/{id}/send", a.send)
	m.HandleFunc("POST /api/v1/integrations/trello/sync", a.trelloSync)
	m.HandleFunc("POST /api/v1/test/agent", a.agentTest)
	m.HandleFunc("POST /api/v1/agent/test", a.agentTest)
	m.Handle("/ws/v1", a.H)
	return cors(requestID(m))
}
func (a API) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]string{"status": "ok", "service": "chat-backend"})
}
func (a API) openapi(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPISpec)
}
func (a API) scalar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<!doctype html><html><head><title>Chat Prospect API</title></head><body><script id="api-reference" data-url="%s/openapi.yaml"></script><script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script></body></html>`, strings.TrimRight(publicBaseURL(r), "/"))
}
func publicBaseURL(r *http.Request) string {
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme + "://" + r.Host
	}
	return "http://" + r.Host
}
func (a API) agentTest(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Message string `json:"message"`
		System  string `json:"system,omitempty"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&payload) != nil {
		problem(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid JSON", nil)
		return
	}
	response, err := a.S.TestAgent(r.Context(), domain.AgentRequest{System: payload.System, User: payload.Message})
	if err != nil {
		problem(w, http.StatusBadGateway, "OMNIROUTE_ERROR", err.Error(), nil)
		return
	}
	write(w, http.StatusOK, map[string]any{"content": response.Content, "tool_calls": response.ToolCalls})
}
func (a API) inbound(w http.ResponseWriter, r *http.Request) {
	var in domain.InboundMessage
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&in) != nil {
		problem(w, 422, "VALIDATION_ERROR", "invalid JSON", nil)
		return
	}
	m, e := a.S.Inbound(r.Context(), in)
	if e == service.ErrDuplicate {
		write(w, 200, map[string]any{"duplicate": true})
		return
	}
	if e != nil {
		problem(w, 422, "VALIDATION_ERROR", e.Error(), nil)
		return
	}
	write(w, 202, m)
}
func (a API) messages(w http.ResponseWriter, r *http.Request) {
	write(w, 200, map[string]any{"messages": a.S.Messages(r.PathValue("id"))})
}
func (a API) conversations(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]any{"conversations": a.S.Conversations()})
}
func (a API) send(w http.ResponseWriter, r *http.Request) {
	var in struct {
		LeadID  string `json:"lead_id"`
		Content string `json:"content"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		problem(w, 422, "VALIDATION_ERROR", "invalid JSON", nil)
		return
	}
	m, e := a.S.Send(r.Context(), domain.OutboundMessage{ConversationID: r.PathValue("id"), LeadID: in.LeadID, Content: in.Content})
	if e != nil {
		problem(w, 422, "VALIDATION_ERROR", e.Error(), nil)
		return
	}
	write(w, 202, m)
}
func (a API) trelloSync(w http.ResponseWriter, r *http.Request) {
	if a.Tools == nil {
		problem(w, http.StatusServiceUnavailable, "COMPOSIO_NOT_CONFIGURED", "Composio executor is not configured", nil)
		return
	}
	var payload map[string]any
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&payload) != nil {
		problem(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid JSON", nil)
		return
	}
	if strings.TrimSpace(fmt.Sprint(payload["lead_id"])) == "" {
		problem(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "lead_id is required", nil)
		return
	}
	result, err := a.Tools.Execute(r.Context(), "trello_sync_lead", payload)
	if err != nil {
		problem(w, http.StatusBadGateway, "COMPOSIO_ERROR", err.Error(), nil)
		return
	}
	write(w, http.StatusAccepted, map[string]any{"status": "accepted", "provider": "composio", "result": result})
}

func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func problem(w http.ResponseWriter, status int, code, msg string, details any) {
	write(w, status, map[string]any{"error": map[string]any{"code": code, "message": msg, "details": details}})
}
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "https://chat-prospect.iainfinito.com.br" || origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = "generated"
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

var _ = fmt.Sprintf
var _ = strings.TrimSpace
