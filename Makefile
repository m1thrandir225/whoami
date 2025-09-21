# Application configuration
SERVICE_NAME := whoami
GO := go
GO_VERSION := $(shell $(GO) version)
GOPATH := $(shell $(GO) env GOPATH)
GOBIN := $(GOPATH)/bin
LDFLAGS := "-w -s" # Strip debug information
DOCKER_COMPOSE := docker compose
COMPOSE_FILE := deployment/compose.yml
COMPOSE_ENV_FILE := deployment/.env
MIGRATE := $(GOBIN)/migrate
MIGRATION_PATH := internal/db/migrations
GOLANGCI_LINT := $(GOBIN)/golangci-lint

# Database configuration
DB_CONTAINER := whoami-db
DB_NAME := whoami_db
DB_TEST_NAME := whoami_test
DB_USER := whoami_user
DB_PASS := secret
DB_PORT := 5432
REDIS_PORT := 6379
HTTP_PORT := 8080

# Database URLs
DB_URL := pgx5://${DB_USER}:${DB_PASS}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable
TEST_DB_URL := pgx5://${DB_USER}:${DB_PASS}@localhost:${DB_PORT}/${DB_TEST_NAME}?sslmode=disable
DOCKER_DB_URL := pgx5://${DB_USER}:${DB_PASS}@${DB_CONTAINER}:${DB_PORT}/${DB_NAME}?sslmode=disable

# Generate comprehensive .env file
$(COMPOSE_ENV_FILE):
	@echo "Generating comprehensive .env file..."
	@echo "# ========================================" > $(COMPOSE_ENV_FILE)
	@echo "# Whoami Service Environment File" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# Global Environment Settings" >> $(COMPOSE_ENV_FILE)
	@echo "ENVIRONMENT=development" >> $(COMPOSE_ENV_FILE)
	@echo "LOG_LEVEL=info" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Database Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "DB_USER=${DB_USER}" >> $(COMPOSE_ENV_FILE)
	@echo "DB_PASSWORD=${DB_PASS}" >> $(COMPOSE_ENV_FILE)
	@echo "DB_NAME=${DB_NAME}" >> $(COMPOSE_ENV_FILE)
	@echo "DB_PORT=${DB_PORT}" >> $(COMPOSE_ENV_FILE)
	@echo "DB_SOURCE=postgres://${DB_USER}:${DB_PASS}@${DB_CONTAINER}:${DB_PORT}/${DB_NAME}?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Redis Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "REDIS_PORT=${REDIS_PORT}" >> $(COMPOSE_ENV_FILE)
	@echo "REDIS_URL=redis://whoami-redis:${REDIS_PORT}" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Server Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "HTTP_SERVER_ADDRESS=0.0.0.0:${HTTP_PORT}" >> $(COMPOSE_ENV_FILE)
	@echo "HTTP_PORT=${HTTP_PORT}" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Authentication Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "ACCESS_TOKEN_DURATION=15m" >> $(COMPOSE_ENV_FILE)
	@echo "REFRESH_TOKEN_DURATION=7d" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Email Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "SMTP_HOST=localhost" >> $(COMPOSE_ENV_FILE)
	@echo "SMTP_PORT=587" >> $(COMPOSE_ENV_FILE)
	@echo "SMTP_USERNAME=" >> $(COMPOSE_ENV_FILE)
	@echo "SMTP_PASSWORD=" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "Environment file generated successfully at $(COMPOSE_ENV_FILE)"

.PHONY: generate-keys
generate-keys: ## Generate secure token symmetric key
	@echo "🔑 Generating secure token symmetric key..."
	@echo "TOKEN_SYMMETRIC_KEY=$$(openssl rand -hex 32)" >> $(COMPOSE_ENV_FILE)
	@echo "✅ Token key generated and added to $(COMPOSE_ENV_FILE)"

.PHONY: setup
setup: $(COMPOSE_ENV_FILE) generate-keys ## Complete setup: generate env file and keys
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Review and update configuration in $(COMPOSE_ENV_FILE)"
	@echo "2. Run 'make docker-up' to start services"
	@echo "3. Run 'make migrate-up-docker' to apply migrations"

.PHONY: build
build: tidy ## Build the application binary
	@echo "Building Go binary..."
	@CGO_ENABLED=0 $(GO) build -ldflags=$(LDFLAGS) -o ./bin/$(SERVICE_NAME) ./cmd/main.go
	@echo "Build complete. Binary in ./bin/$(SERVICE_NAME)"

