# DECISÕES CONFIRMADAS

Preencher durante o discovery.

## Infra
- Diretório local: `/root/prospect`
- API hosting: mesma infraestrutura Swarm
- Frontend hosting:
- Domain:
- PostgreSQL/Supabase: projeto `fsdszcfkjeavuoinyjas` via conexão Composio `supabase_stonen-lapper`
- PostgreSQL local da aplicação: database `prospect` no serviço `postgres_postgres`
- Role da aplicação: `prospect_app`
- Docker secret: `prospect_database_url` — valor não registrar
- Migration aplicada: `chat-backend/migrations/001_initial.sql`
- NATS:
- Temporal:
- Timezone: America/Campo_Grande
- Location: Campo Grande, MS
- UTC offset de referência: -04:00
- Reference time: 2026-09-06T10:00:00-04:00

## Cron
- Scheduler escolhido: Supabase pg_cron via Composio — validação operacional confirmada pelo recebimento da mensagem no número do Joel
- Horário de início: 08:00
- Janela comercial: 08:00–19:00, horário `America/Campo_Grande`
- Dias: segunda a sexta-feira
- Feriados: 2026-09-07 — fila deslocada para 2026-09-08
- Leads/dia:
- Intervalo:
- Tabela: prospect_leads_google
- Filtro: `contact_ready = false` (campo oficial; sem criar `initial_contact_status`) + `do_not_contact = false` + `converted = false` + `processing_status = 'pending'`
- Regra de envio: serializado, um lead por vez (`concurrency = 1`); reler estado e reservar antes de cada envio
- Confirmação do teste: Joel confirmou recebimento de 2 mensagens; teste considerado enviado/recebido
- Modo autorizado: produção controlada, com envio externo habilitado após validação final do `.env` e da fila
- Próxima execução planejada: 2026-09-08 às 08:00
- Envio externo: autorizado pelo Joel; ainda depende da ativação operacional no runtime

## Zernio
- Base URL: descoberto pela CLI oficial
- API Key configured: sim — não registrar chave
- Profile: `Default` (`6a5986cb39199d895a31af96`)
- WhatsApp Account: ativa, plataforma WhatsApp, perfil `Default`
- Template: `oi` — `APPROVED`
- Template ID: `1444129941108858`
- Language: `pt_BR`
- Params: nenhum — body: `Oi tudo bem?`
- Teste executado: sim, somente lead `wdtq75p8up`
- Teste messageId: `wamid.HBgMNTU2NzkyNDY2MzI5FQIAERgUQ0U3N0E4MzIwRDE2MTBDMDUzMDYA`
- Teste conversationId: `6a88bad95e0247d4b83faa2c`
- Teste delivery inicial: aceito pela API; conversa relida com `lastMessage: Oi tudo bem?`

## Trello
- Connection: `trello_lisa-levite` — ACTIVE
- Board: Menu Criativo-CRM
- Board URL: https://trello.com/b/DPwl3zEU/menu-criativo-crm
- List: Cadastro
- Insert position: imediatamente abaixo de `Novo Cliente`
- Custom Fields: PENDING — descobrir via Composio
- Labels: PENDING — descobrir via Composio
- Description: nome; WhatsApp; e-mail/site quando disponível; cidade/endereço; origem básica; `lead_status`; `pipeline_stage`
- Excluded description fields: `processing_status`; qualificações internas; flags técnicas; contadores; follow-up detalhado; IDs de mensagens; erros; timestamps internos; credenciais
- Required comment: `Enviar Tamplete`
- SLA: `created_at + 7 dias úteis`
- Test card: https://trello.com/c/Lxg3MZdI/55-joel-quintana-ia-infinito

## Ryze
- Base URL:
- API Key configured:
- Instance:
- Webhook:
- WebSocket:
- History:

## OmniRoute
- Base URL:
- API Key configured:
- Model:
- Timeout:
- Streaming:

## Composio / Cal.com
- Connection:
- Event Type:
- Duration:
- Timezone:
- Buffers:

## Telegram
- Bot configured:
- Chat ID:
- Critical alerts:

## Follow-up
- Stages:
- FU1:
- FU2:
- FU3:
- FU4:
- FU5:
- FU6:
- FU7:

## Frontend
- Framework:
- Domain:
- Auth:
