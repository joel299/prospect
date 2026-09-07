package service

import (
	"strings"
	"testing"

	"github.com/iainfinito/chat-backend/internal/domain"
)

func TestBuildPromptUserUsesRuntimeContext(t *testing.T) {
	got := BuildPromptUser(PromptContext{
		Lead:                map[string]any{"id": "lead-7", "name": "Ana"},
		State:               "qualified",
		ConversationSummary: "Busca automação",
		RecentMessages:      []domain.Message{{Direction: "inbound", Content: "Quero entender valores"}},
		CurrentMessage:      "Podemos falar amanhã?",
	})
	for _, want := range []string{"lead-7", "qualified", "Busca automação", "Quero entender valores", "Podemos falar amanhã?"} {
		if !strings.Contains(got, want) {
			t.Fatalf("prompt missing %q: %s", want, got)
		}
	}
}
