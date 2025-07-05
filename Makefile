.PHONY: help build clean test run deps lint fmt health

# Default target
help:
	@echo "🐸 fr0g.ai MCP - Model Context Protocol Server"
	@echo ""
	@echo "Available targets:"
	@echo "  build              - Build the MCP server"
	@echo "  clean              - Clean build artifacts"
	@echo "  test               - Run tests"
	@echo "  run                - Run the MCP server"
	@echo "  deps               - Install dependencies"
	@echo "  lint               - Run code linters"
	@echo "  fmt                - Format code"
	@echo "  health             - Check MCP service health"

# Build the MCP server
build:
	@echo "🔨 Building fr0g.ai MCP server..."
	@mkdir -p bin
	go build -o bin/fr0g-ai-mcp ./cmd/mcp


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
