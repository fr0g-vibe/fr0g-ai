#!/bin/bash

# fr0g.ai Deployment Test Script
# Tests the deployment and basic functionality of fr0g.ai services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
AIP_URL="http://localhost:8080"
BRIDGE_URL="http://localhost:8081"
TIMEOUT=30

echo -e "${GREEN}fr0g.ai Deployment Test${NC}"
echo "=================================="

# Function to check if service is responding
check_service() {
    local service_name=$1
    local url=$2
    local endpoint=$3
    
    echo -n "Testing $service_name ($url$endpoint)... "
    
    if curl -s -f --max-time $TIMEOUT "$url$endpoint" > /dev/null; then
        echo -e "${GREEN}PASS${NC}"
        return 0
    else
        echo -e "${RED}FAIL${NC}"
        return 1
    fi
}

# Function to check service health
check_health() {
    local service_name=$1
    local url=$2
    
    echo -n "Checking $service_name health... "
    
    response=$(curl -s --max-time $TIMEOUT "$url/health" || echo "")
    
    if echo "$response" | grep -q "healthy\|ok"; then
        echo -e "${GREEN}HEALTHY${NC}"
        return 0
    else
        echo -e "${RED}UNHEALTHY${NC}"
        echo "Response: $response"
        return 1
    fi
}

# Function to test API endpoint
test_api() {
    local service_name=$1
    local url=$2
    local endpoint=$3
    local method=${4:-GET}
    local data=${5:-""}
    
    echo -n "Testing $service_name API ($method $endpoint)... "
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        response=$(curl -s --max-time $TIMEOUT -X POST \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$url$endpoint" || echo "")
    else
        response=$(curl -s --max-time $TIMEOUT "$url$endpoint" || echo "")
    fi
    
    if [ -n "$response" ]; then
        echo -e "${GREEN}OK${NC}"
        return 0
    else
        echo -e "${RED}FAILED${NC}"
        return 1
    fi
}

# Wait for services to be ready
echo -e "\n${YELLOW}Waiting for services to start...${NC}"
sleep 10

# Test basic connectivity
echo -e "\n${YELLOW}Testing Basic Connectivity${NC}"
echo "--------------------------------"

check_service "AIP Service" "$AIP_URL" "/" || exit 1
check_service "Bridge Service" "$BRIDGE_URL" "/" || exit 1

# Test health endpoints
echo -e "\n${YELLOW}Testing Health Endpoints${NC}"
echo "-----------------------------"

check_health "AIP Service" "$AIP_URL" || exit 1
check_health "Bridge Service" "$BRIDGE_URL" || exit 1

# Test metrics endpoints
echo -e "\n${YELLOW}Testing Metrics Endpoints${NC}"
echo "------------------------------"

check_service "AIP Metrics" "$AIP_URL" "/metrics" || echo -e "${YELLOW}Metrics may not be enabled${NC}"
check_service "Bridge Metrics" "$BRIDGE_URL" "/metrics" || echo -e "${YELLOW}Metrics may not be enabled${NC}"

# Test API endpoints (if available)
echo -e "\n${YELLOW}Testing API Endpoints${NC}"
echo "--------------------------"

# Test basic API endpoints (adjust based on actual API)
test_api "AIP Service" "$AIP_URL" "/api/status" || echo -e "${YELLOW}API endpoint may not be implemented${NC}"
test_api "Bridge Service" "$BRIDGE_URL" "/api/status" || echo -e "${YELLOW}API endpoint may not be implemented${NC}"

# Test inter-service communication
echo -e "\n${YELLOW}Testing Inter-Service Communication${NC}"
echo "---------------------------------------"

# This would test if bridge can communicate with AIP
# Implementation depends on actual API design
echo -e "${YELLOW}Inter-service communication tests not implemented yet${NC}"

# Performance test
echo -e "\n${YELLOW}Basic Performance Test${NC}"
echo "---------------------------"

echo -n "Testing AIP service response time... "
start_time=$(date +%s%N)
curl -s --max-time $TIMEOUT "$AIP_URL/health" > /dev/null
end_time=$(date +%s%N)
duration=$(( (end_time - start_time) / 1000000 ))

if [ $duration -lt 1000 ]; then
    echo -e "${GREEN}${duration}ms (Excellent)${NC}"
elif [ $duration -lt 5000 ]; then
    echo -e "${YELLOW}${duration}ms (Good)${NC}"
else
    echo -e "${RED}${duration}ms (Slow)${NC}"
fi

# Docker health check
echo -e "\n${YELLOW}Testing Docker Health${NC}"
echo "--------------------------"

if command -v docker-compose &> /dev/null; then
    echo "Checking Docker Compose services..."
    docker-compose ps
    
    # Check if all services are healthy
    unhealthy=$(docker-compose ps --services --filter "health=unhealthy" | wc -l)
    if [ "$unhealthy" -eq 0 ]; then
        echo -e "${GREEN}All Docker services are healthy${NC}"
    else
        echo -e "${RED}Some Docker services are unhealthy${NC}"
        docker-compose ps --filter "health=unhealthy"
    fi
else
    echo -e "${YELLOW}Docker Compose not available${NC}"
fi

# Summary
echo -e "\n${GREEN}Deployment Test Complete${NC}"
echo "============================="

# Check if all critical tests passed
if check_service "AIP Service" "$AIP_URL" "/" && \
   check_service "Bridge Service" "$BRIDGE_URL" "/" && \
   check_health "AIP Service" "$AIP_URL" && \
   check_health "Bridge Service" "$BRIDGE_URL"; then
    echo -e "${GREEN}All critical tests passed!${NC}"
    echo -e "${GREEN}fr0g.ai is ready for use${NC}"
    exit 0
else
    echo -e "${RED}Some critical tests failed${NC}"
    echo -e "${RED}Please check the logs and configuration${NC}"
    exit 1
fi
