# Go gRPC Microservice Boilerplate

A minimal, production-ready boilerplate for building gRPC microservices in Go.

This repository provides a clean foundation with:

- Config and validation
- Database setup with connection pooling
- gRPC server with an example RPC to demonstrate service structure
- Multi-stage Docker builds for both development and production
- Makefile with useful developer commands
- Unit and integration tests

Use this as a starting point for your own services — just replace the example RPC with your own application logic.

## Features

- Clean and extensible project structure
- Config package with validation (env-driven)
- Database package with pooling & safe close
- gRPC server with a working example RPC
- Graceful shutdown (cleanly stops gRPC server and background routines on interrupt)
- Multi-stage Dockerfile
  - Development: hot reload with Air
  - Production: optimized binary build
- Makefile for development, testing, and Docker tasks
- Unit & integration tests (using Testify + Testcontainers)
- Migrations & seeders with Makefile commands

## Getting Started

1. Clone the repository

```bash
git clone https://github.com/SagarMaheshwary/go-microservice-boilerplate.git
cd go-microservice-boilerplate
```

2. Setup environment variables (The application falls back to system environment variables if a `.env` file is not found—useful in Kubernetes where variables are mounted via ConfigMaps/Secrets.)

Copy the example environment file and adjust values as needed:

```bash
cp .env.example .env
```

## Requirements

You can run the service either locally or using Docker.

Local requirements

- [Go 1.22+](https://go.dev/dl/)
- [Make](https://www.gnu.org/software/make/)
- (Optional) [Air](https://github.com/air-verse/air?tab=readme-ov-file#via-go-install-recommended) for hot reload in development

Docker requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/)

If you don't have **make** installed on your system, you can install it using:

- **Ubuntu/Debian:** `sudo apt install make`
- **MacOS (Homebrew):** `brew install make`
- **Windows (via Chocolatey):** `choco install make`

## Running the Service

Run locally

```bash
make run     # Production mode, build and run binary
make run-dev # Development mode, reloads application on file change
```

Run inside Docker

```bash
make docker-run     # Production mode
make docker-run-dev # Development mode, reloads application on file change
```

## Testing

- Unit tests with mocks (using Testify).
- Integration tests with [Testcontainers](https://github.com/testcontainers/testcontainers-go):
  - Runs a real Postgres container for end-to-end database testing.
  - Tests ensure migrations + seeders work correctly.

```bash
make test             # all tests
make test-unit        # unit tests only
make test-integration # integration tests only
```

## Database & Migrations

- This boilerplate supports Postgres, SQLite, MySQL (via GORM).
- Uses [golang-migrate](https://github.com/golang-migrate/migrate) CLI for schema migrations.
- Example migration included: `users` table.

#### Install the migrate CLI

By default, this boilerplate is set up for Postgres.
Install the CLI with the Postgres driver:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

If you want to use another database, you’ll need to build the CLI with the corresponding driver tag:

- MySQL:

```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

- SQLite:

```bash
go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Makefile commands:

```bash
make migrate-up dsn="postgres://username:password@localhost:5432/dbname?sslmode=disable"    # Apply migrations
make migrate-down dsn="postgres://username:password@localhost:5432/dbname?sslmode=disable"  # Rollback migrations
make migrate-new name=create_users_table                                                    # Create a new migration file
```

## Seeders

- Seeders are stored in `internal/database/seeders/`.
- Example `users` seeder included.
- Seeders are run via a dedicated CLI:

```bash
make seed
```

- Each seeder logs progress, so you can see which one is running and where it fails.
- CLI entrypoint is under `cmd/cli/` — extensible if you want to add more developer commands later.

## Project Structure

```bash
.
├── proto/          # Protobuf definitions and generated code
├── cmd/            # Service entrypoint (main.go)
├── internal/       # Core application code
│   ├── config/         # Load and manage environment configurations
│   ├── logger/         # Zerolog-based structured logging
│   ├── service/        # Services for application business logic
│   └── database/       # Database initialization and connection handling
│       ├── migrations/     # Database migrations
│       ├── seeder/         # Seeders for generating fake data for dev/test
│       ├── model/          # GORM models
│   └── transports/     # Different communication protocols (e.g grpc, http, websocket). Each protocol can include both server/ and client/ implementations to keep responsibilities organized.
│       ├── grpc/           # gRPC transport
│       │   ├── server/         # gRPC server setup and service registration
│       │   │   ├── handler/         # RPC handlers
│       │   │   ├── interceptor/     # gRPC interceptors
│       │   └── client/         # (Optional) Place for gRPC clients (e.g., microservice-to-microservice communication)
│   └── tests/          # integration tests
│       ├── testutils/      # test helpers
├── Dockerfile      # Multi-stage build for dev/prod
├── Makefile        # Workflow automation (build, run, test, docker)
├── .env.example    # Example environment variables
└── readme.md       # Project documentation
```

## Test gRPC

After running the app (e.g via make docker-run-dev), you can test the example RPC using [grpcurl](https://github.com/fullstorydev/grpcurl)

```bash
grpcurl -d '{"user_id": 1}' -proto ./proto/hello_world/hello_world.proto -plaintext localhost:5000 hello_world.Greeter/SayHello
```

Expected response:

```json
{
  "message": "Hello, World!",
  "user": {
    "id": "1",
    "name": "Alice",
    "email": "alice@example.com"
  }
}
```
