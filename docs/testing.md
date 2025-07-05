# Testing and Build Dispatch Commands

## Build and Test Command Reference

This document contains the standard testing and build commands that can be dispatched to tmux windows or executed directly.

### Build Commands

#### Build All Components
```bash
cd /fr0g-ai && make build-all
```

#### Build Individual Components
```bash
cd /fr0g-ai/fr0g-ai-aip && make build
cd /fr0g-ai/fr0g-ai-bridge && make build
cd /fr0g-ai/fr0g-ai-master-control && make build
cd /fr0g-ai/fr0g-ai-io && make build
```

### Service Registry Verification Commands

#### Find Registry Files
```bash
cd /fr0g-ai && find . -name '*registry*' -type f | head -10
```

#### Find Registry Implementation Files
```bash
cd /fr0g-ai && find . -path '*/registry/server.go' -o -path '*/registry/client.go' | head -5
```

#### Check Registry Integration in Components
```bash
cd /fr0g-ai/fr0g-ai-aip && grep -r 'registry' internal/ | head -3
cd /fr0g-ai/fr0g-ai-bridge && grep -r 'discovery\|registry' internal/ | head -3
```

### Processor Verification Commands

#### Check ESMTP Processor Implementation
```bash
cd /fr0g-ai/fr0g-ai-master-control && find . -path '*/processors/email/*' -name '*.go' | head -5
```

### Health Check Commands

#### System Health Check
```bash
cd /fr0g-ai && make health || echo 'Health check failed - services may not be running'
```

#### Docker Services Status
```bash
cd /fr0g-ai && docker-compose ps
```

### Integration Test Commands

#### Setup Integration Test Framework
```bash
cd /fr0g-ai && mkdir -p tests/integration && echo 'Integration test directory created'
cd /fr0g-ai && chmod +x tests/integration/*.sh
```

#### Run Integration Tests
```bash
cd /fr0g-ai && ./tests/integration/end_to_end_test.sh
cd /fr0g-ai && ./tests/integration/service_registry_test.sh
```

## Tmux Dispatch Commands

### Build Dispatch (Window 7)
```bash
# Build all components
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && make build-all" C-m

# Build individual components
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai/fr0g-ai-aip && make build" C-m
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai/fr0g-ai-bridge && make build" C-m
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai/fr0g-ai-master-control && make build" C-m
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai/fr0g-ai-io && make build" C-m
```

### Verification Dispatch (Window 7)
```bash
# Service registry verification
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && find . -name '*registry*' -type f | head -10" C-m
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && find . -path '*/registry/server.go' -o -path '*/registry/client.go' | head -5" C-m

# Component integration verification
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai/fr0g-ai-aip && grep -r 'registry' internal/ | head -3" C-m
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai/fr0g-ai-bridge && grep -r 'discovery\|registry' internal/ | head -3" C-m

# Processor verification
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai/fr0g-ai-master-control && find . -path '*/processors/email/*' -name '*.go' | head -5" C-m
```

### Health Check Dispatch (Window 7)
```bash
# System health
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && make health || echo 'Health check failed - services may not be running'" C-m

# Docker status
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && docker-compose ps" C-m
```

### Integration Test Dispatch (Window 7)
```bash
# Setup tests
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && mkdir -p tests/integration && echo 'Integration test directory created'" C-m
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && chmod +x tests/integration/*.sh" C-m

# Run tests
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && ./tests/integration/end_to_end_test.sh" C-m
tmux send-keys -t fr0g-ai:7 "cd /fr0g-ai && ./tests/integration/service_registry_test.sh" C-m
```

## Usage Examples

### Quick Build Verification
```bash
# Verify all components build successfully
make build-all

# Check for any build errors
echo $?
```

### Service Registry Status Check
```bash
# Check if service registry is implemented
find . -name '*registry*' -type f

# Verify registry integration
grep -r 'registry' fr0g-ai-*/internal/ | head -10
```

### Integration Test Execution
```bash
# Run full integration test suite
./tests/integration/end_to_end_test.sh

# Run service registry specific tests
./tests/integration/service_registry_test.sh
```

## Expected Outputs

### Successful Build
- All `make build` commands should exit with status 0
- No compilation errors should be reported
- Binary files should be created in expected locations

### Service Registry Verification
- Registry files should be found in multiple components
- Registry integration should show import statements and usage
- Service discovery should show endpoint configuration

### Health Checks
- Services should report healthy status
- Docker containers should be running
- All endpoints should respond to health checks

### Integration Tests
- All services should be discoverable
- Health endpoints should return 200 status
- Service registry should show all registered services
