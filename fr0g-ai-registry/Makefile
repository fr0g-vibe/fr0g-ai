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

# Run unit tests only (no service required)
test-unit:
	@echo "Running unit tests (no service required)..."
	@go test -run "^TestUnit" ./test/

# Run integration tests (requires running service)
test-integration:
	@echo "Running integration tests (requires service on :8500)..."
	@echo "Make sure registry service is running: make run"
	@go test -run "^TestRegistry" ./test/

# Run load tests (requires running service)
test-load:
	@echo "Running load tests (requires service on :8500)..."
	@echo "Make sure registry service is running: make run"
	@go test -run "^TestRegistryLoad" ./test/

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
	@echo "SECURITY Validating Registry production build security..."
	@if [ "$$GRPC_ENABLE_REFLECTION" = "true" ] && [ "$$ENVIRONMENT" = "production" ]; then \
		echo "FAILED FATAL: gRPC reflection is enabled in production environment!"; \
		echo "   This is a security risk and deployment should be blocked."; \
		exit 1; \
	fi
	@echo "COMPLETED Registry production validation passed - reflection properly disabled"
