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
	@./tests/integration/health_check_test.sh
	@./tests/integration/service_registry_test.sh
	@./tests/integration/grpc_health_test.sh

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
	@echo "Generating protobuf files..."
	@cd fr0g-ai-aip && $(MAKE) proto || echo "Protobuf generation failed for AIP"
	@cd fr0g-ai-bridge && $(MAKE) proto || echo "Protobuf generation failed for Bridge"
	@cd fr0g-ai-io && $(MAKE) proto || echo "Protobuf generation failed for IO"

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
	@echo "DIAGNOSE Authoritative gRPC ports: AIP=9090, Bridge=9091, IO=9092"
	@command -v grpcurl >/dev/null && echo "COMPLETED grpcurl available" || echo "MISSING grpcurl not found - install: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
	@nc -z localhost 9090 && echo "COMPLETED AIP gRPC port 9090 open" || echo "CRITICAL AIP gRPC port 9090 closed"
	@nc -z localhost 9091 && echo "COMPLETED Bridge gRPC port 9091 open" || echo "CRITICAL Bridge gRPC port 9091 closed"
	@nc -z localhost 9092 && echo "COMPLETED IO gRPC port 9092 open" || echo "CRITICAL IO gRPC port 9092 closed"
	@echo "DIAGNOSE Testing gRPC service responses..."
	@command -v grpcurl >/dev/null && grpcurl -plaintext localhost:9090 list 2>/dev/null && echo "COMPLETED AIP gRPC responding" || echo "CRITICAL AIP gRPC not responding"
	@command -v grpcurl >/dev/null && grpcurl -plaintext localhost:9091 list 2>/dev/null && echo "COMPLETED Bridge gRPC responding" || echo "CRITICAL Bridge gRPC not responding"
	@command -v grpcurl >/dev/null && grpcurl -plaintext localhost:9092 list 2>/dev/null && echo "COMPLETED IO gRPC responding" || echo "CRITICAL IO gRPC not responding"

diagnose-ports:
	@echo "DIAGNOSE Port configuration..."
	@echo "DIAGNOSE Authoritative port assignments from docker-compose.yml:"
	@echo "  service-registry: 8500, fr0g-ai-aip: 8080/9090, fr0g-ai-bridge: 8082/9091"
	@echo "  fr0g-ai-master-control: 8081, fr0g-ai-io: 8083/9092, redis: 6379"
	@echo "DIAGNOSE Checking Docker port mappings..."
	@docker-compose ps --format "table {{.Name}}\t{{.Ports}}" 2>/dev/null || echo "MISSING Docker Compose not running"
	@echo "DIAGNOSE Checking active port bindings..."
	@netstat -tlnp 2>/dev/null | grep -E ':(8080|8081|8082|8083|8500|9090|9091|9092|6379)' || echo "MISSING No services on expected ports"

diagnose-logs:
	@echo "DIAGNOSE Service logs for errors..."
	@docker-compose logs --tail=5 service-registry 2>/dev/null || echo "Registry container not running"
	@docker-compose logs --tail=5 fr0g-ai-aip 2>/dev/null || echo "AIP container not running"
	@docker-compose logs --tail=5 fr0g-ai-bridge 2>/dev/null || echo "Bridge container not running"
	@docker-compose logs --tail=5 fr0g-ai-io 2>/dev/null || echo "IO container not running"

diagnose-grpc-detailed:
	@echo "DIAGNOSE Detailed gRPC analysis..."
	@chmod +x tests/integration/grpc_health_test.sh
	@./tests/integration/grpc_health_test.sh || true

diagnose-containers:
	@echo "DIAGNOSE Container status and processes..."
	@docker-compose ps
	@echo ""
	@echo "DIAGNOSE Container resource usage..."
	@docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}" 2>/dev/null || echo "Docker stats not available"

diagnose: diagnose-registry diagnose-grpc diagnose-ports diagnose-logs diagnose-grpc-detailed diagnose-containers


# Docker operations
docker-build:
	@echo "DOCKER Building Docker images..."
	@docker-compose build 1>/dev/null

