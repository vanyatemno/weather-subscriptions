GO_PATH := $(shell go env GOPATH)

lint: check-lint dep
	golangci-lint run --timeout=5m -c .golangci.yml

dep:
	@go mod tidy
	@go mod download

build:
	docker compose up --build