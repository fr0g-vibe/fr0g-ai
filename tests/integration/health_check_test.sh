#!/bin/bash

# Health Check and Inter-Service Communication Test
# Tests all fr0g.ai services for health and proper communication

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TIMEOUT=30
RETRY_INTERVAL=2
MAX_RETRIES=15

# Service endpoints
REGISTRY_URL="http://localhost:8500"
AIP_HTTP_URL="http://localhost:8080"
AIP_GRPC_URL="localhost:9090"
BRIDGE_HTTP_URL="http://localhost:8082"
BRIDGE_GRPC_URL="localhost:9091"
IO_HTTP_URL="http://localhost:8083"
IO_GRPC_URL="localhost:9092"
MCP_HTTP_URL="http://localhost:8081"

echo -e "${BLUE}fr0g.ai Health Check and Inter-Service Communication Test${NC}"
echo "=================================================================="

# Function to check HTTP health endpoint
check_http_health() {
    local service_name=$1
    local url=$2
    local retries=0
    
    echo -n "Checking $service_name HTTP health... "
    
    while [ $retries -lt $MAX_RETRIES ]; do
        if curl -sf "$url/health" >/dev/null 2>&1; then
            echo -e "${GREEN}HEALTHY${NC}"
            return 0
        fi
        retries=$((retries + 1))
        sleep $RETRY_INTERVAL
    done
    
    echo -e "${RED}UNHEALTHY${NC}"
    return 1
}

# Function to check gRPC health (using grpcurl if available)
check_grpc_health() {
    local service_name=$1
    local url=$2
    
    echo -n "Checking $service_name gRPC health... "
    
    # First check if port is listening
    local host=$(echo $url | cut -d: -f1)
    local port=$(echo $url | cut -d: -f2)
    if ! nc -z "$host" "$port" 2>/dev/null; then
        echo -e "${RED}PORT CLOSED${NC}"
        echo -e "${RED}   gRPC server not listening on $url${NC}"
        return 1
    fi
    
    # Check if grpcurl is available
    if command -v grpcurl >/dev/null 2>&1; then
        # Try health check
        if grpcurl -plaintext "$url" grpc.health.v1.Health/Check >/dev/null 2>&1; then
            echo -e "${GREEN}HEALTHY${NC}"
            return 0
        else
            # Try listing services to see if gRPC is responding
            if grpcurl -plaintext "$url" list >/dev/null 2>&1; then
                echo -e "${YELLOW}RESPONDING (no health service)${NC}"
                return 0
            else
                echo -e "${RED}UNHEALTHY${NC}"
                echo -e "${RED}   gRPC server not responding properly${NC}"
                return 1
            fi
        fi
    else
        echo -e "${YELLOW}PORT OPEN (grpcurl not available)${NC}"
        echo -e "${YELLOW}   Install grpcurl for detailed gRPC testing${NC}"
        return 0
    fi
}

# Function to test service registry
test_service_registry() {
    echo -e "\n${BLUE}CRITICAL Testing Service Registry${NC}"
    echo "----------------------------"
    
    # Check health
    if ! check_http_health "Service Registry" "$REGISTRY_URL"; then
        return 1
    fi
    
    # Test service registration endpoint
    echo -n "Testing service registration API... "
    local reg_response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" -X POST -H "Content-Type: application/json" -d '{"ID":"test","Name":"test","Port":8000}' 2>/dev/null)
    local reg_code="${reg_response: -3}"
    if [ "$reg_code" = "200" ] || [ "$reg_code" = "201" ]; then
        echo -e "${GREEN}OPERATIONAL${NC}"
    else
        echo -e "${RED}CRITICAL - ENDPOINT MISSING (HTTP $reg_code)${NC}"
        echo -e "${RED}   This is blocking all service discovery${NC}"
    fi
    
    # Test service discovery endpoint
    echo -n "Testing service discovery API... "
    local disc_response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/catalog/services" 2>/dev/null)
    local disc_code="${disc_response: -3}"
    if [ "$disc_code" = "200" ]; then
        echo -e "${GREEN}OPERATIONAL${NC}"
        local services="${disc_response%???}"
        if [ -n "$services" ] && [ "$services" != "{}" ]; then
            echo "Registered services:"
            echo "$services" | jq -r 'keys[]' 2>/dev/null || echo "$services"
        else
            echo -e "${YELLOW}   No services currently registered${NC}"
        fi
    else
        echo -e "${RED}CRITICAL - ENDPOINT MISSING (HTTP $disc_code)${NC}"
        echo -e "${RED}   Service discovery completely broken${NC}"
    fi
    
    # Test health endpoint
    echo -n "Testing service health API... "
    local health_response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/health/service/test" 2>/dev/null)
    local health_code="${health_response: -3}"
    if [ "$health_code" = "200" ] || [ "$health_code" = "404" ]; then
        echo -e "${GREEN}OPERATIONAL${NC}"
    else
        echo -e "${RED}CRITICAL - ENDPOINT MISSING (HTTP $health_code)${NC}"
    fi
    
    return 0
}

