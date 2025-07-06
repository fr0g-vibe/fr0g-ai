#!/bin/bash

# Performance test script for fr0g-ai-registry optimizations
# Tests the <50ms discovery latency target

set -e

REGISTRY_URL="http://localhost:8500"
TEST_SERVICES=100
CONCURRENT_REQUESTS=200

echo "🚀 fr0g-ai-registry Performance Test"
echo "Target: <50ms discovery latency under load"
echo "Registry URL: $REGISTRY_URL"
echo

# Check if registry is running
if ! curl -s "$REGISTRY_URL/health" > /dev/null; then
    echo "MISSING Registry not running at $REGISTRY_URL"
    exit 1
fi

echo "COMPLETED Registry is running"

# Register test services
echo "📝 Registering $TEST_SERVICES test services..."
for i in $(seq 1 $TEST_SERVICES); do
    curl -s -X PUT "$REGISTRY_URL/v1/agent/service/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"id\": \"test-service-$i\",
            \"name\": \"test-service\",
            \"address\": \"127.0.0.1\",
            \"port\": $((8000 + i)),
            \"tags\": [\"test\", \"performance\"],
            \"meta\": {\"version\": \"1.0.0\", \"test\": \"true\"}
        }" > /dev/null
done

echo "COMPLETED Registered $TEST_SERVICES services"

# Test discovery performance under load
echo "⚡ Testing discovery performance with $CONCURRENT_REQUESTS concurrent requests..."

# Create temporary files for results
TEMP_DIR=$(mktemp -d)
RESULTS_FILE="$TEMP_DIR/results.txt"

# Function to test discovery latency
test_discovery() {
    local start_time=$(date +%s%N)
    curl -s "$REGISTRY_URL/v1/catalog/services" > /dev/null
    local end_time=$(date +%s%N)
    local latency=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    echo "$latency" >> "$RESULTS_FILE"
}

# Run concurrent discovery requests
for i in $(seq 1 $CONCURRENT_REQUESTS); do
    test_discovery &
done

# Wait for all requests to complete
wait

# Calculate statistics
TOTAL_REQUESTS=$(wc -l < "$RESULTS_FILE")
AVG_LATENCY=$(awk '{sum+=$1} END {print sum/NR}' "$RESULTS_FILE")
MAX_LATENCY=$(sort -n "$RESULTS_FILE" | tail -1)
MIN_LATENCY=$(sort -n "$RESULTS_FILE" | head -1)
P95_LATENCY=$(sort -n "$RESULTS_FILE" | awk 'BEGIN{p95_line=0} {lines[NR]=$0} END{p95_line=int(NR*0.95); if(p95_line>0) print lines[p95_line]; else print ""}')

echo
echo "📊 Performance Results:"
echo "Total Requests: $TOTAL_REQUESTS"
echo "Average Latency: ${AVG_LATENCY}ms"
echo "Min Latency: ${MIN_LATENCY}ms"
echo "Max Latency: ${MAX_LATENCY}ms"
echo "95th Percentile: ${P95_LATENCY}ms"
echo

# Check if target is met (using awk for arithmetic)
TARGET_LATENCY=50
if [ -n "$AVG_LATENCY" ] && awk "BEGIN {exit !($AVG_LATENCY < $TARGET_LATENCY)}" 2>/dev/null; then
    echo "COMPLETED SUCCESS: Average latency (${AVG_LATENCY}ms) is below target (${TARGET_LATENCY}ms)"
    SUCCESS=true
else
    echo "CRITICAL FAILED: Average latency (${AVG_LATENCY}ms) exceeds target (${TARGET_LATENCY}ms)"
    SUCCESS=false
fi

if [ -n "$P95_LATENCY" ] && awk "BEGIN {exit !($P95_LATENCY < $TARGET_LATENCY)}" 2>/dev/null; then
    echo "COMPLETED SUCCESS: 95th percentile (${P95_LATENCY}ms) is below target (${TARGET_LATENCY}ms)"
else
    echo "CRITICAL FAILED: 95th percentile (${P95_LATENCY}ms) exceeds target (${TARGET_LATENCY}ms)"
    SUCCESS=false
fi

# Test cache effectiveness
echo
echo "🔄 Testing cache effectiveness..."
CACHE_START=$(date +%s%N)
curl -s "$REGISTRY_URL/v1/catalog/services" > /dev/null
CACHE_END=$(date +%s%N)
CACHE_LATENCY=$(( (CACHE_END - CACHE_START) / 1000000 ))

echo "Cache hit latency: ${CACHE_LATENCY}ms"

if (( CACHE_LATENCY < 10 )); then
    echo "COMPLETED SUCCESS: Cache performance is excellent (<10ms)"
else
    echo "WARNING: Cache performance could be improved (${CACHE_LATENCY}ms)"
fi

# Test Redis persistence
echo
echo "💾 Testing Redis persistence..."
REDIS_TEST_SERVICE="redis-test-service"
curl -s -X PUT "$REGISTRY_URL/v1/agent/service/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"id\": \"$REDIS_TEST_SERVICE\",
        \"name\": \"redis-test\",
        \"address\": \"127.0.0.1\",
        \"port\": 9999,
        \"tags\": [\"redis\", \"test\"]
    }" > /dev/null

# Check if service is discoverable
if curl -s "$REGISTRY_URL/v1/catalog/services" | grep -q "$REDIS_TEST_SERVICE"; then
    echo "COMPLETED SUCCESS: Service persisted and discoverable"
else
    echo "CRITICAL FAILED: Service not found after registration"
    SUCCESS=false
fi

# Cleanup test services
echo
echo "🧹 Cleaning up test services..."
for i in $(seq 1 $TEST_SERVICES); do
    curl -s -X PUT "$REGISTRY_URL/v1/agent/service/deregister/test-service-$i" > /dev/null
done
curl -s -X PUT "$REGISTRY_URL/v1/agent/service/deregister/$REDIS_TEST_SERVICE" > /dev/null

# Cleanup temp files
rm -rf "$TEMP_DIR"

echo "COMPLETED Cleanup complete"
echo

# Final result
if [ "$SUCCESS" = true ]; then
    echo "🎉 PERFORMANCE TEST PASSED"
    echo "Registry meets <50ms discovery latency target"
    exit 0
else
    echo "💥 PERFORMANCE TEST FAILED"
    echo "Registry does not meet performance targets"
    exit 1
fi
