.PHONY: test ci api-test api-run web-lint web-build web-dev sqlc-generate db-validate db-test

test: api-test

ci: api-test web-lint web-build

api-test:
	cd apps/api && go test ./...

api-run:
	cd apps/api && go run ./cmd/api

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
