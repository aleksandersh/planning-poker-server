#!/usr/bin/env bash

set -Eeuo pipefail

ROUTE="$POKER_HOST/v1/rooms/$POKER_ROOM/games"
echo "POST $ROUTE Authorization: $POKER_SESSION"

data="{\"name\":\"$1\"}"
curl -X POST "$ROUTE" -H "Content-Type: application/json" -d "$data" -H "Authorization: $POKER_SESSION"
echo ""
