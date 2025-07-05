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

# Build the MCP server
build:
	@echo "🔨 Building fr0g.ai MCP server..."
	@mkdir -p bin
	go build -o bin/fr0g-ai-mcp ./cmd/mcp

# Build AIP only
build-aip: proto-aip
	@echo "Building fr0g-ai-aip..."
	cd fr0g-ai-aip && go build -o bin/fr0g-ai-aip ./main.go

# Run the MCP server
run: build
	@echo "🚀 Starting fr0g.ai MCP server..."
	./bin/fr0g-ai-mcp


# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...


# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/

# Code quality checks
lint:
	@echo "🔍 Running linters..."
	golangci-lint run || echo "⚠️  Install golangci-lint for better linting"

fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Health check MCP service
health:
	@echo "🏥 Checking MCP service health..."
	@curl -sf http://localhost:8081/health && echo "✅ MCP service healthy" || echo "❌ MCP service down"
