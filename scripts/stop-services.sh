#!/bin/bash

# Stop all fr0g-ai services
set -e

echo "=== Stopping fr0g-ai Services ==="

# Function to stop a service
stop_service() {
    local service_name=$1
    local pid_file="logs/${service_name}.pid"
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        echo "Stopping $service_name (PID: $pid)..."
        
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid"
            sleep 2
            
            # Force kill if still running
            if kill -0 "$pid" 2>/dev/null; then
                echo "Force killing $service_name..."
                kill -9 "$pid"
            fi
        fi
        
        rm -f "$pid_file"
        echo "$service_name stopped"
    else
        echo "$service_name PID file not found"
    fi
}

# Stop services in reverse order
stop_service "fr0g-ai-io"
stop_service "fr0g-ai-master-control"
stop_service "fr0g-ai-bridge"
stop_service "fr0g-ai-aip"
stop_service "service-registry"

echo "=== All Services Stopped ==="

# Clean up log files if requested
if [ "$1" = "--clean-logs" ]; then
    echo "Cleaning up log files..."
    rm -f logs/*.log
    echo "Log files cleaned"
fi
