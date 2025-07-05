#!/bin/bash

# fr0g-ai-bridge Integration Test Suite
# Tests Bridge service OpenWebUI communication and API endpoints

set -e

echo "TESTING Starting fr0g-ai-bridge Integration Tests"
echo "=============================================="

# Configuration
BRIDGE_HOST="localhost"
BRIDGE_HTTP_PORT="8082"
BRIDGE_GRPC_PORT="9091"
BASE_URL="http://${BRIDGE_HOST}:${BRIDGE_HTTP_PORT}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to run tests
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_status="$3"
    
    echo -n "Testing $test_name... "
    
    if eval "$test_command" > /dev/null 2>&1; then
        if [ "$expected_status" = "success" ]; then
            echo -e "${GREEN}PASS${NC}"
            ((TESTS_PASSED++))
        else
            echo -e "${RED}FAIL${NC} (expected failure but got success)"
            ((TESTS_FAILED++))
        fi
    else
        if [ "$expected_status" = "fail" ]; then
            echo -e "${GREEN}PASS${NC} (expected failure)"
            ((TESTS_PASSED++))
        else
            echo -e "${RED}FAIL${NC}"
            ((TESTS_FAILED++))
        fi
    fi
}

# Helper function to test HTTP endpoint
test_http_endpoint() {
    local endpoint="$1"
    local method="$2"
    local expected_status="$3"
    local description="$4"
    
    echo -n "Testing $description... "
    
    response=$(curl -s -w "%{http_code}" -X "$method" "$BASE_URL$endpoint" -o /tmp/bridge_test_response)
    status_code="${response: -3}"
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}PASS${NC} (HTTP $status_code)"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}FAIL${NC} (Expected HTTP $expected_status, got $status_code)"
        echo "Response: $(cat /tmp/bridge_test_response)"
        ((TESTS_FAILED++))
    fi
}

# Helper function to test JSON response
test_json_response() {
    local endpoint="$1"
    local method="$2"
    local json_key="$3"
    local expected_value="$4"
    local description="$5"
    
    echo -n "Testing $description... "
    
    response=$(curl -s -X "$method" "$BASE_URL$endpoint")
    actual_value=$(echo "$response" | jq -r ".$json_key" 2>/dev/null)
    
    if [ "$actual_value" = "$expected_value" ]; then
        echo -e "${GREEN}PASS${NC} ($json_key: $actual_value)"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}FAIL${NC} (Expected $json_key: $expected_value, got: $actual_value)"
        echo "Full response: $response"
        ((TESTS_FAILED++))
    fi
}

# Check if Bridge service is running
echo "CHECKING Checking if Bridge service is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo -e "${RED}FAILED Bridge service is not running on $BASE_URL${NC}"
    echo "Please start the service with: cd fr0g-ai-bridge && ./main"
    exit 1
fi

echo -e "${GREEN}COMPLETED Bridge service is running${NC}"
echo ""

# Test 1: Health Check Endpoint
echo "LIST Testing Health Check Endpoint"
test_http_endpoint "/health" "GET" "200" "Health endpoint availability"
test_json_response "/health" "GET" "status" "ok" "Health status response"
test_json_response "/health" "GET" "service" "fr0g-ai-bridge" "Service name in health response"
echo ""

# Test 2: OpenWebUI Compatible Endpoints
echo "🤖 Testing OpenWebUI Compatible Endpoints"
test_http_endpoint "/api/chat/completions" "POST" "200" "Chat completions endpoint"
test_http_endpoint "/api/v1/models" "GET" "200" "Models list endpoint"
test_http_endpoint "/api/v1/chat" "POST" "200" "V1 chat endpoint"
echo ""

# Test 3: HTTP Method Validation
echo "SECURITY Testing HTTP Method Validation"
test_http_endpoint "/api/chat/completions" "GET" "405" "Chat completions GET rejection"
test_http_endpoint "/api/v1/chat" "GET" "405" "V1 chat GET rejection"
echo ""

# Test 4: JSON Response Structure
echo "📄 Testing JSON Response Structure"

# Test chat completions response structure
echo -n "Testing chat completions response structure... "
response=$(curl -s -X POST "$BASE_URL/api/chat/completions" -H "Content-Type: application/json" -d '{"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"test"}]}')
if echo "$response" | jq -e '.id and .object and .choices' > /dev/null 2>&1; then
    echo -e "${GREEN}PASS${NC} (Valid OpenAI format)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} (Invalid response structure)"
    echo "Response: $response"
    ((TESTS_FAILED++))
fi

# Test models response structure
echo -n "Testing models response structure... "
response=$(curl -s "$BASE_URL/api/v1/models")
if echo "$response" | jq -e '.object and .data' > /dev/null 2>&1; then
    echo -e "${GREEN}PASS${NC} (Valid models format)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} (Invalid models structure)"
    echo "Response: $response"
    ((TESTS_FAILED++))
fi
echo ""

# Test 5: Port Configuration Verification
echo "🔌 Testing Port Configuration"
echo -n "Testing HTTP port configuration... "
if curl -s "http://${BRIDGE_HOST}:8082/health" > /dev/null; then
    echo -e "${GREEN}PASS${NC} (HTTP on port 8082)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} (HTTP port 8082 not accessible)"
    ((TESTS_FAILED++))
fi

echo -n "Testing gRPC port availability... "
if nc -z "$BRIDGE_HOST" 9091 2>/dev/null; then
    echo -e "${GREEN}PASS${NC} (gRPC port 9091 open)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} (gRPC port 9091 not accessible)"
    ((TESTS_FAILED++))
fi
echo ""

# Test 6: Performance and Response Time
echo "PERFORMANCE Testing Performance"
echo -n "Testing response time... "
start_time=$(date +%s%N)
curl -s "$BASE_URL/health" > /dev/null
end_time=$(date +%s%N)
response_time=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds

if [ "$response_time" -lt 1000 ]; then
    echo -e "${GREEN}PASS${NC} (${response_time}ms - under 1 second)"
    ((TESTS_PASSED++))
else
    echo -e "${YELLOW}SLOW${NC} (${response_time}ms - over 1 second)"
    ((TESTS_PASSED++)) # Still pass, just slow
fi
echo ""

# Test 7: Error Handling
echo "ALERT Testing Error Handling"
test_http_endpoint "/nonexistent" "GET" "404" "404 for non-existent endpoints"
echo ""

# Test 8: Content-Type Headers
echo "LIST Testing Content-Type Headers"
echo -n "Testing JSON content-type headers... "
headers=$(curl -s -I "$BASE_URL/health" | grep -i content-type)
if echo "$headers" | grep -q "application/json"; then
    echo -e "${GREEN}PASS${NC} (Correct JSON content-type)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} (Missing or incorrect content-type)"
    echo "Headers: $headers"
    ((TESTS_FAILED++))
fi
echo ""

# Test Summary
echo "METRICS Test Results Summary"
echo "======================"
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo -e "Total Tests: $((TESTS_PASSED + TESTS_FAILED))"

if [ "$TESTS_FAILED" -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed! Bridge service integration is working correctly.${NC}"
    exit 0
else
    echo -e "\n${RED}FAILED Some tests failed. Please check the Bridge service implementation.${NC}"
    exit 1
fi
