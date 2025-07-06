#!/bin/bash

# Test script for chat completions endpoint
# Usage: ./test_chat_completions.sh [BASE_URL]

BASE_URL=${1:-"http://localhost:8080"}
ENDPOINT="$BASE_URL/v1/chat/completions"

echo "Testing Chat Completions API at: $ENDPOINT"
echo "================================================"

# Test 1: Basic chat completion
echo "Test 1: Basic chat completion"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 2: Multi-turn conversation
echo "Test 2: Multi-turn conversation"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "What is the capital of France?"},
      {"role": "assistant", "content": "The capital of France is Paris."},
      {"role": "user", "content": "What about Germany?"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 3: With persona prompt
echo "Test 3: With persona prompt"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Explain quantum computing"}
    ],
    "persona_prompt": "You are a physics professor explaining complex topics to undergraduate students."
  }' | jq '.'

echo -e "\n================================================"

# Test 4: With optional parameters
echo "Test 4: With optional parameters"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Write a creative story"}
    ],
    "temperature": 0.8,
    "max_tokens": 150
  }' | jq '.'

echo -e "\n================================================"

# Test 5: Error case - missing model
echo "Test 5: Error case - missing model"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "Hello"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 6: Error case - invalid role
echo "Test 6: Error case - invalid role"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "invalid", "content": "Hello"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 7: Legacy endpoint compatibility
echo "Test 7: Legacy endpoint compatibility"
curl -X POST "$BASE_URL/api/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Test legacy endpoint"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 8: Health check
echo "Test 8: Health check"
curl -X GET "$BASE_URL/health" | jq '.'

echo -e "\n================================================"
echo "Testing complete!"
#!/bin/bash

# Test script for chat completions endpoint
# Usage: ./test_chat_completions.sh [BASE_URL]

BASE_URL=${1:-"http://localhost:8080"}
ENDPOINT="$BASE_URL/v1/chat/completions"

echo "Testing Chat Completions API at: $ENDPOINT"
echo "================================================"

# Test 1: Basic chat completion
echo "Test 1: Basic chat completion"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 2: Multi-turn conversation
echo "Test 2: Multi-turn conversation"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "What is the capital of France?"},
      {"role": "assistant", "content": "The capital of France is Paris."},
      {"role": "user", "content": "What about Germany?"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 3: With persona prompt
echo "Test 3: With persona prompt"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Explain quantum computing"}
    ],
    "persona_prompt": "You are a physics professor explaining complex topics to undergraduate students."
  }' | jq '.'

echo -e "\n================================================"

# Test 4: With optional parameters
echo "Test 4: With optional parameters"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Write a creative story"}
    ],
    "temperature": 0.8,
    "max_tokens": 150
  }' | jq '.'

echo -e "\n================================================"

# Test 5: Error case - missing model
echo "Test 5: Error case - missing model"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "Hello"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 6: Error case - invalid role
echo "Test 6: Error case - invalid role"
curl -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "invalid", "content": "Hello"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 7: Legacy endpoint compatibility
echo "Test 7: Legacy endpoint compatibility"
curl -X POST "$BASE_URL/api/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Test legacy endpoint"}
    ]
  }' | jq '.'

echo -e "\n================================================"

# Test 8: Health check
echo "Test 8: Health check"
curl -X GET "$BASE_URL/health" | jq '.'

echo -e "\n================================================"
echo "Testing complete!"
