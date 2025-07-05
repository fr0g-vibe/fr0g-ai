#!/bin/bash

# OpenWebUI Integration Test for fr0g-ai-bridge
# Tests actual OpenWebUI communication patterns

set -e

echo "NETWORK Testing OpenWebUI Integration with fr0g-ai-bridge"
echo "===================================================="

# Configuration
BRIDGE_HOST="localhost"
BRIDGE_PORT="8082"
BASE_URL="http://${BRIDGE_HOST}:${BRIDGE_PORT}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function
test_openwebui_endpoint() {
    local test_name="$1"
    local endpoint="$2"
    local method="$3"
    local data="$4"
    local expected_status="$5"
    
    echo -n "Testing $test_name... "
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X "$method" "$BASE_URL$endpoint" \
                   -H "Content-Type: application/json" \
                   -d "$data" -o /tmp/openwebui_test_response)
    else
        response=$(curl -s -w "%{http_code}" -X "$method" "$BASE_URL$endpoint" \
                   -o /tmp/openwebui_test_response)
    fi
    
    status_code="${response: -3}"
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}PASS${NC} (HTTP $status_code)"
        ((TESTS_PASSED++))
        
        # Show response for successful tests
        if [ "$status_code" = "200" ]; then
            echo -e "${BLUE}Response:${NC} $(cat /tmp/openwebui_test_response | jq -c . 2>/dev/null || cat /tmp/openwebui_test_response)"
        fi
    else
        echo -e "${RED}FAIL${NC} (Expected HTTP $expected_status, got $status_code)"
        echo "Response: $(cat /tmp/openwebui_test_response)"
        ((TESTS_FAILED++))
    fi
}

# Check if Bridge service is running
echo "CHECKING Checking Bridge service availability..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo -e "${RED}FAILED Bridge service not running on $BASE_URL${NC}"
    echo "Start with: cd fr0g-ai-bridge && ./main"
    exit 1
fi

echo -e "${GREEN}COMPLETED Bridge service is running${NC}"
echo ""

# Test 1: OpenWebUI Chat Completions (Primary Integration Point)
echo "💬 Testing OpenWebUI Chat Completions Integration"
chat_request='{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello, can you help me test the fr0g.ai bridge?"}
  ],
  "temperature": 0.7,
  "max_tokens": 150
}'

test_openwebui_endpoint "Chat completions with system message" "/api/chat/completions" "POST" "$chat_request" "200"
echo ""

# Test 2: OpenWebUI Models Endpoint
echo "LIST Testing OpenWebUI Models Integration"
test_openwebui_endpoint "Models list endpoint" "/api/v1/models" "GET" "" "200"
echo ""

# Test 3: Persona-Aware Chat (fr0g.ai Extension)
echo "🎭 Testing Persona-Aware Chat (fr0g.ai Extension)"
persona_request='{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "user", "content": "Analyze this from a security perspective"}
  ],
  "persona_prompt": "You are a cybersecurity expert with 10 years of experience in threat analysis.",
  "temperature": 0.5
}'

test_openwebui_endpoint "Persona-aware chat completion" "/api/chat/completions" "POST" "$persona_request" "200"
echo ""

# Test 4: OpenWebUI V1 Chat Endpoint
echo "REFRESH Testing OpenWebUI V1 Chat Endpoint"
v1_request='{
  "message": "Test message for v1 endpoint",
  "model": "gpt-3.5-turbo"
}'

test_openwebui_endpoint "V1 chat endpoint" "/api/v1/chat" "POST" "$v1_request" "200"
echo ""

# Test 5: Error Handling for Invalid Requests
echo "ALERT Testing Error Handling"
invalid_request='{"invalid": "request"}'
test_openwebui_endpoint "Invalid request handling" "/api/chat/completions" "POST" "$invalid_request" "200"

# Test method not allowed
test_openwebui_endpoint "Method not allowed handling" "/api/chat/completions" "GET" "" "405"
echo ""

