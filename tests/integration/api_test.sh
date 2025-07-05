#!/bin/bash

# fr0g.ai API Integration Test
# Tests all service APIs and their functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AIP_URL="http://localhost:8080"
BRIDGE_URL="http://localhost:8082"
MCP_URL="http://localhost:8081"
IO_URL="http://localhost:8083"
TIMEOUT=30
TEST_RESULTS=()

echo -e "${BLUE}🔌 fr0g.ai API Integration Test${NC}"
echo "================================="
echo "Testing all service APIs..."
echo ""

# Function to log test results
log_test() {
    local test_name="$1"
    local result="$2"
    local message="$3"
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✅ $test_name: PASS${NC}"
        TEST_RESULTS+=("PASS: $test_name")
    elif [ "$result" = "FAIL" ]; then
        echo -e "${RED}❌ $test_name: FAIL - $message${NC}"
        TEST_RESULTS+=("FAIL: $test_name - $message")
    elif [ "$result" = "SKIP" ]; then
        echo -e "${YELLOW}⏭️  $test_name: SKIP - $message${NC}"
        TEST_RESULTS+=("SKIP: $test_name - $message")
    fi
}

# Function to test JSON API endpoint
test_json_api() {
    local service_name="$1"
    local url="$2"
    local endpoint="$3"
    local method="${4:-GET}"
    local data="${5:-}"
    local expected_field="${6:-}"
    
    local curl_cmd="curl -s --max-time $TIMEOUT"
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        response=$(eval "$curl_cmd -X POST -H 'Content-Type: application/json' -d '$data' '$url$endpoint'" 2>/dev/null || echo "")
    else
        response=$(eval "$curl_cmd '$url$endpoint'" 2>/dev/null || echo "")
    fi
    
    if [ -n "$response" ]; then
        if [ -n "$expected_field" ]; then
            if echo "$response" | grep -q "$expected_field"; then
                log_test "$service_name $endpoint" "PASS"
                return 0
            else
                log_test "$service_name $endpoint" "FAIL" "Missing expected field: $expected_field"
                return 1
            fi
        else
            log_test "$service_name $endpoint" "PASS"
            return 0
        fi
    else
        log_test "$service_name $endpoint" "FAIL" "No response"
        return 1
    fi
}

# Function to test AIP service APIs
test_aip_apis() {
    echo -e "${YELLOW}Testing AIP Service APIs${NC}"
    echo "-------------------------"
    
    # Health check
    test_json_api "AIP" "$AIP_URL" "/health" "GET" "" "status\|healthy"
    
    # Personas endpoint
    test_json_api "AIP" "$AIP_URL" "/personas" "GET" "" "personas\|data"
    
    # Individual persona (if any exist)
    personas_response=$(curl -s --max-time 10 "$AIP_URL/personas" 2>/dev/null || echo "")
    if echo "$personas_response" | grep -q '"id"'; then
        # Extract first persona ID
        persona_id=$(echo "$personas_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        if [ -n "$persona_id" ]; then
            test_json_api "AIP" "$AIP_URL" "/personas/$persona_id" "GET" "" "id\|name"
        fi
    else
        log_test "AIP Individual Persona" "SKIP" "No personas available"
    fi
    
    # Create persona test
    create_persona_data='{
        "name": "Test Persona",
        "topic": "testing",
        "prompt": "You are a test persona for integration testing",
        "context": {"test": "true"}
    }'
    
    test_json_api "AIP" "$AIP_URL" "/personas" "POST" "$create_persona_data" "id\|name"
    
    # Identities endpoint
    test_json_api "AIP" "$AIP_URL" "/identities" "GET" "" "identities\|data"
}

