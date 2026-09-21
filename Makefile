APP=tonalmaster
GO=go

.PHONY: run test build fmt tidy up down logs

run:
	$(GO) run ./cmd/server

test:
	$(GO) test ./...

build:
	$(GO) build ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f api
