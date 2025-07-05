#!/bin/bash

# End-to-end integration test for all fr0g-ai services
set -e

echo "=== End-to-End Integration Test ==="

# Wait for services to be ready
echo "Waiting for services to start..."
sleep 10

# Test AIP service
echo "Testing AIP service..."
curl -f http://localhost:8080/health || {
    echo "ERROR: AIP service health check failed"
    exit 1
}

# Test Bridge service
echo "Testing Bridge service..."
curl -f http://localhost:8082/health || {
    echo "ERROR: Bridge service health check failed"
    exit 1
}

# Test Master Control service
echo "Testing Master Control service..."
curl -f http://localhost:8081/health || {
    echo "ERROR: Master Control service health check failed"
    exit 1
}

# Test I/O service
echo "Testing I/O service..."
curl -f http://localhost:8083/health || {
    echo "ERROR: I/O service health check failed"
    exit 1
}

# Test service registry
echo "Testing Service Registry..."
curl -f http://localhost:8500/health || {
    echo "ERROR: Service Registry health check failed"
    exit 1
}

# Test service discovery integration
echo "Testing service discovery..."
SERVICES=$(curl -s http://localhost:8500/v1/catalog/services | jq -r 'keys[]')
echo "Discovered services: $SERVICES"

# Verify all expected services are registered
for service in "fr0g-ai-aip" "fr0g-ai-bridge" "fr0g-ai-master-control" "fr0g-ai-io"; do
    if echo "$SERVICES" | grep -q "$service"; then
        echo "✓ $service is registered"
    else
        echo "✗ $service is NOT registered"
    fi
done

echo "=== End-to-End Integration Test COMPLETED ==="
