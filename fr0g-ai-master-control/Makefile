.PHONY: build clean test run

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
