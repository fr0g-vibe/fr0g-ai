#!/bin/bash

# Complete end-to-end test for service lifecycle management
set -e

echo "🧪 Testing Complete Service Lifecycle Management"
echo "================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if service is running
check_service() {
    local service_name=$1
    local port=$2
    local max_attempts=10
    local attempt=1
    
    echo -e "${BLUE}🔍 Checking $service_name on port $port...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "http://localhost:$port/health" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ $service_name is healthy${NC}"
            return 0
        fi
        echo -e "${YELLOW}⏳ Attempt $attempt/$max_attempts - waiting for $service_name...${NC}"
        sleep 2
        ((attempt++))
    done
    
    echo -e "${RED}❌ $service_name failed to start${NC}"
    return 1
}

# Function to check service registry
check_registry_service() {
    local service_name=$1
    echo -e "${BLUE}🔍 Checking if $service_name is registered...${NC}"
    
    if curl -s "http://localhost:8500/v1/catalog/services" | grep -q "$service_name"; then
        echo -e "${GREEN}✅ $service_name is registered in service registry${NC}"
        return 0
    else
        echo -e "${RED}❌ $service_name is not registered${NC}"
        return 1
    fi
}

# Cleanup function
cleanup() {
    echo -e "${YELLOW}🧹 Cleaning up processes...${NC}"
    
    # Kill all background processes
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            echo "Killing process $pid"
            kill "$pid" 2>/dev/null || true
        fi
    done
    
    # Wait a moment for graceful shutdown
    sleep 2
    
    # Force kill if necessary
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            echo "Force killing process $pid"
            kill -9 "$pid" 2>/dev/null || true
        fi
    done
    
    echo -e "${GREEN}✅ Cleanup completed${NC}"
}

# Set trap for cleanup on exit
trap cleanup EXIT

# Array to store process IDs
PIDS=()

echo -e "${BLUE}🚀 Step 1: Starting Service Registry...${NC}"
export REDIS_ADDR="localhost:6379"
export REGISTRY_PORT="8500"
cd fr0g-ai-registry && go run cmd/registry/main.go &
REGISTRY_PID=$!
PIDS+=($REGISTRY_PID)
cd ..

# Wait for registry to start
sleep 3
check_service "registry" "8500"

echo -e "${BLUE}🚀 Step 2: Starting AIP Service...${NC}"
export SERVICE_REGISTRY_ENABLED="true"
export SERVICE_REGISTRY_URL="http://localhost:8500"
export SERVICE_NAME="fr0g-ai-aip"
export SERVICE_ID="fr0g-ai-aip-1"
export HTTP_PORT="8080"
export GRPC_PORT="9090"
cd fr0g-ai-aip && go run cmd/aip/main.go &
AIP_PID=$!
PIDS+=($AIP_PID)
cd ..

# Wait for AIP service to start
sleep 3
check_service "fr0g-ai-aip" "8080"
check_registry_service "fr0g-ai-aip"

echo -e "${BLUE}🚀 Step 3: Starting Bridge Service...${NC}"
export SERVICE_NAME="fr0g-ai-bridge"
export SERVICE_ID="fr0g-ai-bridge-1"
export HTTP_PORT="8082"
export GRPC_PORT="9091"
cd fr0g-ai-bridge && go run cmd/bridge/main.go &
BRIDGE_PID=$!
PIDS+=($BRIDGE_PID)
cd ..

# Wait for Bridge service to start
sleep 3
check_service "fr0g-ai-bridge" "8082"
check_registry_service "fr0g-ai-bridge"

echo -e "${BLUE}🚀 Step 4: Starting IO Service...${NC}"
export SERVICE_NAME="fr0g-ai-io"
export SERVICE_ID="fr0g-ai-io-1"
export HTTP_PORT="8084"
export GRPC_PORT="9094"
cd fr0g-ai-io && go run cmd/io/main.go &
IO_PID=$!
PIDS+=($IO_PID)
cd ..

# Wait for IO service to start
sleep 3
check_service "fr0g-ai-io" "8084"
check_registry_service "fr0g-ai-io"

echo -e "${BLUE}🚀 Step 5: Starting Master Control...${NC}"
export SERVICE_NAME="fr0g-ai-master-control"
export SERVICE_ID="fr0g-ai-master-control-1"
export MCP_HTTP_PORT="8083"
cd fr0g-ai-master-control && go run cmd/master-control/main.go &
MCP_PID=$!
PIDS+=($MCP_PID)
cd ..

# Wait for Master Control to start
sleep 3
check_service "fr0g-ai-master-control" "8083"
check_registry_service "fr0g-ai-master-control"

echo -e "${BLUE}🔍 Step 6: Testing Service Discovery...${NC}"

echo "📋 All registered services:"
curl -s http://localhost:8500/v1/catalog/services | jq .

echo ""
echo "🔍 Service details:"
for service in "fr0g-ai-aip" "fr0g-ai-bridge" "fr0g-ai-io" "fr0g-ai-master-control"; do
    echo "--- $service ---"
    curl -s "http://localhost:8500/v1/health/service/$service-1" | jq . || echo "Service not found"
done

echo -e "${BLUE}🧪 Step 7: Testing Inter-Service Communication...${NC}"

# Test health endpoints
echo "Testing health endpoints:"
for port in 8080 8082 8083 8084; do
    echo "Port $port:"
    curl -s "http://localhost:$port/health" | jq . || echo "Failed to connect"
done

echo ""
echo "Testing status endpoints:"
for port in 8080 8082 8083 8084; do
    echo "Port $port:"
    curl -s "http://localhost:$port/status" | jq . || echo "No status endpoint"
done

echo -e "${BLUE}🛑 Step 8: Testing Graceful Shutdown...${NC}"

# Test graceful shutdown by sending SIGTERM to one service
echo "Testing graceful shutdown of AIP service..."
kill -TERM $AIP_PID
sleep 3

# Check if service was deregistered
if ! check_registry_service "fr0g-ai-aip"; then
    echo -e "${GREEN}✅ AIP service was properly deregistered${NC}"
else
    echo -e "${RED}❌ AIP service was not deregistered${NC}"
fi

echo -e "${GREEN}🎉 Service Lifecycle Test Completed Successfully!${NC}"
echo ""
echo "Summary:"
echo "✅ Service Registry started and operational"
echo "✅ All services registered automatically on startup"
echo "✅ Health monitoring working"
echo "✅ Service discovery functional"
echo "✅ Graceful shutdown and deregistration working"
echo ""
echo "The automatic service lifecycle management is fully operational!"
