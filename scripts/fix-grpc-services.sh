#!/bin/bash

# gRPC Services Fix Script
# Addresses gRPC connectivity and protobuf issues

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}fr0g.ai gRPC Services Fix${NC}"
echo "=========================="

# Function to check container logs for specific issues
check_container_issues() {
    local container_name=$1
    local service_name=$2
    
    echo -e "\n${BLUE}Checking $service_name container issues...${NC}"
    
    # Check for protobuf issues
    local protobuf_issues=$(docker logs "$container_name" 2>&1 | grep -i "protobuf\|proto\|grpc.*disabled" || true)
    if [ -n "$protobuf_issues" ]; then
        echo -e "${RED}PROTOBUF ISSUES DETECTED:${NC}"
        echo "$protobuf_issues"
    fi
    
    # Check for gRPC startup issues
    local grpc_issues=$(docker logs "$container_name" 2>&1 | grep -i "grpc.*error\|grpc.*failed\|bind.*failed" || true)
    if [ -n "$grpc_issues" ]; then
        echo -e "${RED}gRPC STARTUP ISSUES:${NC}"
        echo "$grpc_issues"
    fi
    
    # Check for port binding issues
    local port_issues=$(docker logs "$container_name" 2>&1 | grep -i "address.*in use\|bind.*failed\|port.*failed" || true)
    if [ -n "$port_issues" ]; then
        echo -e "${RED}PORT BINDING ISSUES:${NC}"
        echo "$port_issues"
    fi
    
    # Check if gRPC server actually started
    local grpc_started=$(docker logs "$container_name" 2>&1 | grep -i "grpc.*listening\|grpc.*started" || true)
    if [ -n "$grpc_started" ]; then
        echo -e "${GREEN}gRPC SERVER STARTED:${NC}"
        echo "$grpc_started"
    else
        echo -e "${YELLOW}NO gRPC STARTUP CONFIRMATION FOUND${NC}"
    fi
}

# Function to test gRPC connectivity from inside container
test_internal_grpc() {
    local container_name=$1
    local port=$2
    local service_name=$3
    
    echo -e "\n${BLUE}Testing internal gRPC connectivity for $service_name...${NC}"
    
    # Test if gRPC port is listening inside container
    echo -n "Internal port check... "
    if docker exec "$container_name" sh -c "nc -z localhost $port" 2>/dev/null; then
        echo -e "${GREEN}PORT LISTENING${NC}"
    else
        echo -e "${RED}PORT NOT LISTENING${NC}"
        echo "  gRPC server may not be properly bound inside container"
    fi
    
    # Check process listening on port
    echo -n "Process check... "
    local process_info=$(docker exec "$container_name" sh -c "netstat -tlnp 2>/dev/null | grep :$port" || true)
    if [ -n "$process_info" ]; then
        echo -e "${GREEN}PROCESS FOUND${NC}"
        echo "  $process_info"
    else
        echo -e "${RED}NO PROCESS LISTENING${NC}"
    fi
}

# Function to restart specific service
restart_service() {
    local service_name=$1
    
    echo -e "\n${BLUE}Restarting $service_name...${NC}"
    
    docker-compose restart "$service_name"
    sleep 10
    
    # Check if service is healthy after restart
    local container_name="fr0g-ai-${service_name}-1"
    local health=$(docker inspect --format='{{.State.Health.Status}}' "$container_name" 2>/dev/null || echo "no-health-check")
    
    echo -n "Service health after restart... "
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
}

# Function to generate protobuf files
generate_protobuf() {
    echo -e "\n${BLUE}Generating protobuf files...${NC}"
    
    # Check if protoc is available
    if command -v protoc >/dev/null 2>&1; then
        echo -e "${GREEN}protoc available${NC}"
        make proto || echo -e "${YELLOW}Protobuf generation failed - this may be expected${NC}"
    else
        echo -e "${YELLOW}protoc not available - skipping protobuf generation${NC}"
        echo "Install with: apt-get install protobuf-compiler (Ubuntu/Debian)"
        echo "Or: brew install protobuf (macOS)"
    fi
}

main() {
    echo -e "${BLUE}Starting gRPC services diagnostic and fix...${NC}\n"
    
    # Check if Docker is available
    if ! command -v docker >/dev/null 2>&1; then
        echo -e "${RED}ERROR: Docker not available${NC}"
        exit 1
    fi
    
    # Generate protobuf files first
    generate_protobuf
    
    # Check each service
    local services=("fr0g-ai-aip" "fr0g-ai-bridge" "fr0g-ai-io")
    local ports=(9090 9091 9092)
    
    for i in "${!services[@]}"; do
        local service="${services[$i]}"
        local port="${ports[$i]}"
        local container_name="fr0g-ai-${service}-1"
        
        echo -e "\n${BLUE}=== Analyzing $service ===${NC}"
        
        # Check if container is running
        if docker ps --format "{{.Names}}" | grep -q "$container_name"; then
            check_container_issues "$container_name" "$service"
            test_internal_grpc "$container_name" "$port" "$service"
            
            # If gRPC is not working, try restarting
            if ! docker exec "$container_name" sh -c "nc -z localhost $port" 2>/dev/null; then
                echo -e "${YELLOW}gRPC not responding, attempting restart...${NC}"
                restart_service "$service"
            fi
        else
            echo -e "${RED}Container $container_name not running${NC}"
        fi
    done
    
    # Final connectivity test
    echo -e "\n${BLUE}Final gRPC connectivity test...${NC}"
    sleep 5
    
    for i in "${!services[@]}"; do
        local service="${services[$i]}"
        local port="${ports[$i]}"
        
        echo -n "Testing $service gRPC... "
        if timeout 5 grpcurl -plaintext localhost:"$port" list >/dev/null 2>&1; then
            echo -e "${GREEN}RESPONDING${NC}"
        else
            echo -e "${RED}NOT RESPONDING${NC}"
        fi
    done
    
    echo -e "\n${BLUE}gRPC fix attempt completed${NC}"
    echo -e "${BLUE}Run 'make diagnose' to verify improvements${NC}"
}

# Check dependencies
if ! command -v nc >/dev/null 2>&1; then
    echo -e "${YELLOW}WARNING: nc (netcat) not found - some tests will be limited${NC}"
fi

if ! command -v grpcurl >/dev/null 2>&1; then
    echo -e "${YELLOW}WARNING: grpcurl not found - final test will be limited${NC}"
fi

main "$@"
