.PHONY: help setup build build-all clean clean-all test test-all run run-aip run-bridge run-mcp deps lint fmt health proto docker-build docker-up docker-down

# Default target
help:
	@echo "fr0g.ai fr0g.ai - AI-Powered Security Intelligence Platform"
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
	@echo "STARTING Setting up fr0g.ai development environment..."
	@echo "INSTALLING Installing dependencies for all components..."
	@$(MAKE) deps
	@echo "BUILDING Building all components..."
	@$(MAKE) build-all
	@echo "COMPLETED Setup complete!"

# Build all components
build-all:
	@echo "BUILDING Building all fr0g.ai components..."
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
	@echo "COMPLETED All components built successfully!"

# Clean all build artifacts
clean-all:
	@echo "CLEANING Cleaning all build artifacts..."
	@cd fr0g-ai-aip && $(MAKE) clean || true
	@cd fr0g-ai-bridge && $(MAKE) clean || true
	@cd fr0g-ai-master-control && $(MAKE) clean || true
	@cd fr0g-ai-io && $(MAKE) clean || true
	@rm -rf bin/ || true
	@echo "COMPLETED All artifacts cleaned!"

# Run tests for all components
test-all:
	@echo "TESTING Running tests for all components..."
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
	@echo "COMPLETED All tests completed!"

# Run individual services
run-aip:
	@echo "STARTING Starting fr0g-ai-aip service..."
	@cd fr0g-ai-aip && $(MAKE) run

run-bridge:
	@echo "STARTING Starting fr0g-ai-bridge service..."
	@cd fr0g-ai-bridge && $(MAKE) run

run-mcp:
	@echo "STARTING Starting fr0g-ai-master-control service..."
	@cd fr0g-ai-master-control && $(MAKE) run

run-io:
	@echo "STARTING Starting fr0g-ai-io service..."
	@cd fr0g-ai-io && $(MAKE) run

# Install dependencies for all components
deps:
	@echo "INSTALLING Installing dependencies for all components..."
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
	@echo "COMPLETED All dependencies installed!"

# Code quality checks for all components
lint:
	@echo "CHECKING Running linters on all components..."
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
	@echo "HEALTH Checking all service health..."
	@chmod +x tests/integration/health_check_test.sh
	@./tests/integration/health_check_test.sh

# Emergency diagnostic commands
diagnose-registry:
	@echo "DIAGNOSE Service Registry API endpoints..."
	@echo "Testing registry health..."
	@curl -sf http://localhost:8500/health && echo "✓ Registry health OK" || echo "✗ Registry health FAILED"
	@echo "Testing service registration endpoint..."
	@curl -s -w "HTTP %{http_code}\n" http://localhost:8500/v1/agent/service/register -X POST -H "Content-Type: application/json" -d '{"ID":"test","Name":"test","Port":8000}' || echo "✗ Registration endpoint test FAILED"
	@echo "Testing service discovery endpoint..."
	@curl -s -w "HTTP %{http_code}\n" http://localhost:8500/v1/catalog/services || echo "✗ Discovery endpoint test FAILED"
	@echo "Testing service health endpoint..."
	@curl -s -w "HTTP %{http_code}\n" http://localhost:8500/v1/health/service/test || echo "✗ Health endpoint test FAILED"

diagnose-grpc:
	@echo "DIAGNOSE gRPC service health..."
	@echo "Checking if grpcurl is available..."
	@command -v grpcurl >/dev/null && echo "✓ grpcurl available" || echo "✗ grpcurl not found - install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
	@echo "Testing AIP gRPC (port 9090)..."
	@nc -z localhost 9090 && echo "✓ Port 9090 open" || echo "✗ Port 9090 closed"
	@command -v grpcurl >/dev/null && (grpcurl -plaintext localhost:9090 list >/dev/null 2>&1 && echo "✓ AIP gRPC responding" || echo "✗ AIP gRPC not responding") || echo "  (grpcurl not available)"
	@echo "Testing Bridge gRPC (port 9091)..."
	@nc -z localhost 9091 && echo "✓ Port 9091 open" || echo "✗ Port 9091 closed"
	@command -v grpcurl >/dev/null && (grpcurl -plaintext localhost:9091 list >/dev/null 2>&1 && echo "✓ Bridge gRPC responding" || echo "✗ Bridge gRPC not responding") || echo "  (grpcurl not available)"
	@echo "Testing IO gRPC (port 9093)..."
	@nc -z localhost 9093 && echo "✓ Port 9093 open" || echo "✗ Port 9093 closed"
	@command -v grpcurl >/dev/null && (grpcurl -plaintext localhost:9093 list >/dev/null 2>&1 && echo "✓ IO gRPC responding" || echo "✗ IO gRPC not responding") || echo "  (grpcurl not available)"

