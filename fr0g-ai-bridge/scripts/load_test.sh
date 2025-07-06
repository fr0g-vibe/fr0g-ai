#!/bin/bash

# Load test script for chat completions endpoint
# Usage: ./load_test.sh [BASE_URL] [CONCURRENT_REQUESTS] [TOTAL_REQUESTS]

BASE_URL=${1:-"http://localhost:8080"}
CONCURRENT=${2:-5}
TOTAL=${3:-50}
ENDPOINT="$BASE_URL/v1/chat/completions"

echo "Load testing Chat Completions API"
echo "Endpoint: $ENDPOINT"
echo "Concurrent requests: $CONCURRENT"
echo "Total requests: $TOTAL"
echo "================================================"

# Create test payload
TEST_PAYLOAD='{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "user", "content": "Hello, this is a load test message."}
  ]
}'

# Function to make a single request
make_request() {
  local id=$1
  local start_time=$(date +%s.%N)
  
  response=$(curl -s -w "%{http_code}" -X POST "$ENDPOINT" \
    -H "Content-Type: application/json" \
    -d "$TEST_PAYLOAD")
  
  local end_time=$(date +%s.%N)
  local duration=$(echo "$end_time - $start_time" | bc)
  local http_code="${response: -3}"
  
  echo "Request $id: HTTP $http_code, Duration: ${duration}s"
  
  if [ "$http_code" != "200" ]; then
    echo "  Error response: ${response%???}"
  fi
}

# Export function for parallel execution
export -f make_request
export ENDPOINT
export TEST_PAYLOAD

# Run load test
echo "Starting load test..."
seq 1 $TOTAL | xargs -n 1 -P $CONCURRENT -I {} bash -c 'make_request {}'

echo "================================================"
echo "Load test complete!"

# Test health endpoint under load
echo "Testing health endpoint..."
for i in {1..10}; do
  curl -s "$BASE_URL/health" | jq -r '.status' | head -1
done
#!/bin/bash

# Load test script for chat completions endpoint
# Usage: ./load_test.sh [BASE_URL] [CONCURRENT_REQUESTS] [TOTAL_REQUESTS]

BASE_URL=${1:-"http://localhost:8080"}
CONCURRENT=${2:-5}
TOTAL=${3:-50}
ENDPOINT="$BASE_URL/v1/chat/completions"

echo "Load testing Chat Completions API"
echo "Endpoint: $ENDPOINT"
echo "Concurrent requests: $CONCURRENT"
echo "Total requests: $TOTAL"
echo "================================================"

# Create test payload
TEST_PAYLOAD='{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "user", "content": "Hello, this is a load test message."}
  ]
}'

# Function to make a single request
make_request() {
  local id=$1
  local start_time=$(date +%s.%N)
  
  response=$(curl -s -w "%{http_code}" -X POST "$ENDPOINT" \
    -H "Content-Type: application/json" \
    -d "$TEST_PAYLOAD")
  
  local end_time=$(date +%s.%N)
  local duration=$(echo "$end_time - $start_time" | bc)
  local http_code="${response: -3}"
  
  echo "Request $id: HTTP $http_code, Duration: ${duration}s"
  
  if [ "$http_code" != "200" ]; then
    echo "  Error response: ${response%???}"
  fi
}

# Export function for parallel execution
export -f make_request
export ENDPOINT
export TEST_PAYLOAD

# Run load test
echo "Starting load test..."
seq 1 $TOTAL | xargs -n 1 -P $CONCURRENT -I {} bash -c 'make_request {}'

echo "================================================"
echo "Load test complete!"

# Test health endpoint under load
echo "Testing health endpoint..."
for i in {1..10}; do
  curl -s "$BASE_URL/health" | jq -r '.status' | head -1
done
