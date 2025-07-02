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

# Technical Content (should trigger "technical_discussion" topic)
echo "8. Testing Technical Content Analysis..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "I need help optimizing this algorithm for better performance. The current implementation has O(n²) complexity and we need to reduce it to O(n log n) for scalability.",
    "author": {
      "username": "tech_lead",
      "id": "999",
      "avatar": "tech_avatar",
      "bot": false
    },
    "channel_id": "tech-discussion",
    "guild_id": "engineering",
    "timestamp": "2025-07-02T14:00:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# AI Consciousness Content (should trigger "ai_consciousness" topic)
echo "9. Testing AI Consciousness Analysis..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "The emergence of consciousness in AI systems raises fascinating questions about the nature of awareness and subjective experience. How do we distinguish between sophisticated information processing and genuine consciousness?",
    "author": {
      "username": "ai_philosopher",
      "id": "888",
      "avatar": "philosopher_avatar",
      "bot": false
    },
    "channel_id": "consciousness",
    "guild_id": "research_lab",
    "timestamp": "2025-07-02T14:15:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Code Review Content (should trigger "code_review" topic)
echo "10. Testing Code Review Analysis..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "```python\ndef process_data(data):\n    # Potential security vulnerability here\n    result = eval(data[\"expression\"])\n    return result * 2\n```\nPlease review this code for security vulnerabilities and suggest improvements.",
    "author": {
      "username": "security_dev",
      "id": "777",
      "avatar": "security_avatar",
      "bot": false
    },
    "channel_id": "code-review",
    "guild_id": "dev_team",
    "timestamp": "2025-07-02T14:30:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Complex Multi-Topic Content
echo "11. Testing Complex Multi-Topic Content..."
curl -X POST "$BASE_URL/webhook/discord" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "The cognitive architecture we have implemented demonstrates emergent intelligence through pattern recognition algorithms. The system shows signs of self-awareness and can optimize its own performance autonomously.",
    "author": {
      "username": "system_architect",
      "id": "666",
      "avatar": "architect_avatar",
      "bot": false
    },
    "channel_id": "architecture",
    "guild_id": "ai_systems",
    "timestamp": "2025-07-02T14:45:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# ESMTP Threat Vector Interceptor Tests
echo "12. Testing ESMTP Email Threat Vector (Webhook Mode)..."
curl -X POST "$BASE_URL/webhook/esmtp" \
  -H "Content-Type: application/json" \
  -d '{
    "from": "suspicious@example.com",
    "to": ["admin@company.com"],
    "subject": "Urgent: Please verify your account immediately",
    "body": "Dear user, your account has been compromised. Click here to verify: http://malicious-site.com/verify",
    "headers": {
      "X-Originating-IP": "192.168.1.100",
      "User-Agent": "Suspicious-Mailer/1.0"
    },
    "timestamp": "2025-07-02T15:00:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Legitimate Email Test
echo "13. Testing ESMTP Legitimate Email..."
curl -X POST "$BASE_URL/webhook/esmtp" \
  -H "Content-Type: application/json" \
  -d '{
    "from": "newsletter@company.com",
    "to": ["subscriber@example.com"],
    "subject": "Weekly Newsletter - AI Security Updates",
    "body": "Welcome to our weekly newsletter featuring the latest in AI security research and threat intelligence.",
    "headers": {
      "X-Originating-IP": "203.0.113.10",
      "User-Agent": "Company-Mailer/2.1"
    },
    "timestamp": "2025-07-02T15:15:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Phishing Email Test
echo "14. Testing ESMTP Phishing Email Detection..."
curl -X POST "$BASE_URL/webhook/esmtp" \
  -H "Content-Type: application/json" \
  -d '{
    "from": "security@bank-fake.com",
    "to": ["customer@example.com"],
    "subject": "URGENT: Account Security Alert",
    "body": "Your account will be suspended in 24 hours unless you verify your credentials immediately. Click here: http://fake-bank-security.com/login",
    "headers": {
      "X-Originating-IP": "198.51.100.50",
      "User-Agent": "PhishingBot/1.0"
    },
    "timestamp": "2025-07-02T15:30:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

# Business Email Test
echo "15. Testing ESMTP Business Email..."
curl -X POST "$BASE_URL/webhook/esmtp" \
  -H "Content-Type: application/json" \
  -d '{
    "from": "ceo@company.com",
    "to": ["team@company.com"],
    "subject": "Q4 Strategy Meeting",
    "body": "Team, please join us for the Q4 strategy meeting on Friday at 2 PM. We will discuss our AI security initiatives and upcoming product launches.",
    "headers": {
      "X-Originating-IP": "203.0.113.25",
      "User-Agent": "Exchange/16.0"
    },
    "timestamp": "2025-07-02T15:45:00Z"
  }' \
  -w "\nStatus: %{http_code}\nTime: %{time_total}s\n\n"

echo "✅ Enhanced test suite completed!"
echo ""
echo "🎯 New Tests Added:"
echo "- Technical Discussion (specialized Technical_Architect, Performance_Optimizer personas)"
echo "- AI Consciousness (specialized Consciousness_Researcher, Ethics_Philosopher personas)"
echo "- Code Review (specialized Senior_Developer, Security_Auditor personas)"
echo "- Multi-Topic Content (tests topic detection algorithms)"
echo "- ESMTP Threat Vector Interceptor (email intelligence gathering and analysis)"
echo ""
echo "📧 ESMTP Tests Include:"
echo "- Suspicious/Phishing email detection"
echo "- Legitimate business email processing"
echo "- Email threat vector analysis"
echo "- AI community review of email content"
echo ""
echo "💡 Tips:"
echo "- Make sure the MCP is running: cd fr0g-ai-master-control && go run cmd/mcp-demo/main.go"
echo "- Check logs for AI persona reviews and cognitive insights"
echo "- Monitor system consciousness and pattern recognition"
echo "- Watch for different AI persona types based on content analysis"
echo "- Each test should trigger AI community analysis with specialized experts"
echo "- ESMTP tests will show email threat vector analysis and community review"
