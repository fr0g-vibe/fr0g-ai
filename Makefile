.PHONY: help setup init-submodules update-submodules build-all clean test-all run-aip run-bridge run-esmtp docker-build-all deps dev lint fmt health docker-logs docker-clean

# Default target
help:
	@echo "🐸 fr0g.ai - AI-Powered Security Intelligence Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  setup              - Complete development environment setup"
	@echo "  dev                - Start development environment with Docker"
	@echo "  init-submodules    - Initialize and update git submodules"
	@echo "  update-submodules  - Update submodules to latest commits"
	@echo "  build-all          - Build all submodule projects"
	@echo "  build              - Build fr0g-ai-bridge only"
	@echo "  clean              - Clean all build artifacts"
	@echo "  test-all           - Run tests for all projects"
	@echo "  test               - Run bridge tests only"
	@echo "  run-aip            - Run fr0g-ai-aip server locally"
	@echo "  run-bridge         - Run fr0g-ai-bridge server locally"
	@echo "  run                - Run bridge service (alias for run-bridge)"
	@echo "  deps               - Install dependencies for all projects"
	@echo "  proto              - Generate protobuf code for bridge"
	@echo "  docker-build-all   - Build Docker images for all projects"
	@echo "  docker             - Build Docker images"
	@echo "  lint               - Run code linters"
	@echo "  fmt                - Format code"
	@echo "  health             - Check service health"
	@echo "  docker-logs        - View Docker container logs"
	@echo "  docker-clean       - Clean Docker containers and volumes"

# Complete development setup
setup: init-submodules deps
	@echo "🚀 Setting up fr0g.ai development environment..."
	@mkdir -p data/aip data/openwebui config/bridge logs
	@cp .env.example .env 2>/dev/null || true
	@cp fr0g-ai-bridge/config.example.yaml fr0g-ai-bridge/config.yaml 2>/dev/null || true
	@echo "✅ Environment ready!"
	@echo "📝 Edit .env file with your configuration"
	@echo "🐳 Run 'make dev' to start with Docker or 'make run-aip' + 'make run-bridge' for local development"

# Development with Docker
dev: setup
	@echo "🔧 Starting development environment with Docker..."
	docker-compose up --build

# Initialize submodules
init-submodules:
	git submodule init
	git submodule update --recursive --force

# Update submodules to latest
update-submodules:
	git submodule foreach --recursive 'git clean -fd'
	git submodule foreach --recursive 'git reset --hard HEAD'
	git submodule update --remote --recursive --force

# Build all projects (build only, never launch)
build-all: init-submodules deps
	@echo "🔨 Building all fr0g.ai components..."
	@echo "Building fr0g-ai-aip..."
	@cd fr0g-ai-aip && (make build-with-grpc || make build || go build -o bin/fr0g-ai-aip ./cmd/fr0g-ai-aip || echo "❌ AIP build failed")
	@echo "Building fr0g-ai-bridge..."
	@cd fr0g-ai-bridge && (make build-with-grpc || make build || go build -o bin/fr0g-ai-bridge ./cmd/fr0g-ai-bridge || echo "❌ Bridge build failed")
	@echo "Building fr0g-ai-master-control..."
	@cd fr0g-ai-master-control && (make build || go build -o bin/fr0g-ai-master-control ./cmd/master-control || echo "❌ Master-control build failed")
	@echo "Building registry server..."
	@cd fr0g-ai-master-control && (make build-registry || go build -o bin/registry-server ./cmd/registry || echo "❌ Registry build failed")
	@echo "✅ Build process completed"

# Build bridge only
build: proto
	@echo "Building fr0g-ai-bridge..."
	cd fr0g-ai-bridge && go build -o bin/fr0g-ai-bridge cmd/fr0g-ai-bridge/main.go

# Build AIP only
build-aip: proto-aip
	@echo "Building fr0g-ai-aip..."
	cd fr0g-ai-aip && go build -o bin/fr0g-ai-aip ./main.go

# Run the bridge service
run: run-bridge

# Run fr0g-ai-bridge server (build then run)
run-bridge:
	@echo "🚀 Starting fr0g-ai-bridge..."
	@cd fr0g-ai-bridge && (make build || go build -o bin/fr0g-ai-bridge ./cmd/fr0g-ai-bridge || echo "❌ Build failed") && (test -f bin/fr0g-ai-bridge && ./bin/fr0g-ai-bridge || echo "❌ Binary not found")

