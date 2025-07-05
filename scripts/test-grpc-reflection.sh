#!/bin/bash

# gRPC Reflection Testing Script
# Tests gRPC reflection capabilities and security configurations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AIP_GRPC_ENDPOINT="localhost:9090"
AIP_HTTP_ENDPOINT="http://localhost:8080"

echo -e "${BLUE}CHECKING gRPC Reflection Testing Suite${NC}"
echo "================================="
echo "gRPC Endpoint: $AIP_GRPC_ENDPOINT"
echo "HTTP Endpoint: $AIP_HTTP_ENDPOINT"
echo ""

# Function to check if grpcurl is available
check_grpcurl() {
    if ! command -v grpcurl >/dev/null 2>&1; then
        echo -e "${RED}FAILED grpcurl not found${NC}"
        echo "Install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
        exit 1
    fi
    echo -e "${GREEN}COMPLETED grpcurl found${NC}"
}

# Function to test reflection availability
test_reflection_availability() {
    echo -e "\n${BLUE}CHECKING Testing Reflection Availability${NC}"
    echo "-----------------------------------"
    
    echo -n "Testing reflection endpoint... "
    if grpcurl -plaintext "$AIP_GRPC_ENDPOINT" list >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED AVAILABLE${NC}"
        return 0
    else
        echo -e "${RED}FAILED NOT AVAILABLE${NC}"
        return 1
    fi
}

# Function to list services via reflection
list_services() {
    echo -e "\n${BLUE}LIST Listing Available Services${NC}"
    echo "-----------------------------"
    
    local services=$(grpcurl -plaintext "$AIP_GRPC_ENDPOINT" list 2>/dev/null)
    if [ -n "$services" ]; then
        echo "Available services:"
        echo "$services" | while read -r service; do
            echo -e "  ${GREEN}✓${NC} $service"
        done
        echo ""
        return 0
    else
        echo -e "${RED}FAILED No services found${NC}"
        return 1
    fi
}

# Function to test service methods
test_service_methods() {
    echo -e "${BLUE}CONFIG Testing Service Methods${NC}"
    echo "--------------------------"
    
    # Test PersonaService methods
    echo -n "Listing PersonaService methods... "
    if grpcurl -plaintext "$AIP_GRPC_ENDPOINT" list persona.PersonaService >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED AVAILABLE${NC}"
        
        local methods=$(grpcurl -plaintext "$AIP_GRPC_ENDPOINT" list persona.PersonaService 2>/dev/null)
        echo "Available methods:"
        echo "$methods" | while read -r method; do
            echo -e "  ${GREEN}✓${NC} $method"
        done
    else
        echo -e "${RED}FAILED NOT AVAILABLE${NC}"
        return 1
    fi
}

# Function to test method calls
test_method_calls() {
    echo -e "\n${BLUE}TESTING Testing Method Calls${NC}"
    echo "-----------------------"
    
    # Test CreatePersona
    echo -n "Testing CreatePersona method... "
    local test_persona='{
        "persona": {
            "name": "Reflection Test Persona",
            "topic": "gRPC Reflection Testing",
            "prompt": "You are a test persona created via gRPC reflection testing",
            "context": {
                "method": "grpc_reflection",
                "test": "true"
            },
            "rag": ["reflection_test_doc"]
        }
    }'
    
    local response=$(echo "$test_persona" | grpcurl -plaintext -d @ "$AIP_GRPC_ENDPOINT" persona.PersonaService/CreatePersona 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$response" ]; then
        echo -e "${GREEN}COMPLETED SUCCESS${NC}"
        
        # Extract persona ID for cleanup
        local persona_id=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        if [ -n "$persona_id" ]; then
            echo "  Created persona ID: $persona_id"
            
            # Test GetPersona
            echo -n "Testing GetPersona method... "
            local get_response=$(echo "{\"id\":\"$persona_id\"}" | grpcurl -plaintext -d @ "$AIP_GRPC_ENDPOINT" persona.PersonaService/GetPersona 2>/dev/null)
            if [ $? -eq 0 ] && [ -n "$get_response" ]; then
                echo -e "${GREEN}COMPLETED SUCCESS${NC}"
            else
                echo -e "${RED}FAILED FAILED${NC}"
            fi
            
            # Cleanup: delete the test persona
            echo -n "Cleaning up test persona... "
            if echo "{\"id\":\"$persona_id\"}" | grpcurl -plaintext -d @ "$AIP_GRPC_ENDPOINT" persona.PersonaService/DeletePersona >/dev/null 2>&1; then
                echo -e "${GREEN}COMPLETED CLEANED UP${NC}"
            else
                echo -e "${YELLOW}WARNING  CLEANUP FAILED${NC}"
            fi
        fi
    else
        echo -e "${RED}FAILED FAILED${NC}"
        return 1
    fi
}

