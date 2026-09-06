# DATABASE E MÁQUINA DE ESTADOS

## 1. Fonte de verdade

PostgreSQL/Supabase.

---

# 2. Tabela leads

Campos:

```text
id
place_name
whatsapp
source
source_category
source_city
source_state
google_maps_url

qualification_status
qualification_stage
lead_score
needs_enrichment
enrichment_complete
contact_ready

trello_card_id

lead_status
pipeline_stage
status

do_not_contact
converted
responded

initial_contact_status
zernio_profile_id
zernio_account_id
zernio_conversation_id
zernio_message_id

conversation_transport

agent_enabled
agent_pause_reason
conversation_mode

meeting_scheduled
meeting_id
meeting_start_at
meeting_join_url
meeting_provider

followup_enabled
followup_stage
followup_current
followup_count
max_followups
next_followup_at
last_followup_at
followup_completed
followup_completed_at

last_contact_at
last_response_at
last_message
last_message_direction
last_message_id
last_response

processing_status
processing_locked_at
processing_lock_owner
processing_attempts
processing_error

archived
archived_reason
archived_at

created_at
updated_at
```

---

# 3. Conversations

```text
id
lead_id
channel
provider
agent_enabled
summary
last_message_at
created_at
updated_at
```

---

# 4. Messages

```text
id
conversation_id
lead_id
provider
external_message_id
direction
content
content_hash
status
sent_at
delivered_at
read_at
failed_at
created_at
```

---

# 5. Follow-ups

```text
id
lead_id
status
current_followup
followup_1_at ... followup_7_at
followup_1_status ... followup_7_status
next_followup_at
last_followup_at
completed_at
paused_at
error_message
processing_locked_at
processing_attempts
created_at
updated_at
```

---

# 6. Meetings

```text
id
lead_id
provider
external_booking_id
event_type_id
start_at
end_at
timezone
join_url
status
created_at
updated_at
```

---

# 7. Errors

```text
error_events
dead_letters
audit_events
tool_runs
webhook_receipts
idempotency_keys
outbox_events
processed_events
```

---

# 8. Máquina de estados

```text
captured
↓
qualified
↓
ready_for_initial_contact
↓
initial_contact_processing
↓
awaiting_reply
↓
conversation_active
↓
negotiation
↓
meeting
↓
won / lost
```

Alternativos:

```text
paused
do_not_contact
error
needs_review
archived
```

---

# 9. Invariantes

## awaiting_reply

```text
agent_enabled=false
conversation_transport=pending
```

## conversation_active

```text
agent_enabled=true
conversation_transport=ryze
```

## meeting

```text
meeting_scheduled=true
agent_enabled=false
followup_enabled=false
```

## won

```text
converted=true
agent_enabled=false
followup_enabled=false
```

## do_not_contact

Bloqueia qualquer outbound.

---

# 10. Idempotency

Keys:

```text
initial-contact:<lead_id>
agent-reply:<inbound_id>
followup:<lead_id>:<stage>
trello-card:<lead_id>
booking:<lead_id>:<slot>
```

---

# 11. Locks

```text
processing_status
processing_locked_at
processing_lock_owner
```

TTL depende da operação.

Lock expirado exige reconciliação quando houver possível side effect.
