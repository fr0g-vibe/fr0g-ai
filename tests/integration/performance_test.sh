#!/bin/bash

# fr0g.ai Performance Integration Test
# Tests performance and load characteristics of all services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AIP_URL="http://localhost:8080"
BRIDGE_URL="http://localhost:8082"
MCP_URL="http://localhost:8081"
IO_URL="http://localhost:8083"
TIMEOUT=30
CONCURRENT_REQUESTS=10
TEST_DURATION=30
TEST_RESULTS=()

echo -e "${BLUE}fr0g.ai Performance Integration Test${NC}"
echo "======================================="
echo "Testing performance and load characteristics..."
echo ""

# Function to log test results
log_test() {
    local test_name="$1"
    local result="$2"
    local message="$3"
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}PASS $test_name: $message${NC}"
        TEST_RESULTS+=("PASS: $test_name - $message")
    elif [ "$result" = "FAIL" ]; then
        echo -e "${RED}FAIL $test_name: $message${NC}"
        TEST_RESULTS+=("FAIL: $test_name - $message")
    elif [ "$result" = "SKIP" ]; then
        echo -e "${YELLOW}SKIP $test_name: $message${NC}"
        TEST_RESULTS+=("SKIP: $test_name - $message")
    fi
}

# Function to test response time
test_response_time() {
    local service_name="$1"
    local url="$2"
    local endpoint="${3:-/health}"
    local iterations="${4:-10}"
    
    echo -n "Testing $service_name response time... "
    
    local total_time=0
    local successful_requests=0
    local failed_requests=0
    
    for ((i=1; i<=iterations; i++)); do
        start_time=$(date +%s%N)
        if curl -s --max-time 10 "$url$endpoint" > /dev/null 2>&1; then
            end_time=$(date +%s%N)
            duration=$(( (end_time - start_time) / 1000000 ))
            total_time=$((total_time + duration))
            successful_requests=$((successful_requests + 1))
        else
            failed_requests=$((failed_requests + 1))
        fi
    done
    
    if [ $successful_requests -gt 0 ]; then
        avg_time=$((total_time / successful_requests))
        success_rate=$((successful_requests * 100 / iterations))
        
        if [ $avg_time -lt 100 ]; then
            log_test "$service_name Response Time" "PASS" "Avg: ${avg_time}ms, Success: ${success_rate}% (Excellent)"
        elif [ $avg_time -lt 500 ]; then
            log_test "$service_name Response Time" "PASS" "Avg: ${avg_time}ms, Success: ${success_rate}% (Good)"
        elif [ $avg_time -lt 2000 ]; then
            log_test "$service_name Response Time" "PASS" "Avg: ${avg_time}ms, Success: ${success_rate}% (Acceptable)"
        else
            log_test "$service_name Response Time" "FAIL" "Avg: ${avg_time}ms, Success: ${success_rate}% (Slow)"
        fi
    else
        log_test "$service_name Response Time" "FAIL" "No successful requests"
    fi
}

