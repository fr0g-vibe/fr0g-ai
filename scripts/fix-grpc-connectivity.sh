#!/bin/bash

# gRPC Connectivity Fix Script
# Addresses gRPC protocol and service implementation issues

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}fr0g.ai gRPC Connectivity Fix${NC}"
echo "============================="

# Function to test gRPC with different protocols
test_grpc_protocols() {
    local service_name=$1
    local port=$2
    
    echo -e "\n${BLUE}Testing $service_name gRPC protocols...${NC}"
    
    # Test with different grpcurl options
    echo -n "Testing plaintext gRPC... "
    if timeout 10 grpcurl -plaintext -connect-timeout 5 localhost:"$port" list >/dev/null 2>&1; then
        echo -e "${GREEN}SUCCESS${NC}"
        return 0
    else
        echo -e "${RED}FAILED${NC}"
    fi
    
    echo -n "Testing with explicit HTTP/2... "
    if timeout 10 grpcurl -plaintext -proto-set-out /dev/null localhost:"$port" list >/dev/null 2>&1; then
        echo -e "${GREEN}SUCCESS${NC}"
        return 0
    else
        echo -e "${RED}FAILED${NC}"
    fi
    
    echo -n "Testing with verbose output... "
    local verbose_output=$(timeout 10 grpcurl -plaintext -v localhost:"$port" list 2>&1 || true)
    if echo "$verbose_output" | grep -q "connection refused\|timeout\|no such host"; then
        echo -e "${RED}CONNECTION ISSUE${NC}"
        echo "  $verbose_output"
    elif echo "$verbose_output" | grep -q "rpc error\|unimplemented"; then
        echo -e "${YELLOW}PROTOCOL ISSUE${NC}"
        echo "  $verbose_output"
    else
        echo -e "${YELLOW}UNKNOWN ISSUE${NC}"
        echo "  $verbose_output"
    fi
    
    return 1
}

# Function to check gRPC service implementation
check_grpc_implementation() {
    local service_name=$1
    local container_name="fr0g-ai-${service_name}-1"
    
    echo -e "\n${BLUE}Checking $service_name gRPC implementation...${NC}"
    
    # Check if gRPC server is actually implementing services
    echo -n "Checking gRPC service registration... "
    local grpc_logs=$(docker logs "$container_name" 2>&1 | grep -i "grpc.*service\|register.*service\|grpc.*handler" || true)
    if [ -n "$grpc_logs" ]; then
        echo -e "${GREEN}SERVICES REGISTERED${NC}"
        echo "$grpc_logs" | head -3
    else
        echo -e "${RED}NO SERVICES REGISTERED${NC}"
        echo "  gRPC server may be running but not serving any services"
    fi
    
    # Check for gRPC errors
    echo -n "Checking for gRPC errors... "
    local grpc_errors=$(docker logs "$container_name" 2>&1 | grep -i "grpc.*error\|grpc.*fail\|grpc.*panic" || true)
    if [ -n "$grpc_errors" ]; then
        echo -e "${RED}ERRORS FOUND${NC}"
        echo "$grpc_errors" | head -3
    else
        echo -e "${GREEN}NO ERRORS${NC}"
    fi
    
    # Test internal gRPC connectivity with different tools
    echo -n "Testing internal gRPC with nc... "
    if docker exec "$container_name" sh -c "echo 'test' | nc localhost 9090" 2>/dev/null; then
        echo -e "${GREEN}ACCEPTS CONNECTIONS${NC}"
    else
        echo -e "${RED}REJECTS CONNECTIONS${NC}"
    fi
}

# Function to rebuild containers with gRPC fixes
rebuild_grpc_services() {
    echo -e "\n${BLUE}Rebuilding gRPC services with fixes...${NC}"
    
    # Rebuild services that need gRPC fixes
    local services=("fr0g-ai-aip" "fr0g-ai-bridge" "fr0g-ai-io")
    
    for service in "${services[@]}"; do
        echo -e "\n${YELLOW}Rebuilding $service...${NC}"
        docker-compose build --no-cache "$service"
    done
    
    echo -e "\n${BLUE}Restarting services...${NC}"
    docker-compose restart fr0g-ai-aip fr0g-ai-bridge fr0g-ai-io
    
    # Wait for services to start
    sleep 30
    
    echo -e "\n${BLUE}Testing after rebuild...${NC}"
    for service in "${services[@]}"; do
        local port
        case "$service" in
            "fr0g-ai-aip") port=9090 ;;
            "fr0g-ai-bridge") port=9091 ;;
            "fr0g-ai-io") port=9092 ;;
        esac
        
        echo -n "Testing $service gRPC... "
        if timeout 10 grpcurl -plaintext localhost:"$port" list >/dev/null 2>&1; then
            echo -e "${GREEN}SUCCESS${NC}"
        else
            echo -e "${RED}STILL FAILING${NC}"
        fi
    done
}

