#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

DATABASE_URL="${DATABASE_URL:-postgres://salva_food:salva_food@127.0.0.1:55432/salva_food?sslmode=disable}"

go run github.com/pressly/goose/v3/cmd/goose@v3.28.0 \
  -dir "$ROOT_DIR/apps/api/db/migrations" \
  postgres "$DATABASE_URL" up

echo "local postgres migrations: OK"
