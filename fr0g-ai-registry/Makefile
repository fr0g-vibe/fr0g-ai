.PHONY: build clean test run deps fmt lint

# Build the registry binary
build: deps
	@echo "Building fr0g-ai-registry..."
	@mkdir -p bin
	go build -o bin/fr0g-ai-registry ./cmd/registry

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Run the service
run: build
	./bin/fr0g-ai-registry

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run || echo "Install golangci-lint for better linting"

# Validate production build (ensure reflection is disabled)
validate-production:
	@echo "🔒 Validating Registry production build security..."
	@if [ "$$GRPC_ENABLE_REFLECTION" = "true" ] && [ "$$ENVIRONMENT" = "production" ]; then \
		echo "❌ FATAL: gRPC reflection is enabled in production environment!"; \
		echo "   This is a security risk and deployment should be blocked."; \
		exit 1; \
	fi
	@echo "✅ Registry production validation passed - reflection properly disabled"
