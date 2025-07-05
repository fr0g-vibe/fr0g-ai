#!/bin/bash

# FR0G AI - ESMTP Processor & Cognitive Engine Test Runner
echo "🧠 FR0G AI - ESMTP Processor & Cognitive Engine Test"
echo "=================================================="

# Change to the correct directory
cd "$(dirname "$0")/.." || exit 1

# Ensure we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "FAILED Error: Not in fr0g-ai-master-control directory"
    exit 1
fi

echo "📁 Working directory: $(pwd)"

# Build the test program
echo "BUILDING Building test program..."
if ! go build -o bin/test-esmtp ./cmd/test-esmtp/; then
    echo "FAILED Build failed"
    exit 1
fi

echo "COMPLETED Build successful"

# Run the test
echo "STARTING Running ESMTP processor tests..."
echo ""

if ! ./bin/test-esmtp; then
    echo "FAILED Tests failed"
    exit 1
fi

echo ""
echo "🎉 All tests passed successfully!"
echo "📧 ESMTP Processor: OPERATIONAL"
echo "🧠 Cognitive Engine: OPERATIONAL"
