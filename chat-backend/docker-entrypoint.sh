#!/bin/sh
set -eu
if [ -r /run/secrets/prospect_omniroute_api_key ]; then
  export OMNIROUTE_API_KEY="$(cat /run/secrets/prospect_omniroute_api_key)"
fi
exec /app/chat-api
