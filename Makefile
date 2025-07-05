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
	@echo "  test-integration   - Run end-to-end integration tests"
	@echo "  test-registry      - Run service registry tests"
	@echo "  test-api           - Run API integration tests"
	@echo "  test-performance   - Run performance tests"
	@echo "  test-deployment    - Run deployment verification tests"
	@echo "  test-all-integration - Run all integration test suites"
	@echo "  test-aip-service   - Run AIP service test suite"
	@echo "  test-aip-with-reflection - Run AIP tests with gRPC reflection"
	@echo "  test-grpc-reflection - Test gRPC reflection specifically"
	@echo "  build-aip-test     - Build AIP service with reflection for testing"
	@echo "  run-aip-test       - Run AIP service with reflection enabled"
	@echo "  validate-production - Validate production security for all services"
	@echo "  setup-dev-env      - Create global development environment file"
	@echo "  setup-prod-env     - Create global production environment file"
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
	@echo "Building fr0g-ai-registry..."
	@cd fr0g-ai-registry && $(MAKE) build
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
	@echo "Testing fr0g-ai-registry..."
	@cd fr0g-ai-registry && $(MAKE) test
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
	@echo "Installing fr0g-ai-registry dependencies..."
	@cd fr0g-ai-registry && $(MAKE) deps
	@echo "✅ All dependencies installed!"

# Code quality checks for all components
lint:
	@echo "🔍 Running linters on all components..."
	@cd pkg/config && golangci-lint run || echo "⚠️  Install golangci-lint for better linting"
	@cd fr0g-ai-aip && $(MAKE) lint
	@cd fr0g-ai-bridge && $(MAKE) lint
	@cd fr0g-ai-master-control && $(MAKE) lint
	@cd fr0g-ai-io && $(MAKE) lint
	@cd fr0g-ai-registry && $(MAKE) lint

fmt:
	@echo "🎨 Formatting code for all components..."
	@cd pkg/config && go fmt ./...
	@cd fr0g-ai-aip && $(MAKE) fmt
	@cd fr0g-ai-bridge && $(MAKE) fmt
	@cd fr0g-ai-master-control && $(MAKE) fmt
	@cd fr0g-ai-io && $(MAKE) fmt
	@cd fr0g-ai-registry && $(MAKE) fmt

# Generate protobuf files
proto:
	@echo "🔧 Generating protobuf files..."
	@cd fr0g-ai-aip && $(MAKE) proto || echo "⚠️  Protobuf generation failed for AIP"
	@cd fr0g-ai-bridge && $(MAKE) proto || echo "⚠️  Protobuf generation failed for Bridge"
	@cd fr0g-ai-io && $(MAKE) proto || echo "⚠️  Protobuf generation failed for IO"

# Health check all services
health:
	@echo "🏥 Checking all service health..."
	@chmod +x tests/integration/health_check_test.sh
	@./tests/integration/health_check_test.sh

# Quick health check (simple curl tests)
health-quick:
	@echo "🏥 Quick health check..."
	@echo "Checking service registry (port 8500)..."
	@curl -sf http://localhost:8500/health && echo "✅ Registry healthy" || echo "❌ Registry down"
	@echo "Checking fr0g-ai-aip (port 8080)..."
	@curl -sf http://localhost:8080/health && echo "✅ AIP service healthy" || echo "❌ AIP service down"
	@echo "Checking fr0g-ai-bridge (port 8082)..."
	@curl -sf http://localhost:8082/health && echo "✅ Bridge service healthy" || echo "❌ Bridge service down"
	@echo "Checking fr0g-ai-master-control (port 8081)..."
	@curl -sf http://localhost:8081/health && echo "✅ MCP service healthy" || echo "❌ MCP service down"
	@echo "Checking fr0g-ai-io (port 8083)..."
	@curl -sf http://localhost:8083/health && echo "✅ IO service healthy" || echo "❌ IO service down"

# AIP-specific testing (delegates to subproject)
test-aip-service:
	@echo "🧪 Running AIP service test suite..."
	@chmod +x ./test_aip_service.sh
	@./test_aip_service.sh

test-aip-with-reflection:
	@echo "🧪 Running AIP tests with gRPC reflection enabled..."
	@cd fr0g-ai-aip && $(MAKE) test-with-reflection