diagnose-ports:
	@echo "DIAGNOSE Port configuration..."
	@echo "Checking port usage..."
	@echo "Expected ports: Registry(8500), AIP(8080,9090), Bridge(8082,9091), MCP(8081), IO(8083,9093)"
	@netstat -tlnp 2>/dev/null | grep -E ':(8080|8081|8082|8083|8500|9090|9091|9092|9093)' || echo "No services listening on expected ports"
	@echo "Checking Docker port mapping..."
	@docker-compose ps 2>/dev/null || echo "Docker Compose not running"

diagnose-logs:
	@echo "DIAGNOSE Service logs for errors..."
	@echo "=== Registry Logs ==="
	@docker-compose logs --tail=10 service-registry 2>/dev/null || echo "Registry container not running"
	@echo "=== AIP Logs ==="
	@docker-compose logs --tail=10 fr0g-ai-aip 2>/dev/null || echo "AIP container not running"
	@echo "=== Bridge Logs ==="
	@docker-compose logs --tail=10 fr0g-ai-bridge 2>/dev/null || echo "Bridge container not running"
	@echo "=== IO Logs ==="
	@docker-compose logs --tail=10 fr0g-ai-io 2>/dev/null || echo "IO container not running"

diagnose-all: diagnose-registry diagnose-grpc diagnose-ports diagnose-logs

# Test the bridge chat completions endpoint (requires services running)
test-bridge-chat:
	@echo "TESTING Bridge chat completions endpoint..."
	@if ! curl -sf http://localhost:8082/health >/dev/null 2>&1; then \
		echo "ERROR Bridge service not running. Start with: make docker-up"; \
		exit 1; \
	fi
	@curl -X POST http://localhost:8082/v1/chat/completions \
		-H "Content-Type: application/json" \
		-d '{"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"Hello, test message"}]}' \
		|| echo "FAILED Chat completions endpoint not working"

# Test service registry registration (requires services running)
test-registry-register:
	@echo "TESTING Service registry registration..."
	@if ! curl -sf http://localhost:8500/health >/dev/null 2>&1; then \
		echo "ERROR Registry service not running. Start with: make docker-up"; \
		exit 1; \
	fi
	@curl -X POST http://localhost:8500/v1/agent/service/register \
		-H "Content-Type: application/json" \
		-d '{"ID":"test-service","Name":"test","Port":8000,"Address":"localhost"}' \
		|| echo "FAILED Service registration not working"

# Test gRPC connectivity (requires services running)
test-grpc-connectivity:
	@echo "TESTING gRPC service connectivity..."
	@if ! curl -sf http://localhost:8500/health >/dev/null 2>&1; then \
		echo "ERROR Services not running. Start with: make docker-up"; \
		exit 1; \
	fi
	@echo "Testing AIP gRPC (port 9090)..."
	@command -v grpcurl >/dev/null && (grpcurl -plaintext localhost:9090 list >/dev/null 2>&1 && echo "✓ AIP gRPC responding" || echo "✗ AIP gRPC not responding") || echo "grpcurl not available"
	@echo "Testing Bridge gRPC (port 9091)..."
	@command -v grpcurl >/dev/null && (grpcurl -plaintext localhost:9091 list >/dev/null 2>&1 && echo "✓ Bridge gRPC responding" || echo "✗ Bridge gRPC not responding") || echo "grpcurl not available"
	@echo "Testing IO gRPC (port 9093)..."
	@command -v grpcurl >/dev/null && (grpcurl -plaintext localhost:9093 list >/dev/null 2>&1 && echo "✓ IO gRPC responding" || echo "✗ IO gRPC not responding") || echo "grpcurl not available"

