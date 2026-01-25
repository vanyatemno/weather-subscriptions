GO_PATH := $(shell go env GOPATH)

dep:
	$(GO_PATH) mod tidy
	$(GO_PATH) mod download

build:
	docker compose up --build

lint:
	$(GO_PATH)/bin/golangci-lint run --timeout=5m -c .golangci.yml

lint-staged:
	$(GO_PATH)/bin/golangci-lint run --timeout=5m -c .golangci.yml $(FILES)

swagger:
	@swag init -g cmd/main.go

e2e:
	docker compose up --build --abort-on-container-exit --exit-code-from e2e-tests e2e-tests