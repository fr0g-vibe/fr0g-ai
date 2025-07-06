#!/bin/bash

# Service Registry Integration Test
# Verifies all fr0g.ai services properly register with the service registry

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
REGISTRY_URL="http://localhost:8500"
TIMEOUT=60
RETRY_INTERVAL=5
MAX_RETRIES=12

# Expected services
EXPECTED_SERVICES=("fr0g-ai-aip" "fr0g-ai-bridge" "fr0g-ai-io" "fr0g-ai-master-control")

# Test results tracking
TEST_RESULTS=()
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

echo -e "${BLUE}fr0g.ai Service Registry Integration Test${NC}"
echo "=============================================="

# Function to wait for registry to be ready
wait_for_registry() {
    echo -n "Waiting for service registry to be ready... "
    local retries=0
    
    while [ $retries -lt $MAX_RETRIES ]; do
        if curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1; then
            echo -e "${GREEN}COMPLETED READY${NC}"
            log_test "Registry Startup" "PASS" "Registry ready after $((retries * RETRY_INTERVAL)) seconds"
            return 0
        fi
        retries=$((retries + 1))
        sleep $RETRY_INTERVAL
    done
    
    echo -e "${RED}CRITICAL REGISTRY NOT READY${NC}"
    log_test "Registry Startup" "FAIL" "Registry not ready after $((MAX_RETRIES * RETRY_INTERVAL)) seconds"
    return 1
}

# Function to log test results
log_test() {
    local test_name="$1"
    local result="$2"
    local message="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}PASS $test_name: $message${NC}"
        TEST_RESULTS+=("PASS: $test_name - $message")
        PASSED_TESTS=$((PASSED_TESTS + 1))
    elif [ "$result" = "FAIL" ]; then
        echo -e "${RED}FAIL $test_name: $message${NC}"
        TEST_RESULTS+=("FAIL: $test_name - $message")
        FAILED_TESTS=$((FAILED_TESTS + 1))
    elif [ "$result" = "SKIP" ]; then
        echo -e "${YELLOW}SKIP $test_name: $message${NC}"
        TEST_RESULTS+=("SKIP: $test_name - $message")
    fi
}

# Function to check service registration
check_service_registration() {
    local service_name=$1
    echo -n "Checking $service_name registration... "
    
    # Try multiple registry endpoints
    local endpoints=("/v1/catalog/service/$service_name" "/v1/health/service/$service_name" "/services/$service_name")
    
    for endpoint in "${endpoints[@]}"; do
        local response=$(curl -s "$REGISTRY_URL$endpoint" 2>/dev/null)
        if echo "$response" | jq -e '. | length > 0' >/dev/null 2>&1; then
            echo -e "${GREEN}COMPLETED REGISTERED${NC}"
            log_test "$service_name Registration" "PASS" "Service is properly registered"
            return 0
        fi
    done
    
    echo -e "${RED}CRITICAL NOT REGISTERED${NC}"
    log_test "$service_name Registration" "FAIL" "Service not found in registry"
    return 1
}

# Function to verify service health through registry
verify_service_health() {
    local service_name=$1
    echo -n "Verifying $service_name health through registry... "
    
    local health_response=$(curl -s "$REGISTRY_URL/v1/health/service/$service_name" 2>/dev/null)
    if echo "$health_response" | jq -e '.[] | select(.Status == "passing")' >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED HEALTHY${NC}"
        return 0
    else
        echo -e "${RED}CRITICAL UNHEALTHY${NC}"
        return 1
    fi
}

# Function to test service discovery
test_service_discovery() {
    echo -e "\n${BLUE}Testing Service Discovery${NC}"
    echo "----------------------------"
    
    echo -n "Listing all registered services... "
    local services_response=$(curl -s "$REGISTRY_URL/v1/catalog/services" 2>/dev/null)
    if [ -n "$services_response" ] && [ "$services_response" != "{}" ]; then
        echo -e "${GREEN}COMPLETED SERVICES FOUND${NC}"
        log_test "Service Discovery" "PASS" "Service catalog responding with data"
        
        # Count services
        local service_count=$(echo "$services_response" | jq '. | length' 2>/dev/null || echo "0")
        log_test "Service Count" "PASS" "$service_count services in registry"
        
        echo "Registered services:"
        echo "$services_response" | jq -r 'keys[]' 2>/dev/null | sed 's/^/  - /' || echo "$services_response"
        
        # Test performance of discovery
        local start_time=$(date +%s%N)
        curl -s "$REGISTRY_URL/v1/catalog/services" >/dev/null 2>&1
        local end_time=$(date +%s%N)
        local response_time=$(( (end_time - start_time) / 1000000 ))
        
        if [ $response_time -lt 100 ]; then
            log_test "Discovery Performance" "PASS" "Fast discovery: ${response_time}ms"
        elif [ $response_time -lt 500 ]; then
            log_test "Discovery Performance" "PASS" "Good discovery: ${response_time}ms"
        else
            log_test "Discovery Performance" "FAIL" "Slow discovery: ${response_time}ms"
        fi
    else
        echo -e "${RED}CRITICAL NO SERVICES REGISTERED${NC}"
        log_test "Service Discovery" "FAIL" "No services found in registry"
        return 1
    fi
    
    return 0
}

