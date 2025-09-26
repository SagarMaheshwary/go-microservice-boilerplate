# Variables
APP_NAME := go-microservice-boilerplate

.PHONY: help proto build run run-dev test \
        docker-build-dev docker-build-prod \
        docker-run-dev docker-run-prod docker-test

help: ## Show this help
	@echo "Available make commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'

proto: ## Generate probuf code from proto files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/**/*.proto

build: ## Build the service binary
	go build -o bin/$(APP_NAME) ./cmd/server

run: build ## Run the service locally
	./bin/$(APP_NAME)

run-dev: ## Run the app locally with air (requires air installed, https://github.com/air-verse/air?tab=readme-ov-file#via-go-install-recommended)
	air

test: ## Run tests locally
	go test ./internal/... -v

# Docker targets
docker-build-dev: ## Build docker image for development (hot reload via air)
	docker build --target development -t $(APP_NAME):dev .

docker-build-prod: ## Build docker image for production (binary)
	docker build --target production -t $(APP_NAME):latest .

docker-run-dev: docker-build-dev ## Run docker container in development mode
	docker run -it --rm -p 5000:5000 -v $$(pwd):/app $(APP_NAME):dev

docker-run-prod: docker-build-prod ## Run docker container in production mode
	docker run -it --rm -p 5000:5000 -v .env:/app/.env $(APP_NAME):latest

docker-test: docker-build-dev ## Run tests inside docker
	docker run --rm $(APP_NAME):dev go test ./internal/... -v

clean:
	rm -rf bin
	docker image rm $(APP_NAME):latest || true
	docker image rm $(APP_NAME):dev || true
