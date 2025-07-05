#!/bin/bash

# fr0g.ai Comprehensive Test Runner
# Runs all test suites in the correct order

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
WAIT_TIME=10
TEST_RESULTS=()

echo -e "${BLUE}🧪 fr0g.ai Comprehensive Test Suite${NC}"
echo "===================================="
echo "Running all test suites in sequence..."
echo ""

# Function to log test results
log_result() {
    local test_name="$1"
    local result="$2"
    local message="$3"
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✅ $test_name: PASSED${NC}"
        TEST_RESULTS+=("PASS: $test_name")
    elif [ "$result" = "FAIL" ]; then
        echo -e "${RED}❌ $test_name: FAILED - $message${NC}"
        TEST_RESULTS+=("FAIL: $test_name - $message")
    elif [ "$result" = "SKIP" ]; then
        echo -e "${YELLOW}⏭️  $test_name: SKIPPED - $message${NC}"
        TEST_RESULTS+=("SKIP: $test_name - $message")
    fi
}

# Function to run a test suite
run_test_suite() {
    local test_name="$1"
    local test_script="$2"
    local description="$3"
    
    echo -e "\n${BLUE}Running $test_name${NC}"
    echo "$(printf '=%.0s' {1..50})"
    echo "$description"
    echo ""
    
    if [ -f "$test_script" ]; then
        chmod +x "$test_script"
        if "$test_script"; then
            log_result "$test_name" "PASS"
            return 0
        else
            log_result "$test_name" "FAIL" "Test script failed"
            return 1
        fi
    else
        log_result "$test_name" "SKIP" "Test script not found: $test_script"
        return 0
    fi
}

# Function to check prerequisites
check_prerequisites() {
    echo -e "${YELLOW}Checking Prerequisites${NC}"
    echo "----------------------"
    
    # Check if curl is available
    if command -v curl &> /dev/null; then
        echo -e "${GREEN}✅ curl is available${NC}"
    else
        echo -e "${RED}❌ curl is required but not installed${NC}"
        exit 1
    fi
    
    # Check if docker-compose is available
    if command -v docker-compose &> /dev/null; then
        echo -e "${GREEN}✅ docker-compose is available${NC}"
    else
        echo -e "${YELLOW}⚠️  docker-compose not available (Docker tests will be skipped)${NC}"
    fi
    
    # Check if grpcurl is available
    if command -v grpcurl &> /dev/null; then
        echo -e "${GREEN}✅ grpcurl is available${NC}"
    else
        echo -e "${YELLOW}⚠️  grpcurl not available (gRPC tests will be limited)${NC}"
    fi
    
    echo ""
}

# Function to wait for services
wait_for_services() {
    echo -e "${YELLOW}Waiting for services to be ready...${NC}"
    echo "Services need time to initialize properly."
    echo ""
    
    for i in $(seq 1 $WAIT_TIME); do
        echo -ne "\rWaiting... ${i}/${WAIT_TIME}s"
        sleep 1
    done
    echo -e "\n"
}

# Function to check service availability
check_service_availability() {
    echo -e "${YELLOW}Checking Service Availability${NC}"
    echo "------------------------------"
    
    local services_available=0
    local total_services=4
    
    # Check each service
    for service_data in "AIP:http://localhost:8080" "Bridge:http://localhost:8082" "MCP:http://localhost:8081" "I/O:http://localhost:8083"; do
        IFS=':' read -r service_name service_url <<< "$service_data"
        
        if curl -s --max-time 5 "$service_url/health" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ $service_name service is available${NC}"
            services_available=$((services_available + 1))
        else
            echo -e "${RED}❌ $service_name service is not available${NC}"
        fi
    done
    
    echo ""
    echo "Services available: $services_available/$total_services"
    
    if [ $services_available -eq 0 ]; then
        echo -e "${RED}❌ No services are available. Please start the services first.${NC}"
        echo "Try running: make docker-up or ./scripts/start-services.sh"
        exit 1
    elif [ $services_available -lt $total_services ]; then
        echo -e "${YELLOW}⚠️  Some services are not available. Tests may be limited.${NC}"
    else
        echo -e "${GREEN}✅ All services are available!${NC}"
    fi
    
    echo ""
}

