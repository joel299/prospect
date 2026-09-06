package main

import (
	"github.com/iainfinito/chat-backend/internal/httpapi"
	"github.com/iainfinito/chat-backend/internal/providers/composio"
	"github.com/iainfinito/chat-backend/internal/providers/omniroute"
	"github.com/iainfinito/chat-backend/internal/providers/ryze"
	"github.com/iainfinito/chat-backend/internal/realtime"
	"github.com/iainfinito/chat-backend/internal/service"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	h := realtime.New()
	llm := omniroute.Client{BaseURL: os.Getenv("OMNIROUTE_BASE_URL"), APIKey: os.Getenv("OMNIROUTE_API_KEY"), Model: os.Getenv("OMNIROUTE_MODEL"), HTTP: &http.Client{Timeout: 30 * time.Second}}
	s := service.New(ryze.HTTPClient{BaseURL: os.Getenv("RYZE_BASE_URL"), APIKey: os.Getenv("RYZE_API_KEY"), Instance: os.Getenv("RYZE_INSTANCE")}, llm, composio.Client{APIKey: os.Getenv("COMPOSIO_API_KEY")}, h)
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("chat-backend listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, httpapi.API{S: s, H: h}.Handler()))
}
