#!/usr/bin/env bash

set -Eeuo pipefail

ROUTE="$POKER_HOST/v1/rooms/$POKER_ROOM/currentgame/cards"
echo "PUT $ROUTE Authorization: $POKER_SESSION"

data="{\"score\":\"$1\"}"
curl -X PATCH "$ROUTE" -H "Content-Type: application/json" -d "$data" -H "Authorization: $POKER_SESSION"
echo ""
