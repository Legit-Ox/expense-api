.PHONY: help build run test clean docker-build docker-run docker-stop deps

# Default target
help:
	@echo "Available commands:"
	@echo "  deps          - Download Go dependencies"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  test-cover    - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker Compose services"
	@echo "  dev           - Run in development mode with SQLite"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build: deps
	go build -o bin/expense-api main.go

# Run the application locally
run: deps
	go run main.go

# Run tests
test: deps
	go test ./...

# Run tests with coverage
test-cover: deps
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	docker build -t expense-api .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	docker-compose down

# Development mode with SQLite
dev: deps
	DB_URL="sqlite://./expense.db" PORT=8080 go run main.go

# Install testify for testing
install-testify:
	go get github.com/stretchr/testify/assert 