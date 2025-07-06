#!/bin/bash

# gRPC Service Health Check Script
# Comprehensive testing of all gRPC endpoints

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# gRPC Services and Ports
SERVICES=(
    "fr0g-ai-aip:9090"
    "fr0g-ai-bridge:9091"
    "fr0g-ai-io:9092"
)

echo -e "${BLUE}gRPC Service Health Check${NC}"
echo "========================="

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}CRITICAL: grpcurl not installed${NC}"
    echo "Install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Function to test gRPC connectivity
test_grpc_connectivity() {
    local service_name=$1
    local port=$2
    
    echo -e "\n${BLUE}Testing $service_name (port $port)${NC}"
    echo "----------------------------------------"
    
    # Test port connectivity
    echo -n "Port connectivity... "
    if nc -z -w5 localhost "$port" &> /dev/null; then
        echo -e "${GREEN}CONNECTED${NC}"
    else
        echo -e "${RED}DISCONNECTED${NC}"
        return 1
    fi
    
    # Test gRPC reflection
    echo -n "gRPC reflection... "
    if grpcurl -plaintext localhost:"$port" list &> /dev/null; then
        echo -e "${GREEN}AVAILABLE${NC}"
        
        # List available services
        echo "Available services:"
        grpcurl -plaintext localhost:"$port" list | sed 's/^/  - /'
        
    else
        echo -e "${YELLOW}NOT AVAILABLE${NC}"
        echo "  This may be expected if reflection is disabled"
    fi
    
    # Test health check service
    echo -n "Health check service... "
    if grpcurl -plaintext localhost:"$port" grpc.health.v1.Health/Check &> /dev/null; then
        echo -e "${GREEN}RESPONDING${NC}"
    else
        echo -e "${YELLOW}NOT AVAILABLE${NC}"
        echo "  Health check service not implemented"
    fi
    
    # Test basic connectivity with timeout
    echo -n "Basic gRPC connectivity... "
    if timeout 5 grpcurl -plaintext -max-time 3 localhost:"$port" list &> /dev/null; then
        echo -e "${GREEN}RESPONDING${NC}"
    else
        echo -e "${RED}NOT RESPONDING${NC}"
        echo "  gRPC server may not be running or accepting connections"
        return 1
    fi
    
    return 0
}

# Function to diagnose gRPC issues
diagnose_grpc_issues() {
    echo -e "\n${BLUE}gRPC Diagnostic Information${NC}"
    echo "============================"
    
    # Check Docker container gRPC processes
    echo -e "\n${BLUE}Checking gRPC processes in containers:${NC}"
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r service_name port <<< "$service"
        container_name="fr0g-ai-${service_name}-1"
        
        echo -n "Checking $service_name container processes... "
        if docker exec "$container_name" ps aux 2>/dev/null | grep -q "fr0g-ai"; then
            echo -e "${GREEN}PROCESS RUNNING${NC}"
            # Show the actual process
            docker exec "$container_name" ps aux | grep "fr0g-ai" | head -1
        else
            echo -e "${RED}NO PROCESS FOUND${NC}"
        fi
    done
    
    # Check container logs for gRPC startup messages
    echo -e "\n${BLUE}Checking container logs for gRPC startup:${NC}"
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r service_name port <<< "$service"
        container_name="fr0g-ai-${service_name}-1"
        
        echo -e "\n${YELLOW}$service_name logs:${NC}"
        docker logs "$container_name" 2>&1 | grep -i "grpc\|listening\|server\|port" | tail -5 || echo "  No relevant logs found"
    done
}

# Function to test service registration
test_service_registration() {
    echo -e "\n${BLUE}Testing Service Registration${NC}"
    echo "============================="
    
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r service_name port <<< "$service"
        
        echo -n "Checking $service_name registration... "
        if curl -sf "http://localhost:8500/v1/catalog/service/$service_name" >/dev/null 2>&1; then
            echo -e "${GREEN}REGISTERED${NC}"
        else
            echo -e "${RED}NOT REGISTERED${NC}"
            echo "  Service should auto-register with registry on startup"
        fi
    done
}

main() {
    local exit_code=0
    
    echo -e "${BLUE}Starting comprehensive gRPC health check...${NC}\n"
    
    # Test each gRPC service
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r service_name port <<< "$service"
        if ! test_grpc_connectivity "$service_name" "$port"; then
            exit_code=1
        fi
    done
    
    # Run diagnostics
    diagnose_grpc_issues
    
    # Test service registration
    test_service_registration
    
    # Summary
    echo -e "\n${BLUE}gRPC Health Check Summary${NC}"
    echo "=========================="
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}SUCCESS: All gRPC services are healthy${NC}"
    else
        echo -e "${RED}CRITICAL: Some gRPC services have issues${NC}"
        echo -e "\n${BLUE}Troubleshooting Steps:${NC}"
        echo "1. Check if gRPC servers are actually starting in containers"
        echo "2. Verify gRPC server configuration and port binding"
        echo "3. Check for any startup errors in container logs"
        echo "4. Ensure gRPC servers are binding to 0.0.0.0, not localhost"
        echo "5. Verify protobuf definitions are properly compiled"
    fi
    
    exit $exit_code
}

# Check dependencies
if ! command -v nc >/dev/null 2>&1; then
    echo -e "${YELLOW}WARNING: nc (netcat) not found - some tests will be limited${NC}"
fi

if ! command -v docker >/dev/null 2>&1; then
    echo -e "${YELLOW}WARNING: docker not found - container diagnostics will be limited${NC}"
fi

# Run the test
main "$@"
