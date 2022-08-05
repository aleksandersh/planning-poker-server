#!/usr/bin/env bash

set -Eeuo pipefail

ROUTE="$POKER_HOST/v1/rooms/$POKER_ROOM/players"
echo "POST $ROUTE"

data="{\"name\":\"$1\"}"
curl -X POST "$ROUTE" -H "Content-Type: application/json" -d "$data"
echo ""
