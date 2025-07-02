#!/bin/bash

# fr0g.ai Master Control Program - Webhook Test Commands
# Collection of curl commands for testing the MCP webhook system

BASE_URL="http://localhost:8081"

echo "🧪 fr0g.ai MCP Webhook Test Suite"
echo "=================================="
echo ""

# Health Check
echo "1. Testing Health Check..."
curl -X GET "$BASE_URL/health" \
  -H "Accept: application/json" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# System Status
echo "2. Testing System Status..."
curl -X GET "$BASE_URL/status" \
  -H "Accept: application/json" \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Simple Discord Message
echo "3. Testing Simple Discord Message..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello from the AI community! This is a simple test message.",
    "author": {
      "username": "test_user",
      "id": "12345",
      "avatar": "avatar_hash",
      "bot": false
    },
    "channel_id": "general",
    "guild_id": "test_guild",
    "timestamp": "2025-07-02T12:00:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Technical Discussion
echo "4. Testing Technical Discussion..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "I have implemented a new recursive neural network algorithm for pattern recognition. The performance improvements are significant - we are seeing 40% faster processing with 25% better accuracy.",
    "author": {
      "username": "ai_researcher",
      "id": "67890",
      "avatar": "researcher_avatar",
      "bot": false
    },
    "channel_id": "tech-discussion",
    "guild_id": "ai_community",
    "timestamp": "2025-07-02T12:15:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Urgent Message
echo "5. Testing Urgent Message..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "I have an urgent message about system performance - we need to review the cognitive engine patterns immediately!",
    "author": {
      "username": "system_admin",
      "id": "11111",
      "avatar": "admin_avatar",
      "bot": false
    },
    "channel_id": "alerts",
    "guild_id": "operations",
    "timestamp": "2025-07-02T12:30:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# AI Consciousness Discussion
echo "6. Testing AI Consciousness Discussion..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "The AI personas are working together beautifully! This distributed intelligence approach is fascinating. I am observing emergent behaviors that suggest genuine collaborative reasoning.",
    "author": {
      "username": "consciousness_researcher",
      "id": "22222",
      "avatar": "researcher2_avatar",
      "bot": false
    },
    "channel_id": "ai-consciousness",
    "guild_id": "research_lab",
    "timestamp": "2025-07-02T12:45:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Code Review Request
echo "7. Testing Code Review Request..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "```go\nfunc (ce *CognitiveEngine) recognizePatterns() {\n    patterns := ce.analyzeSystemBehavior()\n    for _, pattern := range patterns {\n        ce.storePattern(pattern)\n    }\n}\n```\nCan the AI community review this cognitive pattern recognition code?",
    "author": {
      "username": "developer",
      "id": "33333",
      "avatar": "dev_avatar",
      "bot": false
    },
    "channel_id": "code-review",
    "guild_id": "dev_team",
    "timestamp": "2025-07-02T13:00:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

echo "✅ Test suite completed!"
echo ""
echo "💡 Tips:"
echo "- Make sure the MCP is running: cd fr0g-ai-master-control && go run cmd/mcp-demo/main.go"
echo "- Check logs for AI persona reviews and cognitive insights"
echo "- Monitor system consciousness and pattern recognition"
echo "- Each test should trigger AI community analysis"