# Test 6: Large Request Handling
echo "INSTALLING Testing Large Request Handling"
large_content=$(printf 'A%.0s' {1..1000})  # 1000 character string
large_request="{
  \"model\": \"gpt-3.5-turbo\",
  \"messages\": [
    {\"role\": \"user\", \"content\": \"$large_content\"}
  ]
}"

test_openwebui_endpoint "Large request handling" "/api/chat/completions" "POST" "$large_request" "200"
echo ""

# Test 7: Concurrent Request Handling
echo "PERFORMANCE Testing Concurrent Request Handling"
echo "Sending 5 concurrent requests..."

pids=()
for i in {1..5}; do
    (
        response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/chat/completions" \
                   -H "Content-Type: application/json" \
                   -d "$chat_request" -o "/tmp/concurrent_test_$i")
        echo "$response" > "/tmp/concurrent_result_$i"
    ) &
    pids+=($!)
done

# Wait for all requests to complete
for pid in "${pids[@]}"; do
    wait "$pid"
done

# Check results
concurrent_passed=0
for i in {1..5}; do
    if [ -f "/tmp/concurrent_result_$i" ]; then
        status=$(cat "/tmp/concurrent_result_$i")
        if [ "${status: -3}" = "200" ]; then
            ((concurrent_passed++))
        fi
    fi
done

echo -n "Concurrent request handling... "
if [ "$concurrent_passed" -eq 5 ]; then
    echo -e "${GREEN}PASS${NC} (5/5 requests successful)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} ($concurrent_passed/5 requests successful)"
    ((TESTS_FAILED++))
fi
echo ""

# Test 8: Response Format Validation
echo "LIST Testing Response Format Validation"
echo -n "Validating OpenAI-compatible response format... "

response=$(curl -s -X POST "$BASE_URL/api/chat/completions" \
           -H "Content-Type: application/json" \
           -d "$chat_request")

# Check required OpenAI fields
required_fields=("id" "object" "created" "model" "choices")
format_valid=true

for field in "${required_fields[@]}"; do
    if ! echo "$response" | jq -e ".$field" > /dev/null 2>&1; then
        format_valid=false
        break
    fi
done

if [ "$format_valid" = true ]; then
    echo -e "${GREEN}PASS${NC} (All required OpenAI fields present)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} (Missing required OpenAI fields)"
    echo "Response: $response"
    ((TESTS_FAILED++))
fi
echo ""

# Test 9: Health Check Integration
echo "HEALTH Testing Health Check Integration"
echo -n "Testing health endpoint for OpenWebUI monitoring... "

health_response=$(curl -s "$BASE_URL/health")
if echo "$health_response" | jq -e '.status and .service' > /dev/null 2>&1; then
    echo -e "${GREEN}PASS${NC} (Health check suitable for OpenWebUI monitoring)"
    ((TESTS_PASSED++))
    echo -e "${BLUE}Health Status:${NC} $(echo "$health_response" | jq -c .)"
else
    echo -e "${RED}FAIL${NC} (Health check format invalid)"
    ((TESTS_FAILED++))
fi
echo ""

# Cleanup temporary files
rm -f /tmp/openwebui_test_response /tmp/concurrent_test_* /tmp/concurrent_result_*

# Final Results
echo "METRICS OpenWebUI Integration Test Results"
echo "====================================="
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo -e "Total Tests: $((TESTS_PASSED + TESTS_FAILED))"

if [ "$TESTS_FAILED" -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All OpenWebUI integration tests passed!${NC}"
    echo -e "${GREEN}COMPLETED Bridge service is ready for OpenWebUI integration${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Configure OpenWebUI to use: $BASE_URL"
    echo "2. Set API endpoint in OpenWebUI to: $BASE_URL/api/chat/completions"
    echo "3. Test with actual OpenWebUI instance"
    exit 0
else
    echo -e "\n${RED}FAILED Some OpenWebUI integration tests failed${NC}"
    echo "Please fix the issues before integrating with OpenWebUI"
    exit 1
fi
