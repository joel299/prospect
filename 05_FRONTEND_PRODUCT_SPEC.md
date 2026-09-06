# FRONTEND PRODUCT SPEC

## 1. Objetivo

Criar um painel operacional de atendimento/prospecção em tempo real.

Referências visuais indicadas pelo usuário:

```text
21st.dev — WhatsApp mock chat
21st.dev — chat messages
```

Usar como inspiração de composição, não como motivo para adicionar funções fora do escopo.

---

# 2. Rotas

```text
/login
/app/inbox
/app/leads
/app/followups
/app/agent
/app/automations
/app/integrations
/app/test-lab
/app/observability
/app/settings
```

---

# 3. Login

Tela limpa:

```text
logo
email
senha
entrar
estado de loading
erro
```

Session segura.

---

# 4. Inbox

Composição:

```text
┌──────── Sidebar ───────┬──────── Conversations ────────┬──────── Chat ────────────────┬──── Inspector ────┐
│ Inbox                  │ Search                        │ Header                        │ Lead               │
│ Leads                  │ Avatar / name / preview       │ Messages                      │ Status             │
│ Follow-ups             │ Time / unread                 │ Composer                      │ Trello             │
│ Agent                  │                               │ Agent status                  │ Follow-up          │
│ Errors                 │                               │                               │ Meeting            │
└────────────────────────┴───────────────────────────────┴───────────────────────────────┴───────────────────┘
```

---

# 5. Chat

Manter:
- bolha inbound;
- bolha outbound;
- timestamp;
- sent/delivered/read;
- typing/processing;
- status do agente;
- reunião marcada;
- loading de tool call.

Não incluir:
- telefone/chamada;
- vídeo;
- anexos;
- arquivo;
- câmera;
- gravação de áudio.

---

# 6. Agent Page

Tabs:

```text
Visão Geral
Prompt
Prompt Cache
Ferramentas
Provider
Agenda
Regras
Teste
```

## Prompt

Editor versionado.

Ações:

```text
Salvar rascunho
Publicar
Restaurar versão
Comparar versões
```

## Ferramentas

Toggles:

```text
Consultar agenda
Marcar reunião
Consultar Trello
Mover Card
Criar comentário
Pausar follow-up
```

## Provider

```text
OmniRoute
Base URL configurada
API Key configurada
Model/Route
Timeout
Streaming
```

Nunca exibir API Key real.

---

# 7. Integrations

Cards:

```text
Postgres/Supabase
Zernio
Ryze
OmniRoute
Trello
Composio / Cal.com
Telegram
```

Mostrar:

```text
Configured
Healthy
Latency
Last check
```

Botão:

```text
Testar conexão
```

---

# 8. Follow-ups

Tabela:

```text
Lead
Stage
Next follow-up
Last contact
Status
Agent
```

Ações:
- pause;
- resume;
- inspect.

---

# 9. Errors

Tabs:

```text
Abertos
Críticos
Retries
Reconciliation
Dead Letter
Resolvidos
```

Detalhes:
- service;
- operation;
- provider;
- lead;
- trace;
- analysis;
- retries;
- resolution.

---

# 10. Test Lab

Modos:

```text
Mock
Sandbox
Production
```

Production exige confirmação explícita.

Permitir simular:

```text
Inbound do lead
Resposta do OmniRoute
Tool call
Agenda
Erro
WebSocket
```

---

# 11. Realtime

WebSocket deve atualizar sem refresh:

```text
mensagem
status
agent thinking
tool call
meeting
follow-up
error
integration status
```

---

# 12. Responsividade

Desktop é prioridade operacional.

Tablet suportado.

Mobile deve permitir:
- visualizar conversa;
- responder;
- pausar agente;
- ver lead;
- ver erro crítico.
