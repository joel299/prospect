package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/iainfinito/chat-backend/internal/domain"
	"github.com/iainfinito/chat-backend/internal/realtime"
	"github.com/iainfinito/chat-backend/internal/service"
	"net/http"
	"strings"
)

type API struct {
	S *service.Service
	H *realtime.Hub
}

func (a API) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("GET /health", a.health)
	m.HandleFunc("POST /api/v1/webhooks/ryze", a.inbound)
	m.HandleFunc("GET /api/v1/conversations/{id}/messages", a.messages)
	m.HandleFunc("POST /api/v1/conversations/{id}/send", a.send)
	m.Handle("/ws/v1", a.H)
	return requestID(m)
}
func (a API) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]string{"status": "ok", "service": "chat-backend"})
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
func (a API) send(w http.ResponseWriter, r *http.Request) {
	var in struct{ LeadID, Content string }
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
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func problem(w http.ResponseWriter, status int, code, msg string, details any) {
	write(w, status, map[string]any{"error": map[string]any{"code": code, "message": msg, "details": details}})
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
