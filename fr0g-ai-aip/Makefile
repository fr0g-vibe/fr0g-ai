.PHONY: build test clean run-server run-grpc run-both run-cli proto help run

# Build the application
build: proto-if-needed
	@echo "Building application..."
	go build -o bin/fr0g-ai-aip ./cmd

# Build with local gRPC support (no external dependencies)
build-with-grpc: proto-if-needed build
	@echo "gRPC support built using local JSON-over-HTTP implementation"

# Generate protobuf code (force regeneration)
proto:
	@echo "Generating protobuf code..."
	@if [ ! -f "internal/grpc/proto/persona.proto" ]; then \
		echo "Error: internal/grpc/proto/persona.proto not found"; \
		echo "Please create the proto file first"; \
		exit 1; \
	fi
	@echo "Using protoc: $(shell which protoc)"
	@echo "Using protoc-gen-go: $(shell which protoc-gen-go || echo "$(shell go env GOPATH)/bin/protoc-gen-go")"
	@echo "Using protoc-gen-go-grpc: $(shell which protoc-gen-go-grpc || echo "$(shell go env GOPATH)/bin/protoc-gen-go-grpc")"
	@mkdir -p internal/grpc/pb
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc \
		--proto_path=internal/grpc/proto \
		--go_out=internal/grpc/pb --go_opt=paths=source_relative \
		--go-grpc_out=internal/grpc/pb --go-grpc_opt=paths=source_relative \
		internal/grpc/proto/persona.proto
	@echo "Protobuf code generated successfully in internal/grpc/pb/"

# Generate protobuf code only if files don't exist
proto-if-needed:
	@if [ ! -f "internal/grpc/pb/persona.pb.go" ] || [ ! -f "internal/grpc/pb/persona_grpc.pb.go" ]; then \
		echo "Protobuf files missing, generating..."; \
		$(MAKE) proto; \
	else \
		echo "Protobuf files already exist, skipping generation"; \
	fi

# Run tests
test: proto-if-needed
	go test ./...

# Run tests with coverage
test-coverage: proto-if-needed
	go test -cover ./...

# Run tests with detailed coverage report
test-coverage-detailed: proto-if-needed
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with verbose coverage for specific package
test-coverage-verbose-grpc: proto-if-needed
	go test -v -coverprofile=grpc_coverage.out ./internal/grpc/
	go tool cover -func=grpc_coverage.out
	go tool cover -html=grpc_coverage.out -o grpc_coverage.html
	@echo "gRPC coverage report generated: grpc_coverage.html"

# Run tests with verbose output
test-verbose: proto-if-needed
	go test -v ./...

# Run tests with race detection
test-race: proto-if-needed
	go test -race ./...