# Function to test service registration endpoints
test_registration_endpoints() {
    echo -e "\n${BLUE}Testing Registration Endpoints${NC}"
    echo "--------------------------------"
    
    # Test service registration
    echo -n "Testing service registration endpoint... "
    local test_service='{"ID":"test-service","Name":"test-service","Port":9999,"Check":{"HTTP":"http://localhost:9999/health","Interval":"10s"}}'
    local reg_response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" -X POST -H "Content-Type: application/json" -d "$test_service" 2>/dev/null)
    local reg_code="${reg_response: -3}"
    
    if [ "$reg_code" = "200" ] || [ "$reg_code" = "201" ]; then
        echo -e "${GREEN}COMPLETED OPERATIONAL${NC}"
        log_test "Registration Endpoint" "PASS" "Service registration working (HTTP $reg_code)"
        
        # Test service appears in catalog
        if curl -sf "$REGISTRY_URL/v1/catalog/services" | grep -q "test-service"; then
            log_test "Service Catalog Update" "PASS" "Registered service appears in catalog"
        else
            log_test "Service Catalog Update" "FAIL" "Registered service not found in catalog"
        fi
        
        # Clean up test service
        if curl -s "$REGISTRY_URL/v1/agent/service/deregister/test-service" -X PUT >/dev/null 2>&1; then
            log_test "Service Deregistration" "PASS" "Service deregistration working"
        else
            log_test "Service Deregistration" "FAIL" "Service deregistration failed"
        fi
    else
        echo -e "${RED}CRITICAL ENDPOINT FAILED (HTTP $reg_code)${NC}"
        log_test "Registration Endpoint" "FAIL" "Service registration failed (HTTP $reg_code)"
        return 1
    fi
    
    return 0
}

# Function to generate service registration report
generate_registration_report() {
    echo -e "\n${BLUE}Service Registration Report${NC}"
    echo "============================="
    
    local all_registered=true
    
    for service in "${EXPECTED_SERVICES[@]}"; do
        echo -n "Service: $service ... "
        if check_service_registration "$service" >/dev/null 2>&1; then
            echo -e "${GREEN}REGISTERED${NC}"
        else
            echo -e "${RED}NOT REGISTERED${NC}"
            all_registered=false
        fi
    done
    
    if $all_registered; then
        echo -e "\n${GREEN}SUCCESS: All services properly registered${NC}"
        return 0
    else
        echo -e "\n${RED}CRITICAL: Some services not registered${NC}"
        echo -e "${RED}This will break service discovery and inter-service communication${NC}"
        return 1
    fi
}

# Main test execution
main() {
    echo -e "${BLUE}Starting service registry integration test...${NC}\n"
    
    local exit_code=0
    
    # Wait for registry
    if ! wait_for_registry; then
        echo -e "${RED}CRITICAL: Service registry not available${NC}"
        exit 1
    fi
    
    # Test registration endpoints
    if ! test_registration_endpoints; then
        echo -e "${RED}CRITICAL: Registration endpoints not working${NC}"
        exit_code=1
    fi
    
    # Test service discovery
    if ! test_service_discovery; then
        echo -e "${RED}CRITICAL: Service discovery not working${NC}"
        exit_code=1
    fi
    
    # Wait for services to register
    echo -e "\n${BLUE}Waiting for services to register...${NC}"
    sleep 30
    
    # Generate registration report
    if ! generate_registration_report; then
        exit_code=1
    fi
    
    # Verify service health through registry
    echo -e "\n${BLUE}Verifying Service Health Through Registry${NC}"
    echo "-------------------------------------------"
    for service in "${EXPECTED_SERVICES[@]}"; do
        verify_service_health "$service" || exit_code=1
    done
    
    # Summary
    echo -e "\n${BLUE}Service Registry Test Summary${NC}"
    echo "============================="
    
    echo -e "Total Tests: $TOTAL_TESTS"
    echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
    echo -e "${YELLOW}Skipped: $((TOTAL_TESTS - PASSED_TESTS - FAILED_TESTS))${NC}"
    
    # Show failed tests
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "\n${RED}Failed Tests:${NC}"
        for result in "${TEST_RESULTS[@]}"; do
            if [[ $result == FAIL:* ]]; then
                echo -e "${RED}  $result${NC}"
            fi
        done
    fi
    
    if [ $exit_code -eq 0 ]; then
        echo -e "\n${GREEN}SUCCESS: Service registry integration working properly${NC}"
    else
        echo -e "\n${RED}CRITICAL: Service registry integration has issues${NC}"
        echo -e "${RED}Check service logs and registry configuration${NC}"
        echo -e "\n${BLUE}Troubleshooting Commands:${NC}"
        echo "  docker-compose logs service-registry"
        echo "  curl -v $REGISTRY_URL/health"
        echo "  make diagnose-registry"
    fi
    
    exit $exit_code
}

# Check dependencies
if ! command -v curl >/dev/null 2>&1; then
    echo -e "${RED}ERROR: curl is required${NC}"
    exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
    echo -e "${YELLOW}WARNING: jq not found - JSON parsing will be limited${NC}"
fi

# Run the test
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
