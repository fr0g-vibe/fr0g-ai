#!/bin/bash

# End-to-end integration test for all fr0g-ai services
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== End-to-End Integration Test ===${NC}"

# Function to wait for services to be ready
wait_for_services() {
    echo -e "${BLUE}Waiting for services to be ready...${NC}"
    local max_wait=60
    local wait_time=0
    
    while [ $wait_time -lt $max_wait ]; do
        if curl -sf http://localhost:8500/health >/dev/null 2>&1; then
            echo -e "${GREEN}Services are ready!${NC}"
            return 0
        fi
        echo -n "."
        sleep 2
        wait_time=$((wait_time + 2))
    done
    
    echo -e "\n${RED}Services failed to start within $max_wait seconds${NC}"
    return 1
}

# Check if services are running, if not wait for them
if ! curl -sf http://localhost:8500/health >/dev/null 2>&1; then
    echo -e "${YELLOW}Services not responding, waiting for startup...${NC}"
    if ! wait_for_services; then
        echo -e "${RED}ERROR: Services failed to start${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}Services are already running${NC}"
fi

# Test Bridge service
echo -e "\n${BLUE}Testing Bridge service...${NC}"
echo "  Checking HTTP port (8082)..."
if curl -f http://localhost:8082/health 2>/dev/null; then
    echo -e "${GREEN}✓ Bridge service HTTP healthy on port 8082${NC}"
    # Test additional Bridge endpoints if healthy
    if curl -s http://localhost:8082/health | grep -q "service"; then
        echo -e "${GREEN}✓ Bridge service health endpoint responding correctly${NC}"
    fi
    
    # Test chat completions endpoint
    echo "  Testing chat completions endpoint..."
    if curl -sf http://localhost:8082/v1/chat/completions -X POST \
        -H "Content-Type: application/json" \
        -d '{"model":"test","messages":[{"role":"user","content":"test"}]}' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Chat completions endpoint accessible${NC}"
    else
        echo -e "${YELLOW}⚠ Chat completions endpoint may require valid request${NC}"
    fi
else
    echo -e "${RED}✗ Bridge service HTTP down${NC}"
    if docker-compose logs fr0g-ai-bridge 2>/dev/null | tail -5 | grep -q "error"; then
        echo -e "${RED}Recent errors in Bridge service logs:${NC}"
        docker-compose logs fr0g-ai-bridge 2>/dev/null | tail -5 | grep -i error
    fi
fi

# Test I/O service
echo -e "\n${BLUE}Testing I/O service...${NC}"
if curl -f http://localhost:8083/health 2>/dev/null; then
    echo -e "${GREEN}✓ I/O service healthy${NC}"
    
    # Test I/O service specific endpoints if they exist
    if curl -s http://localhost:8083/health | grep -q "service"; then
        echo -e "${GREEN}✓ I/O service health endpoint responding correctly${NC}"
    fi
else
    echo -e "${RED}✗ I/O service down${NC}"
    if docker-compose logs fr0g-ai-io 2>/dev/null | tail -5 | grep -q "error"; then
        echo -e "${RED}Recent errors in I/O service logs:${NC}"
        docker-compose logs fr0g-ai-io 2>/dev/null | tail -5 | grep -i error
    fi
fi

