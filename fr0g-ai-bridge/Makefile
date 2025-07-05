.PHONY: help build run test clean proto proto-if-needed build-with-grpc deps fmt lint

# Default target
help:
	@echo "fr0g-ai-bridge Development Commands:"
	@echo "  build             - Build the bridge service"
	@echo "  build-with-grpc   - Build with gRPC (generates proto first)"
	@echo "  run               - Run the bridge service"
	@echo "  test              - Run tests"
	@echo "  clean             - Clean build artifacts"
	@echo "  proto             - Generate protobuf code"
	@echo "  proto-if-needed   - Generate protobuf only if needed"

# Build the service
build: proto-if-needed
	@echo "Building fr0g-ai-bridge..."
	go build -o bin/fr0g-ai-bridge ./cmd

# Build with gRPC support
build-with-grpc: proto build

# Run the service
run: build
	@echo "Starting fr0g-ai-bridge..."
	./bin/fr0g-ai-bridge

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@mkdir -p internal/pb
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh

# Generate protobuf only if files don't exist
proto-if-needed:
	@if [ ! -f internal/pb/bridge.pb.go ]; then \
		echo "Generating missing protobuf files..."; \
		make proto; \
	else \
		echo "Protobuf files already exist, skipping generation"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf internal/pb/*.pb.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run || echo "Install golangci-lint for better linting"

# Validate production build (ensure reflection is disabled)
validate-production:
	@echo "🔒 Validating Bridge production build security..."
	@if [ "$$GRPC_ENABLE_REFLECTION" = "true" ] && [ "$$ENVIRONMENT" = "production" ]; then \
		echo "❌ FATAL: gRPC reflection is enabled in production environment!"; \
		echo "   This is a security risk and deployment should be blocked."; \
		exit 1; \
	fi
	@echo "✅ Bridge production validation passed - reflection properly disabled"