# Run fr0g-ai-aip server (build then run)
run-aip:
	@echo "🚀 Starting fr0g-ai-aip..."
	@cd fr0g-ai-aip && (make build || go build -o bin/fr0g-ai-aip ./main.go || echo "❌ Build failed") && (test -f bin/fr0g-ai-aip && ./bin/fr0g-ai-aip || echo "❌ Binary not found")

# Run fr0g-ai-master-control server (build then run)
run-master-control:
	@echo "🚀 Starting fr0g-ai-master-control..."
	@cd fr0g-ai-master-control && (make build || go build -o bin/fr0g-ai-master-control ./cmd/master-control || echo "❌ Build failed") && (test -f bin/fr0g-ai-master-control && ./bin/fr0g-ai-master-control || echo "❌ Binary not found")

# Run service registry server (build then run)
run-registry:
	@echo "🚀 Starting service registry..."
	@cd fr0g-ai-master-control && (make build-registry || go build -o bin/registry-server ./cmd/registry || echo "❌ Build failed") && (test -f bin/registry-server && ./bin/registry-server || echo "❌ Binary not found")

# Build registry only
build-registry:
	@echo "Building service registry..."
	@cd fr0g-ai-master-control && go build -o bin/registry-server ./cmd/registry

# Run fr0g-ai-master-control ESMTP interceptor
run-esmtp: build-all
	@echo "Starting fr0g-ai ESMTP Threat Vector Interceptor..."
	cd fr0g-ai-master-control && ./bin/esmtp-interceptor -config configs/esmtp.yaml

# Run tests for bridge only
test:
	@echo "Running fr0g-ai-bridge tests..."
	cd fr0g-ai-bridge && go test ./...

# Run tests for all projects
test-all:
	@echo "Testing fr0g-ai-aip..."
	cd fr0g-ai-aip && make test
	@echo "Testing fr0g-ai-bridge..."
	cd fr0g-ai-bridge && go test ./...
	@echo "Testing fr0g-ai-master-control..."
	cd fr0g-ai-master-control && go test ./...

# Generate protobuf code for bridge
proto:
	@echo "Generating protobuf code for fr0g-ai-bridge..."
	cd fr0g-ai-bridge && protoc --go_out=. --go-grpc_out=. proto/bridge.proto

# Generate protobuf code for AIP
proto-aip:
	@echo "Generating protobuf code for fr0g-ai-aip..."
	cd fr0g-ai-aip && protoc --proto_path=internal/grpc/proto --go_out=internal/grpc/pb --go_opt=paths=source_relative --go-grpc_out=internal/grpc/pb --go-grpc_opt=paths=source_relative internal/grpc/proto/persona.proto

# Install dependencies for all projects
deps: init-submodules
	@echo "Generating protobuf files for fr0g-ai-aip..."
	make proto-aip || true
	@echo "Installing dependencies for fr0g-ai-aip..."
	cd fr0g-ai-aip && go mod tidy && go mod download
	@echo "Generating protobuf files for fr0g-ai-bridge..."
	cd fr0g-ai-bridge && make proto-if-needed || make proto || true
	@echo "Installing dependencies for fr0g-ai-bridge..."
	cd fr0g-ai-bridge && go mod tidy && go mod download
	@echo "Installing dependencies for fr0g-ai-master-control..."
	cd fr0g-ai-master-control && go mod tidy && go mod download

