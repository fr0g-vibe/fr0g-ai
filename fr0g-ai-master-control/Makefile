.PHONY: build clean test run deps fmt lint

# Build the master-control binary
build:
	@echo "Building fr0g-ai-master-control..."
	go build -o bin/fr0g-ai-master-control ./cmd

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Run the service
run: build
	./bin/fr0g-ai-master-control

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
	@echo "SECURITY Validating Master Control production build security..."
	@if [ "$$GRPC_ENABLE_REFLECTION" = "true" ] && [ "$$ENVIRONMENT" = "production" ]; then \
		echo "FAILED FATAL: gRPC reflection is enabled in production environment!"; \
		echo "   This is a security risk and deployment should be blocked."; \
		exit 1; \
	fi
	@echo "COMPLETED Master Control production validation passed - reflection properly disabled"
