# Chat Backend

Go 1.22+ contract-first API foundation for the chat-prospect conversation vertical slice.

## Run

```bash
cp .env.example .env
go run ./cmd/api
curl http://localhost:8080/health
```

No provider credentials are bundled. Ryze, OmniRoute, Buffer, and Composio are dependency-injected adapters and remain explicitly unconfigured until environment values and provider endpoint mappings are supplied. The API never sends secrets to OmniRoute; only conversation context is passed.

## API

Contract: [`docs/openapi.yaml`](docs/openapi.yaml). Endpoints include health, Ryze inbound webhook, conversation messages, outbound send, and `/ws/v1` reserved for the WebSocket edge adapter. Inbound events are deduplicated by provider external message ID and trigger an asynchronous agent response when OmniRoute is configured. Agent output is validated for non-empty content and a maximum of three lines. Tool calls are executed by Go through the Composio interface in bounded rounds.

## Database

Apply `migrations/001_initial.sql` to PostgreSQL/Supabase. The migration includes conversations, messages, idempotency keys, transactional outbox, constraints, and indexes. The in-memory store is only the development/test vertical slice; production must wire a PostgreSQL repository.

## Scope gates

- No frontend, deployment, DNS, external issue, or Telegram action is included.
- Ryze/OmniRoute/Buffer/Composio clients fail closed when not configured.
- WebSocket handler intentionally returns `501` until a vetted WebSocket implementation is added; persistence and event publication do not depend on it.
- Auth, PostgreSQL repository, NATS/JetStream, provider endpoint discovery, and production E2E remain explicit phase gates.
