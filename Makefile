GO_PATH := $(shell go env GOPATH)

dep:
	$(GO_PATH) mod tidy
	$(GO_PATH) mod download

build:
	docker compose up --build

lint:
	$(GO_PATH)/bin/golangci-lint run --timeout=5m -c .golangci.yml
