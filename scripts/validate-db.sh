#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
API_DIR="$ROOT_DIR/apps/api"
IMAGE="docker.io/library/postgres:17-alpine"

podman run --rm --network=none --user postgres \
  -v "$API_DIR:/src:Z" "$IMAGE" sh -ceu '
    initdb -D /tmp/pgdata >/dev/null
    pg_ctl -D /tmp/pgdata -o "-k /tmp -c listen_addresses=" -w start >/dev/null
    trap "pg_ctl -D /tmp/pgdata -m fast -w stop >/dev/null 2>&1 || true" EXIT

    migration=/src/db/migrations/000001_core.sql
    sed "/-- +goose Down/,\$d" "$migration" \
      | psql -h /tmp -U postgres -d postgres -v ON_ERROR_STOP=1 >/dev/null

    count=$(psql -h /tmp -U postgres -d postgres -Atc \
      "select count(*) from information_schema.tables where table_schema='"'"'public'"'"'")
    test "$count" = "8"

    sed -n "/-- +goose Down/,\$p" "$migration" \
      | psql -h /tmp -U postgres -d postgres -v ON_ERROR_STOP=1 >/dev/null

    remaining=$(psql -h /tmp -U postgres -d postgres -Atc \
      "select count(*) from information_schema.tables where table_schema='"'"'public'"'"'")
    test "$remaining" = "0"
    echo "postgres migration up/down: OK"
  '
