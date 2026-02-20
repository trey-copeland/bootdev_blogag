#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${GATOR_DB_URL:-}" && -f .env ]]; then
  set -a
  source ./.env
  set +a
fi

: "${GATOR_DB_URL:?GATOR_DB_URL is not set}"

action="${1:-}"
case "$action" in
  down)
    goose -dir sql/schema postgres "$GATOR_DB_URL" down
    ;;
  up)
    goose -dir sql/schema postgres "$GATOR_DB_URL" up
    ;;
  status)
    goose -dir sql/schema postgres "$GATOR_DB_URL" status
    ;;
  login)
    psql "$GATOR_DB_URL"
    ;;
  *)
    echo "Usage: $0 {down|up|status|login}" >&2
    exit 1
    ;;
esac