.PHONY: tidy
tidy: ## Tidy Go module files
	@echo "Running go mod tidy..."
	@$(GO) mod tidy

.PHONY: fmt
fmt: ## Format Go source code
	@echo "Formatting Go code..."
	@go fmt ./...

.PHONY: lint
lint: ## Lint Go source code using golangci-lint
	@if ! command -v $(GOLANGCI_LINT) &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run ./...

.PHONY: test
test: tidy ## Run Go tests with coverage
	@echo "Running Go tests..."
	@$(GO) test -v -race -cover ./...

.PHONY: test-short
test-short: tidy ## Run Go tests with coverage (short mode)
	@echo "Running Go tests (short mode)..."
	@$(GO) test -v -cover -short ./...

.PHONY: docker-build
docker-build: ## Build Docker image using Docker Compose
	@echo "Building Docker image defined in $(COMPOSE_FILE)..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) build

.PHONY: docker-up
docker-up: ## Start services using Docker Compose (detached mode)
	@echo "Starting services via Docker Compose..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) up -d --remove-orphans

.PHONY: docker-down
docker-down: ## Stop and remove Docker Compose services
	@echo "Stopping and removing services via Docker Compose..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) down $(options)

.PHONY: docker-stop
docker-stop: ## Stop Docker Compose services without removing them
	@echo "Stopping services via Docker Compose..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) stop

.PHONY: docker-logs
docker-logs: ## Follow logs from Docker Compose services
	@echo "Following logs (Ctrl+C to stop)..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) logs -f

.PHONY: docker-ps
docker-ps: ## List running Docker Compose services
	@echo "Listing running services..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) ps

# Legacy database commands (for local development without Docker)
.PHONY: network
network: ## Create Docker network for the application
	docker network create whoami_network

.PHONY: postgres
postgres: ## Start PostgreSQL container (legacy)
	docker run --name ${DB_CONTAINER} -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASS} -d postgres:17-alpine

.PHONY: createdb
createdb: ## Create main database (legacy)
	docker exec -it ${DB_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

.PHONY: createdb_test
createdb_test: ## Create test database (legacy)
	docker exec -it ${DB_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_TEST_NAME}

.PHONY: dropdb
dropdb: ## Drop main database (legacy)
	docker exec -it ${DB_CONTAINER} dropdb ${DB_NAME}

# Migration commands
.PHONY: migrate-create
migrate-create: ## Create new migration files (e.g., make migrate-create name=add_users_table)
	@$(if $(name),,$(error Please specify migration name, e.g., make migrate-create name=my_migration))
	@echo "Creating migration '$(name)' in $(MIGRATION_PATH)..."
	@mkdir -p $(MIGRATION_PATH)
	@$(MIGRATE) create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

.PHONY: migrate-up
migrate-up: ## Apply all pending UP migrations
	@echo "Applying UP migrations from $(MIGRATION_PATH)..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose up $(steps)

.PHONY: migrate-up-docker
migrate-up-docker: ## Apply migrations using Docker database
	@echo "Applying UP migrations using Docker database..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DOCKER_DB_URL)" -verbose up $(steps)

.PHONY: migrate-down
migrate-down: ## Roll back all migrations (e.g., make migrate-down steps=1 to roll back only 1)
	@$(eval STEPS_ARG := $(if $(steps),$(steps),))
	@echo "Rolling back $(if $(steps),$(steps),all) DOWN migration(s) from $(MIGRATION_PATH)..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose down $(STEPS_ARG)

.PHONY: migrate-down-docker
migrate-down-docker: ## Roll back all migrations using Docker database
	@$(eval STEPS_ARG := $(if $(steps),$(steps),))
	@echo "Rolling back $(if $(steps),$(steps),all) DOWN migration(s) using Docker database..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DOCKER_DB_URL)" -verbose down $(STEPS_ARG)

.PHONY: migrate-status
migrate-status: ## Check migration status
	@echo "Checking migration status..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" version

.PHONY: migrate-status-docker
migrate-status-docker: ## Check migration status using Docker database
	@echo "Checking migration status using Docker database..."
	@$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DOCKER_DB_URL)" version

.PHONY: sqlc
sqlc: ## Generate SQL code using sqlc
	@echo "Generating SQL code using sqlc..."
	@sqlc generate

.PHONY: server
server: ## Start the application server
	@echo "Starting the application server..."
	@$(GO) run ./cmd/main.go

.PHONY: clean
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin
	@$(GO) clean

.PHONY: help
help: ## Display this help screen
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help