# Main test execution
main() {
    echo -e "${BLUE}Starting comprehensive test execution...${NC}"
    echo ""
    
    # Check prerequisites
    check_prerequisites
    
    # Wait for services to be ready
    wait_for_services
    
    # Check service availability
    check_service_availability
    
    # Run test suites in order
    echo -e "${BLUE}Executing Test Suites${NC}"
    echo "====================="
    
    # 1. Basic deployment test
    run_test_suite "Deployment Verification" \
        "scripts/test-deployment.sh" \
        "Verifies basic deployment and service connectivity"
    
    # 2. Service registry tests
    run_test_suite "Service Registry Tests" \
        "tests/integration/service_registry_test.sh" \
        "Tests service discovery and registration functionality"
    
    # 3. End-to-end integration tests
    run_test_suite "End-to-End Integration" \
        "tests/integration/end_to_end_test.sh" \
        "Comprehensive integration testing of all services"
    
    # 4. API integration tests
    run_test_suite "API Integration Tests" \
        "tests/integration/api_test.sh" \
        "Tests all service APIs and their functionality"
    
    # 5. Performance tests
    run_test_suite "Performance Tests" \
        "tests/integration/performance_test.sh" \
        "Tests performance and load characteristics"
    
    # Print final summary
    echo -e "\n${BLUE}Final Test Summary${NC}"
    echo "=================="
    
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    local skipped_tests=0
    
    for result in "${TEST_RESULTS[@]}"; do
        total_tests=$((total_tests + 1))
        if [[ $result == PASS:* ]]; then
            passed_tests=$((passed_tests + 1))
        elif [[ $result == FAIL:* ]]; then
            failed_tests=$((failed_tests + 1))
            echo -e "${RED}$result${NC}"
        elif [[ $result == SKIP:* ]]; then
            skipped_tests=$((skipped_tests + 1))
        fi
    done
    
    echo ""
    echo -e "Total Test Suites: $total_tests"
    echo -e "${GREEN}Passed: $passed_tests${NC}"
    echo -e "${RED}Failed: $failed_tests${NC}"
    echo -e "${YELLOW}Skipped: $skipped_tests${NC}"
    
    # Calculate success rate
    if [ $total_tests -gt 0 ]; then
        success_rate=$((passed_tests * 100 / total_tests))
        echo -e "\nSuccess Rate: ${success_rate}%"
    fi
    
    # Determine overall result
    if [ $failed_tests -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All test suites completed successfully!${NC}"
        echo -e "${GREEN}fr0g.ai is fully operational and ready for production.${NC}"
        exit 0
    else
        echo -e "\n${RED}❌ Some test suites failed.${NC}"
        echo -e "${RED}Please review the failed tests and fix any issues.${NC}"
        exit 1
    fi
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "fr0g.ai Comprehensive Test Runner"
        echo ""
        echo "Usage: $0 [options]"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --quick, -q    Run quick tests only (skip performance tests)"
        echo "  --verbose, -v  Enable verbose output"
        echo ""
        echo "This script runs all fr0g.ai test suites in the correct order:"
        echo "1. Deployment verification"
        echo "2. Service registry tests"
        echo "3. End-to-end integration tests"
        echo "4. API integration tests"
        echo "5. Performance tests"
        echo ""
        echo "Prerequisites:"
        echo "- Services must be running (use 'make docker-up' or './scripts/start-services.sh')"
        echo "- curl must be installed"
        echo "- Optional: docker-compose, grpcurl for enhanced testing"
        exit 0
        ;;
    --quick|-q)
        echo -e "${YELLOW}Running in quick mode (skipping performance tests)${NC}"
        # Set flag to skip performance tests
        SKIP_PERFORMANCE=true
        ;;
    --verbose|-v)
        echo -e "${YELLOW}Running in verbose mode${NC}"
        set -x
        ;;
esac

# Run main function
main "$@"
