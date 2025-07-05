#!/bin/bash

# Integration test for service registry functionality
set -e

echo "=== Service Registry Integration Test ==="

# Test service registry health
echo "Testing service registry health..."
curl -f http://localhost:8500/health || {
    echo "ERROR: Service registry health check failed"
    exit 1
}

# Test service registration
echo "Testing service registration..."
curl -X POST http://localhost:8500/v1/agent/service/register \
    -H "Content-Type: application/json" \
    -d '{
        "ID": "test-service",
        "Name": "test-service",
        "Address": "localhost",
        "Port": 9999,
        "Check": {
            "HTTP": "http://localhost:9999/health",
            "Interval": "10s"
        }
    }' || {
    echo "ERROR: Service registration failed"
    exit 1
}

# Test service discovery
echo "Testing service discovery..."
curl -f http://localhost:8500/v1/catalog/service/test-service || {
    echo "ERROR: Service discovery failed"
    exit 1
}

# Cleanup test service
echo "Cleaning up test service..."
curl -X PUT http://localhost:8500/v1/agent/service/deregister/test-service

echo "=== Service Registry Integration Test PASSED ==="
