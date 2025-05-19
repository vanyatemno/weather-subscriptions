GO_PATH := $(shell go env GOPATH)

dep:
	@go mod tidy
	@go mod download

build:
	docker compose up --build