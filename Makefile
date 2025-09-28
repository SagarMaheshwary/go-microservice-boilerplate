# Variables
APP_NAME := go-microservice-boilerplate
MIGRATIONS_DIR ?= internal/database/migrations

.PHONY: help proto build run run-dev test \
        docker-build-dev docker-build-prod \
        docker-run-dev docker-run-prod clean \
				migrate-up migrate-down migrate-new \
				test-unit test-integration

help: ## Show this help
	@echo "Available make commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'

proto: ## Generate probuf code from proto files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/**/*.proto

build: ## Build the service binary
	go build -o bin/$(APP_NAME) ./cmd/server

run: build ## Run the service locally
	./bin/$(APP_NAME)

run-dev: ## Run the app locally with air
	air

test: ## Run tests locally
	go test ./internal/... -v

test-unit: ## Run only unit tests inside internal/*
	go test $$(go list ./internal/... | grep -v ./internal/tests/) -v

test-integration: ## Run only integration tests inside internal/tests
	go test ./internal/tests/... -v

# Docker targets
docker-build-dev: ## Build docker image for development (hot reload via air)
	docker build --target development -t $(APP_NAME):dev .

docker-build: ## Build docker image for production (binary)
	docker build --target production -t $(APP_NAME):latest .

docker-run-dev: docker-build-dev ## Run docker container in development mode
	docker run -it --rm -p 5000:5000 -v $$(pwd):/app $(APP_NAME):dev

docker-run: docker-build ## Run docker container in production mode
	docker run -it --rm -p 5000:5000 -v .env:/app/.env $(APP_NAME):latest

clean: ## remove generated docker images/binaries
	rm -rf bin
	docker image rm $(APP_NAME):latest || true
	docker image rm $(APP_NAME):dev || true

# Migration targets
migrate-up: ## Run migrations up.
	@if [ -z "$(dsn)" ]; then \
		echo "Usage: make migrate-up dsn=\"postgres://postgres:password@localhost:5432/boilerplate?sslmode=disable\""; \
		exit 1; \
	fi

	@echo "Running migrations up on DSN=$(dsn)"
	@migrate -path $(MIGRATIONS_DIR) -database $(dsn) up

migrate-down: ## Run migrations down (be careful!).
	@if [ -z "$(dsn)" ]; then \
		echo "Usage: make migrate-down dsn=\"postgres://postgres:password@localhost:5432/boilerplate?sslmode=disable\""; \
		exit 1; \
	fi

	@echo "Running migrations down on $(DATABASE_NAME) with DSN=$(DATABASE_DSN)"
	@migrate -path $(MIGRATIONS_DIR) -database $(dsn) down

migrate-new: ## Create a new migration file. 
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-new name=add_users_table"; \
		exit 1; \
	fi

	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
