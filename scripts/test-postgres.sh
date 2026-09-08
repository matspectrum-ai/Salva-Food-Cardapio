#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NAME="salva-food-test-pg-$$"
PORT="${TEST_POSTGRES_PORT:-55432}"
IMAGE="docker.io/library/postgres:17-alpine"

cleanup() {
  podman stop "$NAME" >/dev/null 2>&1 || true
}
trap cleanup EXIT

podman run -d --rm --name "$NAME" --network=host \
  -e POSTGRES_PASSWORD=postgres "$IMAGE" -p "$PORT" >/dev/null

for _ in $(seq 1 30); do
  if podman exec "$NAME" pg_isready -U postgres -p "$PORT" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

for migration in "$ROOT_DIR"/apps/api/db/migrations/*.sql; do
  sed '/-- +goose Down/,$d' "$migration" \
    | podman exec -i "$NAME" psql -U postgres -p "$PORT" -v ON_ERROR_STOP=1 >/dev/null
done

(
  cd "$ROOT_DIR/apps/api"
  TEST_DATABASE_URL="postgres://postgres:postgres@127.0.0.1:${PORT}/postgres?sslmode=disable" \
    go test -count=1 -v ./internal/platform/postgres -run 'Test(Store|IdentityStore|AuthSession|OnboardingStore)Integration'
)
