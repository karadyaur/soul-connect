# 🫂 Soul Connect Social Media (SCSM)

Soul Connect is a peer-support social network where people can safely share personal stories, ask for advice, and receive empathetic feedback from the community. The platform combines a React web client with several Go microservices and supporting infrastructure that securely stores data, processes events, and delivers notifications.

## Key features

- 📝 Share stories and updates tied to personal experiences and emotional well-being.
- 💬 Receive comments, reactions, and suggestions from other community members.
- 🔐 Register, authenticate, and manage your personal profile.
- 📣 Subscribe to updates and get real-time notifications about community activity.
- 🔄 Process events asynchronously through Kafka to keep the platform responsive and scalable.

## Architecture and services

| Component | Description |
| --- | --- |
| `sc-webapp` | React client that interacts with end users. |
| `sc-api-getaway` | API gateway that aggregates requests to microservices and exposes REST/Swagger endpoints. |
| `sc-auth` | gRPC/REST service for registration, login, and authentication. |
| `sc-user`, `sc-post`, `sc-notification` | Domain services that manage profiles, posts, and notifications. |
| `sc-kafka` | Kafka orchestration service that provisions topics and consumers. |
| `postgres` | Persistent storage for users and content. |
| `zookeeper`, `kafka-broker` | Kafka infrastructure required for event streaming. |

## Repository structure

```
├── docker-compose.yml          # Production-ready docker-compose for the entire stack
├── local.docker-compose.yml    # Simplified compose file for local development
├── postgres/                   # Database configuration and Dockerfile
├── proto/                      # gRPC contract definitions
├── sc-api-getaway/             # API gateway service written in Go
├── sc-auth/                    # Authentication microservice
├── sc-post/, sc-user/, ...     # Additional domain services
├── sc-kafka/                   # Kafka integration utilities
└── sc-webapp/                  # React client
```

## Prerequisites

Make sure the following tools are installed:

- [Git](https://git-scm.com/) for cloning the repository.
- [Docker](https://www.docker.com/products/docker-desktop) and Docker Compose v2.20 or newer.
- [Go 1.23+](https://go.dev/dl/) if you plan to run Go services outside Docker.
- [Node.js 18+](https://nodejs.org/en/download) and a package manager (npm, pnpm, or yarn) for running the web client locally.

## Quick start with Docker Compose

1. **Clone the repository:**
   ```bash
   git clone https://github.com/<your-account>/soul-connect.git
   cd soul-connect
   ```
2. **(Optional) Configure environment variables:**
   - Copy `example.env` from each service into `.env` and adjust values if necessary.
   - For the web client: `cp sc-webapp/example.env sc-webapp/.env`.
3. **Build and start the backend stack:**
   ```bash
   docker compose up --build
   ```
   This will start Postgres, Kafka, and all core microservices.
4. **Verify everything is running:**
   - API Gateway: http://localhost:8000
   - Swagger (if enabled in the API Gateway): http://localhost:8000/swagger/index.html
   - Launch the web client separately (see “Local setup without Docker”).
5. **Shut down the services:**
   ```bash
   docker compose down
   ```

### Local docker-compose for development

The `local.docker-compose.yml` file spins up a minimal stack (database, API gateway, Kafka). Start it the same way:

```bash
docker compose -f local.docker-compose.yml up --build
```

Attach additional services and the web client manually when needed.

## Local setup without Docker

1. **Start PostgreSQL:**
   ```bash
   docker compose up postgres -d
   ```
   Create the `sc_db` database if it is not provisioned automatically.
2. **Prepare environment variables:**
   ```bash
   cp sc-auth/example.env sc-auth/.env
   cp sc-api-getaway/example.env sc-api-getaway/.env
   # repeat for other services when necessary
   ```
   Update connection parameters (`DB_SOURCE`, `GRPC_AUTH_PORT`, `WEBAPP_BASE_URL`, etc.) to match local ports.
3. **Install dependencies and run services:**
   ```bash
   cd sc-auth && go mod tidy && go run cmd/general/main.go
   cd sc-api-getaway && go mod tidy && go run cmd/general/main.go
   # do the same for sc-user, sc-post, sc-notification
   ```
   Use `make proto-generate` or Makefile commands when regenerating gRPC code from the `proto` directory.
4. **Start Kafka (if required by the services):**
   ```bash
   docker compose up zookeeper kafka-broker sc-kafka -d
   ```
5. **Launch the web client:**
   ```bash
   cd sc-webapp
   npm install    # or pnpm install / yarn
   npm start
   ```
   The application runs at http://localhost:3000 and calls the API gateway at http://localhost:8000.

## Testing

- **Go services:**
  ```bash
  cd sc-auth && go test ./...
  cd sc-api-getaway && go test ./...
  ```
- **Web client:**
  ```bash
  cd sc-webapp
  npm test
  ```

## Helpful Makefile commands

Each service exposes utility commands. Example for `sc-auth`:

```bash
make install        # Install dependencies (go mod tidy)
make start          # Run the service (go run ...)
make test           # Execute tests (go test ./...)
make proto-generate # Generate gRPC code from proto/auth.proto
```

## Additional notes

- All proto files live in the `proto/` directory. Regenerate gRPC clients and servers after modifying contracts.
- The default Kafka setup creates `post.created`, `subscription.created`, and `notification.created` topics. Override them via `sc-kafka` environment variables if required.
- For easier debugging, use `docker compose logs -f <service>` and `docker compose exec <service> sh` to inspect logs and access running containers.

Happy building and sharing on Soul Connect! 🫶
