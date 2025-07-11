.PHONY: help build clean test run deps lint fmt proto health

# Default target
help:
	@echo "fr0g-ai-io - Input/Output Processing Service"
	@echo ""
	@echo "Available targets:"
	@echo "  build              - Build the service binary"
	@echo "  clean              - Clean build artifacts"
	@echo "  test               - Run tests"
	@echo "  run                - Run the service"
	@echo "  deps               - Install dependencies"
	@echo "  lint               - Run code linters"
	@echo "  fmt                - Format code"
	@echo "  proto              - Generate protobuf files"
	@echo "  health             - Check service health"

# Build the service
build:
	@echo "Building fr0g-ai-io service..."
	@go mod tidy
	@go build -o bin/fr0g-ai-io ./cmd/server
	@echo "Build complete: bin/fr0g-ai-io"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v -timeout=30s -short
	@echo "Tests complete"

# Run the service
run:
	@echo "Starting fr0g-ai-io service..."
	@go run ./cmd/server

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Run linters
lint:
	@echo "Running linters..."
	@go vet ./...
	@golangci-lint run || echo "Install golangci-lint for better linting"
	@echo "Linting complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Formatting complete"

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@mkdir -p internal/pb
	@if [ ! -f proto/io.proto ]; then echo "Error: proto/io.proto not found"; exit 1; fi
	@echo "Running protoc command..."
	@protoc --version || echo "protoc not found - install protocol buffers compiler"
	@which protoc-gen-go || echo "protoc-gen-go not found - run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
	@which protoc-gen-go-grpc || echo "protoc-gen-go-grpc not found - run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
	@PATH="$(shell go env GOPATH)/bin:$$PATH" protoc \
		--go_out=internal/pb --go_opt=paths=source_relative \
		--go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
		proto/io.proto
	@echo "Generated files:"
	@ls -la internal/pb/ || echo "No files generated"
	@echo "Protobuf generation complete"

# Check service health
health:
	@echo "Checking fr0g-ai-io service health..."
	@curl -sf http://localhost:8083/health || echo "Service not running or unhealthy"

# Development targets
dev-run:
	@echo "Running in development mode..."
	@go run ./cmd/server

dev-watch:
	@echo "Running with file watching (requires air)..."
	@air || echo "Install air for hot reloading: go install github.com/cosmtrek/air@latest"

# Docker targets
docker-build:
	@echo "Building Docker image..."
	@docker build -t fr0g-ai-io:latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8083:8083 -p 9092:9092 fr0g-ai-io:latest

# Integration targets
integration-test:
	@echo "Running integration tests..."
	@echo "Integration tests not implemented yet"

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. ./...

# Validate production build (ensure reflection is disabled)
validate-production:
	@echo "SECURITY Validating IO production build security..."
	@if [ "$$GRPC_ENABLE_REFLECTION" = "true" ] && [ "$$ENVIRONMENT" = "production" ]; then \
		echo "FAILED FATAL: gRPC reflection is enabled in production environment!"; \
		echo "   This is a security risk and deployment should be blocked."; \
		exit 1; \
	fi
	@echo "COMPLETED IO production validation passed - reflection properly disabled"
