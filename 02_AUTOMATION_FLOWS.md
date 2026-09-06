# FLUXOS DE AUTOMAÇÃO — ESPECIFICAÇÃO COMPLETA

# FLUXO A — CRON DIÁRIO DE PROSPECÇÃO

## Objetivo

Uma vez por dia montar a fila dos leads da aplicação e enviar o primeiro template.

## Scheduler

Preferência definida:

```text
Hermes → Composio → ferramenta de cron/scheduler disponível
```

Antes de implementar:
- descobrir se a conexão Composio possui uma ferramenta compatível;
- se existir, usar;
- se não existir, pedir autorização para Temporal Schedule ou cron interno.

## Passos

```text
CRON
↓
validar janela comercial
↓
consultar PostgreSQL/Supabase
↓
selecionar leads elegíveis
↓
ordenar fila
↓
pegar 1 lead
↓
reler estado
↓
lock/reserva
↓
Zernio template "oi"
↓
persistir messageId/conversationId
↓
Composio → Trello
↓
persistir Card ID
↓
aguardar intervalo
↓
próximo lead
```

## Seleção

Campos conceituais:

```text
contact_ready = false
do_not_contact = false
converted = false
processing_status = pending
```

## Reserva

Transação atômica:

```text
processing_status = initial_contact_processing
processing_locked_at = now
```

Se outro worker já reservou, ignorar.

## Zernio

Usar:

```text
templateName = oi
```

O Hermes deve solicitar:
- Base URL;
- API Key.

Depois descobrir Profile/Account.

## Rate limit

Não definir número arbitrário.

O sistema deve:
- respeitar limites documentados/configurados;
- serializar envios inicialmente;
- aplicar `concurrency=1` permanentemente neste fluxo inicial: um lead por vez; só iniciar próximo após confirmação/persistência do anterior;
- interpretar 429;
- honrar Retry-After;
- usar backoff;
- não tentar contornar limite.

## Regras de WhatsApp

A implementação deve operar apenas dentro das políticas aplicáveis ao template oficial e ao consentimento/base legal definido pelo usuário.

Não criar mecanismos para burlar políticas, qualidade, limites ou bloqueios do provider.

## Trello

Após envio confirmado:

```text
Hermes
↓
Composio Trello tool
↓
criar/atualizar Card
```

Campos recomendados:

```text
Lead ID
Nome
WhatsApp
Cidade
Origem
Lead Score
Initial Contact Status
```

Checklist:

```text
Prospecção
- Template enviado
```

Due date:

```text
created_at + 7 dias úteis
```

## Finalização

```text
lead_status = awaiting_reply
contact_ready = true
processing_status = pending
```

---

# FLUXO B — RECEBIMENTO DA RESPOSTA

```text
Ryze WebSocket/Webhook
↓
Go Provider Gateway
↓
dedupe
↓
persist inbound
↓
identificar lead
↓
conversation_transport = ryze
agent_enabled = true
↓
publicar NATS event
↓
Hermes
```

Antes da primeira resposta:
- consultar histórico Ryze;
- carregar prompt;
- prompt cache;
- dados do lead;
- estado.

---

# FLUXO C — RESPOSTA DO SDR

```text
Hermes
↓
load context
↓
OmniRoute
↓
validator
↓
Ryze send text
↓
persist outbound
↓
Trello update quando necessário
↓
WebSocket frontend
```

## Validator

Bloquear:
- vazio;
- analysis;
- reasoning;
- metadata;
- safety labels;
- JSON indevido;
- múltiplas respostas.

Regras:
- máximo 3 linhas;
- uma resposta;
- apresentação somente uma vez.

---

# FLUXO D — AGENDA

## Intenção

Exemplos:

```text
"podemos marcar"
"qual horário"
"amanhã funciona"
"terça à tarde"
```

Fluxo:

```text
Hermes
↓
calendar_get_availability
↓
Composio
↓
Cal.com
↓
slots
↓
Hermes
↓
OmniRoute
↓
Ryze
```

## Booking

```text
lead escolhe slot
↓
Hermes
↓
calendar_create_booking
↓
Composio
↓
Cal.com
↓
success
↓
persist
↓
OmniRoute gera confirmação final
↓
Ryze envia
↓
agent_enabled=false
followup_enabled=false
```

Nunca inventar horário.

Nunca afirmar booking antes de sucesso.

---

# FLUXO E — FOLLOW-UP

## Elegibilidade

```text
followup_enabled = true
meeting_scheduled = false
converted = false
do_not_contact = false
next_followup_at <= now
```

## Scheduler

Preferir Temporal para timers duráveis.

## Geração

```text
load history
↓
Hermes
↓
OmniRoute
↓
gate
↓
Ryze
```

## SKIP

OmniRoute pode retornar:

```text
SKIP_FOLLOWUP
```

Nesse caso não envia.

## Fim

Após último estágio:

```text
followup_completed = true
followup_enabled = false
```

---

# FLUXO F — ERROS

```text
erro
↓
error_events
↓
fingerprint
↓
classificação
↓
retry? reconcile? DLQ?
↓
Analista de Erros
↓
validator
↓
Telegram
↓
WebSocket frontend
```

## Telegram

O Hermes deve pedir:
- Bot Token;
- Chat ID/Channel ID.

Formato:

```text
Erro no workflow
Serviço: <service> | Link: <dashboard>
Sugestão: <ação>
```

## Timeout

Se uma chamada pode ter produzido side effect:

```text
delivery_unknown
↓
reconciliation
```

Não retry cego.

---

# FLUXO G — MAINTENANCE

Jobs:

```text
lock_watchdog
provider_reconciliation
data_hygiene
trello_reconciliation
meeting_reconciliation
followup_reconciliation
```

## Data hygiene

Detectar:
- duplicados;
- lead já enviado voltando para fila;
- Card órfão;
- follow-up ativo após reunião;
- lock antigo.

Preferir soft delete/archive.

---

# FLUXO H — REALTIME

Todos os fluxos publicam eventos.

```text
NATS
↓
Go WebSocket Hub
↓
frontend
```

Eventos:

```text
lead.updated
initial_contact.sent
message.inbound
message.outbound
agent.thinking
agent.tool.started
agent.tool.completed
meeting.scheduled
followup.updated
error.created
error.resolved
```
