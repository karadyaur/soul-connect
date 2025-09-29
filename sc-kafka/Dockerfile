# syntax=docker/dockerfile:1

FROM golang:1.23 as build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/sc-kafka ./sc-kafka/cmd/main

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /bin/sc-kafka /usr/bin/sc-kafka

USER nonroot

ENTRYPOINT ["/usr/bin/sc-kafka"]
