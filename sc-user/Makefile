build: ## Build the project
	@echo "Building the project..."
	go build -o main cmd/general/main.go

start: ## Run the project
	go run cmd/general/main.go

test: ## Run all tests
	go test ./...

install: ## Install dependencies
	go mod tidy

swag: ##  Generate swagger specification
	swag init -g cmd/general/main.go

db-schema: ## Generate db schema
	sqlc generate

proto: ## Generate gRPC stubs
	protoc --proto_path=../proto --go_out=./internal --go-grpc_out=./internal ../proto/user.proto

help: ## Display help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' Makefile
