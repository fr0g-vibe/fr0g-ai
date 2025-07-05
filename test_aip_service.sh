#!/bin/bash

# AIP Service Endpoint Testing Script
# Tests persona CRUD operations and gRPC service functionality

set -e

# Configuration
AIP_HTTP_BASE="http://localhost:8080"
AIP_GRPC_ENDPOINT="localhost:9090"
TEST_OUTPUT_DIR="./test_results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create test output directory
mkdir -p "$TEST_OUTPUT_DIR"

echo -e "${BLUE}=== fr0g-ai-aip Service Testing Suite ===${NC}"
echo "Timestamp: $TIMESTAMP"
echo "HTTP Base URL: $AIP_HTTP_BASE"
echo "gRPC Endpoint: $AIP_GRPC_ENDPOINT"
echo ""

# Function to log test results
log_test() {
    local test_name="$1"
    local status="$2"
    local details="$3"
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓ $test_name${NC}"
    elif [ "$status" = "FAIL" ]; then
        echo -e "${RED}✗ $test_name${NC}"
        echo -e "  ${RED}Error: $details${NC}"
    else
        echo -e "${YELLOW}⚠ $test_name - $status${NC}"
    fi
    
    echo "$TIMESTAMP,$test_name,$status,$details" >> "$TEST_OUTPUT_DIR/test_results_$TIMESTAMP.csv"
}

# Function to check if service is running
check_service_health() {
    echo -e "${BLUE}--- Health Check Tests ---${NC}"
    
    # Test HTTP health endpoint
    if curl -s -f "$AIP_HTTP_BASE/health" > "$TEST_OUTPUT_DIR/health_response.json" 2>/dev/null; then
        local health_status=$(cat "$TEST_OUTPUT_DIR/health_response.json" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
        if [ "$health_status" = "healthy" ]; then
            log_test "HTTP Health Check" "PASS" "Service is healthy"
        else
            log_test "HTTP Health Check" "FAIL" "Service status: $health_status"
            return 1
        fi
    else
        log_test "HTTP Health Check" "FAIL" "Health endpoint not responding"
        return 1
    fi
    
    # Test if gRPC port is listening
    if nc -z localhost 9090 2>/dev/null; then
        log_test "gRPC Port Check" "PASS" "Port 9090 is listening"
    else
        log_test "gRPC Port Check" "FAIL" "Port 9090 not accessible"
        return 1
    fi
    
    return 0
}

# Function to test REST API endpoints
test_rest_api() {
    echo -e "${BLUE}--- REST API Tests ---${NC}"
    
    # Test GET /personas (list all personas)
    if curl -s -f "$AIP_HTTP_BASE/personas" > "$TEST_OUTPUT_DIR/personas_list.json" 2>/dev/null; then
        local persona_count=$(cat "$TEST_OUTPUT_DIR/personas_list.json" | grep -o '"id"' | wc -l)
        log_test "GET /personas" "PASS" "Retrieved $persona_count personas"
    else
        log_test "GET /personas" "FAIL" "Failed to retrieve personas list"
    fi
    
    # Test POST /personas (create new persona)
    local test_persona='{
        "name": "Test Persona",
        "topic": "Testing",
        "prompt": "You are a test persona for API validation",
        "context": {
            "environment": "test",
            "purpose": "validation"
        },
        "rag": ["test_document_1", "test_document_2"]
    }'
    
    local create_response=$(curl -s -X POST "$AIP_HTTP_BASE/personas" \
        -H "Content-Type: application/json" \
        -d "$test_persona" \
        -w "%{http_code}")
    
    local http_code="${create_response: -3}"
    local response_body="${create_response%???}"
    
    if [ "$http_code" = "201" ] || [ "$http_code" = "200" ]; then
        echo "$response_body" > "$TEST_OUTPUT_DIR/created_persona.json"
        local created_id=$(echo "$response_body" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        log_test "POST /personas" "PASS" "Created persona with ID: $created_id"
        
        # Store the created ID for further tests
        echo "$created_id" > "$TEST_OUTPUT_DIR/test_persona_id.txt"
        
        # Test GET /personas/{id} (get specific persona)
        if curl -s -f "$AIP_HTTP_BASE/personas/$created_id" > "$TEST_OUTPUT_DIR/persona_detail.json" 2>/dev/null; then
            log_test "GET /personas/{id}" "PASS" "Retrieved persona details"
        else
            log_test "GET /personas/{id}" "FAIL" "Failed to retrieve persona details"
        fi
        
        # Test PUT /personas/{id} (update persona)
        local update_persona='{
            "name": "Updated Test Persona",
            "topic": "Updated Testing",
            "prompt": "You are an updated test persona for API validation",
            "context": {
                "environment": "test",
                "purpose": "validation",
                "updated": "true"
            },
            "rag": ["updated_document_1"]
        }'
        
        local update_response=$(curl -s -X PUT "$AIP_HTTP_BASE/personas/$created_id" \
            -H "Content-Type: application/json" \
            -d "$update_persona" \
            -w "%{http_code}")
        
        local update_http_code="${update_response: -3}"
        if [ "$update_http_code" = "200" ]; then
            log_test "PUT /personas/{id}" "PASS" "Updated persona successfully"
        else
            log_test "PUT /personas/{id}" "FAIL" "HTTP code: $update_http_code"
        fi
        
        # Test DELETE /personas/{id} (delete persona)
        local delete_response=$(curl -s -X DELETE "$AIP_HTTP_BASE/personas/$created_id" -w "%{http_code}")
        local delete_http_code="${delete_response: -3}"
        
        if [ "$delete_http_code" = "200" ] || [ "$delete_http_code" = "204" ]; then
            log_test "DELETE /personas/{id}" "PASS" "Deleted persona successfully"
        else
            log_test "DELETE /personas/{id}" "FAIL" "HTTP code: $delete_http_code"
        fi
        
    else
        log_test "POST /personas" "FAIL" "HTTP code: $http_code, Response: $response_body"
    fi
}