# Test with automatic service startup
test-bridge-chat-auto:
	@echo "TESTING Bridge chat completions endpoint (auto-start services)..."
	@if ! curl -sf http://localhost:8082/health >/dev/null 2>&1; then \
		echo "Starting services..."; \
		docker-compose up -d >/dev/null 2>&1; \
		sleep 15; \
	fi
	@$(MAKE) test-bridge-chat

test-registry-register-auto:
	@echo "TESTING Service registry registration (auto-start services)..."
	@if ! curl -sf http://localhost:8500/health >/dev/null 2>&1; then \
		echo "Starting services..."; \
		docker-compose up -d >/dev/null 2>&1; \
		sleep 15; \
	fi
	@$(MAKE) test-registry-register

test-grpc-connectivity-auto:
	@echo "TESTING gRPC service connectivity (auto-start services)..."
	@if ! curl -sf http://localhost:8500/health >/dev/null 2>&1; then \
		echo "Starting services..."; \
		docker-compose up -d >/dev/null 2>&1; \
		sleep 15; \
	fi
	@$(MAKE) test-grpc-connectivity

# Health check with clean service restart
health-clean:
	@echo "HEALTH Performing clean health check with service restart..."
	@echo "Stopping any running services..."
	@docker-compose down >/dev/null 2>&1 || true
	@echo "Starting services fresh..."
	@docker-compose up -d >/dev/null 2>&1
	@echo "Waiting for services to be ready..."
	@sleep 15
	@echo "Running health checks..."
	@chmod +x tests/integration/health_check_test.sh
	@./tests/integration/health_check_test.sh

# Run clean integration test (stops services, starts fresh, tests, then cleans up)
test-clean:
	@echo "TESTING Running clean integration test with fresh services..."
	@chmod +x tests/integration/clean_test.sh
	@./tests/integration/clean_test.sh

# Quick health check (simple curl tests)
health-quick:
	@echo "HEALTH Quick health check..."
	@echo "Checking service registry (port 8500)..."
	@curl -sf http://localhost:8500/health && echo "COMPLETED Registry healthy" || echo "FAILED Registry down"
	@echo "Checking fr0g-ai-aip (port 8080)..."
	@curl -sf http://localhost:8080/health && echo "COMPLETED AIP service healthy" || echo "FAILED AIP service down"
	@echo "Checking fr0g-ai-bridge (port 8082)..."
	@curl -sf http://localhost:8082/health && echo "COMPLETED Bridge service healthy" || echo "FAILED Bridge service down"
	@echo "Checking fr0g-ai-master-control (port 8081)..."
	@curl -sf http://localhost:8081/health && echo "COMPLETED MCP service healthy" || echo "FAILED MCP service down"
	@echo "Checking fr0g-ai-io (port 8083)..."
	@curl -sf http://localhost:8083/health && echo "COMPLETED IO service healthy" || echo "FAILED IO service down"


# AIP-specific testing (delegates to subproject)
test-aip-service:
	@echo "TESTING Running AIP service test suite..."
	@chmod +x ./test_aip_service.sh
	@./test_aip_service.sh

test-aip-with-reflection:
	@echo "TESTING Running AIP tests against running service with gRPC reflection..."
	@echo "🔧 gRPC reflection enables MCP integration for dynamic service discovery"
	@cd fr0g-ai-aip && $(MAKE) test-with-reflection

test-grpc-reflection:
	@echo "CHECKING Testing gRPC reflection for MCP compatibility..."
	@echo "🔧 This enables Model Context Protocol exposure for other gRPC services"
	@cd fr0g-ai-aip && $(MAKE) test-grpc-reflection

build-aip-test:
	@echo "BUILDING Building AIP service with gRPC reflection for testing and MCP integration..."
	@cd fr0g-ai-aip && $(MAKE) build-test

run-aip-test:
	@echo "STARTING Starting AIP service with reflection enabled for MCP exposure..."
	@echo "🔧 Other gRPC clients can now discover AIP services dynamically"
	@cd fr0g-ai-aip && $(MAKE) run-test