# Function to check health endpoint reflection status
check_health_reflection_status() {
    echo -e "\n${BLUE}HEALTH Checking Health Endpoint Reflection Status${NC}"
    echo "---------------------------------------------"
    
    echo -n "Fetching health status... "
    local health_response=$(curl -s "$AIP_HTTP_ENDPOINT/health" 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$health_response" ]; then
        echo -e "${GREEN}COMPLETED SUCCESS${NC}"
        
        # Parse reflection status
        local reflection_status=$(echo "$health_response" | jq -r '.grpc_reflection // "unknown"' 2>/dev/null)
        echo "gRPC Reflection Status: $reflection_status"
        
        # Check for warnings
        local reflection_warning=$(echo "$health_response" | jq -r '.grpc_reflection_warning // ""' 2>/dev/null)
        if [ -n "$reflection_warning" ]; then
            echo -e "${YELLOW}WARNING  Warning: $reflection_warning${NC}"
        fi
        
        # Display full health response
        echo "Full health response:"
        echo "$health_response" | jq . 2>/dev/null || echo "$health_response"
    else
        echo -e "${RED}FAILED FAILED${NC}"
        return 1
    fi
}

# Function to test security implications
test_security_implications() {
    echo -e "\n${BLUE}SECURITY Security Implications${NC}"
    echo "------------------------"
    
    echo -e "${YELLOW}WARNING  Security Considerations:${NC}"
    echo "1. Reflection exposes all service definitions"
    echo "2. Attackers can discover available methods"
    echo "3. Should be disabled in production environments"
    echo "4. Useful for development and testing only"
    echo ""
    
    echo -e "${BLUE}TIP Recommendations:${NC}"
    echo "• Use GRPC_ENABLE_REFLECTION=false in production"
    echo "• Set ENVIRONMENT=production to enforce security"
    echo "• Monitor reflection status in health checks"
    echo "• Use make validate-production before deployment"
}

# Function to generate summary
generate_summary() {
    echo -e "\n${BLUE}METRICS Reflection Test Summary${NC}"
    echo "=========================="
    
    local reflection_available=false
    if grpcurl -plaintext "$AIP_GRPC_ENDPOINT" list >/dev/null 2>&1; then
        reflection_available=true
    fi
    
    if [ "$reflection_available" = true ]; then
        echo -e "${YELLOW}WARNING  gRPC Reflection is ENABLED${NC}"
        echo -e "${YELLOW}   This is suitable for development/testing${NC}"
        echo -e "${YELLOW}   Should be DISABLED for production${NC}"
        echo ""
        echo -e "${BLUE}To disable reflection:${NC}"
        echo "  export GRPC_ENABLE_REFLECTION=false"
        echo "  export ENVIRONMENT=production"
        echo "  make validate-production"
    else
        echo -e "${GREEN}COMPLETED gRPC Reflection is DISABLED${NC}"
        echo -e "${GREEN}   This is the recommended production setting${NC}"
        echo ""
        echo -e "${BLUE}To enable reflection for testing:${NC}"
        echo "  make run-aip-test"
        echo "  make test-aip-with-reflection"
    fi
}

# Main execution
main() {
    check_grpcurl
    
    if test_reflection_availability; then
        list_services
        test_service_methods
        test_method_calls
        check_health_reflection_status
        test_security_implications
    else
        echo -e "\n${BLUE}INFO  Reflection is disabled${NC}"
        echo "This is the recommended production setting."
        echo ""
        echo "To test with reflection enabled:"
        echo "  make run-aip-test"
        echo "  make test-grpc-reflection"
        
        check_health_reflection_status
        test_security_implications
    fi
    
    generate_summary
}

# Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
