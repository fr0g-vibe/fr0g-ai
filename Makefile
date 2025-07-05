.PHONY: help setup build build-all clean clean-all test test-all run run-aip run-bridge run-mcp deps lint fmt health proto docker-build docker-up docker-down

# Default target
help:
	@echo "🐸 fr0g.ai - AI-Powered Security Intelligence Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  setup              - Initialize development environment"
	@echo "  build-all          - Build all fr0g.ai components"
	@echo "  clean-all          - Clean all build artifacts"
	@echo "  test-all           - Run tests for all components"
	@echo "  run-aip            - Run fr0g-ai-aip service"
	@echo "  run-bridge         - Run fr0g-ai-bridge service"
	@echo "  run-mcp            - Run fr0g-ai-master-control service"
	@echo "  run-io             - Run fr0g-ai-io service"
	@echo "  deps               - Install dependencies for all components"
	@echo "  lint               - Run code linters on all components"
	@echo "  fmt                - Format code for all components"
	@echo "  proto              - Generate protobuf files"
	@echo "  health             - Check all service health"
	@echo "  docker-build       - Build Docker images"
	@echo "  docker-up          - Start services with Docker Compose"
	@echo "  docker-down        - Stop Docker Compose services"

# Setup development environment
setup:
	@echo "🚀 Setting up fr0g.ai development environment..."
	@echo "📦 Installing dependencies for all components..."
	@$(MAKE) deps
	@echo "🔨 Building all components..."
	@$(MAKE) build-all
	@echo "✅ Setup complete!"

# Build all components
build-all:
	@echo "🔨 Building all fr0g.ai components..."
	@echo "Building shared configuration..."
	@cd pkg/config && go build ./...
	@echo "Building fr0g-ai-aip..."
	@cd fr0g-ai-aip && $(MAKE) build
	@echo "Building fr0g-ai-bridge..."
	@cd fr0g-ai-bridge && $(MAKE) build
	@echo "Building fr0g-ai-master-control..."
	@cd fr0g-ai-master-control && $(MAKE) build
	@echo "Building fr0g-ai-io..."
	@cd fr0g-ai-io && $(MAKE) build
	@echo "✅ All components built successfully!"

# Clean all build artifacts
clean-all:
	@echo "🧹 Cleaning all build artifacts..."
	@cd fr0g-ai-aip && $(MAKE) clean || true
	@cd fr0g-ai-bridge && $(MAKE) clean || true
	@cd fr0g-ai-master-control && $(MAKE) clean || true
	@cd fr0g-ai-io && $(MAKE) clean || true
	@rm -rf bin/ || true
	@echo "✅ All artifacts cleaned!"

# Run tests for all components
test-all:
	@echo "🧪 Running tests for all components..."
	@echo "Testing shared configuration..."
	@cd pkg/config && go test ./...
	@echo "Testing fr0g-ai-aip..."
	@cd fr0g-ai-aip && $(MAKE) test
	@echo "Testing fr0g-ai-bridge..."
	@cd fr0g-ai-bridge && $(MAKE) test
	@echo "Testing fr0g-ai-master-control..."
	@cd fr0g-ai-master-control && $(MAKE) test
	@echo "Testing fr0g-ai-io..."
	@cd fr0g-ai-io && $(MAKE) test
	@echo "✅ All tests completed!"

# Run individual services
run-aip:
	@echo "🚀 Starting fr0g-ai-aip service..."
	@cd fr0g-ai-aip && $(MAKE) run

run-bridge:
	@echo "🚀 Starting fr0g-ai-bridge service..."
	@cd fr0g-ai-bridge && $(MAKE) run

run-mcp:
	@echo "🚀 Starting fr0g-ai-master-control service..."
	@cd fr0g-ai-master-control && $(MAKE) run

run-io:
	@echo "🚀 Starting fr0g-ai-io service..."
	@cd fr0g-ai-io && $(MAKE) run

# Install dependencies for all components
deps:
	@echo "📦 Installing dependencies for all components..."
	@echo "Installing shared config dependencies..."
	@cd pkg/config && go mod tidy && go mod download
	@echo "Installing fr0g-ai-aip dependencies..."
	@cd fr0g-ai-aip && $(MAKE) deps
	@echo "Installing fr0g-ai-bridge dependencies..."
	@cd fr0g-ai-bridge && $(MAKE) deps
	@echo "Installing fr0g-ai-master-control dependencies..."
	@cd fr0g-ai-master-control && $(MAKE) deps
	@echo "Installing fr0g-ai-io dependencies..."
	@cd fr0g-ai-io && $(MAKE) deps
	@echo "✅ All dependencies installed!"

# Code quality checks for all components
lint:
	@echo "🔍 Running linters on all components..."
	@cd pkg/config && golangci-lint run || echo "⚠️  Install golangci-lint for better linting"
	@cd fr0g-ai-aip && $(MAKE) lint
	@cd fr0g-ai-bridge && $(MAKE) lint
	@cd fr0g-ai-master-control && $(MAKE) lint
	@cd fr0g-ai-io && $(MAKE) lint

fmt:
	@echo "🎨 Formatting code for all components..."
	@cd pkg/config && go fmt ./...
	@cd fr0g-ai-aip && $(MAKE) fmt
	@cd fr0g-ai-bridge && $(MAKE) fmt
	@cd fr0g-ai-master-control && $(MAKE) fmt
	@cd fr0g-ai-io && $(MAKE) fmt

# Generate protobuf files
proto:
	@echo "🔧 Generating protobuf files..."
	@cd fr0g-ai-aip && $(MAKE) proto || echo "⚠️  Protobuf generation failed for AIP"
	@cd fr0g-ai-bridge && $(MAKE) proto || echo "⚠️  Protobuf generation failed for Bridge"
	@cd fr0g-ai-io && $(MAKE) proto || echo "⚠️  Protobuf generation failed for IO"

# Health check all services
health:
	@echo "🏥 Checking all service health..."
	@echo "Checking fr0g-ai-aip (port 8080)..."
	@curl -sf http://localhost:8080/health && echo "✅ AIP service healthy" || echo "❌ AIP service down"
	@echo "Checking fr0g-ai-bridge (port 8082)..."
	@curl -sf http://localhost:8082/health && echo "✅ Bridge service healthy" || echo "❌ Bridge service down"
	@echo "Checking fr0g-ai-master-control (port 8081)..."
	@curl -sf http://localhost:8081/health && echo "✅ MCP service healthy" || echo "❌ MCP service down"
	@echo "Checking fr0g-ai-io (port 8083)..."
	@curl -sf http://localhost:8083/health && echo "✅ IO service healthy" || echo "❌ IO service down"

# Docker operations
docker-build:
	@echo "🐳 Building Docker images..."
	docker-compose build

docker-up:
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "🐳 Stopping Docker Compose services..."
	docker-compose down

# Legacy targets for backward compatibility
build: build-all
clean: clean-all
test: test-all
run: run-mcp
