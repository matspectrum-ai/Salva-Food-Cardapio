.PHONY: test ci api-test api-run sqlc-generate db-validate db-test

test: api-test

ci: test

api-test:
	cd apps/api && go test ./...

api-run:
	cd apps/api && go run ./cmd/api

sqlc-generate:
	cd apps/api && podman run --rm --network=none -v "$$PWD:/src:Z" -w /src docker.io/sqlc/sqlc:1.30.0 generate

db-validate:
	./scripts/validate-db.sh

db-test:
	./scripts/test-postgres.sh
