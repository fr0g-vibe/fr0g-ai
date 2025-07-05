#!/bin/bash

# Generate protobuf files for fr0g-ai-bridge
set -e

echo "Generating protobuf files..."

# Create the pb directory if it doesn't exist
mkdir -p internal/pb

# Clean any existing generated files
rm -f internal/pb/*.pb.go

# Debug: Check if protoc is available
if ! command -v protoc &> /dev/null; then
    echo "ERROR: protoc is not installed or not in PATH"
    echo "Please install protobuf compiler:"
    echo "  Ubuntu/Debian: sudo apt install protobuf-compiler"
    echo "  macOS: brew install protobuf"
    echo "  Arch Linux: sudo pacman -S protobuf"
    exit 1
fi

# Debug: Check if proto file exists
if [ ! -f proto/bridge.proto ]; then
    echo "ERROR: proto/bridge.proto not found"
    exit 1
fi

echo "Generating Go code from protobuf..."

# Debug: Show protoc version
echo "Protoc version:"
protoc --version

# Generate Go code from protobuf with explicit output directory
echo "Running protoc command..."
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go_opt=Mproto/bridge.proto=github.com/fr0g-vibe/fr0g-ai-bridge/internal/pb \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --go-grpc_opt=Mproto/bridge.proto=github.com/fr0g-vibe/fr0g-ai-bridge/internal/pb \
    proto/bridge.proto

echo "Protoc command completed with exit code: $?"

# Check if files were generated in proto directory and move them
if [ -f proto/bridge.pb.go ]; then
    echo "Moving bridge.pb.go to internal/pb/"
    mv proto/bridge.pb.go internal/pb/
fi

if [ -f proto/bridge_grpc.pb.go ]; then
    echo "Moving bridge_grpc.pb.go to internal/pb/"
    mv proto/bridge_grpc.pb.go internal/pb/
fi

echo "Protobuf generation complete!"
echo "Generated files:"
ls -la internal/pb/

# Verify the files were created
if [ ! -f internal/pb/bridge.pb.go ]; then
    echo "ERROR: bridge.pb.go was not generated"
    echo "Trying alternative generation method..."
    
    # Alternative method: generate directly to internal/pb
    protoc \
        --go_out=internal/pb \
        --go_opt=paths=source_relative \
        --go-grpc_out=internal/pb \
        --go-grpc_opt=paths=source_relative \
        proto/bridge.proto
    
    if [ ! -f internal/pb/bridge.pb.go ]; then
        echo "ERROR: Still failed to generate bridge.pb.go"
        exit 1
    fi
fi

if [ ! -f internal/pb/bridge_grpc.pb.go ]; then
    echo "ERROR: bridge_grpc.pb.go was not generated"
    exit 1
fi

echo "COMPLETED Protobuf generation successful"
