# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files from root
COPY go.mod go.sum ./
COPY pkg/ ./pkg/

# Copy fr0g-ai-io source
COPY fr0g-ai-io/ ./fr0g-ai-io/

# Download dependencies
RUN go mod download

# Build the application
WORKDIR /app/fr0g-ai-io
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fr0g-ai-io ./cmd

# Final stage
FROM alpine:latest

# Install ca-certificates, curl, and netcat for health checks
RUN apk --no-cache add ca-certificates curl netcat-openbsd

# Create non-root user
RUN adduser -D -s /bin/sh fr0g

# Set working directory
WORKDIR /app

# Copy the binary
COPY --from=builder /app/fr0g-ai-io/fr0g-ai-io .

# Create data and logs directories
RUN mkdir -p /app/data /app/logs && chown -R fr0g:fr0g /app

# Switch to non-root user
USER fr0g

# Expose ports  
EXPOSE 8083 9092

# Enhanced health check with HTTP and gRPC support
HEALTHCHECK --interval=15s --timeout=5s --start-period=45s --retries=5 \
  CMD curl -f http://localhost:8083/health && nc -z localhost 9092 || exit 1

# Run the application
CMD ["./fr0g-ai-io"]
