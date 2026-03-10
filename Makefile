.PHONY: all build run test clean docker-build docker-up docker-down migrate lint fmt help

# Application name
APP_NAME=pet-log-api
MAIN_PATH=./cmd/api

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(GOBIN)/$(APP_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)/main.go

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@echo "Running $(APP_NAME) with hot reload..."
	@air

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	@docker-compose logs -f

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/pet_log?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/pet_log?sslmode=disable" down

migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Generate API documentation (requires swag)
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs

# Install development tools
tools:
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# Show help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run with hot reload (requires air)"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter (requires golangci-lint)"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - Show Docker logs"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make migrate-create - Create new migration"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make tools          - Install development tools"
