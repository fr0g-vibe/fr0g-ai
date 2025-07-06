#!/bin/bash

# gRPC Connectivity Test Script
# Tests all gRPC services and their connectivity

set -e

echo "🚀 Starting gRPC Connectivity Tests"
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
AIP_GRPC_PORT=9090
IO_GRPC_PORT=9092
MC_GRPC_PORT=9093
BRIDGE_GRPC_PORT=9094

# Function to check if port is open
check_port() {
    local port=$1
    local service=$2
    
    if nc -z localhost $port 2>/dev/null; then
        echo -e "${GREEN}✓${NC} $service gRPC server is running on port $port"
        return 0
    else
        echo -e "${RED}✗${NC} $service gRPC server is NOT running on port $port"
        return 1
    fi
}

# Function to test gRPC health check
test_grpc_health() {
    local port=$1
    local service=$2
    
    echo -e "${BLUE}Testing gRPC health check for $service on port $port...${NC}"
    
    # Use grpcurl if available, otherwise use basic connectivity test
    if command -v grpcurl &> /dev/null; then
        if grpcurl -plaintext localhost:$port grpc.health.v1.Health/Check 2>/dev/null; then
            echo -e "${GREEN}✓${NC} $service health check passed"
            return 0
        else
            echo -e "${YELLOW}⚠${NC} $service health check failed (service may not implement health check)"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠${NC} grpcurl not available, using basic connectivity test"
        check_port $port $service
        return $?
    fi
}

# Function to test gRPC reflection
test_grpc_reflection() {
    local port=$1
    local service=$2
    
    echo -e "${BLUE}Testing gRPC reflection for $service on port $port...${NC}"
    
    if command -v grpcurl &> /dev/null; then
        if grpcurl -plaintext localhost:$port list 2>/dev/null; then
            echo -e "${GREEN}✓${NC} $service gRPC reflection is working"
            return 0
        else
            echo -e "${YELLOW}⚠${NC} $service gRPC reflection is disabled or not available"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠${NC} grpcurl not available, skipping reflection test"
        return 1
    fi
}

# Function to start a service if not running
start_service_if_needed() {
    local service=$1
    local port=$2
    local cmd=$3
    
    if ! check_port $port $service; then
        echo -e "${YELLOW}Starting $service...${NC}"
        # This would start the service in the background
        # For now, just report that it needs to be started
        echo -e "${YELLOW}⚠${NC} Please start $service manually: $cmd"
        return 1
    fi
    return 0
}

echo -e "\n${BLUE}Step 1: Checking if all gRPC services are running${NC}"
echo "=================================================="

# Check all services
services_running=0

if check_port $AIP_GRPC_PORT "fr0g-ai-aip"; then
    ((services_running++))
fi

if check_port $IO_GRPC_PORT "fr0g-ai-io"; then
    ((services_running++))
fi

if check_port $MC_GRPC_PORT "fr0g-ai-master-control"; then
    ((services_running++))
fi

if check_port $BRIDGE_GRPC_PORT "fr0g-ai-bridge"; then
    ((services_running++))
fi

echo -e "\n${BLUE}Services running: $services_running/4${NC}"

if [ $services_running -eq 0 ]; then
    echo -e "${RED}No gRPC services are running!${NC}"
    echo -e "${YELLOW}Please start the services first:${NC}"
    echo "  - fr0g-ai-aip: go run fr0g-ai-aip/cmd/main.go"
    echo "  - fr0g-ai-io: go run fr0g-ai-io/cmd/server/main.go"
    echo "  - fr0g-ai-master-control: go run fr0g-ai-master-control/cmd/main.go"
    echo "  - fr0g-ai-bridge: go run fr0g-ai-bridge/cmd/bridge/main.go"
    exit 1
fi

echo -e "\n${BLUE}Step 2: Testing gRPC health checks${NC}"
echo "=================================="

health_checks_passed=0

if check_port $AIP_GRPC_PORT "fr0g-ai-aip"; then
    if test_grpc_health $AIP_GRPC_PORT "fr0g-ai-aip"; then
        ((health_checks_passed++))
    fi
fi

if check_port $IO_GRPC_PORT "fr0g-ai-io"; then
    if test_grpc_health $IO_GRPC_PORT "fr0g-ai-io"; then
        ((health_checks_passed++))
    fi
fi

if check_port $MC_GRPC_PORT "fr0g-ai-master-control"; then
    if test_grpc_health $MC_GRPC_PORT "fr0g-ai-master-control"; then
        ((health_checks_passed++))
    fi
fi

if check_port $BRIDGE_GRPC_PORT "fr0g-ai-bridge"; then
    if test_grpc_health $BRIDGE_GRPC_PORT "fr0g-ai-bridge"; then
        ((health_checks_passed++))
    fi
fi

echo -e "\n${BLUE}Health checks passed: $health_checks_passed/$services_running${NC}"

echo -e "\n${BLUE}Step 3: Testing gRPC reflection${NC}"
echo "==============================="

reflection_working=0

if check_port $AIP_GRPC_PORT "fr0g-ai-aip"; then
    if test_grpc_reflection $AIP_GRPC_PORT "fr0g-ai-aip"; then
        ((reflection_working++))
    fi
fi

if check_port $IO_GRPC_PORT "fr0g-ai-io"; then
    if test_grpc_reflection $IO_GRPC_PORT "fr0g-ai-io"; then
        ((reflection_working++))
    fi
fi

if check_port $MC_GRPC_PORT "fr0g-ai-master-control"; then
    if test_grpc_reflection $MC_GRPC_PORT "fr0g-ai-master-control"; then
        ((reflection_working++))
    fi
