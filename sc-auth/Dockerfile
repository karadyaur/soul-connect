# syntax=docker/dockerfile:1.4

FROM golang:1.22 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/auth-service ./sc-auth/cmd/general
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/migrate github.com/golang-migrate/migrate/v4/cmd/migrate

FROM busybox:1.36.1 AS busybox

FROM gcr.io/distroless/base-debian12:nonroot AS final
WORKDIR /app

COPY --from=builder /out/auth-service /app/auth-service
COPY --from=builder /out/migrate /bin/migrate
COPY --from=builder /src/sc-auth/internal/db/migration /app/internal/db/migration
COPY --from=busybox /bin/busybox /busybox
COPY --chmod=0755 sc-auth/entrypoint.sh /entrypoint.sh

ENV PATH="/bin:${PATH}"

ENTRYPOINT ["/entrypoint.sh"]
