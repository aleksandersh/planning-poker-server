#!/usr/bin/env bash

set -Eeuo pipefail

ROUTE="$POKER_HOST/v1/rooms/$POKER_ROOM"
echo "GET $ROUTE Authorization: $POKER_SESSION"

curl -X GET "$ROUTE" -H "Content-Type: application/json" -H "Authorization: $POKER_SESSION"
echo ""
