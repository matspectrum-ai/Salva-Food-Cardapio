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

    for migration in /src/db/migrations/*.sql; do
      sed "/-- +goose Down/,\$d" "$migration" \
        | psql -h /tmp -U postgres -d postgres -v ON_ERROR_STOP=1 >/dev/null
    done

    for table in tenants catalog_items orders app_users tenant_memberships auth_sessions; do
      test "$(psql -h /tmp -U postgres -d postgres -Atc "select to_regclass('"'"'public.$table'"'"') is not null")" = "t"
    done

    for migration in $(ls -r /src/db/migrations/*.sql); do
      sed -n "/-- +goose Down/,\$p" "$migration" \
        | psql -h /tmp -U postgres -d postgres -v ON_ERROR_STOP=1 >/dev/null
    done

    remaining=$(psql -h /tmp -U postgres -d postgres -Atc \
      "select count(*) from information_schema.tables where table_schema='"'"'public'"'"'")
    test "$remaining" = "0"
    echo "postgres migrations up/down: OK"
  '
