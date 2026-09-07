package ryze

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iainfinito/chat-backend/internal/domain"
)

type EventListener struct {
	BaseURL, APIKey, Instance string
	HTTP                      *http.Client
	Dialer                    *websocket.Dialer
}

func (l EventListener) Configure(ctx context.Context) error {
	if l.BaseURL == "" || l.APIKey == "" || l.Instance == "" {
		return fmt.Errorf("ryze websocket not configured")
	}
	body := strings.NewReader(`{"enabled":true,"events":["message.exchange","message.status"],"mediaBase64":true}`)
	u := strings.TrimRight(l.BaseURL, "/") + "/api/events/websocket/" + url.PathEscape(l.Instance)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return err
	}
	req.Header.Set("token", l.APIKey)
	req.Header.Set("Content-Type", "application/json")
	h := l.HTTP
	if h == nil {
		h = http.DefaultClient
	}
	resp, err := h.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ryze websocket config status %d", resp.StatusCode)
	}
	log.Printf("ryze websocket configured instance=%s", l.Instance)
	return nil
}

func (l EventListener) Listen(ctx context.Context, onMessage func(domain.InboundMessage)) error {
	if l.APIKey == "" || l.Instance == "" {
		return fmt.Errorf("ryze websocket not configured")
	}
	dialer := l.Dialer
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}
	backoff := time.Second
	for {
		u := "wss://ryzeapi.cloud/ws/" + url.PathEscape(l.Instance) + "?token=" + url.QueryEscape(l.APIKey)
		if strings.HasPrefix(strings.ToLower(l.BaseURL), "http://") {
			u = "ws://" + strings.TrimPrefix(strings.TrimPrefix(l.BaseURL, "http://"), "https://") + "/ws/" + url.PathEscape(l.Instance) + "?token=" + url.QueryEscape(l.APIKey)
		}
		conn, _, err := dialer.DialContext(ctx, u, nil)
		if err == nil {
			log.Printf("ryze websocket connected instance=%s", l.Instance)
			backoff = time.Second
			err = l.read(ctx, conn, onMessage)
			_ = conn.Close()
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("ryze websocket read stopped: %v", err)
		} else {
			log.Printf("ryze websocket dial failed: %v", err)
		}
		jitter := time.Duration(rand.Int63n(int64(backoff/2 + 1)))
		wait := backoff + jitter
		if wait > 30*time.Second {
			wait = 30 * time.Second
		}
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
		backoff *= 2
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
		_ = err
	}
}

func (l EventListener) read(ctx context.Context, conn *websocket.Conn, onMessage func(domain.InboundMessage)) error {
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		var envelope struct {
			Event string          `json:"event"`
			Data  json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			log.Printf("ryze websocket invalid frame bytes=%d", len(raw))
			continue
		}
		log.Printf("ryze websocket frame event=%s bytes=%d", envelope.Event, len(raw))
		if envelope.Event != "message.exchange" {
			continue
		}
		var data struct {
			ID        string `json:"id"`
			FromMe    *bool  `json:"fromMe"`
			IsFromMe  *bool  `json:"isFromMe"`
			Direction string `json:"direction"`
			Message   struct {
				ID        string `json:"id"`
				Direction string `json:"direction"`
				FromMe    *bool  `json:"fromMe"`
				Chat      struct {
					JID string `json:"jid"`
				} `json:"chat"`
				Content struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(envelope.Data, &data) != nil {
			continue
		}
		m := data.Message
		id := m.ID
		if id == "" {
			id = data.ID
		}
		fromMe := false
		if data.FromMe != nil {
			fromMe = *data.FromMe
		}
		if m.FromMe != nil {
			fromMe = *m.FromMe
		}
		if strings.EqualFold(m.Direction, "outgoing") || strings.EqualFold(data.Direction, "outgoing") {
			fromMe = true
		}
		if data.IsFromMe != nil {
			fromMe = *data.IsFromMe
		}
		phone := strings.TrimSpace(m.Chat.JID)
		if phone == "" {
			continue
		}
		phone = strings.Split(phone, "@")[0]
		content := strings.TrimSpace(m.Content.Text)
		if content == "" || id == "" {
			continue
		}
		log.Printf("ryze websocket event id=%s lead=%s fromMe=%t", id, phone, fromMe)
		onMessage(domain.InboundMessage{ConversationID: "ryze:" + phone, LeadID: phone, ExternalID: id, Content: content, FromMe: fromMe})
		_ = ctx
	}
}
