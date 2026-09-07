.PHONY: api-test api-run

api-test:
	cd apps/api && go test ./...

api-run:
	cd apps/api && go run ./cmd/api