# Function to test concurrent load
test_concurrent_load() {
    local service_name="$1"
    local url="$2"
    local endpoint="${3:-/health}"
    local concurrent="${4:-5}"
    
    echo -n "Testing $service_name concurrent load ($concurrent requests)... "
    
    # Create temporary directory for results
    local temp_dir=$(mktemp -d)
    
    # Start concurrent requests
    for ((i=1; i<=concurrent; i++)); do
        {
            start_time=$(date +%s%N)
            if curl -s --max-time 15 "$url$endpoint" > /dev/null 2>&1; then
                end_time=$(date +%s%N)
                duration=$(( (end_time - start_time) / 1000000 ))
                echo "SUCCESS:$duration" > "$temp_dir/result_$i"
            else
                echo "FAIL:0" > "$temp_dir/result_$i"
            fi
        } &
    done
    
    # Wait for all requests to complete
    wait
    
    # Analyze results
    local successful=0
    local failed=0
    local total_time=0
    local max_time=0
    local min_time=999999
    
    for result_file in "$temp_dir"/result_*; do
        if [ -f "$result_file" ]; then
            result=$(cat "$result_file")
            status=$(echo "$result" | cut -d':' -f1)
            time=$(echo "$result" | cut -d':' -f2)
            
            if [ "$status" = "SUCCESS" ]; then
                successful=$((successful + 1))
                total_time=$((total_time + time))
                if [ $time -gt $max_time ]; then max_time=$time; fi
                if [ $time -lt $min_time ]; then min_time=$time; fi
            else
                failed=$((failed + 1))
            fi
        fi
    done
    
    # Clean up
    rm -rf "$temp_dir"
    
    if [ $successful -gt 0 ]; then
        avg_time=$((total_time / successful))
        success_rate=$((successful * 100 / concurrent))
        
        if [ $success_rate -ge 90 ] && [ $avg_time -lt 1000 ]; then
            log_test "$service_name Concurrent Load" "PASS" "Success: ${success_rate}%, Avg: ${avg_time}ms, Max: ${max_time}ms"
        elif [ $success_rate -ge 70 ]; then
            log_test "$service_name Concurrent Load" "PASS" "Success: ${success_rate}%, Avg: ${avg_time}ms (Acceptable)"
        else
            log_test "$service_name Concurrent Load" "FAIL" "Success: ${success_rate}%, Avg: ${avg_time}ms (Poor)"
        fi
    else
        log_test "$service_name Concurrent Load" "FAIL" "No successful requests"
    fi
}

# Function to test memory usage (if available)
test_memory_usage() {
    local service_name="$1"
    local process_pattern="$2"
    
    if command -v ps &> /dev/null; then
        # Get memory usage for the service process
        memory_kb=$(ps aux | grep "$process_pattern" | grep -v grep | awk '{sum += $6} END {print sum}' || echo "0")
        
        if [ "$memory_kb" -gt 0 ]; then
            memory_mb=$((memory_kb / 1024))
            
            if [ $memory_mb -lt 100 ]; then
                log_test "$service_name Memory Usage" "PASS" "${memory_mb}MB (Efficient)"
            elif [ $memory_mb -lt 500 ]; then
                log_test "$service_name Memory Usage" "PASS" "${memory_mb}MB (Normal)"
            elif [ $memory_mb -lt 1000 ]; then
                log_test "$service_name Memory Usage" "PASS" "${memory_mb}MB (High)"
            else
                log_test "$service_name Memory Usage" "FAIL" "${memory_mb}MB (Excessive)"
            fi
        else
            log_test "$service_name Memory Usage" "SKIP" "Process not found"
        fi
    else
        log_test "$service_name Memory Usage" "SKIP" "ps command not available"
    fi
}

# Function to test API throughput
test_api_throughput() {
    local service_name="$1"
    local url="$2"
    local endpoint="$3"
    local duration="${4:-10}"
    
    echo -n "Testing $service_name API throughput (${duration}s)... "
    
    local start_time=$(date +%s)
    local end_time=$((start_time + duration))
    local request_count=0
    local error_count=0
    
    while [ $(date +%s) -lt $end_time ]; do
        if curl -s --max-time 5 "$url$endpoint" > /dev/null 2>&1; then
            request_count=$((request_count + 1))
        else
            error_count=$((error_count + 1))
        fi
    done
    
    local actual_duration=$(($(date +%s) - start_time))
    local rps=$((request_count / actual_duration))
    local error_rate=0
    
    if [ $((request_count + error_count)) -gt 0 ]; then
        error_rate=$((error_count * 100 / (request_count + error_count)))
    fi
    
    if [ $rps -ge 10 ] && [ $error_rate -le 5 ]; then
        log_test "$service_name API Throughput" "PASS" "${rps} RPS, ${error_rate}% errors (Excellent)"
    elif [ $rps -ge 5 ] && [ $error_rate -le 10 ]; then
        log_test "$service_name API Throughput" "PASS" "${rps} RPS, ${error_rate}% errors (Good)"
    elif [ $rps -ge 1 ] && [ $error_rate -le 20 ]; then
        log_test "$service_name API Throughput" "PASS" "${rps} RPS, ${error_rate}% errors (Acceptable)"
    else
        log_test "$service_name API Throughput" "FAIL" "${rps} RPS, ${error_rate}% errors (Poor)"
    fi
}