# Run benchmarks
test-bench: proto-if-needed
	go test -bench=. ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf internal/grpc/pb/*.pb.go

# Run HTTP REST API server
run-server: proto-if-needed
	go run ./cmd -server

# Run gRPC server
run-grpc: proto-if-needed
	go run ./cmd -grpc

# Run both HTTP and gRPC servers
run-both: proto-if-needed
	go run ./cmd -server -grpc

# Run CLI help
run-cli:
	go run ./cmd -help

# Install dependencies
deps:
	go mod tidy

# Install protobuf tools (optional - only needed for full gRPC implementation)
install-proto-tools:
	@echo "Installing protobuf tools (this will add external dependencies)..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Make sure $(shell go env GOPATH)/bin is in your PATH"
	@echo "Note: You'll also need to add gRPC dependencies to go.mod for full gRPC support"

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Generate documentation
docs: proto-if-needed
	@echo "Generating documentation..."
	@mkdir -p docs/generated
	@echo "Generating package documentation..."
	@for pkg in $$(go list ./...); do \
		echo "=== $$pkg ===" >> docs/generated/godoc.txt; \
		go doc -all $$pkg >> docs/generated/godoc.txt 2>/dev/null || echo "Failed to generate docs for $$pkg" >> docs/generated/godoc.txt; \
		echo "" >> docs/generated/godoc.txt; \
	done
	@echo "Documentation generated in docs/generated/"


# Generate OpenAPI documentation
docs-openapi:
	@echo "OpenAPI specification available at docs/OPENAPI_SPEC.yaml"
	@echo "View with: swagger-ui-serve docs/OPENAPI_SPEC.yaml"

# Generate MCP documentation
docs-mcp:
	@echo "MCP integration documentation available at docs/MCP_INTEGRATION.md"

# Run the service (for main project integration)
run: build
	@echo "🚀 Starting fr0g-ai-aip service..."
	@echo "HTTP server will start on port 8080"
	@echo "gRPC server will start on port 9090"
	@echo "Service registry will register at http://localhost:8500"
	@./bin/fr0g-ai-aip

# gRPC Reflection and Testing Support
test-grpc-reflection:
	@echo "🔍 Testing gRPC reflection..."
	@if ! command -v grpcurl >/dev/null 2>&1; then \
		echo "❌ grpcurl not found. Install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"; \
		exit 1; \
	fi
	@echo "✅ grpcurl found"
	@echo "Testing reflection endpoint..."
	@if grpcurl -plaintext localhost:9090 list >/dev/null 2>&1; then \
		echo "✅ gRPC reflection is working"; \
		grpcurl -plaintext localhost:9090 list; \
	else \
		echo "❌ gRPC reflection not available"; \
		echo "💡 Start service with reflection: make run-test"; \
	fi

# Build and run with reflection enabled for testing
build-test:
	@echo "🔨 Building AIP service with gRPC reflection enabled for testing..."
	@echo "⚠️  WARNING: gRPC reflection is ENABLED - DO NOT use in production!"
	@export GRPC_ENABLE_REFLECTION=true && \
	export ENVIRONMENT=development && \
	export AIP_HTTP_PORT=8080 && \
	export AIP_GRPC_PORT=9090 && \
	go build -o bin/fr0g-ai-aip ./cmd

run-test: build-test
	@echo "🚀 Starting AIP service with reflection enabled..."
	@echo "⚠️  WARNING: gRPC reflection is ENABLED - DO NOT use in production!"
	@export GRPC_ENABLE_REFLECTION=true && \
	export ENVIRONMENT=development && \
	export AIP_HTTP_PORT=8080 && \
	export AIP_GRPC_PORT=9090 && \
	./bin/fr0g-ai-aip

# Run tests with reflection enabled
test-with-reflection:
	@echo "🧪 Running AIP tests against running service..."
	@echo "💡 Make sure the service is running with reflection enabled:"
	@echo "   GRPC_ENABLE_REFLECTION=true ENVIRONMENT=development docker-compose up fr0g-ai-aip -d"
	@echo ""
	@chmod +x ../test_aip_service.sh
	@../test_aip_service.sh

# Validate production build (ensure reflection is disabled)
validate-production:
	@echo "🔒 Validating AIP production build security..."
	@if [ "$$GRPC_ENABLE_REFLECTION" = "true" ] && [ "$$ENVIRONMENT" = "production" ]; then \
		echo "❌ FATAL: gRPC reflection is enabled in production environment!"; \
		echo "   This is a security risk and deployment should be blocked."; \
		exit 1; \
	fi
	@echo "✅ AIP production validation passed - reflection properly disabled"

# Development environment setup for AIP
setup-dev:
	@echo "🛠️  Setting up AIP development environment..."
	@echo "export GRPC_ENABLE_REFLECTION=true" > .env.aip.development
	@echo "export ENVIRONMENT=development" >> .env.aip.development
	@echo "export LOG_LEVEL=debug" >> .env.aip.development
	@echo "export AIP_HTTP_PORT=8080" >> .env.aip.development
	@echo "export AIP_GRPC_PORT=9090" >> .env.aip.development
	@echo "✅ AIP development environment configured (.env.aip.development)"
	@echo "💡 Source with: source .env.aip.development"

# Production environment setup for AIP
setup-prod:
	@echo "🔒 Setting up AIP production environment..."
	@echo "export GRPC_ENABLE_REFLECTION=false" > .env.aip.production
	@echo "export ENVIRONMENT=production" >> .env.aip.production
	@echo "export LOG_LEVEL=info" >> .env.aip.production
	@echo "export AIP_HTTP_PORT=8080" >> .env.aip.production
	@echo "export AIP_GRPC_PORT=9090" >> .env.aip.production
	@echo "✅ AIP production environment configured (.env.aip.production)"
	@echo "💡 Source with: source .env.aip.production"

# Show help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Building:"
	@echo "  build              - Build the application (no external deps)"
	@echo "  build-with-grpc    - Build with full gRPC support"
	@echo "  build-test         - Build with gRPC reflection for testing"
	@echo "  clean              - Clean build artifacts"
	@echo ""
	@echo "Protocol Buffers:"
	@echo "  proto              - Force generate protobuf code"
	@echo "  proto-if-needed    - Generate protobuf code only if missing"
	@echo ""
	@echo "Testing:"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-coverage-detailed - Generate HTML coverage report"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-race          - Run tests with race detection"
	@echo "  test-bench         - Run benchmarks"
	@echo "  test-grpc-reflection - Test gRPC reflection specifically"
	@echo "  test-with-reflection - Run tests with reflection enabled"
	@echo ""
	@echo "Running:"
	@echo "  run-server         - Run HTTP REST API server"
	@echo "  run-grpc           - Run gRPC server"
	@echo "  run-both           - Run both HTTP and gRPC servers"
	@echo "  run-test           - Run with gRPC reflection enabled"
	@echo "  run-cli            - Show CLI help"
	@echo ""
	@echo "Security:"
	@echo "  validate-production - Ensure production build is secure"
	@echo "  setup-dev          - Create development environment file"
	@echo "  setup-prod         - Create production environment file"
	@echo ""
	@echo "Documentation:"
	@echo "  docs               - Generate all documentation"
	@echo "  docs-openapi       - Show OpenAPI documentation info"
	@echo "  docs-mcp           - Show MCP integration documentation info"
	@echo ""
	@echo "Development:"
	@echo "  deps               - Install/update dependencies"
	@echo "  install-proto-tools - Install protobuf generation tools"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code (requires golangci-lint)"
	@echo ""
	@echo "Quick Start Examples:"
	@echo "  # Generate sample data"
	@echo "  go run ./cmd generate-identities"
	@echo ""
	@echo "  # Create a tech community"
	@echo "  go run ./cmd generate-random-community -size 20 \\"
	@echo "    -name \"Tech Startup\" -type \"professional\" -location \"San Francisco\""
	@echo ""
	@echo "  # Start API server"
	@echo "  make run-server"
	@echo ""
	@echo "  # Test with reflection"
	@echo "  make run-test"
	@echo "  make test-grpc-reflection"
	@echo ""
	@echo "Help:"
	@echo "  help               - Show this help"
