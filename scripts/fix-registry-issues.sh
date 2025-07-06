#!/bin/bash

# Service Registry Fix Script
# Addresses JSON decode EOF and registration issues

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

REGISTRY_URL="http://localhost:8500"

echo -e "${BLUE}fr0g.ai Service Registry Fix${NC}"
echo "============================"

# Function to test registry endpoints
test_registry_endpoints() {
    echo -e "\n${BLUE}Testing registry endpoints...${NC}"
    
    # Test health endpoint
    echo -n "Health endpoint... "
    if curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}FAILED${NC}"
        return 1
    fi
    
    # Test catalog endpoint
    echo -n "Catalog endpoint... "
    local catalog_response=$(curl -s "$REGISTRY_URL/v1/catalog/services" 2>/dev/null)
    if [ -n "$catalog_response" ]; then
        echo -e "${GREEN}OK${NC}"
        echo "  Current services: $catalog_response"
    else
        echo -e "${RED}FAILED${NC}"
    fi
    
    # Test registration endpoint with minimal payload
    echo -n "Registration endpoint (minimal test)... "
    local test_response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"ID":"test-minimal","Name":"test"}' 2>/dev/null)
    
    local test_code="${test_response: -3}"
    if [ "$test_code" = "200" ] || [ "$test_code" = "201" ]; then
        echo -e "${GREEN}OK${NC}"
        # Clean up
        curl -s "$REGISTRY_URL/v1/agent/service/deregister/test-minimal" -X PUT >/dev/null 2>&1
    else
        echo -e "${RED}FAILED (HTTP $test_code)${NC}"
        echo "  Response: ${test_response%???}"
    fi
}

# Function to check registry logs for specific issues
check_registry_logs() {
    echo -e "\n${BLUE}Checking registry logs for issues...${NC}"
    
    # Get recent logs
    local logs=$(docker logs fr0g-ai-service-registry-1 --tail=20 2>&1)
    
    # Check for JSON decode errors
    local json_errors=$(echo "$logs" | grep -i "json.*error\|decode.*error\|EOF" || true)
    if [ -n "$json_errors" ]; then
        echo -e "${RED}JSON DECODE ERRORS:${NC}"
        echo "$json_errors"
    fi
    
    # Check for registration handler issues
    local reg_errors=$(echo "$logs" | grep -i "register.*error\|register.*failed" || true)
    if [ -n "$reg_errors" ]; then
        echo -e "${RED}REGISTRATION ERRORS:${NC}"
        echo "$reg_errors"
    fi
    
    # Check for successful registrations
    local reg_success=$(echo "$logs" | grep -i "register.*success\|service.*registered" || true)
    if [ -n "$reg_success" ]; then
        echo -e "${GREEN}SUCCESSFUL REGISTRATIONS:${NC}"
        echo "$reg_success"
    fi
    
    # Show recent registration attempts
    local reg_attempts=$(echo "$logs" | grep -i "register.*handler" || true)
    if [ -n "$reg_attempts" ]; then
        echo -e "${BLUE}RECENT REGISTRATION ATTEMPTS:${NC}"
        echo "$reg_attempts"
    fi
}

# Function to test different registration payloads
test_registration_payloads() {
    echo -e "\n${BLUE}Testing different registration payloads...${NC}"
    
    # Test 1: Minimal payload
    echo -n "Minimal payload... "
    local response1=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"ID":"test1","Name":"test1"}' 2>/dev/null)
    local code1="${response1: -3}"
    echo "HTTP $code1"
    
    # Test 2: Standard payload
    echo -n "Standard payload... "
    local response2=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"ID":"test2","Name":"test2","Port":8000,"Address":"localhost"}' 2>/dev/null)
    local code2="${response2: -3}"
    echo "HTTP $code2"
    
    # Test 3: Full payload with health check
    echo -n "Full payload... "
    local response3=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"ID":"test3","Name":"test3","Port":8000,"Address":"localhost","Check":{"HTTP":"http://localhost:8000/health","Interval":"10s"}}' 2>/dev/null)
    local code3="${response3: -3}"
    echo "HTTP $code3"
    
    # Clean up test services
    curl -s "$REGISTRY_URL/v1/agent/service/deregister/test1" -X PUT >/dev/null 2>&1
    curl -s "$REGISTRY_URL/v1/agent/service/deregister/test2" -X PUT >/dev/null 2>&1
    curl -s "$REGISTRY_URL/v1/agent/service/deregister/test3" -X PUT >/dev/null 2>&1
}

