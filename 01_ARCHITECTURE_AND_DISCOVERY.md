# ARQUITETURA, DISCOVERY E ORDEM DE CONSTRUÇÃO

## 1. Princípio

Este sistema deve ser construído de forma orientada a eventos e estados persistidos.

Não criar uma automação monolítica.

Domínios:

```text
Lead Intake
Initial Contact
CRM
Conversation
Agent
Calendar
Follow-up
Errors
Realtime
Frontend
```

---

# 2. Fluxo principal

```text
PostgreSQL/Supabase
       ↓
Daily Prospecting Cron
       ↓
Lead Reservation
       ↓
Zernio Template "oi"
       ↓
Trello Sync
       ↓
Awaiting Reply
       ↓
Ryze Inbound
       ↓
Conversation API
       ↓
Hermes
       ↓
OmniRoute
       ↓
Ryze Outbound
```

---

# 3. Agenda

```text
Conversation
↓
scheduling intent
↓
Hermes Tool Call
↓
Composio
↓
Cal.com
↓
availability
↓
Hermes
↓
OmniRoute
↓
Ryze
```

Booking:

```text
selected slot
↓
Hermes
↓
Composio
↓
Cal.com
↓
booking confirmed
↓
final confirmation
↓
agent OFF
follow-up OFF
```

---

# 4. Erro

```text
ANY COMPONENT
↓
Error Event
↓
Postgres
↓
Classification
↓
Retry / Reconciliation / DLQ
↓
Error Analyst
↓
Telegram
↓
WebSocket
↓
Frontend
```

---

# 5. Ordem de implementação detalhada

## Fase 0 — Discovery

Criar:

```text
DECISÕES CONFIRMADAS
```

Perguntar infraestrutura e ambiente.

## Fase 1 — Repository / Foundation

Estrutura recomendada:

```text
/apps/frontend
/apps/api
/apps/worker
/internal
/docs
/deploy
```

## Fase 2 — Banco

Criar schema e migrations.

## Fase 3 — Cron de prospecção

Construir sem ainda enviar mensagem real.

Primeiro rodar em `dry_run`.

## Fase 4 — Zernio

Discovery:
- Base URL;
- API Key;
- Profile;
- Account;
- template `oi`.

Depois sandbox.

## Fase 5 — Trello

Usar Composio.

Configurar Board/Lists/Fields.

## Fase 6 — Go API

REST + auth + services + provider adapters.

## Fase 7 — WebSocket

Eventos realtime.

## Fase 8 — Ryze

Inbound/history/outbound.

## Fase 9 — Agent Runtime

Hermes + OmniRoute + prompts.

## Fase 10 — Agenda

MCP/Composio/Cal.com.

## Fase 11 — Follow-up

Timers e regras comerciais.

## Fase 12 — Erros

Error Engine + Telegram.

## Fase 13 — Frontend

Implementar telas.

## Fase 14 — Seed / Test Lab

Mock e sandbox.

## Fase 15 — Hardening

Idempotência, rate limit, observabilidade, retries.

## Fase 16 — Produção

Feature flags e rollout controlado.

---

# 6. Perguntas de discovery

## Infra

- Onde o Go API será hospedado?
- Onde o frontend será hospedado?
- PostgreSQL/Supabase já existe?
- Qual domínio será usado?
- TLS/reverse proxy?
- Docker?
- Temporal será self-hosted ou cloud?
- NATS será self-hosted?

## Prospecção

- Nome da tabela?
- Campos atuais?
- Qual coluna indica elegibilidade?
- Quantos leads/dia?
- Horário do cron?
- Timezone?
- Dias permitidos?
- Feriados?
- Qual intervalo entre envios?

## Zernio

- Base URL?
- API Key?
- qual Profile?
- qual account?
- `oi` está aprovado?
- idioma?
- parâmetros?

## Trello

- credenciais?
- Board?
- Lists?
- Fields?
- Labels?

## Ryze

- Base URL?
- key?
- instance?
- webhook?
- WebSocket?
- formato dos eventos?

## OmniRoute

- URL?
- key?
- model?
- timeout?
- streaming?

## Agenda

- Composio connection?
- Cal.com Event Type?
- duração?
- timezone?
- buffers?

## Telegram

- Bot token?
- chat/channel ID?
- tipos de alerta?

---

# 7. Regras de performance

- operações de banco com índices;
- HTTP keep-alive;
- connection pools;
- payloads mínimos;
- histórico resumido;
- NATS para eventos;
- Temporal para waits;
- WebSocket separado de persistência;
- evitar consultas repetidas;
- tool calls somente quando necessárias.

---

# 8. Regra de source of truth

```text
PostgreSQL/Supabase = verdade
Trello = representação
WebSocket = atualização UI
NATS = transporte interno
```
