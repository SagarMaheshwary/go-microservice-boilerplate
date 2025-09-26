# Go gRPC Microservice Boilerplate

A minimal, production-ready boilerplate for building gRPC microservices in Go.

This repository provides a clean foundation with:

- Config and validation
- Database setup with connection pooling
- gRPC server with an example RPC to demonstrate service structure
- Multi-stage Docker builds for both development and production
- Makefile with useful developer commands
- Unit tests for core components

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
- Unit tests

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

Option 1: Run locally with hot reload (Air required)

```bash
make run-dev
```

Option 2: Run inside Docker (development mode with hot reload)

```bash
make docker-run-dev
```

Option 3: Run inside Docker (production mode)

```bash
make docker-run-prod
```

## Running Tests

Run unit tests locally:

```bash
make test
```

Or run them inside Docker:

```bash
make docker-test
```

## Project Structure

```bash
.
├── proto/          # Protobuf definitions and generated code
├── cmd/            # Service entrypoint (main.go)
├── internal/           # Core application code
│   ├── config/         # Load and manage environment configurations
│   ├── database/       # Database initialization and connection handling
│   ├── logger/         # Zerolog-based structured logging
│   └── transports/         # Different communication protocols (e.g grpc, http, websocket). Each protocol can include both server/ and client/ implementations to keep responsibilities organized.
│       ├── grpc/           # gRPC transport
│       │   ├── server/     # gRPC server setup and service registration
│       │   │   ├── handler/         # RPC handlers
│       │   │   ├── interceptor/     # gRPC interceptors
│       │   └── client/     # (Optional) Place for gRPC clients (e.g., microservice-to-microservice communication)
├── Dockerfile      # Multi-stage build for dev/prod
├── Makefile        # Workflow automation (build, run, test, docker)
├── .env.example    # Example environment variables
└── README.md       # Project documentation
```

## Test gRPC

After running the app (e.g via make docker-run-dev), you can test the example RPC using [grpcurl](https://github.com/fullstorydev/grpcurl)

```bash
grpcurl -proto ./proto/hello_world/hello_world.proto -plaintext localhost:5000 hello_world.Greeter/SayHello
```

Expected response:

```json
{
  "message": "Hello World"
}
```

### Database

This boilerplate includes a database package (Postgres/MySQL/SQLite supported).

In Part One, database connection is disabled by default to keep the service runnable out-of-the-box.

If you’d like to enable database integration:

- Set up a Postgres instance (local or Docker).
- Update .env with valid DB credentials.
- Uncomment the database initialization code in cmd/server/main.go.
- Database-backed RPCs will be added in the next part of this series.

## Next Steps

This boilerplate is part of a series. In the next part, we’ll add:

- A real database-backed RPC
- Database migrations
- Logging improvements
- Middleware (interceptors)

Stay tuned! 🚀
