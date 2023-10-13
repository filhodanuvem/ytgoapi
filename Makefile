.PHONY: dev
dev:
	docker compose -f docker-compose.dev.yaml up -d

.PHONY: up
up:
	docker compose -f docker-compose.yaml up -d --build

.PHONY: down
down:
	docker compose down

.PHONY: e2e
e2e:
	go build -o ./tmp/e2e ./cmd/e2e/main.go
	./tmp/e2e

.PHONY: runapi
runapi:
	go run cmd/api/main.go

.PHONY: logdb
logdb:
	docker logs -f ytgoapi-db-1
