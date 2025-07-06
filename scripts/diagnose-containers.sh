#!/bin/bash

echo "DIAGNOSE Container Health and gRPC Status"
echo "========================================"

# Function to check container health
check_container_health() {
    local container_name=$1
    local http_port=$2
    local grpc_port=$3
    
    echo ""
    echo "DIAGNOSE Checking $container_name..."
    
    # Container status
    status=$(docker inspect --format='{{.State.Health.Status}}' $container_name 2>/dev/null || echo "not_found")
    echo "  Container health: $status"
    
    # Process check
    if [ "$status" != "not_found" ]; then
        echo "  Running processes:"
        docker exec $container_name ps aux 2>/dev/null | grep -E "(fr0g|main)" || echo "    No fr0g processes found"
        
        # Port binding check
        if [ ! -z "$http_port" ]; then
            docker exec $container_name netstat -tlnp 2>/dev/null | grep ":$http_port" && echo "  HTTP port $http_port: BOUND" || echo "  HTTP port $http_port: NOT BOUND"
        fi
        
        if [ ! -z "$grpc_port" ]; then
            docker exec $container_name netstat -tlnp 2>/dev/null | grep ":$grpc_port" && echo "  gRPC port $grpc_port: BOUND" || echo "  gRPC port $grpc_port: NOT BOUND"
        fi
        
        # Recent logs
        echo "  Recent logs (last 5 lines):"
        docker logs --tail=5 $container_name 2>&1 | sed 's/^/    /'
    fi
}

# Check all fr0g-ai containers
check_container_health "fr0g-ai-service-registry-1" "8500" ""
check_container_health "fr0g-ai-fr0g-ai-aip-1" "8080" "9090"
check_container_health "fr0g-ai-fr0g-ai-bridge-1" "8082" "9091"
check_container_health "fr0g-ai-fr0g-ai-master-control-1" "8081" ""
check_container_health "fr0g-ai-fr0g-ai-io-1" "8083" "9092"

echo ""
echo "DIAGNOSE Network connectivity test..."
echo "===================================="

# Test inter-container connectivity
containers=("fr0g-ai-fr0g-ai-aip-1" "fr0g-ai-fr0g-ai-bridge-1" "fr0g-ai-fr0g-ai-io-1")
for container in "${containers[@]}"; do
    if docker exec $container nc -z service-registry 8500 2>/dev/null; then
        echo "  $container -> service-registry:8500: CONNECTED"
    else
        echo "  $container -> service-registry:8500: FAILED"
    fi
done

echo ""
echo "COMPLETED Container diagnostics complete"