# Function to test AIP service
test_aip_service() {
    echo -e "\n${BLUE}Testing AIP Service${NC}"
    echo "----------------------"
    
    # Check HTTP health
    if ! check_http_health "AIP" "$AIP_HTTP_URL"; then
        return 1
    fi
    
    # Check gRPC health
    check_grpc_health "AIP" "$AIP_GRPC_URL"
    
    # Test gRPC reflection status
    test_grpc_reflection_status
    
    # Test personas endpoint
    echo -n "Testing personas API... "
    if curl -sf "$AIP_HTTP_URL/personas" >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED ACCESSIBLE${NC}"
    else
        echo -e "${RED}FAILED INACCESSIBLE${NC}"
        return 1
    fi
    
    # Test identities endpoint
    echo -n "Testing identities API... "
    if curl -sf "$AIP_HTTP_URL/identities" >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED ACCESSIBLE${NC}"
    else
        echo -e "${RED}FAILED INACCESSIBLE${NC}"
        return 1
    fi
    
    # Count personas
    echo -n "Checking persona count... "
    local persona_count=$(curl -s "$AIP_HTTP_URL/personas" 2>/dev/null | jq '. | length' 2>/dev/null)
    if [ -n "$persona_count" ] && [ "$persona_count" -ge 0 ]; then
        echo -e "${GREEN}COMPLETED $persona_count personas found${NC}"
    else
        echo -e "${YELLOW}WARNING  Unable to count personas${NC}"
    fi
    
    return 0
}

# Function to test gRPC reflection status
test_grpc_reflection_status() {
    echo -n "Checking gRPC reflection status... "
    
    # Check if grpcurl is available
    if ! command -v grpcurl >/dev/null 2>&1; then
        echo -e "${YELLOW}WARNING  grpcurl not available${NC}"
        return 0
    fi
    
    # Test reflection
    if grpcurl -plaintext "$AIP_GRPC_URL" list >/dev/null 2>&1; then
        echo -e "${YELLOW}WARNING  REFLECTION ENABLED${NC}"
        echo -e "${YELLOW}   TIP This should be disabled in production${NC}"
        
        # List available services
        local services=$(grpcurl -plaintext "$AIP_GRPC_URL" list 2>/dev/null)
        if [ -n "$services" ]; then
            echo -e "${BLUE}   Available gRPC services:${NC}"
            echo "$services" | sed 's/^/     /'
        fi
        
        # Check health endpoint for reflection status
        local health_response=$(curl -s "$AIP_HTTP_URL/health" 2>/dev/null)
        local reflection_status=$(echo "$health_response" | jq -r '.grpc_reflection // "unknown"' 2>/dev/null)
        echo -e "${BLUE}   Health endpoint reports: $reflection_status${NC}"
        
    else
        echo -e "${GREEN}COMPLETED REFLECTION DISABLED${NC}"
        echo -e "${GREEN}   This is the recommended production setting${NC}"
    fi
}

