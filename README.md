# HERMES — Plataforma de Prospecção, Conversação e Operação Comercial

[![Status](https://img.shields.io/badge/status-discovery-orange)](./PACKAGE_STATUS.md)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](./03_GO_API_AND_REALTIME.md)
[![React](https://img.shields.io/badge/React-TypeScript-61DAFB?logo=react&logoColor=black)](./05_FRONTEND_PRODUCT_SPEC.md)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL%2FSupabase-database-4169E1?logo=postgresql&logoColor=white)](./06_DATABASE_AND_STATE.md)
[![WebSocket](https://img.shields.io/badge/realtime-WebSocket-010101)](./03_GO_API_AND_REALTIME.md)
[![License](https://img.shields.io/badge/license-TBD-lightgrey)](./DECISOES_CONFIRMADAS.md)

> Plataforma de prospecção, conversação, agenda e observabilidade comercial, orquestrada pelo Hermes.

Status atual: discovery concluído parcialmente; implementação aguardando fechamento das credenciais e decisões de infraestrutura.

## Runtime infrastructure

- PostgreSQL database: `prospect` on the existing Swarm service `postgres_postgres`.
- Application role: `prospect_app`.
- Docker secret: `prospect_database_url` (value never committed or printed).
- PostgreSQL schema migration: `chat-backend/migrations/001_initial.sql`.
- Realtime API endpoint: `/ws/v1`.
- Trello/Cal.com operations use the installed Composio CLI through the server-side adapter.

## 1. Visão geral

Este é um projeto novo.

O Hermes não deve assumir conhecimento de sistemas anteriores. Toda a implementação deve ser construída a partir desta especificação e das respostas fornecidas pelo usuário durante o discovery.

O sistema terá quatro blocos principais:

```text
1. PROSPECÇÃO
2. CONVERSAÇÃO
3. AGENDA
4. ERROS / OBSERVABILIDADE
```

Todos os blocos serão integrados por uma API em Go, PostgreSQL/Supabase, WebSocket e ferramentas acessadas pelo Hermes por meio de Composio/MCP.

---

# 2. Arquitetura principal

```text
POSTGRESQL / SUPABASE
        ↓
CRON DIÁRIO DO HERMES
        ↓
LISTA DE LEADS ELEGÍVEIS
        ↓
1 LEAD POR VEZ
        ↓
ZERNIO
        ↓
TEMPLATE WHATSAPP "oi"
        ↓
TRELLO
        ↓
AGUARDA RESPOSTA
        ↓
RYZE WEBSOCKET / WEBHOOK
        ↓
GO CONVERSATION API
        ↓
HERMES AGENT RUNTIME
        ↓
OMNIROUTE
        ↓
VALIDAÇÃO
        ↓
RYZE SEND TEXT
        ↓
LEAD

Quando houver intenção de agenda:

HERMES
  ↓
COMPOSIO
  ↓
CAL.COM
  ↓
slots disponíveis
  ↓
OMNIROUTE
  ↓
RYZE
  ↓
LEAD

Quando booking for confirmado:

agent_enabled = false
followup_enabled = false
meeting_scheduled = true
```

---

# 3. Ordem obrigatória de desenvolvimento

O Hermes deve seguir esta ordem:

```text
FASE 0 — DISCOVERY
FASE 1 — FUNDAÇÃO DO PROJETO
FASE 2 — BANCO / ESTADOS
FASE 3 — CRON DE PROSPECÇÃO
FASE 4 — ZERNIO
FASE 5 — TRELLO
FASE 6 — API GO
FASE 7 — WEBSOCKET
FASE 8 — RYZE
FASE 9 — OMNIROUTE / SDR
FASE 10 — COMPOSIO / CAL.COM
FASE 11 — FOLLOW-UP
FASE 12 — ERROR ENGINE / TELEGRAM
FASE 13 — FRONTEND
FASE 14 — TEST LAB / SEED
FASE 15 — TESTES E HARDENING
FASE 16 — PRODUÇÃO
```

Nenhuma fase deve ser implementada sem que as dependências da fase anterior estejam confirmadas.

---

# 4. Regra de discovery

Antes de desenvolver, o Hermes deve perguntar e registrar:

```text
DECISÕES CONFIRMADAS
```

O Hermes deve perguntar somente o necessário para a próxima etapa.

Ele não deve inventar:

- Base URL;
- API Key;
- IDs;
- nomes de Boards;
- Lists;
- Profiles;
- Accounts;
- Event Types;
- endpoints privados;
- credenciais;
- nomes de instâncias;
- schema não confirmado.

---

# 5. Stack principal

## Backend

```text
Go
PostgreSQL / Supabase
NATS JetStream
Temporal
WebSocket
REST API
```

## Agente

```text
Hermes
OmniRoute
MCP
Composio
Cal.com
```

## Providers

```text
Zernio = primeiro template oficial
Ryze = conversa depois da resposta
Trello = CRM visual
Telegram = canal de erro
```

## Frontend

```text
React
TypeScript
componentes inspirados nas referências fornecidas
WebSocket para realtime
```

---

# 6. Regras centrais

## Primeiro contato

Sempre:

```text
Zernio
→ template aprovado "oi"
```

Nunca iniciar conversa livre pela Ryze.

## Conversação

Somente depois da resposta:

```text
Ryze inbound
→ Hermes
→ OmniRoute
→ Ryze outbound
```

## Reunião

Quando booking for confirmado:

```text
meeting_scheduled = true
agent_enabled = false
followup_enabled = false
next_followup_at = null
conversation_mode = meeting_hold
```

## CRM

PostgreSQL/Supabase é a fonte de verdade.

Trello apenas representa visualmente o processo.

## Erro

Erro em qualquer fluxo:

```text
persistir
→ analisar
→ notificar Telegram
→ retry/reconcile/DLQ conforme regra
```

## Realtime

Frontend recebe:

```text
mensagens
estado do agente
tool calls
reuniões
follow-ups
erros
```

via WebSocket da API Go.

---

# 7. Arquivos do projeto

```text
CLAUDE.md
README.md
DESIGN_SYSTEM.md
MASTER_PROMPT_HERMES.md
01_ARCHITECTURE_AND_DISCOVERY.md
02_AUTOMATION_FLOWS.md
03_GO_API_AND_REALTIME.md
04_AGENT_CONVERSATION_AND_PROMPTS.md
05_FRONTEND_PRODUCT_SPEC.md
06_DATABASE_AND_STATE.md
07_OPERATIONS_ERRORS_AND_TESTS.md
ENV_EXAMPLE.md
```

Esses arquivos são autoritativos e devem ser lidos antes da implementação.
