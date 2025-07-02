#!/bin/bash

# Generate protobuf files for fr0g-ai-bridge
set -e

echo "Generating protobuf files..."

# Create the pb directory if it doesn't exist
mkdir -p internal/pb

# Clean any existing generated files
rm -f internal/pb/*.pb.go

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

# Verify the files were created
if [ ! -f internal/pb/bridge.pb.go ]; then
    echo "ERROR: bridge.pb.go was not generated"
    exit 1
fi

if [ ! -f internal/pb/bridge_grpc.pb.go ]; then
    echo "ERROR: bridge_grpc.pb.go was not generated"
    exit 1
fi

echo "✅ Protobuf generation successful"
