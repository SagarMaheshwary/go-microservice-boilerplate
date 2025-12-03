# Go gRPC Microservice Boilerplate

A minimal, production-ready boilerplate for building gRPC microservices in Go.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Running the Service](#running-the-service)
- [Database & Migrations](#database--migrations)
- [Seeders](#seeders)
- [Testing](#testing)
- [Observability](#observability)
- [Examples branch](#examples-branch)
- [Tutorial Series](#tutorial-series)

## Features

- **Modular, extensible architecture** — clean separation of concerns for scalability and maintainability
- **Config-driven setup** — flexible environment-based configuration management
- **Database integration** — built-in migrations and seeders for easy setup
- **Redis caching layer** — ready-to-use caching for performance optimization
- **Graceful shutdown** — proper cleanup of resources on service termination
- **gRPC server** — includes structured logging interceptor
- **Observability ready** — integrated logging, Prometheus metrics, and OpenTelemetry tracing
- **Health checks** — `/livez` and `/readyz` endpoints for service monitoring
- **Testing ready** — unit and integration testing with Testify and Testcontainers
- **Developer tooling** — Makefile commands, hot reload (Air), and multi-stage Docker builds

> **Looking for a complete working example?**
> Check out the [examples](https://github.com/SagarMaheshwary/go-microservice-boilerplate/tree/examples) branch — includes sample gRPC service, Redis cache usage, migrations, seeders, metrics, tracing, and a full Docker-based observability stack.

## Project Structure

```bash
.
├── proto/          # Protobuf definitions and generated code
├── cmd/            # Service entrypoint (main.go)
├── internal/       # Core application code
│   ├── config/       # Load and manage environment configurations
│   ├── logger/       # Zerolog-based structured logging
│   ├── service/      # Services for application business logic
│   ├── cache/        # Redis
│   └── database/     # Database initialization and connection handling
│       ├── migrations/   # Database migrations
│       ├── seeder/       # Seeders for generating fake data for dev/test
│       └── model/        # GORM models
│   └── transports/   # Different communication protocols (e.g grpc, http, websocket). Each protocol can include both server/ and client/ implementations to keep responsibilities organized.
│       └── grpc/         # gRPC transport
│           ├── server/         # gRPC server setup and service registration
│           │   ├── handler/         # RPC handlers
│           │   └── interceptor/     # gRPC interceptors
│           └── client/         # (Optional) Place for gRPC clients (e.g., microservice-to-microservice communication)
│       └── http/         # HTTP transport
│           └── server/         # HTTP server setup, api routes (healthchecks/metrics)
│               └── handler/         # Route handlers
│   └── observability/
│       ├── metrics/      # Prometheus metrics
│       └── tracing/      # OpenTelemetry tracing
│   └── tests/
│       ├── integration/  # Integration tests
│       ├── mock/         # Mocks for unit tests
│       └── testutils/    # Test helpers
├── Dockerfile         # Multi-stage build for dev/prod
├── Makefile           # Workflow automation (build, run, test, docker)
├── docker-compose.yml # Postgres/Redis
└── readme.md          # Project documentation
```

## Requirements

You can run the service either locally or using Docker.

#### Local requirements

- [Go 1.22+](https://go.dev/dl/)
- [Make](https://www.gnu.org/software/make/)
- (Optional) [Air](https://github.com/air-verse/air?tab=readme-ov-file#via-go-install-recommended) for hot reload in development

#### Docker requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/)

#### Makefile

This project comes with a Makefile to simplify common workflows (building, running, migrations, tests, etc.).
Run the following to see all available commands:

```bash
make help
```

> Tip: Not every command is listed in the README — use `make help` to see all available workflows (build, lint, test, migrate-up, etc.).

#### Installing Make

If you don't have **make** installed on your system, you can install it using:

- **Ubuntu/Debian:** `sudo apt install make`
- **MacOS (Homebrew):** `brew install make`
- **Windows (via Chocolatey):** `choco install make`

## Getting Started

Clone the repository

```bash
git clone https://github.com/SagarMaheshwary/go-microservice-boilerplate.git
cd go-microservice-boilerplate
```

Setup environment variables (The application falls back to system environment variables if a `.env` file is not found—useful in Kubernetes where variables are mounted via ConfigMaps/Secrets.)

Copy the example environment file and adjust values as needed:

```bash
cp .env.example .env
```

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

This boilerplate also includes a `docker-compose.yml` with Postgres and Redis services. You can bring it up with:

```bash
docker compose up
```

Or, if you already have your own Postgres/Redis instance running, update the configuration accordingly.

## Database & Migrations

This boilerplate supports Postgres, SQLite, MySQL (via GORM) and uses [golang-migrate](https://github.com/golang-migrate/migrate) CLI for schema migrations.

By default, this boilerplate is set up for Postgres.

#### Install the migrate CLI

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

#### Migration commands

```bash
# Apply migrations
make migrate-up dsn="postgres://postgres:password@localhost:5432/boilerplate?sslmode=disable"

# Rollback migrations
make migrate-down dsn="postgres://postgres:password@localhost:5432/boilerplate?sslmode=disable"

# Create a new migration file
make migrate-new name=create_users_table
```

## Seeders

Seeders are stored in `internal/database/seeder/` and run via a dedicated CLI `cmd/cli/` (extensible if you want to add more developer commands later):

```bash
make seed
```

Each seeder logs progress, so you can see which one is running and where it fails.

## Testing

- Unit tests with mocks (using Testify).
- Integration tests with [Testcontainers](https://github.com/testcontainers/testcontainers-go):
  - The boilerplate is set up to support integration testing with real services.
  - Checkout `HealthService` integration tests in `internal/tests/integrations/service/health.go` that spins up Postgres and Redis containers and test against them.

You can use below make commands to run tests:

```bash
make test             # all tests
make test-unit        # unit tests only
make test-integration # integration tests only
```

## Observability

This boilerplate includes built-in observability tools to help you monitor, debug, and understand your services in production.

### Logging (Structured JSON)

- Centralized, structured logging using zerolog.
- Outputs clean, machine-readable **JSON logs**.
- Easily integrate with log aggregators like **Loki**, **ELK**, or **CloudWatch**.

### Metrics (Prometheus)

- Integrated Prometheus metrics via a dedicated `MetricsService`.
- Exposes `/metrics` endpoint for HTTP scraping.
- `example` branch includes custom **gRPC metrics interceptor** which tracks request count and latency per RPC method.
- Easily extendable with custom business metrics.

### Tracing (OpenTelemetry + OTLP)

- OpenTelemetry-based tracing using the **OTLP HTTP exporter**
- Vendor-neutral — works with Jaeger, Grafana Tempo, Datadog, New Relic, etc.
- Automatic trace propagation between gRPC services via `grpc.StatsHandler`

Together, **metrics**, **logs**, and **tracing** provide full visibility into your system’s behavior — helping you detect latency issues, understand dependencies, and debug bottlenecks efficiently.

## Examples Branch

A complete working demo is available in the [examples](https://github.com/SagarMaheshwary/go-microservice-boilerplate/tree/examples) branch. It includes:

- Sample gRPC service & handler (`SayHello`)
- User service example with DB and Redis
- Migrations & seeders
- Prometheus metrics + OpenTelemetry tracing
- docker-compose setup for Grafana, Prometheus, Jaeger
- Example dashboard & trace (SayHello → UserService)

> See the `examples` branch README for setup and usage instructions.

## Tutorial Series

This boilerplate is built as part of the **Designing Microservices in Go Series**.

- [Part 1](https://dev.to/sagarmaheshwary/go-microservices-boilerplate-series-from-hello-world-to-production-part-1-46k5) – Project setup: config, logging, gRPC server, graceful shutdowns, Dockerfile, and Makefile.
- [Part 2](https://dev.to/sagarmaheshwary/go-microservices-boilerplate-series-from-hello-world-to-production-part-2-428b) – Database integration: Postgres with GORM, migrations, seeders, service layer, and integration tests.
- [Part 3](https://dev.to/sagarmaheshwary/go-microservices-boilerplate-series-part-3-redis-healthchecks-observability-prometheus-metrics-32jo) – Redis caching, observability (Prometheus metrics + OpenTelemetry tracing), and health checks.

---

## Support & Contributions

If you find this project useful, consider giving it a ⭐, it helps others discover it.

Contributions, feedback, and suggestions are always welcome.
Feel free to open an issue or submit a PR anytime.
