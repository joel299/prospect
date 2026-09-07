package main

import (
	"github.com/iainfinito/chat-backend/internal/httpapi"
	"github.com/iainfinito/chat-backend/internal/providers/buffer"
	"github.com/iainfinito/chat-backend/internal/providers/composio"
	"github.com/iainfinito/chat-backend/internal/providers/lead"
	"github.com/iainfinito/chat-backend/internal/providers/omniroute"
	"github.com/iainfinito/chat-backend/internal/providers/ryze"
	"github.com/iainfinito/chat-backend/internal/realtime"
	"github.com/iainfinito/chat-backend/internal/service"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func envFirst(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func main() {
	h := realtime.New()
	omniBaseURL := envFirst("OMNIROUTE_BASE_URL", "OMNIROUTER_BASE_URL")
	if omniBaseURL == "" {
		omniBaseURL = "https://omnirouter.iainfinito.com.br"
	}
	omniModel := envFirst("OMNIROUTE_MODEL", "OMNIROUTER_MODEL")
	if omniModel == "" {
		omniModel = "auto"
	}
	llm := omniroute.Client{
		BaseURL: omniBaseURL,
		APIKey:  envFirst("OMNIROUTE_API_KEY", "OMNIROUTER_APIKEY"),
		Model:   omniModel,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
	ryzeBaseURL := envFirst("RYZE_BASE_URL", "RYZE_API_BASE_URL")
	if ryzeBaseURL == "" {
		ryzeBaseURL = "https://ryzeapi.cloud"
	}
	ryzeClient := ryze.HTTPClient{
		BaseURL:  ryzeBaseURL,
		APIKey:   envFirst("RYZE_API_KEY", "RAYZE_APIKEY", "Token_Instance"),
		Instance: envFirst("RYZE_INSTANCE", "Instance_Name"),
		HTTP:     &http.Client{Timeout: 30 * time.Second},
	}
	composioClient := composio.Client{BaseURL: envFirst("COMPOSIO_BASE_URL"), APIKey: envFirst("COMPOSIO_API_KEY", "COMPOSIO_APIKEY")}
	s := service.New(ryzeClient, llm, composioClient, h)
	if historyURL := envFirst("BUFFER_HISTORY_URL"); historyURL != "" {
		s.SetContextProviders(
			lead.SupabaseClient{BaseURL: envFirst("SUPABASE_URL"), APIKey: envFirst("SUPABASE_SERVICE_ROLE_KEY", "SUPABASE_ANON_KEY"), Table: envFirst("SUPABASE_LEADS_TABLE"), HTTP: &http.Client{Timeout: 10 * time.Second}},
			buffer.HistoryClient{HistoryURL: historyURL, AccessToken: envFirst("BUFFER_ACCESS_TOKEN"), HTTP: &http.Client{Timeout: 10 * time.Second}},
		)
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("chat-backend listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, httpapi.API{S: s, H: h, Tools: composioClient}.Handler()))
}