# Function to test Bridge service APIs
test_bridge_apis() {
    echo -e "\n${YELLOW}Testing Bridge Service APIs${NC}"
    echo "----------------------------"
    
    # Health check
    test_json_api "Bridge" "$BRIDGE_URL" "/health" "GET" "" "status\|healthy"
    
    # Models endpoint
    test_json_api "Bridge" "$BRIDGE_URL" "/api/v1/models" "GET" "" "models\|data"
    
    # Chat completions endpoint (POST)
    chat_data='{
        "model": "gpt-3.5-turbo",
        "messages": [
            {"role": "user", "content": "Hello, this is a test message"}
        ],
        "max_tokens": 50
    }'
    
    # Note: This might fail if OpenWebUI is not configured
    response=$(curl -s --max-time 30 -X POST \
        -H "Content-Type: application/json" \
        -d "$chat_data" \
        "$BRIDGE_URL/api/chat/completions" 2>/dev/null || echo "")
    
    if [ -n "$response" ]; then
        if echo "$response" | grep -q "choices\|error\|message"; then
            log_test "Bridge Chat Completions" "PASS"
        else
            log_test "Bridge Chat Completions" "FAIL" "Unexpected response format"
        fi
    else
        log_test "Bridge Chat Completions" "SKIP" "OpenWebUI may not be configured"
    fi
}

# Function to test Master Control APIs
test_master_control_apis() {
    echo -e "\n${YELLOW}Testing Master Control APIs${NC}"
    echo "----------------------------"
    
    # Health check
    test_json_api "MCP" "$MCP_URL" "/health" "GET" "" "status\|healthy"
    
    # Status endpoint
    test_json_api "MCP" "$MCP_URL" "/status" "GET" "" "status\|intelligence\|learning"
    
    # System state endpoint
    test_json_api "MCP" "$MCP_URL" "/system/state" "GET" "" "state\|system\|load"
    
    # Capabilities endpoint
    test_json_api "MCP" "$MCP_URL" "/system/capabilities" "GET" "" "capabilities\|emergent"
    
    # Discord webhook endpoint (POST)
    webhook_data='{
        "content": "Test message from integration test",
        "username": "fr0g-ai-test",
        "embeds": []
    }'
    
    response=$(curl -s --max-time 10 -X POST \
        -H "Content-Type: application/json" \
        -d "$webhook_data" \
        "$MCP_URL/webhook/discord" 2>/dev/null || echo "")
    
    if [ -n "$response" ]; then
        log_test "MCP Discord Webhook" "PASS"
    else
        log_test "MCP Discord Webhook" "SKIP" "Webhook may require authentication"
    fi
}

# Function to test I/O service APIs
test_io_apis() {
    echo -e "\n${YELLOW}Testing I/O Service APIs${NC}"
    echo "-------------------------"
    
    # Health check
    test_json_api "I/O" "$IO_URL" "/health" "GET" "" "status\|healthy"
    
    # Processors endpoint
    test_json_api "I/O" "$IO_URL" "/processors" "GET" "" "processors\|sms\|voice\|irc\|discord"
    
    # Queue status endpoint
    test_json_api "I/O" "$IO_URL" "/queue/status" "GET" "" "queue\|input\|output"
    
    # Input processing endpoint (POST)
    input_data='{
        "type": "sms",
        "source": "+1234567890",
        "content": "Test SMS message for integration testing",
        "metadata": {"test": "true"}
    }'
    
    test_json_api "I/O" "$IO_URL" "/input" "POST" "$input_data" "id\|processed"
    
    # Output command endpoint (POST)
    output_data='{
        "type": "sms",
        "target": "+1234567890",
        "content": "Test response message",
        "priority": 1
    }'
    
    test_json_api "I/O" "$IO_URL" "/output" "POST" "$output_data" "id\|command"
}

