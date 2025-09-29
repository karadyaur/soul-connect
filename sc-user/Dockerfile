# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /app

# Download module dependencies first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy only the service source needed for the build
COPY sc-user ./sc-user

# Build the service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/general ./sc-user/cmd/general

# Build the migrate CLI with postgres support
RUN GOBIN=/build go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy binaries and application files
COPY --from=builder /build/general /usr/local/bin/general
COPY --from=builder /build/migrate /usr/local/bin/migrate
COPY sc-user/internal/db/migration ./internal/db/migration
COPY sc-user/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
