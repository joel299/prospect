# OPERAÇÕES, ERROS, TESTES E PRODUÇÃO

## 1. Error Engine

Todo componente produz erro normalizado.

Campos:

```text
trace_id
service
operation
provider
lead_id
conversation_id
severity
category
error_code
http_status
message
retryable
attempt
occurred_at
```

---

# 2. Categorias

```text
AUTH
VALIDATION
DATABASE
PROVIDER_4XX
PROVIDER_5XX
PROVIDER_TIMEOUT
RATE_LIMIT
NETWORK
DUPLICATE
STATE_CONFLICT
LOCK_TIMEOUT
INVALID_AI_OUTPUT
TOOL_ERROR
DELIVERY_UNKNOWN
CONFIGURATION
INTERNAL
```

---

# 3. Retry

Retry somente quando seguro.

Backoff inicial sugerido:

```text
5s
30s
2m
10m
DLQ
```

Confirmar valores finais.

Adicionar jitter.

---

# 4. Reconciliation

Obrigatório para:

```text
Zernio timeout
Ryze outbound timeout
Cal.com booking timeout
Trello create uncertain
```

Resultado:

```text
confirmed_success
confirmed_failure
needs_review
```

---

# 5. Dead Letter

Depois de tentativas esgotadas ou estado não resolvido.

Painel permite:
- reprocess;
- reconcile;
- resolve;
- discard com justificativa.

---

# 6. Telegram

Fluxo:

```text
error
↓
Analista de Erros
↓
validator
↓
Telegram
```

Anti-spam:
- fingerprint;
- agrupamento;
- contador.

---

# 7. Lock Watchdog

Jobs monitorados:

```text
initial_contact_processing
followup_dispatching
booking_in_progress
reconciling
trello_sync
```

Não liberar lock cegamente.

---

# 8. Data Hygiene

Diariamente:
- duplicados;
- estados impossíveis;
- Card órfão;
- follow-up órfão;
- leads enviados voltando para fila;
- booking com agent still enabled.

---

# 9. Delivery Status

Normalizar:

```text
queued
accepted
sent
delivered
read
failed
unknown
```

---

# 10. Webhook Security

- HTTPS;
- signatures/secrets quando suportados;
- body limit;
- dedupe;
- replay protection;
- persist before process;
- async processing.

---

# 11. Rate Limit

Por provider.

Nunca contornar.

429:
- Retry-After;
- pause;
- backoff.

---

# 12. Test Lab

## Mock

Nenhuma chamada externa.

## Sandbox

Somente allowlist.

## Production

Requer confirmação.

---

# 13. Testes obrigatórios

## Cron

- query;
- zero leads;
- vários leads;
- lock;
- duplicate;
- 429;
- Zernio timeout;
- Trello error.

## Zernio

- auth;
- profile;
- account;
- template oi;
- idioma;
- params;
- success;
- timeout;
- reconciliation.

## Ryze

- inbound;
- duplicate inbound;
- history;
- outbound;
- timeout.

## Agent

- sem histórico;
- com histórico;
- 3 linhas;
- 4 perguntas;
- invalid output;
- tool call.

## Agenda

- slots;
- slot expirado;
- booking;
- booking timeout;
- agent stops.

## Follow-up

- business window;
- holiday;
- SKIP;
- invalid output;
- meeting cancels.

## Errors

- fingerprint;
- Telegram;
- retry;
- DLQ;
- reconciliation.

## Realtime

- reconnect;
- multiple clients;
- event order;
- slow client.

## Auth

- login;
- expired session;
- unauthorized WebSocket.

---

# 14. Seed

Criar usuário dev via env:

```text
SEED_ENABLED
SEED_ADMIN_EMAIL
SEED_ADMIN_PASSWORD
```

Nunca senha fixa no código.

---

# 15. Feature Flags

```text
INITIAL_CONTACT_ENABLED
AGENT_ENABLED
FOLLOWUP_ENABLED
TRELLO_SYNC_ENABLED
ERROR_NOTIFICATIONS_ENABLED
```

---

# 16. Produção

Rollout:

```text
mock
↓
sandbox
↓
allowlist production
↓
limited batch
↓
full production
```

Antes de ativar:
- health;
- metrics;
- trace;
- Telegram;
- database backup;
- rollback plan.

---

# 17. E2E principal

```text
lead elegível
↓
cron
↓
Zernio "oi"
↓
Trello
↓
lead responde
↓
Ryze
↓
Hermes
↓
OmniRoute
↓
Ryze
↓
lead pede reunião
↓
Hermes
↓
Composio
↓
Cal.com
↓
booking
↓
confirmação
↓
agent OFF
↓
painel realtime
```

---

# 18. E2E erro

```text
provider falha
↓
error event
↓
retry/reconcile
↓
Analista
↓
Telegram
↓
frontend realtime
```
