# syntax=docker/dockerfile:1.4

FROM golang:1.22 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api-gateway ./sc-api-getaway/cmd/general

FROM gcr.io/distroless/base-debian12:nonroot AS final
WORKDIR /app

COPY --from=builder /out/api-gateway /app/api-gateway

EXPOSE 8000

ENTRYPOINT ["/app/api-gateway"]
