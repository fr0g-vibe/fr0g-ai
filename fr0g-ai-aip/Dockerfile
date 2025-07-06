# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./
COPY pkg/ ./pkg/

# Copy fr0g.ai-aip source
COPY fr0g.ai-aip/ ./fr0g.ai-aip/

# Download dependencies
RUN go mod download

# Build the application
WORKDIR /app/fr0g.ai-aip
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fr0g.ai-aip ./cmd/fr0g.ai-aip

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh fr0g

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/fr0g.ai-aip/fr0g.ai-aip .

# Create data directory
RUN mkdir -p /app/data /app/logs && chown -R fr0g:fr0g /app

# Switch to non-root user
USER fr0g

# Expose ports
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Default command - start both HTTP and gRPC servers
CMD ["./fr0g.ai-aip", "--server", "--grpc"]
