#!/usr/bin/env bash

set -Eeuo pipefail

ROUTE="$POKER_HOST/v1/rooms/$POKER_ROOM"
echo "DELETE $ROUTE"

curl -X DELETE "$ROUTE" -H "Content-Type: application/json" -H "Authorization: $POKER_SESSION"
echo ""
