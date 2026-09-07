# E2E TEST REPORT — 2026-09-06

## Executed

- `go test ./...`
- `go vet ./...`
- `go build ./cmd/api`
- Ryze adapter contract tests
- conversation vertical-slice E2E: inbound → agent → outbound
- live Ryze outbound + history read-back
- WebSocket connect → publish → receive E2E
- API health smoke test
- Trello bridge through the local API → Composio CLI
- direct Trello CLI read-back

## Evidence

- Build/tests/vet: PASS
- Live Ryze send HTTP: 200
- Live Ryze history HTTP: 200
- Ryze read-back found the exact E2E message: true
- Local conversation E2E: PASS
- WebSocket E2E: PASS
- API health: HTTP 200
- Trello API bridge: HTTP 202
- Composio CLI account: `trello_lisa-levite`
- Existing card resolved: `6a9ce557959dbdf409743e93`
- Composio Trello comment operation: successful
- Trello read-back: successful
- Exact comment `Enviar Tamplete`: confirmed

## Code changes

- Implemented Ryze text/history adapters.
- Wired outbound service delivery to Ryze.
- Implemented WebSocket transport in the same Go API.
- Added `POST /api/v1/integrations/trello/sync`.
- Changed the Composio bridge to use the installed Composio CLI, not MCP.
- Added fixed operation allowlist for Trello and Cal.com.
- Added Trello board/list discovery, duplicate card resolution, comment write, and read-back.
- Added Cal.com CLI operation mappings for availability and booking.
- Added environment aliases for existing Ryze and OmniRoute names.
- Fixed OmniRoute OpenAI response and SSE parsing in the Go client.

## OmniRoute health validation — 2026-09-07

- Swarm task: running on the manager node; service converged.
- SQLite `integrity_check`: `ok`.
- SQLite `quick_check`: `ok`.
- Public liveness endpoint: `GET /api/health/ping` → HTTP 200.
- Public completion probe with explicit connected provider/model → HTTP 200.
- Probe response: `HEALTH_OK`.
- Provider/model used for health probe: `antigravity/gemini-3.6-flash-low`.
- Root cause of previous timeout: default `auto` route attempted unavailable/ambiguous candidates; logs showed empty pools and abandoned upstream selection. Explicit connected model responds normally.
- API key setup intentionally not changed or recorded.
- Configured/tested `antigravity/gemini-3.8-flash-low`: upstream returned HTTP 404 `model_not_found`; it is not currently available through the connected Antigravity account.
- Known working probe model remains `antigravity/gemini-3.6-flash-low`.
- OmniRoute skill routes exist, but are management/authenticated APIs: `/api/skills`, `/api/skills/install`, and `/api/skills/executions`. They are not exposed as unauthenticated `/v1` model tools.

## Remaining gates

- Real inbound webhook/WebSocket receipt after the recipient replies.
- Configure and validate the production client API key separately.
- PostgreSQL repository/read-back must be wired into the application runtime.
- Cal.com real availability/booking requires the confirmed event-type payload.

No secrets are recorded in this report.
