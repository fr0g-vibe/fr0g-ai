# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./
COPY pkg/ ./pkg/

# Copy fr0g-ai-registry source
COPY fr0g-ai-registry/ ./fr0g-ai-registry/

# Download dependencies
RUN go mod download

# Build the registry service from cmd/registry
WORKDIR /app/fr0g-ai-registry
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o registry-server ./cmd/registry

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh registry

WORKDIR /app

# Copy the binary
COPY --from=builder /app/fr0g-ai-registry/registry-server .

# Create logs directory
RUN mkdir -p /app/logs && chown -R registry:registry /app

# Switch to non-root user
USER registry

# Expose port
EXPOSE 8500

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8500/health || exit 1

# Run the registry server
CMD ["./registry-server"]