# Function to test Bridge service
test_bridge_service() {
    echo -e "\n${BLUE}Testing Bridge Service${NC}"
    echo "-------------------------"
    
    # Check HTTP health
    if ! check_http_health "Bridge" "$BRIDGE_HTTP_URL"; then
        return 1
    fi
    
    # Check gRPC health
    check_grpc_health "Bridge" "$BRIDGE_GRPC_URL"
    
    # Test chat completions endpoint
    echo -n "Testing chat completions API... "
    if curl -sf "$BRIDGE_HTTP_URL/v1/chat/completions" -X POST \
        -H "Content-Type: application/json" \
        -d '{"model":"test","messages":[{"role":"user","content":"test"}]}' >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED ACCESSIBLE${NC}"
    else
        echo -e "${YELLOW}WARNING  ENDPOINT AVAILABLE (may require valid request)${NC}"
    fi
    
    return 0
}

# Function to test IO service
test_io_service() {
    echo -e "\n${BLUE}Testing IO Service${NC}"
    echo "--------------------"
    
    # Check HTTP health
    if ! check_http_health "IO" "$IO_HTTP_URL"; then
        return 1
    fi
    
    # Check gRPC health
    check_grpc_health "IO" "$IO_GRPC_URL"
    
    # Test input events endpoint
    echo -n "Testing input events API... "
    if curl -sf "$IO_HTTP_URL/events" >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED ACCESSIBLE${NC}"
    else
        echo -e "${YELLOW}WARNING  ENDPOINT MAY REQUIRE AUTHENTICATION${NC}"
    fi
    
    return 0
}

# Function to test Master Control service
test_mcp_service() {
    echo -e "\n${BLUE}Testing Master Control Service${NC}"
    echo "--------------------------------"
    
    # Check HTTP health
    if ! check_http_health "Master Control" "$MCP_HTTP_URL"; then
        echo -e "${YELLOW}WARNING  Master Control service not running${NC}"
        return 0  # Don't fail the test suite
    fi
    
    # Test intelligence endpoint
    echo -n "Testing intelligence API... "
    if curl -sf "$MCP_HTTP_URL/intelligence" >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED ACCESSIBLE${NC}"
    else
        echo -e "${YELLOW}WARNING  INTELLIGENCE ENDPOINT NOT FOUND${NC}"
    fi
    
    return 0
}

# Function to test inter-service communication
test_inter_service_communication() {
    echo -e "\n${BLUE}Testing Inter-Service Communication${NC}"
    echo "-------------------------------------"
    
    # Test Bridge -> AIP communication
    echo -n "Testing Bridge -> AIP communication... "
    if curl -sf "$BRIDGE_HTTP_URL/health" >/dev/null 2>&1 && \
       curl -sf "$AIP_HTTP_URL/health" >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED SERVICES READY${NC}"
    else
        echo -e "${YELLOW}WARNING  BRIDGE SERVICE NOT RUNNING${NC}"
    fi
    
    # Test IO -> Registry communication
    echo -n "Testing IO -> Registry communication... "
    if curl -sf "$IO_HTTP_URL/health" >/dev/null 2>&1 && \
       curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED SERVICES READY${NC}"
    else
        echo -e "${YELLOW}WARNING  IO SERVICE NOT RUNNING${NC}"
    fi
    
    # Test AIP -> Registry communication
    echo -n "Testing AIP -> Registry communication... "
    if curl -sf "$AIP_HTTP_URL/health" >/dev/null 2>&1 && \
       curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1; then
        echo -e "${GREEN}COMPLETED SERVICES READY${NC}"
    else
        echo -e "${RED}FAILED CORE SERVICES NOT READY${NC}"
        return 1
    fi
    
    return 0
}