# Function to restart registry service
restart_registry() {
    echo -e "\n${BLUE}Restarting registry service...${NC}"
    
    docker-compose restart service-registry
    sleep 10
    
    # Wait for registry to be ready
    local retries=0
    while [ $retries -lt 30 ]; do
        if curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1; then
            echo -e "${GREEN}Registry restarted successfully${NC}"
            return 0
        fi
        retries=$((retries + 1))
        sleep 2
    done
    
    echo -e "${RED}Registry restart failed${NC}"
    return 1
}

# Function to check Redis connection
check_redis_connection() {
    echo -e "\n${BLUE}Checking Redis connection...${NC}"
    
    # Check if Redis container is running
    if docker ps --format "{{.Names}}" | grep -q "fr0g-ai-redis-1"; then
        echo -e "${GREEN}Redis container running${NC}"
        
        # Test Redis connectivity from registry container
        echo -n "Registry -> Redis connectivity... "
        if docker exec fr0g-ai-service-registry-1 sh -c "nc -z redis 6379" 2>/dev/null; then
            echo -e "${GREEN}OK${NC}"
        else
            echo -e "${RED}FAILED${NC}"
            echo "  Registry cannot connect to Redis"
        fi
    else
        echo -e "${RED}Redis container not running${NC}"
        return 1
    fi
}

main() {
    echo -e "${BLUE}Starting registry diagnostic and fix...${NC}\n"
    
    # Check if registry is accessible
    if ! curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1; then
        echo -e "${RED}Registry not accessible, attempting restart...${NC}"
        restart_registry
    fi
    
    # Run diagnostics
    check_redis_connection
    test_registry_endpoints
    check_registry_logs
    test_registration_payloads
    
    # If issues persist, restart registry
    local reg_response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" \
        -X POST \
        -H "Content-Type: application/json" \
        -d '{"ID":"final-test","Name":"final-test"}' 2>/dev/null)
    local reg_code="${reg_response: -3}"
    
    if [ "$reg_code" != "200" ] && [ "$reg_code" != "201" ]; then
        echo -e "\n${YELLOW}Registration still failing, restarting registry...${NC}"
        restart_registry
        
        # Test again after restart
        sleep 5
        local final_response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" \
            -X POST \
            -H "Content-Type: application/json" \
            -d '{"ID":"final-test-2","Name":"final-test-2"}' 2>/dev/null)
        local final_code="${final_response: -3}"
        
        if [ "$final_code" = "200" ] || [ "$final_code" = "201" ]; then
            echo -e "${GREEN}Registration fixed after restart${NC}"
        else
            echo -e "${RED}Registration still failing after restart${NC}"
            echo "Manual intervention may be required"
        fi
    else
        echo -e "\n${GREEN}Registration endpoint working${NC}"
    fi
    
    # Clean up test service
    curl -s "$REGISTRY_URL/v1/agent/service/deregister/final-test" -X PUT >/dev/null 2>&1
    curl -s "$REGISTRY_URL/v1/agent/service/deregister/final-test-2" -X PUT >/dev/null 2>&1
    
    echo -e "\n${BLUE}Registry fix attempt completed${NC}"
    echo -e "${BLUE}Run './scripts/fix-service-registration.sh' to test service registration${NC}"
}

# Check dependencies
if ! command -v curl >/dev/null 2>&1; then
    echo -e "${RED}ERROR: curl is required${NC}"
    exit 1
fi

main "$@"