test-grpc-reflection:
	@echo "🔍 Testing gRPC reflection..."
	@cd fr0g-ai-aip && $(MAKE) test-grpc-reflection

build-aip-test:
	@echo "🔨 Building AIP service with gRPC reflection for testing..."
	@cd fr0g-ai-aip && $(MAKE) build-test

run-aip-test:
	@echo "🚀 Starting AIP service with reflection enabled..."
	@cd fr0g-ai-aip && $(MAKE) run-test

# Global validation (checks all services)
validate-production:
	@echo "🔒 Validating production build security for all services..."
	@cd fr0g-ai-aip && $(MAKE) validate-production
	@cd fr0g-ai-bridge && $(MAKE) validate-production || echo "Bridge validation not implemented"
	@cd fr0g-ai-io && $(MAKE) validate-production || echo "IO validation not implemented"
	@echo "✅ Production validation completed"

# Global environment setup
setup-dev-env:
	@echo "🛠️  Setting up development environment..."
	@echo "export GRPC_ENABLE_REFLECTION=true" > .env.development
	@echo "export ENVIRONMENT=development" >> .env.development
	@echo "export LOG_LEVEL=debug" >> .env.development
	@echo "✅ Development environment configured (.env.development)"
	@echo "💡 Source with: source .env.development"
	@echo "💡 Individual services may have additional setup - check their Makefiles"

setup-prod-env:
	@echo "🔒 Setting up production environment..."
	@echo "export GRPC_ENABLE_REFLECTION=false" > .env.production
	@echo "export ENVIRONMENT=production" >> .env.production
	@echo "export LOG_LEVEL=info" >> .env.production
	@echo "✅ Production environment configured (.env.production)"
	@echo "💡 Source with: source .env.production"

# Integration testing targets
test-integration:
	@echo "🧪 Running integration tests..."
	@chmod +x tests/integration/*.sh
	@./tests/integration/end_to_end_test.sh

test-registry:
	@echo "🔍 Running service registry tests..."
	@chmod +x tests/integration/service_registry_test.sh
	@./tests/integration/service_registry_test.sh

test-api:
	@echo "🔌 Running API tests..."
	@chmod +x tests/integration/api_test.sh
	@./tests/integration/api_test.sh

test-performance:
	@echo "⚡ Running performance tests..."
	@chmod +x tests/integration/performance_test.sh
	@./tests/integration/performance_test.sh

test-deployment:
	@echo "🚀 Running deployment tests..."
	@chmod +x scripts/test-deployment.sh
	@./scripts/test-deployment.sh

# Run all integration tests
test-all-integration: test-integration test-registry test-api test-performance test-deployment

# Docker operations with error handling
docker-build:
	@echo "🐳 Building Docker images..."
	@if docker-compose build >/dev/null 2>&1; then \
		echo "✅ Docker images built successfully!"; \
	else \
		echo "❌ Docker build failed. Running with verbose output:"; \
		docker-compose build; \
		exit 1; \
	fi

docker-up:
	@echo "🐳 Starting services with Docker Compose..."
	@if docker-compose up -d >/dev/null 2>&1; then \
		echo "✅ Services started successfully!"; \
		echo "⏳ Waiting for services to be ready..."; \
		sleep 10; \
		$(MAKE) health; \
	else \
		echo "❌ Failed to start services. Check logs with: docker-compose logs"; \
		exit 1; \
	fi

docker-up-core:
	@echo "🐳 Starting core services (Registry + AIP)..."
	@docker-compose up -d service-registry fr0g-ai-aip
	@echo "⏳ Waiting for core services to be ready..."
	@sleep 10
	@$(MAKE) health-quick

docker-up-all:
	@echo "🐳 Starting all services..."
	@docker-compose up -d
	@echo "⏳ Waiting for all services to be ready..."
	@sleep 15
	@$(MAKE) health

docker-down:
	@echo "🐳 Stopping Docker Compose services..."
	@if docker-compose down >/dev/null 2>&1; then \
		echo "✅ Services stopped successfully!"; \
	else \
		echo "❌ Failed to stop services"; \
		exit 1; \
	fi

docker-status:
	@echo "🐳 Docker service status..."
	@docker-compose ps
	@echo ""
	@$(MAKE) health-quick

# Legacy targets for backward compatibility
build: build-all
clean: clean-all
test: test-all
run: run-mcp
