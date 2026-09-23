.PHONY: test ci api-test api-run web-lint web-build web-dev sqlc-generate db-validate db-test infra-up infra-down db-migrate redis-cli

test: api-test

ci: api-test web-lint web-build

api-test:
	cd apps/api && go test ./...

api-run:
	DATABASE_URL="$${DATABASE_URL:-postgres://salva_food:salva_food@127.0.0.1:55432/salva_food?sslmode=disable}" cd apps/api && go run ./cmd/api

web-lint:
	cd apps/web && pnpm lint

web-build:
	cd apps/web && pnpm build

web-dev:
	cd apps/web && pnpm dev

sqlc-generate:
	cd apps/api && podman run --rm --network=none -v "$$PWD:/src:Z" -w /src docker.io/sqlc/sqlc:1.30.0 generate

db-validate:
	./scripts/validate-db.sh

db-test:
	./scripts/test-postgres.sh

infra-up:
	./scripts/compose-local.sh up -d

infra-down:
	./scripts/compose-local.sh down

db-migrate:
	./scripts/migrate-local.sh

redis-cli:
	./scripts/compose-local.sh exec redis redis-cli -h 127.0.0.1 -p 56379
