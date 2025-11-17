.PHONY: build test run clean docker-build docker-run migrate-up migrate-down migrate-create help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Test the project
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run unit tests only
test-unit:
	$(GOTEST) -v ./internal/...

# Run integration tests
test-integration:
	$(GOTEST) -v ./tests/integration/

# Run E2E tests
test-e2e:
	$(GOTEST) -v ./tests/e2e/

# Run the application
run:
	$(GOCMD) run ./cmd/server

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Download dependencies
deps:
	$(GOMOD) download

# Docker commands
docker-build:
	docker build -t pr-reviewer-assignment-service .

docker-run:
	docker run -p 8080:8080 --env-file .env pr-reviewer-assignment-service

# Docker Compose commands
docker-compose-up:
	docker-compose up --build

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f app

# Migration commands (requires golang-migrate installed locally)
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/reviewer_assigner?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/reviewer_assigner?sslmode=disable" down

migrate-create:
	@echo "Usage: make migrate-create name=<migration_name>"
	@if [ -z "$(name)" ]; then echo "Error: migration name is required. Use: make migrate-create name=your_migration_name"; exit 1; fi
	migrate create -ext sql -dir migrations -seq $(name)

# Development setup
dev-setup: deps
	@echo "Development environment setup complete"

# CI/CD commands
ci-test: test-coverage
ci-build: build-linux

# Show help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application binary"
	@echo "  build-linux     - Build for Linux platform"
	@echo "  test            - Run all tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-integration- Run integration tests"
	@echo "  test-e2e        - Run E2E tests"
	@echo "  run             - Run the application"
	@echo "  clean           - Clean build files"
	@echo "  tidy            - Tidy Go modules"
	@echo "  deps            - Download dependencies"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  docker-compose-up    - Start services with docker-compose"
	@echo "  docker-compose-down  - Stop services with docker-compose"
	@echo "  docker-compose-logs  - Show docker-compose logs"
	@echo "  migrate-up      - Run database migrations up"
	@echo "  migrate-down    - Run database migrations down"
	@echo "  migrate-create  - Create new migration file"
	@echo "  dev-setup       - Setup development environment"
	@echo "  ci-test         - Run tests for CI"
	@echo "  ci-build        - Build for CI"
	@echo "  help            - Show this help message"