# Clean all build artifacts
clean:
	@echo "Cleaning fr0g-ai-aip..."
	cd fr0g-ai-aip && make clean || true
	@echo "Cleaning fr0g-ai-bridge..."
	rm -rf fr0g-ai-bridge/bin/
	rm -rf fr0g-ai-bridge/internal/pb/*.pb.go
	@echo "Cleaning fr0g-ai-master-control..."
	rm -rf fr0g-ai-master-control/bin/

# Build Docker images
docker:
	@echo "Building Docker images..."
	docker-compose build

# Build Docker images for all
docker-build-all:
	@echo "Building Docker image for fr0g-ai-aip..."
	cd fr0g-ai-aip && docker build -t fr0g-ai-aip .
	@echo "Building Docker image for fr0g-ai-bridge..."
	cd fr0g-ai-bridge && docker build -t fr0g-ai-bridge .

# Code quality checks
lint:
	@echo "🔍 Running linters..."
	@cd fr0g-ai-aip && golangci-lint run || echo "⚠️  Install golangci-lint for better linting"
	@cd fr0g-ai-bridge && golangci-lint run || echo "⚠️  Install golangci-lint for better linting"
	@cd fr0g-ai-master-control && golangci-lint run || echo "⚠️  Install golangci-lint for better linting"

fmt:
	@echo "🎨 Formatting code..."
	@cd fr0g-ai-aip && go fmt ./...
	@cd fr0g-ai-bridge && go fmt ./...
	@cd fr0g-ai-master-control && go fmt ./...

# Health check all services
health:
	@echo "🏥 Checking service health..."
	@curl -sf http://localhost:8500/health && echo "✅ Service registry healthy" || echo "❌ Service registry down"
	@curl -sf http://localhost:8080/health && echo "✅ AIP service healthy" || echo "❌ AIP service down"
	@curl -sf http://localhost:8082/health && echo "✅ Bridge service healthy" || echo "❌ Bridge service down"
	@curl -sf http://localhost:8081/health && echo "✅ Master-control service healthy" || echo "❌ Master-control service down"
	@(command -v nc >/dev/null 2>&1 && nc -z localhost 2525 && echo "✅ ESMTP Interceptor healthy") || echo "❌ ESMTP Interceptor down"

# Quick health summary
health-summary:
	@echo "🎉 fr0g.ai Service Status Summary:"
	@echo "=================================="
	@curl -sf http://localhost:8080/health | jq -r '"AIP: \(.status) - \(.persona_count) personas loaded"' 2>/dev/null || echo "❌ AIP: Down"
	@curl -sf http://localhost:8082/health | jq -r '"Bridge: \(.status) - \(.service)"' 2>/dev/null || echo "❌ Bridge: Down"
	@curl -sf http://localhost:8081/health | jq -r '"Master-Control: \(.status) - Intelligence: \(.intelligence.status)"' 2>/dev/null || echo "❌ Master-Control: Down"
	@echo "=================================="
	@echo "✅ All core services operational!"

# Detailed health check with verbose output
health-verbose:
	@echo "🏥 Detailed service health check..."
	@echo "Checking AIP service (port 8080)..."
	@curl -v http://localhost:8080/health 2>&1 || echo "❌ AIP service connection failed"
	@echo "Checking Bridge service (port 8082)..."
	@curl -v http://localhost:8082/health 2>&1 || echo "❌ Bridge service connection failed"
	@echo "Checking Master-control service (port 8081)..."
	@curl -v http://localhost:8081/health 2>&1 || echo "❌ Master-control service connection failed"
	@echo "Checking ESMTP port (2525)..."
	@nc -z localhost 2525 && echo "✅ ESMTP port open" || echo "❌ ESMTP port closed"

# Start all services in background for testing
start-services:
	@echo "🚀 Starting all fr0g.ai services..."
	@cd fr0g-ai-aip && ./bin/fr0g-ai-aip -server &
	@cd fr0g-ai-bridge && ./bin/fr0g-ai-bridge &
	@cd fr0g-ai-master-control && ./bin/fr0g-ai-master-control &
	@echo "✅ All services started in background"
	@echo "💡 Use 'make health' to check status"
	@echo "💡 Use 'make stop-services' to stop all services"

# Stop all services
stop-services:
	@echo "🛑 Stopping all fr0g.ai services..."
	@pkill -f fr0g-ai-aip || true
	@pkill -f fr0g-ai-bridge || true
	@pkill -f fr0g-ai-master-control || true
	@echo "✅ All services stopped"

# Development helpers
dev-deps:
	@echo "Installing development dependencies..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Test endpoints
test-rest:
	@echo "Testing REST endpoints..."
	curl -X POST http://localhost:8081/api/v1/chat \
		-H "Content-Type: application/json" \
		-d '{"message": "Hello from fr0g.ai!", "model": "gpt-3.5-turbo"}' | jq .

test-grpc:
	@echo "Testing gRPC endpoints..."
	grpcurl -plaintext localhost:9091 list

# Docker utilities
docker-logs:
	@echo "📋 Viewing Docker container logs..."
	docker-compose logs -f

docker-clean:
	@echo "🧹 Cleaning Docker containers and volumes..."
	docker-compose down -v
	docker system prune -f

# Production deployment
prod-up:
	@echo "🚀 Starting production environment..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

prod-down:
	@echo "🛑 Stopping production environment..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# Testing
test-deployment:
	@echo "🧪 Running deployment tests..."
	@chmod +x scripts/test-deployment.sh
	@./scripts/test-deployment.sh

# Monitoring
monitoring-up:
	@echo "📊 Starting monitoring stack..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d prometheus grafana

monitoring-down:
	@echo "📊 Stopping monitoring stack..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml stop prometheus grafana
