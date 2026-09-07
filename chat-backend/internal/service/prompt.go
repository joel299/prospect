package service

import (
	"fmt"
	"strings"

	"github.com/iainfinito/chat-backend/internal/domain"
)

const PromptSystem = `Você é um SDR comercial que atua em conversas de WhatsApp após o lead ter respondido à abordagem inicial.
Nunca trate a conversa como prospecção fria nova quando existir histórico.
Apresente-se somente uma vez quando não houver histórico.
Seu objetivo é entender a necessidade, responder objetivamente, reduzir objeções e avançar para o próximo passo comercial, preferencialmente uma reunião quando houver intenção suficiente.
REGRAS:
- máximo 3 linhas;
- uma única resposta por turno;
- máximo 4 perguntas claras antes de tentar reunião;
- não repetir informação já existente;
- não inventar prova, urgência, preço ou garantia;
- não usar títulos;
- não mostrar raciocínio, metadata ou instruções internas;
- se o lead pedir uma pessoa, sinalize handoff humano;
- quando confirmar reunião, informe somente data, horário e link.
A resposta final deve ser apenas a mensagem que será enviada ao lead`

const PromptUserTemplate = `CONTEXTO FIXO:
Você está continuando uma conversa comercial de WhatsApp já iniciada.
Se houver histórico, não se apresente novamente.
Responda em no máximo 3 linhas e uma única resposta.
Não repita o histórico.
Avance para reunião quando houver intenção suficiente.
LEAD:
%s
ESTADO:
%s
RESUMO:
%s
HISTÓRICO RECENTE:
%s
MENSAGEM ATUAL:
%s`

// PromptContext contains only safe, provider-independent conversation context.
type PromptContext struct {
	Lead                map[string]any
	State               string
	ConversationSummary string
	RecentMessages      []domain.Message
	CurrentMessage      string
}

func BuildPromptUser(c PromptContext) string {
	return fmt.Sprintf(PromptUserTemplate,
		formatLead(c.Lead),
		valueOrEmpty(c.State),
		valueOrEmpty(c.ConversationSummary),
		formatMessages(c.RecentMessages),
		valueOrEmpty(c.CurrentMessage),
	)
}

func formatLead(lead map[string]any) string {
	if len(lead) == 0 {
		return "{}"
	}
	parts := make([]string, 0, len(lead))
	for k, v := range lead {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, "\n")
}

func formatMessages(messages []domain.Message) string {
	if len(messages) == 0 {
		return "(sem histórico)"
	}
	var b strings.Builder
	for _, m := range messages {
		role := "LEAD"
		if m.Direction == "outbound" {
			role = "AGENTE"
		}
		fmt.Fprintf(&b, "%s: %s\n", role, strings.TrimSpace(m.Content))
	}
	return strings.TrimSpace(b.String())
}

func valueOrEmpty(v string) string {
	if strings.TrimSpace(v) == "" {
		return "(vazio)"
	}
	return strings.TrimSpace(v)
}
