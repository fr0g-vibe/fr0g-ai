.PHONY: build clean test run deps fmt lint

# Build the registry binary
build:
	@echo "Building fr0g-ai-registry..."
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