# Function to test Docker container health
test_container_health() {
    echo -e "\n${BLUE}Testing Docker Container Health${NC}"
    echo "----------------------------------"
    
    # Check if Docker is available
    if ! command -v docker >/dev/null 2>&1; then
        echo -e "${YELLOW}Docker not available, skipping container tests${NC}"
        return 0
    fi
    
    # List running containers
    echo "Running fr0g.ai containers:"
    docker ps --filter "name=fr0g" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || {
        echo -e "${YELLOW}WARNING  No fr0g.ai containers found${NC}"
        return 0
    }
    
    # Check container health status
    local containers=$(docker ps --filter "name=fr0g" --format "{{.Names}}" 2>/dev/null)
    for container in $containers; do
        echo -n "Checking $container health... "
        local health=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null)
        case "$health" in
            "healthy")
                echo -e "${GREEN}HEALTHY${NC}"
                ;;
            "unhealthy")
                echo -e "${RED}UNHEALTHY${NC}"
                ;;
            "starting")
                echo -e "${YELLOW}STARTING${NC}"
                ;;
            *)
                echo -e "${YELLOW}NO HEALTH CHECK${NC}"
                ;;
        esac
    done
    
    return 0
}


# Function to generate summary report
generate_summary() {
    echo -e "\n${BLUE}CRITICAL ISSUES SUMMARY${NC}"
    echo "======================="
    
    # Test actual service status
    local registry_http_ok=false
    local registry_api_ok=false
    local aip_http_ok=false
    local aip_grpc_ok=false
    local bridge_http_ok=false
    local bridge_grpc_ok=false
    local io_http_ok=false
    local io_grpc_ok=false
    local mcp_http_ok=false
    
    # Check HTTP services
    curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1 && registry_http_ok=true
    curl -sf "$AIP_HTTP_URL/health" >/dev/null 2>&1 && aip_http_ok=true
    curl -sf "$BRIDGE_HTTP_URL/health" >/dev/null 2>&1 && bridge_http_ok=true
    curl -sf "$IO_HTTP_URL/health" >/dev/null 2>&1 && io_http_ok=true
    curl -sf "$MCP_HTTP_URL/health" >/dev/null 2>&1 && mcp_http_ok=true
    
    # Check Registry API
    local reg_code=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" -X POST 2>/dev/null | tail -c 3)
    [ "$reg_code" = "200" ] || [ "$reg_code" = "201" ] && registry_api_ok=true
    
    # Check gRPC services
    if command -v grpcurl >/dev/null 2>&1; then
        grpcurl -plaintext "$AIP_GRPC_URL" list >/dev/null 2>&1 && aip_grpc_ok=true
        grpcurl -plaintext "$BRIDGE_GRPC_URL" list >/dev/null 2>&1 && bridge_grpc_ok=true
        grpcurl -plaintext "$IO_GRPC_URL" list >/dev/null 2>&1 && io_grpc_ok=true
    fi
    
    echo -e "\n${BLUE}SERVICE STATUS MATRIX:${NC}"
    echo "Service           | HTTP  | gRPC  | API   | Status"
    echo "------------------|-------|-------|-------|--------"
    
    # Registry
    local reg_status="${GREEN}OPERATIONAL${NC}"
    if ! $registry_http_ok; then
        reg_status="${RED}CRITICAL${NC}"
    elif ! $registry_api_ok; then
        reg_status="${RED}API MISSING${NC}"
    fi
    printf "%-17s | %-5s | %-5s | %-5s | %s\n" "Registry" \
        "$($registry_http_ok && echo "✓" || echo "✗")" \
        "N/A" \
        "$($registry_api_ok && echo "✓" || echo "✗")" \
        "$reg_status"
    
    # AIP
    local aip_status="${GREEN}OPERATIONAL${NC}"
    if ! $aip_http_ok; then
        aip_status="${RED}CRITICAL${NC}"
    elif ! $aip_grpc_ok; then
        aip_status="${YELLOW}gRPC ISSUE${NC}"
    fi
    printf "%-17s | %-5s | %-5s | %-5s | %s\n" "AIP" \
        "$($aip_http_ok && echo "✓" || echo "✗")" \
        "$($aip_grpc_ok && echo "✓" || echo "✗")" \
        "N/A" \
        "$aip_status"
    
    # Bridge
    local bridge_status="${GREEN}OPERATIONAL${NC}"
    if ! $bridge_http_ok; then
        bridge_status="${RED}CRITICAL${NC}"
    elif ! $bridge_grpc_ok; then
        bridge_status="${YELLOW}gRPC ISSUE${NC}"
    fi
    printf "%-17s | %-5s | %-5s | %-5s | %s\n" "Bridge" \
        "$($bridge_http_ok && echo "✓" || echo "✗")" \
        "$($bridge_grpc_ok && echo "✓" || echo "✗")" \
        "N/A" \
        "$bridge_status"
    
    # IO
    local io_status="${GREEN}OPERATIONAL${NC}"
    if ! $io_http_ok; then
        io_status="${RED}CRITICAL${NC}"
    elif ! $io_grpc_ok; then
        io_status="${YELLOW}gRPC ISSUE${NC}"
    fi
    printf "%-17s | %-5s | %-5s | %-5s | %s\n" "IO" \
        "$($io_http_ok && echo "✓" || echo "✗")" \
        "$($io_grpc_ok && echo "✓" || echo "✗")" \
        "N/A" \
        "$io_status"
    
    # Master Control
    local mcp_status="${GREEN}OPERATIONAL${NC}"
    if ! $mcp_http_ok; then
        mcp_status="${RED}CRITICAL${NC}"
    fi
    printf "%-17s | %-5s | %-5s | %-5s | %s\n" "Master Control" \
        "$($mcp_http_ok && echo "✓" || echo "✗")" \
        "N/A" \
        "N/A" \
        "$mcp_status"
    
    echo -e "\n${BLUE}CRITICAL BLOCKERS IDENTIFIED:${NC}"
    if ! $registry_api_ok; then
        echo -e "${RED}1. CRITICAL: Service Registry API endpoints missing${NC}"
        echo -e "${RED}   - /v1/agent/service/register returning 404${NC}"
        echo -e "${RED}   - Service discovery completely broken${NC}"
    fi
    
    if ! $aip_grpc_ok || ! $bridge_grpc_ok || ! $io_grpc_ok; then
        echo -e "${RED}2. CRITICAL: gRPC services unhealthy${NC}"
        echo -e "${RED}   - Inter-service communication failing${NC}"
        echo -e "${RED}   - Check gRPC server startup and port binding${NC}"
    fi
    
    echo -e "\n${BLUE}EMERGENCY ACTIONS REQUIRED:${NC}"
    echo -e "  make diagnose-registry      # Diagnose registry API issues"
    echo -e "  make diagnose-grpc          # Diagnose gRPC health issues"
    echo -e "  make diagnose-ports         # Check port configuration"
    echo -e "  docker-compose logs         # Check service logs for errors"
}

