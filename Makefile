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

# Build all projects
build-all: init-submodules deps
	@echo "Building fr0g-ai-aip..."
	cd fr0g-ai-aip && make build-with-grpc
	@echo "Building fr0g-ai-bridge..."
	cd fr0g-ai-bridge && make build-with-grpc
	@echo "Building fr0g-ai-master-control..."
	cd fr0g-ai-master-control && go build -o bin/esmtp-interceptor ./cmd/esmtp-interceptor

# Build bridge only
build: proto
	@echo "Building fr0g-ai-bridge..."
	cd fr0g-ai-bridge && go build -o bin/fr0g-ai-bridge cmd/fr0g-ai-bridge/main.go

# Run the bridge service
run: run-bridge

# Run fr0g-ai-bridge server
run-bridge: build
	@echo "Starting fr0g-ai-bridge..."
	cd fr0g-ai-bridge && ./bin/fr0g-ai-bridge

# Run fr0g-ai-aip server
run-aip: build-all
	cd fr0g-ai-aip && ./bin/fr0g-ai-aip -server -grpc

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

# Install dependencies for all projects
deps: init-submodules
	@echo "Generating protobuf files for fr0g-ai-aip..."
	cd fr0g-ai-aip && make proto-if-needed || make proto || true
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
	@curl -sf http://localhost:8080/health && echo "✅ AIP service healthy" || echo "❌ AIP service down"
	@curl -sf http://localhost:8081/health && echo "✅ Bridge service healthy" || echo "❌ Bridge service down"
	@nc -z localhost 2525 && echo "✅ ESMTP Interceptor healthy" || echo "❌ ESMTP Interceptor down"

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
