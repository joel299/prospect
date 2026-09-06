package realtime

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
type Hub struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

func New() *Hub { return &Hub{clients: map[chan []byte]struct{}{}} }
func (h *Hub) Publish(e Event) {
	b, _ := json.Marshal(e)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		select {
		case c <- b:
		default:
		}
	}
}
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "websocket upgrade required", 426)
		return
	}
	http.Error(w, "websocket transport requires a WebSocket-capable edge adapter", 501)
}