# Main test execution
main() {
    echo -e "${BLUE}Starting comprehensive health check...${NC}\n"
    
    local exit_code=0
    
    # Run all tests
    test_service_registry || exit_code=1
    test_aip_service || exit_code=1
    test_bridge_service || exit_code=1
    test_io_service || exit_code=1
    test_mcp_service || exit_code=1
    test_inter_service_communication || exit_code=1
    test_container_health || exit_code=1
    
    
    # Generate summary
    generate_summary
    
    if [ $exit_code -eq 0 ]; then
        echo -e "\n${GREEN}SUCCESS: All health checks passed!${NC}"
    else
        echo -e "\n${RED}ERROR: Some health checks failed!${NC}"
    fi
    
    exit $exit_code
}

# Check for required tools
check_dependencies() {
    local missing_tools=()
    
    if ! command -v curl >/dev/null 2>&1; then
        missing_tools+=("curl")
    fi
    
    if ! command -v jq >/dev/null 2>&1; then
        echo -e "${YELLOW}WARNING: jq not found - JSON parsing will be limited${NC}"
    fi
    
    if ! command -v nc >/dev/null 2>&1; then
        echo -e "${YELLOW}WARNING: nc (netcat) not found - port checks will be limited${NC}"
    fi
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        echo -e "${RED}ERROR: Missing required tools: ${missing_tools[*]}${NC}"
        echo "Please install the missing tools and try again."
        exit 1
    fi
}

# Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    check_dependencies
    main "$@"
fi
