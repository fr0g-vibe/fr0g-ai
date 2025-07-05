#!/bin/bash

# Integration test script for fr0g-ai-registry
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REGISTRY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$REGISTRY_DIR")"

echo "🧪 Starting fr0g-ai-registry integration tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if registry is running
check_registry() {
    print_status "Checking if registry service is running..."
    
    if curl -s http://localhost:8500/health > /dev/null 2>&1; then
        print_success "Registry service is running"
        return 0
    else
        print_error "Registry service is not running on port 8500"
        return 1
    fi
}

# Start registry if not running
start_registry() {
    print_status "Starting registry service..."
    
    cd "$REGISTRY_DIR"
    
    # Build the registry if needed
    if [ ! -f "bin/fr0g-ai-registry" ]; then
        print_status "Building registry service..."
        make build
    fi
    
    # Start registry in background
    ./bin/fr0g-ai-registry &
    REGISTRY_PID=$!
    
    # Wait for registry to start
    print_status "Waiting for registry to start..."
    for i in {1..30}; do
        if curl -s http://localhost:8500/health > /dev/null 2>&1; then
            print_success "Registry started successfully (PID: $REGISTRY_PID)"
            return 0
        fi
        sleep 1
    done
    
    print_error "Registry failed to start within 30 seconds"
    kill $REGISTRY_PID 2>/dev/null || true
    return 1
}

# Run integration tests
run_tests() {
    print_status "Running integration tests..."
    
    cd "$REGISTRY_DIR"
    
    # Run the integration tests
    if go test -v ./test/... -run TestRegistryIntegration; then
        print_success "Integration tests passed"
    else
        print_error "Integration tests failed"
        return 1
    fi
    
    # Run performance tests
    print_status "Running performance tests..."
    if go test -v ./test/... -run TestRegistryPerformance; then
        print_success "Performance tests passed"
    else
        print_warning "Performance tests failed or had issues"
    fi
    
    # Run benchmarks
    print_status "Running benchmarks..."
    if go test -bench=. ./test/... -benchmem; then
        print_success "Benchmarks completed"
    else
        print_warning "Benchmarks failed or had issues"
    fi
}

# Test service registration workflow
test_service_workflow() {
    print_status "Testing service registration workflow..."
    
    # Test service registration
    print_status "Registering test service..."
    curl -X PUT http://localhost:8500/v1/agent/service/register \
        -H "Content-Type: application/json" \
        -d '{
            "id": "test-workflow-service",
            "name": "test-workflow-service",
            "address": "localhost",
            "port": 9999,
            "tags": ["test", "workflow"],
            "meta": {
                "version": "1.0.0",
                "env": "test"
            },
            "check": {
                "http": "http://localhost:9999/health",
                "interval": "10s",
                "timeout": "3s"
            }
        }'
    
    if [ $? -eq 0 ]; then
        print_success "Service registered successfully"
    else
        print_error "Service registration failed"
        return 1
    fi
    
    # Test service discovery
    print_status "Testing service discovery..."
    SERVICES=$(curl -s http://localhost:8500/v1/catalog/services)
    
    if echo "$SERVICES" | grep -q "test-workflow-service"; then
        print_success "Service discovered successfully"
    else
        print_error "Service not found in discovery"
        return 1
    fi
    
    # Test service deregistration
    print_status "Deregistering test service..."
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/test-workflow-service
    
    if [ $? -eq 0 ]; then
        print_success "Service deregistered successfully"
    else
        print_error "Service deregistration failed"
        return 1
    fi
}

