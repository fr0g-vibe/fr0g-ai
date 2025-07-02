#!/bin/bash

# Generate protobuf files for fr0g-ai-bridge
set -e

echo "Generating protobuf files..."

# Create the pb directory if it doesn't exist
mkdir -p internal/pb

# Generate Go code from protobuf
protoc \
    --go_out=internal/pb \
    --go_opt=paths=source_relative \
    --go-grpc_out=internal/pb \
    --go-grpc_opt=paths=source_relative \
    proto/bridge.proto

echo "Protobuf generation complete!"
echo "Generated files:"
ls -la internal/pb/
