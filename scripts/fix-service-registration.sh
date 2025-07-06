#!/bin/bash

# Service Registration Fix Script
# Manually registers services with the registry to test registration endpoint

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

REGISTRY_URL="http://localhost:8500"

echo -e "${BLUE}fr0g.ai Service Registration Fix${NC}"
echo "================================="

# Function to register a service
register_service() {
    local service_id=$1
    local service_name=$2
    local port=$3
    local health_url=$4
    
    echo -n "Registering $service_name... "
    
    local service_def=$(cat <<EOF
{
    "ID": "$service_id",
    "Name": "$service_name",
    "Address": "localhost",
    "Port": $port,
    "Check": {
        "HTTP": "$health_url",
        "Interval": "10s",
        "Timeout": "5s"
    },
    "Tags": ["fr0g-ai", "microservice"]
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" "$REGISTRY_URL/v1/agent/service/register" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$service_def" 2>/dev/null)
    
    local status_code="${response: -3}"
    local response_body="${response%???}"
    
    if [ "$status_code" = "200" ] || [ "$status_code" = "201" ]; then
        echo -e "${GREEN}SUCCESS${NC}"
        # Wait a moment for registration to propagate
        sleep 2
        return 0
    else
        echo -e "${RED}FAILED (HTTP $status_code)${NC}"
        echo "Response: $response_body"
        echo "Request: $service_def"
        return 1
    fi
}

# Function to verify registration
verify_registration() {
    local service_name=$1
    
    echo -n "Verifying $service_name registration... "
    
    local response=$(curl -s "$REGISTRY_URL/v1/catalog/service/$service_name" 2>/dev/null)
    
    if echo "$response" | jq -e '. | length > 0' >/dev/null 2>&1; then
        echo -e "${GREEN}VERIFIED${NC}"
        return 0
    else
        echo -e "${RED}NOT FOUND${NC}"
        echo "  Catalog response: $response"
        
        # Try alternative endpoints
        local alt_response=$(curl -s "$REGISTRY_URL/services/$service_name" 2>/dev/null)
        if [ -n "$alt_response" ] && [ "$alt_response" != "null" ]; then
            echo -e "${YELLOW}FOUND VIA ALTERNATIVE ENDPOINT${NC}"
            return 0
        fi
        
        return 1
    fi
}

main() {
    echo -e "${BLUE}Attempting to manually register all services...${NC}\n"
    
    local exit_code=0
    
    # Check if registry is available
    if ! curl -sf "$REGISTRY_URL/health" >/dev/null 2>&1; then
        echo -e "${RED}ERROR: Service registry not available at $REGISTRY_URL${NC}"
        exit 1
    fi
    
    # Register each service
    register_service "aip-001" "fr0g-ai-aip" 8080 "http://localhost:8080/health" || exit_code=1
    register_service "bridge-001" "fr0g-ai-bridge" 8082 "http://localhost:8082/health" || exit_code=1
    register_service "io-001" "fr0g-ai-io" 8083 "http://localhost:8083/health" || exit_code=1
    register_service "mcp-001" "fr0g-ai-master-control" 8081 "http://localhost:8081/health" || exit_code=1
    
    echo ""
    
    # Verify registrations
    echo -e "${BLUE}Verifying service registrations...${NC}"
    verify_registration "fr0g-ai-aip" || exit_code=1
    verify_registration "fr0g-ai-bridge" || exit_code=1
    verify_registration "fr0g-ai-io" || exit_code=1
    verify_registration "fr0g-ai-master-control" || exit_code=1
    
    echo ""
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}SUCCESS: All services registered successfully${NC}"
        echo -e "${BLUE}You can now test service discovery with:${NC}"
        echo "  curl http://localhost:8500/v1/catalog/services"
    else
        echo -e "${RED}ERROR: Some services failed to register${NC}"
        echo -e "${BLUE}Check the registry logs for more details:${NC}"
        echo "  docker-compose logs service-registry"
    fi
    
    exit $exit_code
}

# Check dependencies
if ! command -v curl >/dev/null 2>&1; then
    echo -e "${RED}ERROR: curl is required${NC}"
    exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
    echo -e "${YELLOW}WARNING: jq not found - verification will be limited${NC}"
fi

main "$@"
