FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY go.mod go.sum *.go ./
COPY deploy ./deploy
COPY api ./api
RUN ls

# Toggle CGO based on your app requirement. CGO_ENABLED=1 for enabling CGO
RUN #CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /cmd/main.go
# Use below if using vendor
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/appbin ./cmd/main.go

FROM alpine:latest
LABEL MAINTAINER = <vanya04400@gmail.com>

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
# RUN apk --update add ca-certificates && \
#    rm -rf /var/cache/apk/*

## Add new user 'appuser'
#RUN adduser -D appuser
#USER appuser

COPY --from=builder /app /home/appuser/app

WORKDIR /home/appuser/app

# Since running as a non-root user, port bindings < 1024 is not possible
# 8000 for HTTP
EXPOSE 8000

CMD ["./appbin"]