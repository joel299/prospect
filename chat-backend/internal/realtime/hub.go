package realtime

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]struct{}
}

func New() *Hub { return &Hub{clients: map[*websocket.Conn]struct{}{}} }

func (h *Hub) Publish(e Event) {
	b, err := json.Marshal(e)
	if err != nil {
		return
	}
	h.mu.RLock()
	clients := make([]*websocket.Conn, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()
	for _, c := range clients {
		_ = c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
			h.remove(c)
		}
	}
}

func (h *Hub) remove(c *websocket.Conn) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		_ = c.Close()
	}
	h.mu.Unlock()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
	defer h.remove(c)
	for {
		if _, _, err := c.ReadMessage(); err != nil {
			return
		}
	}
}
