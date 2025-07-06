#!/bin/bash
# fr0g.ai Comprehensive Test Runner
# Runs all test suites with proper reporting and coverage

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TEST_TIMEOUT=${TEST_TIMEOUT:-"10m"}
COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD:-80}
PARALLEL_TESTS=${PARALLEL_TESTS:-true}

# Directories
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COVERAGE_DIR="${PROJECT_ROOT}/coverage"
TEST_RESULTS_DIR="${PROJECT_ROOT}/test-results"

# Components to test
COMPONENTS=(
    "pkg/config"
    "fr0g-ai-aip"
    "fr0g-ai-bridge"
    "fr0g-ai-master-control"
    "fr0g-ai-io"
    "fr0g-ai-registry"
)

# Test types
TEST_TYPES=(
    "unit"
    "race"
    "coverage"
    "benchmark"
)

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

setup_test_environment() {
    log "Setting up test environment..."
    
    # Create directories
    mkdir -p "${COVERAGE_DIR}"
    mkdir -p "${TEST_RESULTS_DIR}"
    
    # Clean previous results
    rm -f "${COVERAGE_DIR}"/*.out
    rm -f "${COVERAGE_DIR}"/*.html
    rm -f "${TEST_RESULTS_DIR}"/*.xml
    
    success "Test environment ready"
}

run_unit_tests() {
    local component=$1
    log "Running unit tests for ${component}..."
    
    cd "${PROJECT_ROOT}/${component}"
    
    if go test -v -race -cover -timeout="${TEST_TIMEOUT}" ./...; then
        success "Unit tests passed for ${component}"
        return 0
    else
        error "Unit tests failed for ${component}"
        return 1
    fi
}

run_coverage_tests() {
    local component=$1
    local coverage_file="${COVERAGE_DIR}/$(basename ${component}).out"
    
    log "Running coverage tests for ${component}..."
    
    cd "${PROJECT_ROOT}/${component}"
    
    if go test -coverprofile="${coverage_file}" -covermode=atomic -timeout="${TEST_TIMEOUT}" ./...; then
        # Generate HTML report
        go tool cover -html="${coverage_file}" -o "${coverage_file%.out}.html"
        
        # Check coverage threshold
        local coverage=$(go tool cover -func="${coverage_file}" | grep total | awk '{print $3}' | sed 's/%//')
        if (( $(echo "${coverage} >= ${COVERAGE_THRESHOLD}" | bc -l) )); then
            success "Coverage tests passed for ${component} (${coverage}%)"
            return 0
        else
            warning "Coverage below threshold for ${component} (${coverage}% < ${COVERAGE_THRESHOLD}%)"
            return 1
        fi
    else
        error "Coverage tests failed for ${component}"
        return 1
    fi
}

run_race_tests() {
    local component=$1
    log "Running race detection tests for ${component}..."
    
    cd "${PROJECT_ROOT}/${component}"
    
    if go test -race -timeout="${TEST_TIMEOUT}" ./...; then
        success "Race detection tests passed for ${component}"
        return 0
    else
        error "Race detection tests failed for ${component}"
        return 1
    fi
}

run_benchmark_tests() {
    local component=$1
    local bench_file="${TEST_RESULTS_DIR}/$(basename ${component})-bench.txt"
    
    log "Running benchmark tests for ${component}..."
    
    cd "${PROJECT_ROOT}/${component}"
    
    if go test -bench=. -benchmem -timeout="${TEST_TIMEOUT}" ./... > "${bench_file}" 2>&1; then
        success "Benchmark tests completed for ${component}"
        return 0
    else
        warning "Benchmark tests had issues for ${component}"
        return 0  # Don't fail on benchmark issues
    fi
}

run_component_tests() {
    local component=$1
    local failed_tests=0
    
    log "Testing component: ${component}"
    
    # Check if component exists
    if [[ ! -d "${PROJECT_ROOT}/${component}" ]]; then
        warning "Component directory not found: ${component}"
        return 0
    fi
    
    # Run different test types
    for test_type in "${TEST_TYPES[@]}"; do
        case $test_type in
            "unit")
                run_unit_tests "$component" || ((failed_tests++))
                ;;
            "coverage")
                run_coverage_tests "$component" || ((failed_tests++))
                ;;
            "race")
                run_race_tests "$component" || ((failed_tests++))
                ;;
            "benchmark")
                run_benchmark_tests "$component" || ((failed_tests++))
                ;;
        esac
    done
    
    return $failed_tests
}

generate_test_report() {
    log "Generating test report..."
    
    local report_file="${TEST_RESULTS_DIR}/test-summary.txt"
    
    cat > "$report_file" << EOF
fr0g.ai Test Summary Report
Generated: $(date)
========================================

Test Configuration:
- Timeout: ${TEST_TIMEOUT}
- Coverage Threshold: ${COVERAGE_THRESHOLD}%
- Parallel Tests: ${PARALLEL_TESTS}

Components Tested:
EOF
    
    for component in "${COMPONENTS[@]}"; do
        echo "- ${component}" >> "$report_file"
    done
    
    echo "" >> "$report_file"
    echo "Coverage Reports:" >> "$report_file"
    
    for coverage_file in "${COVERAGE_DIR}"/*.html; do
        if [[ -f "$coverage_file" ]]; then
            echo "- $(basename "$coverage_file")" >> "$report_file"
        fi
    done
    
    success "Test report generated: ${report_file}"
}

main() {
    log "Starting comprehensive test suite for fr0g.ai"
    
    setup_test_environment
    
    local total_failures=0
    
    # Run tests for each component
    for component in "${COMPONENTS[@]}"; do
        if ! run_component_tests "$component"; then
            ((total_failures++))
        fi
        echo ""  # Add spacing between components
    done
    
    generate_test_report
    
    # Summary
    if [[ $total_failures -eq 0 ]]; then
        success "All tests completed successfully!"
        exit 0
    else
        error "Tests completed with ${total_failures} component(s) having failures"
        exit 1
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --timeout)
            TEST_TIMEOUT="$2"
            shift 2
            ;;
        --coverage-threshold)
            COVERAGE_THRESHOLD="$2"
            shift 2
            ;;
        --no-parallel)
            PARALLEL_TESTS=false
            shift
            ;;
        --component)
            COMPONENTS=("$2")
            shift 2
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --timeout DURATION          Set test timeout (default: 10m)"
            echo "  --coverage-threshold PERCENT Set coverage threshold (default: 80)"
            echo "  --no-parallel               Disable parallel test execution"
            echo "  --component NAME             Test only specific component"
            echo "  --help                       Show this help message"
            exit 0
            ;;
        *)
            error "Unknown option: $1"
            exit 1
            ;;
    esac
done

main "$@"
