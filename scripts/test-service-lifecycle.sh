#!/bin/bash

# Test script for service lifecycle management
set -e

echo "🧪 Testing Service Lifecycle Management"
echo "======================================="

# Start registry in background
echo "🚀 Starting service registry..."
export REDIS_ADDR="localhost:6379"
export REGISTRY_PORT="8500"
cd fr0g-ai-registry && go run cmd/registry/main.go &
REGISTRY_PID=$!
sleep 3

# Test service registration
echo "📋 Testing service registration..."
curl -X PUT http://localhost:8500/v1/agent/service/register \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-service-1",
    "name": "test-service",
    "address": "localhost",
    "port": 8080,
    "tags": ["test", "ai"],
    "meta": {"version": "1.0.0"},
    "check": {
      "http": "http://localhost:8080/health",
      "interval": "30s",
      "timeout": "10s"
    }
  }'

echo ""
echo "🔍 Checking registered services..."
curl -s http://localhost:8500/v1/catalog/services | jq .

echo ""
echo "❤️  Checking service health..."
curl -s http://localhost:8500/v1/health/service/test-service-1 | jq .

echo ""
echo "🗑️  Testing service deregistration..."
curl -X PUT http://localhost:8500/v1/agent/service/deregister/test-service-1

echo ""
echo "🔍 Verifying service removed..."
curl -s http://localhost:8500/v1/catalog/services | jq .

# Cleanup
echo ""
echo "🧹 Cleaning up..."
kill $REGISTRY_PID 2>/dev/null || true

echo "✅ Service lifecycle test completed successfully!"