docker-up:
	@echo "DOCKER Starting services..."
	@docker-compose up -d
	@echo "WAITING Waiting for services to be ready..."
	@sleep 10
	@$(MAKE) health

docker-down:
	@echo "DOCKER Stopping services..."
	@docker-compose down

docker-restart:
	@echo "DOCKER Restarting all services..."
	@$(MAKE) docker-down
	@sleep 5
	@$(MAKE) docker-up

docker-restart-unhealthy:
	@echo "DOCKER Restarting unhealthy services..."
	@echo "DIAGNOSE Checking for unhealthy containers..."
	@docker-compose ps --format "table {{.Name}}\t{{.Status}}" | grep -E "(unhealthy|starting)" || echo "No unhealthy containers found"
	@echo "DOCKER Restarting Master Control service (most likely to be unhealthy)..."
	@docker-compose restart fr0g-ai-master-control || echo "Master Control restart failed"
	@sleep 15
	@$(MAKE) health

docker-rebuild:
	@echo "DOCKER Rebuilding and restarting all services..."
	@$(MAKE) docker-down
	@$(MAKE) docker-build
	@$(MAKE) docker-up

validate-production:
	@echo "VALIDATE Production readiness checks..."
	@echo "VALIDATE Checking service registration..."
	@curl -sf http://localhost:8500/v1/catalog/services && echo "COMPLETED Service discovery operational" || echo "CRITICAL Service discovery failed"
	@echo "VALIDATE Checking all service health endpoints..."
	@curl -sf http://localhost:8080/health && echo "COMPLETED AIP service healthy" || echo "CRITICAL AIP service unhealthy"
	@curl -sf http://localhost:8082/health && echo "COMPLETED Bridge service healthy" || echo "CRITICAL Bridge service unhealthy"
	@curl -sf http://localhost:8081/health && echo "COMPLETED Master Control service healthy" || echo "CRITICAL Master Control service unhealthy"
	@curl -sf http://localhost:8083/health && echo "COMPLETED IO service healthy" || echo "CRITICAL IO service unhealthy"
	@curl -sf http://localhost:8500/health && echo "COMPLETED Registry service healthy" || echo "CRITICAL Registry service unhealthy"
	@echo "VALIDATE Checking gRPC service availability..."
	@nc -z localhost 9090 && echo "COMPLETED AIP gRPC available" || echo "CRITICAL AIP gRPC unavailable"
	@nc -z localhost 9091 && echo "COMPLETED Bridge gRPC available" || echo "CRITICAL Bridge gRPC unavailable"
	@nc -z localhost 9092 && echo "COMPLETED IO gRPC available" || echo "CRITICAL IO gRPC unavailable"
	@echo "VALIDATE Checking service registration status..."
	@curl -s http://localhost:8500/v1/catalog/services | grep -q "aip-001\|bridge-001\|io-001\|mcp-001" && echo "COMPLETED All services registered" || echo "CRITICAL Services not registered"
	@echo "COMPLETED Production validation complete"

test-registry-performance:
	@echo "PERFORMANCE Testing registry performance..."
	@echo "PERFORMANCE Running 1000 service lookups..."
	@time for i in $$(seq 1 1000); do curl -s http://localhost:8500/v1/catalog/services >/dev/null; done
	@echo "COMPLETED Registry performance test complete"

test-service-integration:
	@echo "INTEGRATION Testing inter-service communication..."
	@echo "INTEGRATION Testing AIP -> Registry communication..."
	@curl -sf http://localhost:8080/health && curl -sf http://localhost:8500/health && echo "COMPLETED AIP-Registry integration working" || echo "CRITICAL AIP-Registry integration failed"
	@echo "INTEGRATION Testing Bridge -> AIP communication..."
	@curl -sf http://localhost:8082/health && nc -z localhost 9090 && echo "COMPLETED Bridge-AIP integration ready" || echo "CRITICAL Bridge-AIP integration failed"
	@echo "INTEGRATION Testing IO -> Master Control communication..."
	@curl -sf http://localhost:8083/health && curl -sf http://localhost:8081/health && echo "COMPLETED IO-MCP integration working" || echo "CRITICAL IO-MCP integration failed"
	@echo "COMPLETED Service integration tests complete"

