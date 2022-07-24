#!/usr/bin/env bash

set -Eeuo pipefail

ROUTE="$POKER_HOST/v1/rooms"
echo "POST $ROUTE"

curl -X POST "$ROUTE" -H "Content-Type: application/json"
echo ""
