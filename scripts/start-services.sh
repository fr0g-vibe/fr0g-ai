#!/bin/bash

# Start all fr0g-ai services for development
set -e

echo "=== Starting fr0g-ai Services ==="

# Function to start a service in background
start_service() {
    local service_name=$1
    local service_dir=$2
    local service_port=$3
    
    echo "Starting $service_name on port $service_port..."
    cd "$service_dir"
    
    # Start service in background
    ./bin/$service_name > "../logs/${service_name}.log" 2>&1 &
    local pid=$!
    echo "$pid" > "../logs/${service_name}.pid"
    
    echo "$service_name started with PID $pid"
    cd ..
}

# Create logs directory
mkdir -p logs

# Start services in dependency order
echo "1. Starting Service Registry..."
cd fr0g-ai-master-control
./bin/fr0g-ai-master-control --registry-mode > ../logs/service-registry.log 2>&1 &
echo $! > ../logs/service-registry.pid
cd ..
sleep 2

echo "2. Starting AIP Service..."
start_service "fr0g-ai-aip" "fr0g-ai-aip" "8080"
sleep 2

echo "3. Starting Bridge Service..."
start_service "fr0g-ai-bridge" "fr0g-ai-bridge" "8082"
sleep 2

echo "4. Starting Master Control..."
start_service "fr0g-ai-master-control" "fr0g-ai-master-control" "8081"
sleep 2

echo "5. Starting I/O Service..."
start_service "fr0g-ai-io" "fr0g-ai-io" "8083"
sleep 2

echo "=== All Services Started ==="
echo "Waiting for services to initialize..."
sleep 5

echo "=== Service Health Check ==="
make health

echo "=== Service Status ==="
echo "Service logs available in logs/ directory"
echo "To stop services: ./scripts/stop-services.sh"
echo "To view logs: tail -f logs/[service-name].log"
