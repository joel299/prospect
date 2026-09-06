# API GO, EVENT BUS E REALTIME

## 1. Responsabilidade

A API Go é o núcleo de acesso da aplicação.

Ela não substitui o Hermes.

Ela oferece:
- REST;
- WebSocket;
- auth;
- provider adapters;
- persistence;
- state transitions;
- event publishing;
- dashboard APIs.

---

# 2. Estrutura sugerida

```text
/apps/api/cmd/api
/apps/worker/cmd/worker

/internal/auth
/internal/leads
/internal/prospecting
/internal/conversations
/internal/agents
/internal/followup
/internal/calendar
/internal/errors
/internal/notifications
/internal/realtime
/internal/persistence

/internal/providers/zernio
/internal/providers/ryze
/internal/providers/omniroute
/internal/providers/trello
/internal/providers/composio
```

---

# 3. REST API

## Auth

```text
POST /api/v1/auth/login
POST /api/v1/auth/logout
GET  /api/v1/auth/me
```

## Leads

```text
GET /api/v1/leads
GET /api/v1/leads/:id
PATCH /api/v1/leads/:id
```

## Conversations

```text
GET  /api/v1/conversations
GET  /api/v1/conversations/:id
GET  /api/v1/conversations/:id/messages
POST /api/v1/conversations/:id/send
POST /api/v1/conversations/:id/agent/pause
POST /api/v1/conversations/:id/agent/resume
```

## Agent

```text
GET /api/v1/agent/config
PUT /api/v1/agent/config
GET /api/v1/agent/prompts
POST /api/v1/agent/prompts
POST /api/v1/agent/test
```

## Integrations

```text
GET /api/v1/integrations
POST /api/v1/integrations/:provider/test
PUT /api/v1/integrations/:provider
```

## Follow-up

```text
GET /api/v1/followups
POST /api/v1/followups/:id/pause
POST /api/v1/followups/:id/resume
```

## Errors

```text
GET /api/v1/errors
GET /api/v1/errors/:id
POST /api/v1/errors/:id/retry
POST /api/v1/errors/:id/reconcile
POST /api/v1/errors/:id/resolve
```

## Test Lab

```text
POST /api/v1/test/agent
POST /api/v1/test/zernio
POST /api/v1/test/ryze
POST /api/v1/test/calendar
```

---

# 4. WebSocket

Endpoint:

```text
GET /ws/v1
```

Auth:
- token da aplicação;
- nunca token de provider.

Eventos:

```text
conversation.message.inbound
conversation.message.outbound
conversation.updated
agent.thinking
agent.tool.started
agent.tool.completed
agent.paused
meeting.scheduled
followup.updated
error.created
error.updated
error.resolved
integration.health
```

---

# 5. WebSocket Gateway

Fluxo:

```text
provider event
↓
persist
↓
NATS
↓
WebSocket Hub
↓
authenticated clients
```

Nunca depender do WebSocket para persistir evento.

---

# 6. NATS

Subjects sugeridos:

```text
lead.*
conversation.*
agent.*
meeting.*
followup.*
error.*
integration.*
```

JetStream para eventos que não podem se perder.

---

# 7. Temporal

Usar para:

```text
Daily Prospecting
Follow-up timers
Retry schedules
Reconciliation
Maintenance
```

Se o cron principal for construído via Composio scheduler, Temporal ainda pode ser usado internamente para processamento da fila.

---

# 8. Idempotência

Chaves:

```text
initial-contact:<lead_id>
agent-reply:<inbound_message_id>
followup:<lead_id>:<stage>
trello-card:<lead_id>
booking:<lead_id>:<slot>
```

---

# 9. Transactional Outbox

Toda mudança crítica:

```text
transaction
├── state update
└── outbox event
```

Worker publica no NATS.

---

# 10. Circuit breakers

Separados:

```text
Zernio
Ryze
OmniRoute
Trello
Composio
```

---

# 11. Observabilidade

Por request:

```text
trace_id
request_id
lead_id
conversation_id
provider
duration_ms
status
```

Métricas:
- API p50/p95/p99;
- provider p95;
- WebSocket clients;
- queue depth;
- error rate;
- 429;
- retries;
- DLQ.
