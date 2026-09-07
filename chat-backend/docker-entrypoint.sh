#!/bin/sh
set -eu
if [ -r /run/secrets/prospect_omniroute_api_key ]; then
  export OMNIROUTE_API_KEY="$(cat /run/secrets/prospect_omniroute_api_key)"
fi
if [ -r /run/secrets/prospect_supabase_service_role_key ]; then
  export SUPABASE_SERVICE_ROLE_KEY="$(cat /run/secrets/prospect_supabase_service_role_key)"
fi
if [ -r /run/secrets/prospect_buffer_api_key ]; then
  export BUFFER_API_KEY="$(cat /run/secrets/prospect_buffer_api_key)"
fi
exec /app/chat-api