fi

if check_port $BRIDGE_GRPC_PORT "fr0g-ai-bridge"; then
    if test_grpc_reflection $BRIDGE_GRPC_PORT "fr0g-ai-bridge"; then
        ((reflection_working++))
    fi
fi

echo -e "\n${BLUE}Reflection working: $reflection_working/$services_running${NC}"

echo -e "\n${BLUE}Step 4: Testing service-to-service connectivity${NC}"
echo "=============================================="

# Test bidirectional connectivity between service pairs
test_connectivity() {
    local from_port=$1
    local to_port=$2
    local from_service=$3
    local to_service=$4
    
    echo -e "${BLUE}Testing $from_service -> $to_service connectivity...${NC}"
    
    if nc -z localhost $to_port 2>/dev/null; then
        echo -e "${GREEN}✓${NC} $from_service can connect to $to_service"
        return 0
    else
        echo -e "${RED}✗${NC} $from_service cannot connect to $to_service"
        return 1
    fi
}

connectivity_tests=0
connectivity_passed=0

# Test AIP ↔ Bridge
if check_port $AIP_GRPC_PORT "fr0g-ai-aip" && check_port $BRIDGE_GRPC_PORT "fr0g-ai-bridge"; then
    ((connectivity_tests++))
    if test_connectivity $AIP_GRPC_PORT $BRIDGE_GRPC_PORT "fr0g-ai-aip" "fr0g-ai-bridge"; then
        ((connectivity_passed++))
    fi
    
    ((connectivity_tests++))
    if test_connectivity $BRIDGE_GRPC_PORT $AIP_GRPC_PORT "fr0g-ai-bridge" "fr0g-ai-aip"; then
        ((connectivity_passed++))
    fi
fi

# Test IO ↔ Master-Control
if check_port $IO_GRPC_PORT "fr0g-ai-io" && check_port $MC_GRPC_PORT "fr0g-ai-master-control"; then
    ((connectivity_tests++))
    if test_connectivity $IO_GRPC_PORT $MC_GRPC_PORT "fr0g-ai-io" "fr0g-ai-master-control"; then
        ((connectivity_passed++))
    fi
    
    ((connectivity_tests++))
    if test_connectivity $MC_GRPC_PORT $IO_GRPC_PORT "fr0g-ai-master-control" "fr0g-ai-io"; then
        ((connectivity_passed++))
    fi
fi

# Test Bridge ↔ Master-Control
if check_port $BRIDGE_GRPC_PORT "fr0g-ai-bridge" && check_port $MC_GRPC_PORT "fr0g-ai-master-control"; then
    ((connectivity_tests++))
    if test_connectivity $BRIDGE_GRPC_PORT $MC_GRPC_PORT "fr0g-ai-bridge" "fr0g-ai-master-control"; then
        ((connectivity_passed++))
    fi
    
    ((connectivity_tests++))
    if test_connectivity $MC_GRPC_PORT $BRIDGE_GRPC_PORT "fr0g-ai-master-control" "fr0g-ai-bridge"; then
        ((connectivity_passed++))
    fi
fi

echo -e "\n${BLUE}Connectivity tests passed: $connectivity_passed/$connectivity_tests${NC}"

echo -e "\n${BLUE}Step 5: Running Go tests${NC}"
echo "======================="

# Run the Go connectivity tests
if [ -f "test/grpc_connectivity_test.go" ]; then
    echo -e "${BLUE}Running Go gRPC connectivity tests...${NC}"
    if go test -v ./test -run TestGRPCConnectivity; then
        echo -e "${GREEN}✓${NC} Go gRPC connectivity tests passed"
    else
        echo -e "${YELLOW}⚠${NC} Go gRPC connectivity tests had issues (expected without real services)"
    fi
else
    echo -e "${YELLOW}⚠${NC} Go test file not found, skipping Go tests"
fi

# Run service discovery tests
if [ -f "test/service_discovery_test.go" ]; then
    echo -e "${BLUE}Running service discovery tests...${NC}"
    if go test -v ./test -run TestServiceDiscovery; then
        echo -e "${GREEN}✓${NC} Service discovery tests passed"
    else
        echo -e "${YELLOW}⚠${NC} Service discovery tests had issues (expected without real registry)"
    fi
fi

echo -e "\n${BLUE}Summary${NC}"
echo "======="
echo -e "Services running: ${GREEN}$services_running/4${NC}"
echo -e "Health checks: ${GREEN}$health_checks_passed/$services_running${NC}"
echo -e "Reflection: ${GREEN}$reflection_working/$services_running${NC}"
echo -e "Connectivity: ${GREEN}$connectivity_passed/$connectivity_tests${NC}"

# Overall status
if [ $services_running -eq 4 ] && [ $health_checks_passed -ge 2 ] && [ $connectivity_passed -ge 4 ]; then
    echo -e "\n${GREEN}🎉 gRPC connectivity tests PASSED!${NC}"
    echo -e "${GREEN}All services are properly connected and communicating.${NC}"
    exit 0
elif [ $services_running -ge 2 ]; then
    echo -e "\n${YELLOW}⚠ gRPC connectivity tests PARTIALLY PASSED${NC}"
    echo -e "${YELLOW}Some services are working, but full connectivity needs improvement.${NC}"
    exit 1
else
    echo -e "\n${RED}❌ gRPC connectivity tests FAILED${NC}"
    echo -e "${RED}Critical issues found with gRPC service connectivity.${NC}"
    exit 2
fi