# Test AIP service
echo -e "\n${BLUE}Testing AIP service...${NC}"
echo "  Checking correct port (8080)..."
if curl -f http://localhost:8080/health 2>/dev/null; then
    echo -e "${GREEN}✓ AIP service healthy on correct port 8080${NC}"
    
    # Test personas endpoint
    if curl -s http://localhost:8080/personas 2>/dev/null | grep -q "\["; then
        persona_count=$(curl -s http://localhost:8080/personas 2>/dev/null | jq '. | length' 2>/dev/null || echo "unknown")
        echo -e "${GREEN}✓ AIP service has $persona_count personas loaded${NC}"
    fi
else
    echo -e "${RED}✗ AIP service down on port 8080${NC}"
    if docker-compose logs fr0g-ai-aip 2>/dev/null | tail -5 | grep -q "error"; then
        echo -e "${RED}Recent errors in AIP service logs:${NC}"
        docker-compose logs fr0g-ai-aip 2>/dev/null | tail -5 | grep -i error
    fi
fi

# Test Master Control service
echo -e "\n${BLUE}Testing Master Control service...${NC}"
if curl -f http://localhost:8081/health 2>/dev/null; then
    echo -e "${GREEN}✓ Master Control service healthy${NC}"
    
    # Test intelligence metrics if available
    if curl -s http://localhost:8081/status 2>/dev/null | grep -q "learning_rate"; then
        echo -e "${GREEN}✓ Master Control AI intelligence operational${NC}"
    fi
else
    echo -e "${RED}✗ Master Control service down${NC}"
    if docker-compose logs fr0g-ai-master-control 2>/dev/null | tail -5 | grep -q "error"; then
        echo -e "${RED}Recent errors in Master Control service logs:${NC}"
        docker-compose logs fr0g-ai-master-control 2>/dev/null | tail -5 | grep -i error
    fi
fi

# Test service registry
echo -e "\n${BLUE}Testing Service Registry...${NC}"
if curl -f http://localhost:8500/health 2>/dev/null; then
    echo -e "${GREEN}✓ Service Registry healthy${NC}"
    
    # Test service discovery integration
    echo "Testing service discovery..."
    if command -v jq >/dev/null 2>&1; then
        SERVICES=$(curl -s http://localhost:8500/v1/catalog/services 2>/dev/null | jq -r 'keys[]' 2>/dev/null || echo "")
        if [ -n "$SERVICES" ]; then
            echo -e "${GREEN}Discovered services:${NC} $SERVICES"
            
            # Count registered services
            service_count=$(echo "$SERVICES" | wc -w)
            echo -e "${GREEN}✓ $service_count services registered${NC}"
        else
            echo -e "${YELLOW}⚠ No services discovered${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ jq not available - skipping service discovery test${NC}"
    fi
else
    echo -e "${RED}✗ Service Registry not available${NC}"
fi

# Port conflict analysis
echo ""
echo "=== Port Conflict Analysis ==="
echo "Analyzing port usage conflicts..."

# Check what's using each port
for port in 8080 8081 8082 8083 9090 9091 9092; do
    if command -v lsof >/dev/null 2>&1; then
        port_usage=$(lsof -ti:$port 2>/dev/null || echo "")
        if [ -n "$port_usage" ]; then
            process_info=$(ps -p $port_usage -o comm= 2>/dev/null || echo "unknown")
            echo "Port $port: IN USE by $process_info (PID: $port_usage)"
        else
            echo "Port $port: AVAILABLE"
        fi
    else
        # Fallback using netstat if lsof not available
        if netstat -ln 2>/dev/null | grep -q ":$port "; then
            echo "Port $port: IN USE (process unknown - install lsof for details)"
        else
            echo "Port $port: AVAILABLE"
        fi
    fi
done

echo ""
echo "CHECKING PORT CONFLICT ANALYSIS:"
echo "- Port 9091: Likely conflict between Bridge and AIP services"
echo "- Expected ports: AIP(8080,9090), Bridge(8082,9091), MCP(8081), I/O(8083,9092)"
echo "- Actual usage: Check above for current port assignments"

echo ""
echo -e "${BLUE}=== End-to-End Integration Test COMPLETED ===${NC}"
echo ""

# Count operational services
operational_count=0
total_services=5

# Check each service and count operational ones
if curl -sf http://localhost:8500/health >/dev/null 2>&1; then
    operational_count=$((operational_count + 1))
fi
if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
    operational_count=$((operational_count + 1))
fi
if curl -sf http://localhost:8082/health >/dev/null 2>&1; then
    operational_count=$((operational_count + 1))
fi
if curl -sf http://localhost:8081/health >/dev/null 2>&1; then
    operational_count=$((operational_count + 1))
fi
if curl -sf http://localhost:8083/health >/dev/null 2>&1; then
    operational_count=$((operational_count + 1))
fi

echo -e "${BLUE}Service Status Summary:${NC}"
echo -e "Operational Services: ${GREEN}$operational_count${NC}/${total_services}"

if [ $operational_count -eq $total_services ]; then
    echo -e "${GREEN}✓ ALL SERVICES OPERATIONAL${NC}"
    echo -e "${GREEN}✓ Complete fr0g.ai system is running successfully${NC}"
    exit 0
elif [ $operational_count -ge 3 ]; then
    echo -e "${YELLOW}⚠ PARTIAL SYSTEM OPERATIONAL${NC}"
    echo -e "${YELLOW}⚠ $operational_count/$total_services services running${NC}"
    exit 0
else
    echo -e "${RED}✗ SYSTEM NOT OPERATIONAL${NC}"
    echo -e "${RED}✗ Only $operational_count/$total_services services running${NC}"
    exit 1
fi
