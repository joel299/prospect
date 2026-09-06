package realtime

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketReceivesPublishedEvent(t *testing.T) {
	h := New()
	srv := httptest.NewServer(h)
	defer srv.Close()
	url := "ws" + srv.URL[len("http"):]
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	h.Publish(Event{Type: "test.event", Data: map[string]string{"ok": "true"}})
	_ = c.SetReadDeadline(time.Now().Add(time.Second))
	_, msg, err := c.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) != `{"type":"test.event","data":{"ok":"true"}}` {
		t.Fatalf("unexpected event: %s", msg)
	}
}