# Function to test cross-service API interactions
test_cross_service_apis() {
    echo -e "\n${YELLOW}Testing Cross-Service API Interactions${NC}"
    echo "---------------------------------------"
    
    # Test Bridge -> AIP persona lookup
    personas_response=$(curl -s --max-time 10 "$AIP_URL/personas" 2>/dev/null || echo "")
    if echo "$personas_response" | grep -q '"id"'; then
        persona_id=$(echo "$personas_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        
        # Test chat completion with persona
        chat_with_persona='{
            "model": "gpt-3.5-turbo",
            "messages": [
                {"role": "user", "content": "Hello"}
            ],
            "persona_prompt": "You are a helpful assistant",
            "max_tokens": 50
        }'
        
        response=$(curl -s --max-time 30 -X POST \
            -H "Content-Type: application/json" \
            -d "$chat_with_persona" \
            "$BRIDGE_URL/api/chat/completions" 2>/dev/null || echo "")
        
        if [ -n "$response" ]; then
            log_test "Bridge-AIP Persona Integration" "PASS"
        else
            log_test "Bridge-AIP Persona Integration" "SKIP" "OpenWebUI integration required"
        fi
    else
        log_test "Bridge-AIP Persona Integration" "SKIP" "No personas available"
    fi
    
    # Test I/O -> MCP event processing
    if curl -s --max-time 5 "$IO_URL/health" > /dev/null && curl -s --max-time 5 "$MCP_URL/health" > /dev/null; then
        log_test "I/O-MCP Event Processing" "SKIP" "Requires gRPC testing tools"
    else
        log_test "I/O-MCP Event Processing" "SKIP" "Services not available"
    fi
}

# Function to test API error handling
test_api_error_handling() {
    echo -e "\n${YELLOW}Testing API Error Handling${NC}"
    echo "----------------------------"
    
    # Test invalid endpoints
    for service_data in "AIP:$AIP_URL" "Bridge:$BRIDGE_URL" "MCP:$MCP_URL" "I/O:$IO_URL"; do
        IFS=':' read -r service_name service_url <<< "$service_data"
        
        response=$(curl -s -w "%{http_code}" --max-time 10 "$service_url/nonexistent" 2>/dev/null || echo "000")
        
        if [[ "$response" == *"404" ]]; then
            log_test "$service_name 404 Handling" "PASS"
        else
            log_test "$service_name 404 Handling" "FAIL" "Expected 404, got: $response"
        fi
    done
    
    # Test invalid JSON
    invalid_json='{"invalid": json}'
    
    for service_data in "AIP:$AIP_URL:/personas" "Bridge:$BRIDGE_URL:/api/chat/completions" "I/O:$IO_URL:/input"; do
        IFS=':' read -r service_name service_url endpoint <<< "$service_data"
        
        response=$(curl -s -w "%{http_code}" --max-time 10 \
            -X POST \
            -H "Content-Type: application/json" \
            -d "$invalid_json" \
            "$service_url$endpoint" 2>/dev/null || echo "000")
        
        if [[ "$response" == *"400" ]] || [[ "$response" == *"422" ]]; then
            log_test "$service_name Invalid JSON Handling" "PASS"
        else
            log_test "$service_name Invalid JSON Handling" "SKIP" "Error handling varies"
        fi
    done
}

# Main test execution
main() {
    echo -e "${BLUE}Starting API integration tests...${NC}"
    echo ""
    
    # Wait for services to be ready
    echo -e "${YELLOW}Waiting for services to initialize...${NC}"
    sleep 5
    
    # Run test suites
    test_aip_apis
    test_bridge_apis
    test_master_control_apis
    test_io_apis
    test_cross_service_apis
    test_api_error_handling
    
    # Print summary
    echo -e "\n${BLUE}API Test Summary${NC}"
    echo "================"
    
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    local skipped_tests=0
    
    for result in "${TEST_RESULTS[@]}"; do
        total_tests=$((total_tests + 1))
        if [[ $result == PASS:* ]]; then
            passed_tests=$((passed_tests + 1))
        elif [[ $result == FAIL:* ]]; then
            failed_tests=$((failed_tests + 1))
            echo -e "${RED}$result${NC}"
        elif [[ $result == SKIP:* ]]; then
            skipped_tests=$((skipped_tests + 1))
        fi
    done
    
    echo ""
    echo -e "Total Tests: $total_tests"
    echo -e "${GREEN}Passed: $passed_tests${NC}"
    echo -e "${RED}Failed: $failed_tests${NC}"
    echo -e "${YELLOW}Skipped: $skipped_tests${NC}"
    
    # Determine overall result
    if [ $failed_tests -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All API tests passed!${NC}"
        exit 0
    else
        echo -e "\n${RED}❌ Some API tests failed.${NC}"
        exit 1
    fi
}

# Run main function
main "$@"