# Test with actual fr0g.ai services if they're running
test_real_services() {
    print_status "Testing integration with real fr0g.ai services..."
    
    # Check which services are running
    SERVICES_TO_TEST=()
    
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        SERVICES_TO_TEST+=("fr0g-ai-aip:8080")
        print_status "Found fr0g-ai-aip running on port 8080"
    fi
    
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        SERVICES_TO_TEST+=("fr0g-ai-bridge:8081")
        print_status "Found fr0g-ai-bridge running on port 8081"
    fi
    
    if curl -s http://localhost:8082/health > /dev/null 2>&1; then
        SERVICES_TO_TEST+=("fr0g-ai-io:8082")
        print_status "Found fr0g-ai-io running on port 8082"
    fi
    
    if curl -s http://localhost:8083/health > /dev/null 2>&1; then
        SERVICES_TO_TEST+=("fr0g-ai-master-control:8083")
        print_status "Found fr0g-ai-master-control running on port 8083"
    fi
    
    if [ ${#SERVICES_TO_TEST[@]} -eq 0 ]; then
        print_warning "No fr0g.ai services found running - skipping real service tests"
        return 0
    fi
    
    # Test registration with real services
    for service_info in "${SERVICES_TO_TEST[@]}"; do
        IFS=':' read -r service_name service_port <<< "$service_info"
        
        print_status "Testing registration for $service_name..."
        
        curl -X PUT http://localhost:8500/v1/agent/service/register \
            -H "Content-Type: application/json" \
            -d "{
                \"id\": \"${service_name}-test\",
                \"name\": \"${service_name}\",
                \"address\": \"localhost\",
                \"port\": ${service_port},
                \"tags\": [\"fr0g-ai\", \"production\"],
                \"check\": {
                    \"http\": \"http://localhost:${service_port}/health\",
                    \"interval\": \"10s\",
                    \"timeout\": \"3s\"
                }
            }"
        
        if [ $? -eq 0 ]; then
            print_success "Successfully registered $service_name"
        else
            print_error "Failed to register $service_name"
        fi
    done
    
    # Verify all services are discoverable
    print_status "Verifying service discovery..."
    DISCOVERED_SERVICES=$(curl -s http://localhost:8500/v1/catalog/services)
    
    for service_info in "${SERVICES_TO_TEST[@]}"; do
        IFS=':' read -r service_name service_port <<< "$service_info"
        
        if echo "$DISCOVERED_SERVICES" | grep -q "$service_name"; then
            print_success "$service_name is discoverable"
        else
            print_error "$service_name not found in service discovery"
        fi
    done
}

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    
    # Kill registry if we started it
    if [ ! -z "$REGISTRY_PID" ]; then
        print_status "Stopping registry service (PID: $REGISTRY_PID)..."
        kill $REGISTRY_PID 2>/dev/null || true
        wait $REGISTRY_PID 2>/dev/null || true
    fi
    
    # Clean up any test services
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/test-workflow-service 2>/dev/null || true
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/fr0g-ai-aip-test 2>/dev/null || true
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/fr0g-ai-bridge-test 2>/dev/null || true
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/fr0g-ai-io-test 2>/dev/null || true
    curl -X PUT http://localhost:8500/v1/agent/service/deregister/fr0g-ai-master-control-test 2>/dev/null || true
}

# Set up cleanup trap
trap cleanup EXIT

# Main execution
main() {
    print_status "Starting fr0g-ai-registry integration test suite..."
    
    # Check if registry is already running
    if ! check_registry; then
        # Start registry if not running
        if ! start_registry; then
            print_error "Failed to start registry service"
            exit 1
        fi
        REGISTRY_STARTED_BY_SCRIPT=true
    else
        REGISTRY_STARTED_BY_SCRIPT=false
    fi
    
    # Run tests
    if ! run_tests; then
        print_error "Integration tests failed"
        exit 1
    fi
    
    # Test service workflow
    if ! test_service_workflow; then
        print_error "Service workflow tests failed"
        exit 1
    fi
    
    # Test with real services
    test_real_services
    
    print_success "All integration tests completed successfully!"
    
    # Only cleanup if we started the registry
    if [ "$REGISTRY_STARTED_BY_SCRIPT" = true ]; then
        cleanup
    else
        print_status "Registry was already running - leaving it running"
        # Still clean up test services
        curl -X PUT http://localhost:8500/v1/agent/service/deregister/test-workflow-service 2>/dev/null || true
    fi
}

# Run main function
main "$@"
