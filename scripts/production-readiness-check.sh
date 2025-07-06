#!/bin/bash

echo "PRODUCTION READINESS CHECK"
echo "========================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
PASSED=0
FAILED=0
WARNINGS=0

# Function to check and report
check_status() {
    local test_name="$1"
    local command="$2"
    local critical="$3"
    
    echo -n "CHECKING $test_name... "
    
    if eval "$command" >/dev/null 2>&1; then
        echo -e "${GREEN}PASSED${NC}"
        ((PASSED++))
        return 0
    else
        if [ "$critical" = "true" ]; then
            echo -e "${RED}CRITICAL FAILED${NC}"
            ((FAILED++))
        else
            echo -e "${YELLOW}WARNING${NC}"
            ((WARNINGS++))
        fi
        return 1
    fi
}

echo -e "${BLUE}1. CONTAINER HEALTH CHECKS${NC}"
echo "----------------------------"
check_status "Registry container health" "docker inspect --format='{{.State.Health.Status}}' fr0g-ai-service-registry-1 | grep -q healthy" true
check_status "AIP container health" "docker inspect --format='{{.State.Health.Status}}' fr0g-ai-fr0g-ai-aip-1 | grep -q healthy" true
check_status "Bridge container health" "docker inspect --format='{{.State.Health.Status}}' fr0g-ai-fr0g-ai-bridge-1 | grep -q healthy" true
check_status "Master Control container health" "docker inspect --format='{{.State.Health.Status}}' fr0g-ai-fr0g-ai-master-control-1 | grep -q healthy" true
check_status "IO container health" "docker inspect --format='{{.State.Health.Status}}' fr0g-ai-fr0g-ai-io-1 | grep -q healthy" true

echo ""
echo -e "${BLUE}2. SERVICE ENDPOINT AVAILABILITY${NC}"
echo "-----------------------------------"
check_status "Registry HTTP endpoint" "curl -sf http://localhost:8500/health" true
check_status "AIP HTTP endpoint" "curl -sf http://localhost:8080/health" true
check_status "Bridge HTTP endpoint" "curl -sf http://localhost:8082/health" true
check_status "Master Control HTTP endpoint" "curl -sf http://localhost:8081/health" true
check_status "IO HTTP endpoint" "curl -sf http://localhost:8083/health" true

echo ""
echo -e "${BLUE}3. GRPC SERVICE AVAILABILITY${NC}"
echo "------------------------------"
check_status "AIP gRPC port" "nc -z localhost 9090" true
check_status "Bridge gRPC port" "nc -z localhost 9091" true
check_status "IO gRPC port" "nc -z localhost 9092" true

echo ""
echo -e "${BLUE}4. SERVICE REGISTRATION STATUS${NC}"
echo "--------------------------------"
check_status "Service discovery API" "curl -sf http://localhost:8500/v1/catalog/services" true
check_status "AIP service registered" "curl -s http://localhost:8500/v1/catalog/services | grep -q aip" true
check_status "Bridge service registered" "curl -s http://localhost:8500/v1/catalog/services | grep -q bridge" true
check_status "Master Control service registered" "curl -s http://localhost:8500/v1/catalog/services | grep -q mcp" true
check_status "IO service registered" "curl -s http://localhost:8500/v1/catalog/services | grep -q io" true

echo ""
echo -e "${BLUE}5. API FUNCTIONALITY TESTS${NC}"
echo "----------------------------"
check_status "AIP personas API" "curl -sf http://localhost:8080/personas" false
check_status "AIP identities API" "curl -sf http://localhost:8080/identities" false
check_status "Bridge chat completions API" "curl -sf http://localhost:8082/v1/chat/completions -X POST -H 'Content-Type: application/json' -d '{}'" false
check_status "IO input events API" "curl -sf http://localhost:8083/events" false

echo ""
echo -e "${BLUE}6. PERFORMANCE BENCHMARKS${NC}"
echo "----------------------------"
echo -n "CHECKING Registry response time... "
RESPONSE_TIME=$(curl -o /dev/null -s -w '%{time_total}' http://localhost:8500/health)
# Use awk instead of bc for floating point comparison
if awk "BEGIN {exit !($RESPONSE_TIME < 0.1)}"; then
    echo -e "${GREEN}PASSED${NC} (${RESPONSE_TIME}s)"
    ((PASSED++))
else
    echo -e "${YELLOW}WARNING${NC} (${RESPONSE_TIME}s - should be <0.1s)"
    ((WARNINGS++))
fi

echo -n "CHECKING Service discovery performance... "
START_TIME=$(date +%s.%N)
for i in {1..100}; do
    curl -s http://localhost:8500/v1/catalog/services >/dev/null
done
END_TIME=$(date +%s.%N)
# Use awk for floating point arithmetic instead of bc
TOTAL_TIME=$(awk "BEGIN {print $END_TIME - $START_TIME}")
AVG_TIME=$(awk "BEGIN {printf \"%.4f\", $TOTAL_TIME / 100}")

if awk "BEGIN {exit !($AVG_TIME < 0.01)}"; then
    echo -e "${GREEN}PASSED${NC} (avg ${AVG_TIME}s per lookup)"
    ((PASSED++))
else
    echo -e "${YELLOW}WARNING${NC} (avg ${AVG_TIME}s per lookup - should be <0.01s)"
    ((WARNINGS++))
fi

echo ""
echo -e "${BLUE}7. SECURITY CHECKS${NC}"
echo "-------------------"
check_status "Non-root container users" "docker exec fr0g-ai-fr0g-ai-aip-1 whoami | grep -v root" false
check_status "No privileged containers" "docker inspect fr0g-ai-fr0g-ai-aip-1 | grep -q '\"Privileged\": false'" false
check_status "Health check endpoints secure" "curl -sf http://localhost:8080/health | grep -qv 'error\\|exception\\|stack'" false

echo ""
echo -e "${BLUE}PRODUCTION READINESS SUMMARY${NC}"
echo "=============================="
echo -e "Tests Passed: ${GREEN}$PASSED${NC}"
echo -e "Critical Failures: ${RED}$FAILED${NC}"
echo -e "Warnings: ${YELLOW}$WARNINGS${NC}"

if [ $FAILED -eq 0 ]; then
    if [ $WARNINGS -eq 0 ]; then
        echo -e "\n${GREEN}✓ PRODUCTION READY${NC} - All checks passed!"
        exit 0
    else
        echo -e "\n${YELLOW}⚠ PRODUCTION READY WITH WARNINGS${NC} - Review warnings before deployment"
        exit 0
    fi
else
    echo -e "\n${RED}✗ NOT PRODUCTION READY${NC} - Critical failures must be resolved"
    exit 1
fi