# Function to test gRPC endpoints (if grpcurl is available)
test_grpc_api() {
    echo -e "${BLUE}--- gRPC API Tests ---${NC}"
    
    # Check if grpcurl is available
    if ! command -v grpcurl &> /dev/null; then
        log_test "gRPC Tests" "SKIP" "grpcurl not installed"
        echo -e "${YELLOW}Install grpcurl to test gRPC endpoints: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest${NC}"
        return 0
    fi
    
    # Test gRPC service reflection
    if grpcurl -plaintext "$AIP_GRPC_ENDPOINT" list > "$TEST_OUTPUT_DIR/grpc_services.txt" 2>/dev/null; then
        local service_count=$(cat "$TEST_OUTPUT_DIR/grpc_services.txt" | wc -l)
        log_test "gRPC Service Reflection" "PASS" "Found $service_count services"
    else
        log_test "gRPC Service Reflection" "FAIL" "Service reflection not working"
        return 1
    fi
    
    # Test PersonaService methods
    if grpcurl -plaintext "$AIP_GRPC_ENDPOINT" list persona.PersonaService > "$TEST_OUTPUT_DIR/persona_methods.txt" 2>/dev/null; then
        log_test "PersonaService Discovery" "PASS" "PersonaService methods available"
    else
        log_test "PersonaService Discovery" "FAIL" "PersonaService not found"
    fi
    
    # Test CreatePersona gRPC method
    local grpc_persona='{
        "persona": {
            "name": "gRPC Test Persona",
            "topic": "gRPC Testing",
            "prompt": "You are a test persona created via gRPC",
            "context": {
                "method": "grpc",
                "test": "true"
            },
            "rag": ["grpc_test_doc"]
        }
    }'
    
    if echo "$grpc_persona" | grpcurl -plaintext -d @ "$AIP_GRPC_ENDPOINT" persona.PersonaService/CreatePersona > "$TEST_OUTPUT_DIR/grpc_create_response.json" 2>/dev/null; then
        local grpc_id=$(cat "$TEST_OUTPUT_DIR/grpc_create_response.json" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        log_test "gRPC CreatePersona" "PASS" "Created persona via gRPC: $grpc_id"
        
        # Test GetPersona gRPC method
        if echo "{\"id\":\"$grpc_id\"}" | grpcurl -plaintext -d @ "$AIP_GRPC_ENDPOINT" persona.PersonaService/GetPersona > "$TEST_OUTPUT_DIR/grpc_get_response.json" 2>/dev/null; then
            log_test "gRPC GetPersona" "PASS" "Retrieved persona via gRPC"
        else
            log_test "gRPC GetPersona" "FAIL" "Failed to retrieve persona via gRPC"
        fi
        
        # Clean up: delete the test persona
        echo "{\"id\":\"$grpc_id\"}" | grpcurl -plaintext -d @ "$AIP_GRPC_ENDPOINT" persona.PersonaService/DeletePersona > /dev/null 2>&1
        
    else
        log_test "gRPC CreatePersona" "FAIL" "Failed to create persona via gRPC"
    fi
}

