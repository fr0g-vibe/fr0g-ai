#!/bin/bash

# VERIFICATION TASK: Monitor critical fixes across all services
# TRACK: 1) AIP configuration migration progress
#        2) Master-control storage validation fix  
#        3) I/O API integration status
# VERIFY: Build success, service startup, health checks
# COORDINATE: Report any infrastructure issues blocking fixes
# TARGET: All services operational after fixes

set -e

echo "=========================================="
echo "VERIFICATION Critical Fixes Monitoring"
echo "=========================================="

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Status tracking
ISSUES_FOUND=0

log_status() {
    local status=$1
    local message=$2
    case $status in
        "COMPLETED")
            echo -e "${GREEN}COMPLETED${NC} $message"
            ;;
        "CRITICAL")
            echo -e "${RED}CRITICAL${NC} $message"
            ISSUES_FOUND=$((ISSUES_FOUND + 1))
            ;;
        "MISSING")
            echo -e "${YELLOW}MISSING${NC} $message"
            ISSUES_FOUND=$((ISSUES_FOUND + 1))
            ;;
        *)
            echo "$status $message"
            ;;
    esac
}

# 1) AIP Configuration Migration Progress
echo ""
echo "TRACKING 1) AIP configuration migration progress..."
echo "=================================================="

# Check AIP build
if [ -d "fr0g-ai-aip" ]; then
    cd fr0g-ai-aip
    if go build -v ./... >/dev/null 2>&1; then
        log_status "COMPLETED" "AIP builds successfully"
    else
        log_status "CRITICAL" "AIP build failed"
        echo "Build errors:"
        go build -v ./... 2>&1 | tail -10
    fi
    cd ..
else
    log_status "MISSING" "AIP directory not found"
fi

