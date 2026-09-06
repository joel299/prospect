# CLAUDE.md — Instruções do Projeto Hermes

## Missão

Construir uma plataforma nova de prospecção e atendimento comercial com:

- PostgreSQL/Supabase;
- cron diário;
- Zernio para primeiro template;
- Trello;
- Ryze para conversa;
- Hermes como orquestrador;
- OmniRoute como provider LLM;
- Composio + Cal.com;
- API em Go;
- WebSocket;
- frontend realtime;
- Error Engine;
- alertas Telegram.

---

# REGRA 1 — NÃO COMEÇAR CODIFICANDO

Antes de implementar qualquer etapa:

1. leia todos os arquivos `.md` deste pacote;
2. identifique a próxima fase;
3. liste exatamente o que está faltando;
4. faça perguntas ao usuário;
5. aguarde respostas;
6. escreva `DECISÕES CONFIRMADAS`;
7. somente depois desenvolva.

Não invente valores.

---

# REGRA 2 — PEDIR CREDENCIAIS NO MOMENTO CERTO

## PostgreSQL/Supabase

Perguntar:

```text
Database URL ou Supabase Project URL
método de autenticação
schema/tabela
se a tabela já existe
```

Nunca imprimir password/token depois de salvo.

## Zernio

Quando chegar ao fluxo de primeiro template:

```text
Base URL
API Key
```

Depois descobrir Profile e WhatsApp Account pela API.

Template padrão:

```text
oi
```

Confirmar idioma e parâmetros antes de enviar.

## Ryze

Quando iniciar o fluxo conversacional:

```text
Base URL
API Key/token
Instance
Account, se aplicável
Webhook
WebSocket
```

## OmniRoute

Perguntar:

```text
Base URL
API Key
Model/Route
Auth format
Timeout
Streaming
```

## Trello

Perguntar:

```text
API Key
Token
Board
```

Depois descobrir/confirmar:

```text
Lists
Labels
Custom Fields
Members
```

## Composio / Cal.com

Perguntar:

```text
Composio API Key / conexão
Connection ID
Cal.com Event Type
Timezone
Duração
Buffers
Antecedência mínima
```

## Telegram

Perguntar:

```text
Bot Token
Chat ID / Channel ID
```

---

# REGRA 3 — COMPOSIO

O Hermes deve usar Composio sempre que a especificação indicar uma ferramenta externa suportada.

Casos obrigatórios definidos:

```text
Cal.com
Trello
cron/scheduler, se a conexão/ferramenta disponível suportar
```

Para o cron:

1. primeiro inspecionar as tools disponíveis no Composio;
2. se houver scheduler/cron adequado, usar;
3. se não houver uma ferramenta compatível, NÃO substituir silenciosamente;
4. perguntar ao usuário se deve usar Temporal/cron interno da aplicação.

---

# REGRA 4 — NUNCA EXPOR SECRETS AO OMNIROUTE

OmniRoute recebe somente contexto conversacional.

Nunca enviar:

```text
API Keys
tokens
authorization headers
database credentials
```

---

# REGRA 5 — PRIMEIRO CONTATO

Fluxo obrigatório:

```text
Postgres/Supabase
↓
cron diário
↓
lead elegível
↓
Zernio
↓
template "oi"
↓
Trello
```

Um lead por vez.

Respeitar:

- rate limit documentado;
- resposta 429;
- Retry-After quando houver;
- políticas do WhatsApp;
- nenhuma tentativa de evasão de limites;
- janela comercial definida;
- idempotência.

---

# REGRA 6 — CONVERSAÇÃO

Somente depois de inbound real:

```text
Ryze
→ Go API
→ Hermes
→ OmniRoute
→ validator
→ Ryze
```

---

# REGRA 7 — AGENDA

Quando houver intenção de agenda:

```text
Hermes
→ Composio
→ Cal.com
```

Availability sempre live ou cache curtíssimo.

Booking sempre live.

Nunca dizer que marcou sem confirmação.

Após booking:

```text
agent_enabled=false
followup_enabled=false
```

---

# REGRA 8 — ERROS

Todo erro deve:

```text
ser persistido
ter trace_id
ser classificado
ser analisado
ser enviado ao Telegram quando notificável
```

Timeout com side effect potencial:

```text
não retry cego
→ reconciliação
```

---

# REGRA 9 — FRONTEND

Seguir `DESIGN_SYSTEM.md`.

Componentes de chat devem manter o padrão visual das referências fornecidas pelo usuário.

Não adicionar:

```text
ligação
videochamada
upload de arquivo
câmera
áudio
```

---

# REGRA 10 — TESTE

Antes de produção:

```text
mock
sandbox
allowlist
E2E
```

Nunca usar cliente real em testes iniciais.

---

# REGRA 11 — DOCUMENTAÇÃO

Toda mudança arquitetural precisa atualizar:

```text
README.md
arquivo de fluxo correspondente
DECISÕES CONFIRMADAS
```

---

# REGRA 12 — CRITÉRIO DE ENTREGA

Não declarar “pronto” sem:

```text
build
tests
healthcheck
provider test
database validation
WebSocket validation
E2E
```
