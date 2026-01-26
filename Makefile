GO_PATH := $(shell go env GOPATH)
include .env

dep:
	$(GO_PATH) mod tidy
	$(GO_PATH) mod download

build:
	docker compose up --build

lint:
	$(GO_PATH)/bin/golangci-lint run --timeout=5m -c .golangci.yml

lint-staged:
	$(GO_PATH)/bin/golangci-lint run --timeout=5m -c ./.golangci.yml $(FILES)

swagger:
	@swag init -g ./cmd/main.go

postgres-up:
	docker compose up --build -d postgres

e2e: postgres-up
	DNS=$(DNS) go test -v ./test/...