FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY cmd ./cmd
COPY docs ./docs
COPY internal ./internal
COPY go.mod go.sum *.go ./
COPY deploy ./deploy
COPY api ./api
COPY test ./test

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/appbin ./cmd/main.go

FROM builder AS tester
WORKDIR /app
COPY . .
ENV DSN=""

FROM alpine:latest
LABEL MAINTAINER = <vanya04400@gmail.com>

COPY --from=builder /app /home/appuser/app

WORKDIR /home/appuser/app

# Since running as a non-root user, port bindings < 1024 is not possible
# 8000 for HTTP
EXPOSE 8000

CMD ["./appbin"]