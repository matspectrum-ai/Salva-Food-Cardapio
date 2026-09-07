.PHONY: test ci api-test api-run

test: api-test

ci: test

api-test:
	cd apps/api && go test ./...

api-run:
	cd apps/api && go run ./cmd/api
