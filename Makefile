.PHONY: build clean run test docker-up docker-down seed help

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the application binaries"
	@echo "  clean       - Clean build artifacts"
	@echo "  run         - Run the application locally"
	@echo "  test        - Run tests"
	@echo "  docker-up   - Start the application with Docker Compose"
	@echo "  docker-down - Stop the Docker Compose stack"
	@echo "  seed        - Create sample data (requires database)"
	@echo "  help        - Show this help message"

# Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	go build -o bin/api cmd/api/main.go
	go build -o bin/seed cmd/seed/main.go
	@echo "Build completed!"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	go clean
	@echo "Clean completed!"

# Run the application locally
run: build
	@echo "Starting the application..."
	./bin/api

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Start with Docker Compose
docker-up:
	@echo "Starting with Docker Compose..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services are ready!"

# Stop Docker Compose
docker-down:
	@echo "Stopping Docker Compose stack..."
	docker-compose down

# Create sample data
seed: build
	@echo "Creating sample data..."
	./bin/seed

# Full setup: build, start docker, wait, and seed
setup: build docker-up
	@sleep 15
	@echo "Creating sample data..."
	./bin/seed
	@echo "Setup completed! API is running on http://localhost:8080"
	@echo "Try: curl http://localhost:8080/health"

# Development workflow
dev: clean build run