# Check AIP service health
if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
    log_status "COMPLETED" "AIP service operational"
    
    # Check configuration validation in health response
    aip_health=$(curl -s http://localhost:8080/health 2>/dev/null)
    if echo "$aip_health" | grep -q "config.*valid\|validation.*success" 2>/dev/null; then
        log_status "COMPLETED" "AIP configuration validation working"
    else
        log_status "MISSING" "AIP configuration validation status unclear"
    fi
else
    log_status "MISSING" "AIP service not responding"
fi

# 2) Master-control Storage Validation Fix
echo ""
echo "TRACKING 2) Master-control storage validation fix..."
echo "===================================================="

# Check MCP build
if [ -d "fr0g-ai-master-control" ]; then
    cd fr0g-ai-master-control
    if go build -v ./... >/dev/null 2>&1; then
        log_status "COMPLETED" "MCP builds successfully"
    else
        log_status "CRITICAL" "MCP build failed"
        echo "Build errors:"
        go build -v ./... 2>&1 | tail -10
    fi
    cd ..
else
    log_status "MISSING" "MCP directory not found"
fi

# Check MCP service health
if curl -sf http://localhost:8081/health >/dev/null 2>&1; then
    log_status "COMPLETED" "MCP service operational"
    
    # Check storage validation in health response
    mcp_health=$(curl -s http://localhost:8081/health 2>/dev/null)
    if echo "$mcp_health" | grep -q "storage.*valid\|validation.*success" 2>/dev/null; then
        log_status "COMPLETED" "MCP storage validation working"
    else
        log_status "MISSING" "MCP storage validation status unclear"
    fi
else
    log_status "MISSING" "MCP service not responding"
fi

# 3) I/O API Integration Status
echo ""
echo "TRACKING 3) I/O API integration status..."
echo "========================================="

# Check IO build
if [ -d "fr0g-ai-io" ]; then
    cd fr0g-ai-io
    if go build -v ./... >/dev/null 2>&1; then
        log_status "COMPLETED" "IO builds successfully"
    else
        log_status "CRITICAL" "IO build failed"
        echo "Build errors:"
        go build -v ./... 2>&1 | tail -10
    fi
    cd ..
else
    log_status "MISSING" "IO directory not found"
fi

# Check IO service health
if curl -sf http://localhost:8083/health >/dev/null 2>&1; then
    log_status "COMPLETED" "IO service operational"
    
    # Check API integration in health response
    io_health=$(curl -s http://localhost:8083/health 2>/dev/null)
    if echo "$io_health" | grep -q "api.*integrated\|integration.*success" 2>/dev/null; then
        log_status "COMPLETED" "IO API integration working"
    else
        log_status "MISSING" "IO API integration status unclear"
    fi
else
    log_status "MISSING" "IO service not responding"
fi

# VERIFY: Build Success
echo ""
echo "VERIFY Build success across all services..."
echo "==========================================="

if command -v make >/dev/null 2>&1 && make build-all >/dev/null 2>&1; then
    log_status "COMPLETED" "All services build successfully"
else
    log_status "CRITICAL" "Build failures detected or make not available"
fi

# VERIFY: Service Startup
echo ""
echo "VERIFY Service startup verification..."
echo "====================================="

if command -v docker-compose >/dev/null 2>&1; then
    if docker-compose ps --format "table {{.Name}}\t{{.Status}}" 2>/dev/null | grep -v "Exit" >/dev/null 2>&1; then
        log_status "COMPLETED" "All containers running"
        echo "Container status:"
        docker-compose ps --format "table {{.Name}}\t{{.Status}}" 2>/dev/null || echo "Unable to get container status"
    else
        log_status "CRITICAL" "Container startup issues"
        docker-compose ps 2>/dev/null || echo "Unable to get container status"
    fi
else
    log_status "MISSING" "Docker Compose not available"
fi

# VERIFY: Health Checks
echo ""
echo "VERIFY Health checks verification..."
echo "==================================="

# Service Registry
if curl -sf http://localhost:8500/health >/dev/null 2>&1; then
    log_status "COMPLETED" "Registry health check passing"
else
    log_status "CRITICAL" "Registry health check failing"
fi

# Core services health check
services_healthy=0
total_services=4

for service_port in "8080:AIP" "8081:MCP" "8082:Bridge" "8083:IO"; do
    port=$(echo $service_port | cut -d: -f1)
    name=$(echo $service_port | cut -d: -f2)
    
    if curl -sf "http://localhost:$port/health" >/dev/null 2>&1; then
        services_healthy=$((services_healthy + 1))
    fi
done

if [ $services_healthy -eq $total_services ]; then
    log_status "COMPLETED" "All health checks passing ($services_healthy/$total_services)"
else
    log_status "CRITICAL" "Health check failures detected ($services_healthy/$total_services healthy)"
fi

# COORDINATE: Infrastructure Issues
echo ""
echo "COORDINATE Infrastructure issues blocking fixes..."
echo "================================================="

if command -v docker-compose >/dev/null 2>&1; then
    echo "Docker container status:"
    docker-compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || echo "Unable to get container status"
    
    echo ""
    echo "Recent service logs (errors only):"
    
    for service in "fr0g-ai-aip:AIP" "fr0g-ai-master-control:MCP" "fr0g-ai-io:IO"; do
        container=$(echo $service | cut -d: -f1)
        name=$(echo $service | cut -d: -f2)
        
        echo "$name errors:"
        if docker-compose logs --tail=5 "$container" 2>/dev/null | grep -i error; then
            :  # Errors found and displayed
        else
            echo "No $name errors found"
        fi
    done
else
    echo "Docker Compose not available - cannot check container status"
fi

echo ""
echo "Network connectivity:"
if command -v docker >/dev/null 2>&1 && docker network ls | grep -q fr0g-ai-network; then
    log_status "COMPLETED" "Network operational"
else
    log_status "CRITICAL" "Network missing or Docker not available"
fi

# Final Summary
echo ""
echo "=========================================="
echo "VERIFICATION SUMMARY"
echo "=========================================="

if [ $ISSUES_FOUND -eq 0 ]; then
    log_status "COMPLETED" "All services operational after fixes - TARGET ACHIEVED"
    exit 0
else
    log_status "CRITICAL" "$ISSUES_FOUND issues found blocking fixes"
    echo ""
    echo "RECOMMENDED ACTIONS:"
    echo "1. Check service logs: docker-compose logs [service-name]"
    echo "2. Rebuild services: make build-all"
    echo "3. Restart services: docker-compose restart"
    echo "4. Check configuration: verify environment variables"
    exit 1
fi
