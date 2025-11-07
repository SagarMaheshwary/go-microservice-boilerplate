# Go gRPC Microservice Boilerplate — Examples Branch

A production-ready boilerplate for building scalable gRPC microservices in Go.

This **examples branch** includes:

- Example gRPC service (SayHello)
- Users table migration and seeder
- Observability setup (Prometheus, Grafana, Jaeger)
- Example metrics and traces

It’s meant to **demonstrate how the boilerplate works in action** —
the master branch, on the other hand, is a **clean version** ready for your own service code and configurations.

## Table of Contents

- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Test gRPC](#test-grpc)
- [Test APIs](#test-apis)
- [Observability](#observability)

## Project Structure

```bash
.
├── proto/          # Protobuf definitions and generated code
├── cmd/            # Service entrypoint (main.go)
├── internal/       # Core application code
│   ├── config/       # Load and manage environment configurations
│   ├── logger/       # Zerolog-based structured logging
│   ├── service/      # Services for application business logic
│   ├── cache/        # Redis caching
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
├── docker-compose.observability.yml
└── readme.md          # Project documentation
```

## Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/)

#### Makefile

This project comes with a Makefile to simplify common workflows (building, running, migrations, tests, etc.).
Run the following to see all available commands:

```bash
make help
```

> Not every command is documented in the README, so `make help` is the best way to explore what’s available.

#### Installing Make

If you don't have **make** installed on your system, you can install it using:

- **Ubuntu/Debian:** `sudo apt install make`
- **MacOS (Homebrew):** `brew install make`
- **Windows (via Chocolatey):** `choco install make`

## Getting Started

Clone the repository and switch to the **examples** branch:

```bash
git clone https://github.com/SagarMaheshwary/go-microservice-boilerplate.git
git checkout examples
cd go-microservice-boilerplate
```

Copy the example environment file:

```bash
cp .env.example .env
```

Start the core stack (Application, Postgres, Redis):

```bash
docker compose up
```

In another terminal, start the observability stack (Grafana, Prometheus, Jaeger):

```bash
docker compose -f docker-compose.observability.yml up
```

Install `migrate` CLI and run database migrations:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

make migrate-up dsn="postgres://postgres:password@localhost:5432/boilerplate?sslmode=disable"
```

Seed test data:

```bash
make seed dsn="postgres://postgres:password@localhost:5432/boilerplate?sslmode=disable"
```

> See `internal/database/seeder/users.go` to check what’s being seeded.

## Test gRPC

You can test the example **SayHello** RPC using [grpcurl](https://github.com/fullstorydev/grpcurl):

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

## Test APIs

Check service health:

**Livez:**

```bash
curl localhost:4000/livez
```

Response:

```json
{ "status": "ok" }
```

**Readyz:**

```bash
curl localhost:4000/readyz
```

Response:

```json
{ "status": "ready", "details": { "cache": "ok", "database": "ok" } }
```

**Metrics:**

```bash
curl localhost:4000/metrics
```

Example output:

```bash
# HELP grpc_requests_total Total number of gRPC requests received, labeled by method and status.
# TYPE grpc_requests_total counter
grpc_requests_total{method="/hello_world.Greeter/SayHello",status="OK"} 1
```

> Default Go runtime metrics are disabled — set `METRICS_ENABLE_DEFAULT_METRICS=true` to enable them.

## Observability

### Prometheus + Grafana

When you start `docker-compose.observability.yml`, Grafana automatically mounts:

- a Prometheus data source, and
- a sample gRPC dashboard visualizing request counts and latencies.

Generate some gRPC traffic and open Grafana to view the metrics:

- URL: [http://localhost:3000](http://localhost:3000)
- Username: **admin**
- Password: **admin**

![Grafana Sample Dashboard](./assets/grafana-sample-dashboard.png)

### OpenTelemetry Tracing

Jaeger automatically collects traces from the SayHello RPC example.

Visit [http://localhost:16686](http://localhost:16686), select the service `go-microservice-boilerplate`, and explore traces:

**Example — search results:**

![Jaeger Tracing Search Traces](./assets/jaeger-tracing-search-traces.png)

**Example — selected trace details:**

![Jaeger Tracing Search Traces](.//assets/jaeger-tracing-selected-trace.png)

---

You now have a fully functional example showing:

- gRPC service layer
- Redis caching
- Prometheus metrics
- Jaeger tracing
- Grafana visualization
