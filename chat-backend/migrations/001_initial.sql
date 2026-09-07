CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS conversations (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), lead_id text NOT NULL, channel text NOT NULL DEFAULT 'whatsapp', provider text NOT NULL DEFAULT 'ryze', agent_enabled boolean NOT NULL DEFAULT false, meeting_scheduled boolean NOT NULL DEFAULT false, summary text NOT NULL DEFAULT '', last_message_at timestamptz, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS messages (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), conversation_id uuid NOT NULL REFERENCES conversations(id), lead_id text NOT NULL, provider text NOT NULL, external_message_id text, direction text NOT NULL CHECK (direction IN ('inbound','outbound')), content text NOT NULL, content_hash text NOT NULL, status text NOT NULL DEFAULT 'queued', created_at timestamptz NOT NULL DEFAULT now(), UNIQUE(provider, external_message_id)
);
CREATE TABLE IF NOT EXISTS idempotency_keys (key text PRIMARY KEY, response_hash text NOT NULL, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS outbox_events (id uuid PRIMARY KEY DEFAULT gen_random_uuid(), topic text NOT NULL, aggregate_id text NOT NULL, payload jsonb NOT NULL, published_at timestamptz, created_at timestamptz NOT NULL DEFAULT now());
CREATE INDEX IF NOT EXISTS messages_conversation_created_idx ON messages(conversation_id, created_at);
CREATE INDEX IF NOT EXISTS outbox_unpublished_idx ON outbox_events(created_at) WHERE published_at IS NULL;
CREATE TABLE IF NOT EXISTS agent_memory (
 lead_id text PRIMARY KEY,
 summary text NOT NULL DEFAULT '',
 facts jsonb NOT NULL DEFAULT '{}'::jsonb,
 updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS agent_memory_updated_idx ON agent_memory(updated_at);
