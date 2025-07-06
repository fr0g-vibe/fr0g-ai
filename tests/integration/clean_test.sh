#!/bin/bash

# Clean Integration Test - Stops services, starts fresh, and runs comprehensive tests
# This ensures accurate testing without interference from previous service state

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}fr0g.ai Clean Integration Test${NC}"
echo "=================================="
echo -e "${BLUE}This test will stop all services, start fresh, and run comprehensive tests${NC}\n"

# Function to cleanup on exit
cleanup() {
    echo -e "\n${BLUE}Cleaning up...${NC}"
    docker-compose down >/dev/null 2>&1 || true
}

# Set trap for cleanup
trap cleanup EXIT

# Step 1: Stop any running services
echo -e "${BLUE}Step 1: Stopping any running services...${NC}"
docker-compose down >/dev/null 2>&1 || true
echo -e "${GREEN}Services stopped${NC}"

# Step 2: Start services fresh
echo -e "\n${BLUE}Step 2: Starting services fresh...${NC}"
if docker-compose up -d >/dev/null 2>&1; then
    echo -e "${GREEN}Services started successfully${NC}"
else
    echo -e "${RED}Failed to start services${NC}"
    exit 1
fi

# Step 3: Wait for services to be ready
echo -e "\n${BLUE}Step 3: Waiting for services to be ready...${NC}"
echo -n "Waiting"
for i in {1..30}; do
    if curl -sf http://localhost:8500/health >/dev/null 2>&1; then
        echo -e "\n${GREEN}Services are ready!${NC}"
        break
    fi
    echo -n "."
    sleep 2
    if [ $i -eq 30 ]; then
        echo -e "\n${RED}Services failed to start within 60 seconds${NC}"
        exit 1
    fi
done

# Step 4: Run health checks
echo -e "\n${BLUE}Step 4: Running health checks...${NC}"
if make health; then
    echo -e "\n${GREEN}Health checks passed!${NC}"
else
    echo -e "\n${RED}Health checks failed!${NC}"
    exit 1
fi

# Step 5: Run integration tests
echo -e "\n${BLUE}Step 5: Running integration tests...${NC}"
if make test-integration; then
    echo -e "\n${GREEN}Integration tests passed!${NC}"
else
    echo -e "\n${RED}Integration tests failed!${NC}"
    exit 1
fi

echo -e "\n${GREEN}SUCCESS: All clean integration tests passed!${NC}"
echo -e "${BLUE}Services will be stopped during cleanup...${NC}"
