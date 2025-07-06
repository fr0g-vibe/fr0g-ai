.PHONY: help setup build-all clean-all test test-integration test-performance test-all run-aip run-bridge run-mcp run-io health diagnose proto docker-build docker-up docker-down

# Default target
help:
	@echo "fr0g.ai - AI-Powered Security Intelligence Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  setup              - Initialize development environment"
	@echo "  build-all          - Build all components"
	@echo "  clean-all          - Clean all build artifacts"
	@echo "  test               - Run unit tests"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-performance   - Run performance tests"
	@echo "  test-all           - Run all test suites"
	@echo "  run-aip            - Run AIP service"
	@echo "  run-bridge         - Run Bridge service"
	@echo "  run-mcp            - Run Master Control service"
	@echo "  run-io             - Run I/O service"
	@echo "  health             - Check service health"
	@echo "  diagnose           - Run diagnostic checks"
	@echo "  proto              - Generate protobuf files"
	@echo "  docker-build       - Build Docker images"
	@echo "  docker-up          - Start services"
	@echo "  docker-down        - Stop services"

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

# Core test targets
test:
	@echo "TESTING Running unit tests for all components..."
	@cd pkg/config && go test ./...
	@cd fr0g-ai-aip && $(MAKE) test
	@cd fr0g-ai-bridge && $(MAKE) test
	@cd fr0g-ai-master-control && $(MAKE) test
	@cd fr0g-ai-io && $(MAKE) test
	@cd fr0g-ai-registry && $(MAKE) test
	@echo "COMPLETED All unit tests completed!"

test-integration:
	@echo "TESTING Running integration tests..."
	@chmod +x tests/integration/*.sh
	@./tests/integration/end_to_end_test.sh
	@./tests/integration/service_registry_test.sh
	@./tests/integration/api_test.sh

test-performance:
	@echo "PERFORMANCE Running performance tests..."
	@chmod +x tests/integration/performance_test.sh
	@./tests/integration/performance_test.sh

test-all: test test-integration test-performance
	@echo "COMPLETED All test suites completed!"

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

# Diagnostic commands
diagnose-registry:
	@echo "DIAGNOSE Service Registry API endpoints..."
	@curl -sf http://localhost:8500/health && echo "✓ Registry health OK" || echo "✗ Registry health FAILED"
	@curl -s -w "HTTP %{http_code}\n" http://localhost:8500/v1/catalog/services || echo "✗ Discovery endpoint test FAILED"

diagnose-grpc:
	@echo "DIAGNOSE gRPC service health..."
	@command -v grpcurl >/dev/null && echo "✓ grpcurl available" || echo "✗ grpcurl not found"
	@nc -z localhost 9090 && echo "✓ AIP gRPC port open" || echo "✗ AIP gRPC port closed"
	@nc -z localhost 9091 && echo "✓ Bridge gRPC port open" || echo "✗ Bridge gRPC port closed"

diagnose-ports:
	@echo "DIAGNOSE Port configuration..."
	@echo "Expected ports: Registry(8500), AIP(8080,9090), Bridge(8082,9091), MCP(8081), IO(8083,9093)"
	@netstat -tlnp 2>/dev/null | grep -E ':(8080|8081|8082|8083|8500|9090|9091|9093)' || echo "No services on expected ports"

diagnose-logs:
	@echo "DIAGNOSE Service logs for errors..."
	@docker-compose logs --tail=5 service-registry 2>/dev/null || echo "Registry container not running"
	@docker-compose logs --tail=5 fr0g-ai-aip 2>/dev/null || echo "AIP container not running"
	@docker-compose logs --tail=5 fr0g-ai-bridge 2>/dev/null || echo "Bridge container not running"
	@docker-compose logs --tail=5 fr0g-ai-io 2>/dev/null || echo "IO container not running"

diagnose: diagnose-registry diagnose-grpc diagnose-ports diagnose-logs


# Docker operations
docker-build:
	@echo "DOCKER Building Docker images..."
	@docker-compose build

docker-up:
	@echo "DOCKER Starting services..."
	@docker-compose up -d
	@echo "WAITING Waiting for services to be ready..."
	@sleep 10
	@$(MAKE) health

docker-down:
	@echo "DOCKER Stopping services..."
	@docker-compose down

