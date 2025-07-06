# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./
COPY pkg/ ./pkg/

# Copy fr0g-ai-bridge source
COPY fr0g-ai-bridge/ ./fr0g-ai-bridge/

# Download dependencies
RUN go mod download

# Build the application
WORKDIR /app/fr0g-ai-bridge
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fr0g-ai-bridge .

# Runtime stage
FROM alpine:latest

# Install ca-certificates, curl, and netcat for health checks
RUN apk --no-cache add ca-certificates curl netcat-openbsd

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/fr0g-ai-bridge/fr0g-ai-bridge .

# Create config and logs directories
RUN mkdir -p /app/config /app/logs && chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8082 9091

# Enhanced health check with HTTP and gRPC support
HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=5 \
  CMD curl -f http://localhost:8082/health && nc -z localhost 9091 || exit 1

# Run the application
CMD ["./fr0g-ai-bridge"]
