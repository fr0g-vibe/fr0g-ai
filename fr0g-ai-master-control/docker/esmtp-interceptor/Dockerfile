FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o esmtp-interceptor ./cmd/esmtp-interceptor

FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary and config
COPY --from=builder /app/esmtp-interceptor .
COPY --from=builder /app/configs/esmtp.yaml ./configs/

# Create directories for logs and data
RUN mkdir -p /var/log/fr0g-ai /var/lib/fr0g-ai

# Expose ESMTP ports
EXPOSE 2525 2465

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD nc -z localhost 2525 || exit 1

# Run the interceptor
CMD ["./esmtp-interceptor", "-config", "configs/esmtp.yaml"]
