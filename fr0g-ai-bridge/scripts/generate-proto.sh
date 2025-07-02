#!/bin/bash

# Generate protobuf files for fr0g-ai-bridge
set -e

echo "Generating protobuf files..."

# Create the pb directory if it doesn't exist
mkdir -p internal/pb

# Generate Go code from protobuf with correct output paths
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --go_opt=Mproto/bridge.proto=github.com/fr0g-vibe/fr0g-ai-bridge/internal/pb \
    --go-grpc_opt=Mproto/bridge.proto=github.com/fr0g-vibe/fr0g-ai-bridge/internal/pb \
    proto/bridge.proto

# Move generated files to correct location if they're in the wrong place
if [ -f proto/bridge.pb.go ]; then
    mv proto/bridge.pb.go internal/pb/
fi
if [ -f proto/bridge_grpc.pb.go ]; then
    mv proto/bridge_grpc.pb.go internal/pb/
fi

echo "Protobuf generation complete!"
echo "Generated files:"
ls -la internal/pb/