# Function to test all services
test_all_services() {
    echo -e "${YELLOW}Testing Service Performance${NC}"
    echo "---------------------------"
    
    # Test each service if it's available
    for service_data in "AIP:$AIP_URL:fr0g-ai-aip" "Bridge:$BRIDGE_URL:fr0g-ai-bridge" "MCP:$MCP_URL:fr0g-ai-master-control" "I/O:$IO_URL:fr0g-ai-io"; do
        IFS=':' read -r service_name service_url process_pattern <<< "$service_data"
        
        # Check if service is available
        if curl -s --max-time 5 "$service_url/health" > /dev/null 2>&1; then
            echo -e "\n${BLUE}Testing $service_name Service${NC}"
            echo "$(printf '=%.0s' {1..30})"
            
            # Response time test
            test_response_time "$service_name" "$service_url" "/health" 10
            
            # Concurrent load test
            test_concurrent_load "$service_name" "$service_url" "/health" 5
            
            # Memory usage test
            test_memory_usage "$service_name" "$process_pattern"
            
            # API throughput test
            test_api_throughput "$service_name" "$service_url" "/health" 10
            
        else
            log_test "$service_name Performance" "SKIP" "Service not available"
        fi
    done
}

# Function to test system resources
test_system_resources() {
    echo -e "\n${YELLOW}Testing System Resources${NC}"
    echo "-------------------------"
    
    # CPU usage
    if command -v top &> /dev/null; then
        cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1 || echo "0")
        if [ "${cpu_usage%.*}" -lt 80 ]; then
            log_test "System CPU Usage" "PASS" "${cpu_usage}% (Normal)"
        else
            log_test "System CPU Usage" "FAIL" "${cpu_usage}% (High)"
        fi
    else
        log_test "System CPU Usage" "SKIP" "top command not available"
    fi
    
    # Memory usage
    if command -v free &> /dev/null; then
        memory_info=$(free -m | grep "Mem:")
        total_mem=$(echo $memory_info | awk '{print $2}')
        used_mem=$(echo $memory_info | awk '{print $3}')
        mem_percent=$((used_mem * 100 / total_mem))
        
        if [ $mem_percent -lt 80 ]; then
            log_test "System Memory Usage" "PASS" "${mem_percent}% (${used_mem}MB/${total_mem}MB)"
        else
            log_test "System Memory Usage" "FAIL" "${mem_percent}% (${used_mem}MB/${total_mem}MB)"
        fi
    else
        log_test "System Memory Usage" "SKIP" "free command not available"
    fi
    
    # Disk usage
    if command -v df &> /dev/null; then
        disk_usage=$(df -h / | tail -1 | awk '{print $5}' | cut -d'%' -f1)
        if [ "$disk_usage" -lt 90 ]; then
            log_test "System Disk Usage" "PASS" "${disk_usage}% (Normal)"
        else
            log_test "System Disk Usage" "FAIL" "${disk_usage}% (High)"
        fi
    else
        log_test "System Disk Usage" "SKIP" "df command not available"
    fi
}

# Main test execution
main() {
    echo -e "${BLUE}Starting performance tests...${NC}"
    echo ""
    
    # Wait for services to be ready
    echo -e "${YELLOW}Waiting for services to initialize...${NC}"
    sleep 5
    
    # Run test suites
    test_all_services
    test_system_resources
    
    # Print summary
    echo -e "\n${BLUE}Performance Test Summary${NC}"
    echo "========================"
    
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
    echo -e "Total Tests: $total_tests"
    echo -e "${GREEN}Passed: $passed_tests${NC}"
    echo -e "${RED}Failed: $failed_tests${NC}"
    echo -e "${YELLOW}Skipped: $skipped_tests${NC}"
    
    # Determine overall result
    if [ $failed_tests -eq 0 ]; then
        echo -e "\n${GREEN}SUCCESS: All performance tests passed!${NC}"
        exit 0
    else
        echo -e "\n${RED}ERROR: Some performance tests failed.${NC}"
        exit 1
    fi
}

# Run main function
main "$@"