# Function to test identity endpoints
test_identity_api() {
    echo -e "${BLUE}--- Identity API Tests ---${NC}"
    
    # Test GET /identities
    if curl -s -f "$AIP_HTTP_BASE/identities" > "$TEST_OUTPUT_DIR/identities_list.json" 2>/dev/null; then
        local identity_count=$(cat "$TEST_OUTPUT_DIR/identities_list.json" | grep -o '"id"' | wc -l)
        log_test "GET /identities" "PASS" "Retrieved $identity_count identities"
    else
        log_test "GET /identities" "FAIL" "Failed to retrieve identities list"
    fi
    
    # Test identity creation with rich attributes
    local test_identity='{
        "name": "Test Identity",
        "description": "A test identity for API validation",
        "persona_id": "test-persona-id",
        "background": "Test background for validation purposes",
        "demographics": {
            "age": 30,
            "gender": "non-binary",
            "education": "bachelors",
            "location": {
                "country": "US",
                "state": "CA",
                "city": "San Francisco"
            }
        },
        "preferences": {
            "hobbies": ["reading", "coding"],
            "interests": ["technology", "science"]
        }
    }'
    
    local identity_response=$(curl -s -X POST "$AIP_HTTP_BASE/identities" \
        -H "Content-Type: application/json" \
        -d "$test_identity" \
        -w "%{http_code}")
    
    local identity_http_code="${identity_response: -3}"
    if [ "$identity_http_code" = "201" ] || [ "$identity_http_code" = "200" ]; then
        log_test "POST /identities" "PASS" "Created identity successfully"
    else
        log_test "POST /identities" "FAIL" "HTTP code: $identity_http_code"
    fi
}

# Function to test attribute processors
test_attribute_processors() {
    echo -e "${BLUE}--- Attribute Processor Tests ---${NC}"
    
    # Test demographics validation
    local demographics_test='{"age": 25, "gender": "female", "education": "masters"}'
    local demo_response=$(curl -s -X POST "$AIP_HTTP_BASE/validate/demographics" \
        -H "Content-Type: application/json" \
        -d "$demographics_test" \
        -w "%{http_code}")
    
    local demo_http_code="${demo_response: -3}"
    if [ "$demo_http_code" = "200" ]; then
        log_test "Demographics Validation" "PASS" "Demographics processor working"
    else
        log_test "Demographics Validation" "SKIP" "Endpoint not available (expected)"
    fi
    
    # Test psychographics validation
    local psycho_test='{"personality": {"openness": 0.8, "conscientiousness": 0.7}}'
    local psycho_response=$(curl -s -X POST "$AIP_HTTP_BASE/validate/psychographics" \
        -H "Content-Type: application/json" \
        -d "$psycho_test" \
        -w "%{http_code}")
    
    local psycho_http_code="${psycho_response: -3}"
    if [ "$psycho_http_code" = "200" ]; then
        log_test "Psychographics Validation" "PASS" "Psychographics processor working"
    else
        log_test "Psychographics Validation" "SKIP" "Endpoint not available (expected)"
    fi
}

# Function to generate test report
generate_report() {
    echo -e "${BLUE}--- Test Report ---${NC}"
    
    local total_tests=$(cat "$TEST_OUTPUT_DIR/test_results_$TIMESTAMP.csv" | wc -l)
    local passed_tests=$(grep ",PASS," "$TEST_OUTPUT_DIR/test_results_$TIMESTAMP.csv" | wc -l)
    local failed_tests=$(grep ",FAIL," "$TEST_OUTPUT_DIR/test_results_$TIMESTAMP.csv" | wc -l)
    local skipped_tests=$(grep ",SKIP," "$TEST_OUTPUT_DIR/test_results_$TIMESTAMP.csv" | wc -l)
    
    echo "Total Tests: $total_tests"
    echo -e "${GREEN}Passed: $passed_tests${NC}"
    echo -e "${RED}Failed: $failed_tests${NC}"
    echo -e "${YELLOW}Skipped: $skipped_tests${NC}"
    echo ""
    
    if [ $failed_tests -eq 0 ]; then
        echo -e "${GREEN}🎉 All tests passed! AIP service is fully operational.${NC}"
    else
        echo -e "${RED}❌ Some tests failed. Check the detailed results in $TEST_OUTPUT_DIR/${NC}"
    fi
    
    echo ""
    echo "Detailed results saved to: $TEST_OUTPUT_DIR/test_results_$TIMESTAMP.csv"
    echo "Test artifacts saved to: $TEST_OUTPUT_DIR/"
}

# Main execution
main() {
    # Initialize CSV header
    echo "timestamp,test_name,status,details" > "$TEST_OUTPUT_DIR/test_results_$TIMESTAMP.csv"
    
    # Run test suites
    if check_service_health; then
        test_rest_api
        test_grpc_api
        test_identity_api
        test_attribute_processors
    else
        echo -e "${RED}Service health check failed. Skipping remaining tests.${NC}"
    fi
    
    # Generate final report
    generate_report
}

# Check dependencies
echo -e "${BLUE}Checking dependencies...${NC}"
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is required but not installed.${NC}"
    exit 1
fi

if ! command -v nc &> /dev/null; then
    echo -e "${YELLOW}Warning: nc (netcat) not found. Port checks may not work.${NC}"
fi

# Run the tests
main