# Function to test alternative gRPC tools
test_alternative_grpc_tools() {
    echo -e "\n${BLUE}Testing with alternative gRPC tools...${NC}"
    
    # Test with curl (HTTP/2)
    echo -n "Testing with curl HTTP/2... "
    if curl -s --http2-prior-knowledge http://localhost:9090 >/dev/null 2>&1; then
        echo -e "${GREEN}HTTP/2 RESPONDING${NC}"
    else
        echo -e "${RED}HTTP/2 NOT RESPONDING${NC}"
    fi
    
    # Test with telnet
    echo -n "Testing with telnet... "
    if timeout 5 telnet localhost 9090 </dev/null >/dev/null 2>&1; then
        echo -e "${GREEN}PORT ACCEPTS TELNET${NC}"
    else
        echo -e "${RED}PORT REJECTS TELNET${NC}"
    fi
}

# Function to check for common gRPC issues
check_common_grpc_issues() {
    echo -e "\n${BLUE}Checking for common gRPC issues...${NC}"
    
    # Check for reflection service
    echo -n "Checking gRPC reflection service... "
    local reflection_test=$(timeout 5 grpcurl -plaintext localhost:9090 grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo 2>&1 || true)
    if echo "$reflection_test" | grep -q "Unimplemented"; then
        echo -e "${YELLOW}REFLECTION NOT IMPLEMENTED${NC}"
        echo "  This is expected if reflection is disabled"
    elif echo "$reflection_test" | grep -q "connection refused"; then
        echo -e "${RED}CONNECTION REFUSED${NC}"
    else
        echo -e "${GREEN}REFLECTION AVAILABLE${NC}"
    fi
    
    # Check for health service
    echo -n "Checking gRPC health service... "
    local health_test=$(timeout 5 grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check 2>&1 || true)
    if echo "$health_test" | grep -q "Unimplemented"; then
        echo -e "${YELLOW}HEALTH SERVICE NOT IMPLEMENTED${NC}"
    elif echo "$health_test" | grep -q "connection refused"; then
        echo -e "${RED}CONNECTION REFUSED${NC}"
    else
        echo -e "${GREEN}HEALTH SERVICE AVAILABLE${NC}"
    fi
    
    # Check for any implemented services
    echo -n "Checking for any gRPC services... "
    local any_service=$(timeout 5 grpcurl -plaintext localhost:9090 list 2>&1 || true)
    if echo "$any_service" | grep -q "Failed to list services"; then
        echo -e "${RED}NO SERVICES AVAILABLE${NC}"
        echo "  gRPC server running but no services registered"
    elif echo "$any_service" | grep -q "connection refused"; then
        echo -e "${RED}CONNECTION REFUSED${NC}"
    else
        echo -e "${GREEN}SERVICES AVAILABLE${NC}"
        echo "$any_service" | head -3
    fi
}

main() {
    echo -e "${BLUE}Starting gRPC connectivity diagnostic and fix...${NC}\n"
    
    # Check if grpcurl is available
    if ! command -v grpcurl >/dev/null 2>&1; then
        echo -e "${RED}ERROR: grpcurl not available${NC}"
        echo "Install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
        exit 1
    fi
    
    # Test each service
    local services=("fr0g-ai-aip:9090" "fr0g-ai-bridge:9091" "fr0g-ai-io:9092")
    local any_working=false
    
    for service in "${services[@]}"; do
        IFS=':' read -r service_name port <<< "$service"
        
        if test_grpc_protocols "$service_name" "$port"; then
            any_working=true
        fi
        
        check_grpc_implementation "$service_name"
    done
    
    # Check for common issues
    check_common_grpc_issues
    
    # Test alternative tools
    test_alternative_grpc_tools
    
    # If nothing is working, try rebuilding
    if ! $any_working; then
        echo -e "\n${YELLOW}No gRPC services responding, attempting rebuild...${NC}"
        rebuild_grpc_services
    fi
    
    echo -e "\n${BLUE}gRPC connectivity fix completed${NC}"
    echo -e "${BLUE}Run 'make diagnose' to verify improvements${NC}"
}

# Check dependencies
if ! command -v docker >/dev/null 2>&1; then
    echo -e "${RED}ERROR: Docker not available${NC}"
    exit 1
fi

main "$@"