# Global validation (checks all services)
validate-production:
	@echo "SECURITY Validating production build security for all services..."
	@cd fr0g-ai-aip && $(MAKE) validate-production
	@cd fr0g-ai-bridge && $(MAKE) validate-production
	@cd fr0g-ai-io && $(MAKE) validate-production
	@cd fr0g-ai-master-control && $(MAKE) validate-production
	@cd fr0g-ai-registry && $(MAKE) validate-production
	@echo "COMPLETED Production validation completed for all services"

# Global environment setup
setup-dev-env:
	@echo "SETUP️  Setting up development environment..."
	@echo "export GRPC_ENABLE_REFLECTION=true" > .env.development
	@echo "export ENVIRONMENT=development" >> .env.development
	@echo "export LOG_LEVEL=debug" >> .env.development
	@echo "COMPLETED Development environment configured (.env.development)"
	@echo "TIP Source with: source .env.development"
	@echo "TIP Individual services may have additional setup - check their Makefiles"

setup-prod-env:
	@echo "SECURITY Setting up production environment..."
	@echo "export GRPC_ENABLE_REFLECTION=false" > .env.production
	@echo "export ENVIRONMENT=production" >> .env.production
	@echo "export LOG_LEVEL=info" >> .env.production
	@echo "COMPLETED Production environment configured (.env.production)"
	@echo "TIP Source with: source .env.production"

# Integration testing targets
test-integration:
	@echo "TESTING Running integration tests..."
	@chmod +x tests/integration/*.sh
	@./tests/integration/end_to_end_test.sh

test-registry:
	@echo "CHECKING Running service registry tests..."
	@chmod +x tests/integration/service_registry_test.sh
	@./tests/integration/service_registry_test.sh

test-registry-unit:
	@echo "TESTING Running registry unit tests..."
	@cd fr0g-ai-registry && $(MAKE) test-unit

test-registry-integration:
	@echo "TESTING Running registry integration tests..."
	@cd fr0g-ai-registry && $(MAKE) test-integration

test-registry-load:
	@echo "PERFORMANCE Running registry load tests..."
	@cd fr0g-ai-registry && $(MAKE) test-load

test-api:
	@echo "🔌 Running API tests..."
	@chmod +x tests/integration/api_test.sh
	@./tests/integration/api_test.sh

test-performance:
	@echo "PERFORMANCE Running performance tests..."
	@chmod +x tests/integration/performance_test.sh
	@./tests/integration/performance_test.sh

test-deployment:
	@echo "STARTING Running deployment tests..."
	@chmod +x scripts/test-deployment.sh
	@./scripts/test-deployment.sh

# Run all integration tests
test-all-integration: test-integration test-registry test-api test-performance test-deployment

# Docker operations with error handling
docker-build:
	@echo "DOCKER Building Docker images..."
	@if docker-compose build >/dev/null 2>&1; then \
		echo "COMPLETED Docker images built successfully!"; \
	else \
		echo "FAILED Docker build failed. Running with verbose output:"; \
		docker-compose build; \
		exit 1; \
	fi

docker-up:
	@echo "DOCKER Starting services with Docker Compose..."
	@if docker-compose up -d >/dev/null 2>&1; then \
		echo "COMPLETED Services started successfully!"; \
		echo "WAITING Waiting for services to be ready..."; \
		sleep 10; \
		$(MAKE) health; \
	else \
		echo "FAILED Failed to start services. Check logs with: docker-compose logs"; \
		exit 1; \
	fi

docker-up-core:
	@echo "DOCKER Starting core services (Registry + AIP)..."
	@docker-compose up -d service-registry fr0g-ai-aip
	@echo "WAITING Waiting for core services to be ready..."
	@sleep 10
	@$(MAKE) health-quick

docker-up-all:
	@echo "DOCKER Starting all services..."
	@docker-compose up -d
	@echo "WAITING Waiting for all services to be ready..."
	@sleep 15
	@$(MAKE) health

docker-down:
	@echo "DOCKER Stopping Docker Compose services..."
	@if docker-compose down >/dev/null 2>&1; then \
		echo "COMPLETED Services stopped successfully!"; \
	else \
		echo "FAILED Failed to stop services"; \
		exit 1; \
	fi

docker-status:
	@echo "DOCKER Docker service status..."
	@docker-compose ps
	@echo ""
	@$(MAKE) health-quick

# Legacy targets for backward compatibility
build: build-all
clean: clean-all
run: run-mcp
