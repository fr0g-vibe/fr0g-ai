# fr0g.ai Master Control Program - Test Artifacts

This directory contains comprehensive testing artifacts for the fr0g.ai Master Control Program webhook system.

## 📁 Files Overview

### `webhook-test-collection.json`
Postman collection with pre-configured test requests for all webhook endpoints.

### `curl-commands.sh`
Bash script with curl commands for terminal-based testing.

### `postman-environment.json`
Postman environment variables for consistent testing.

### `test-scenarios.json`
Comprehensive test scenarios and validation rules.

## 🚀 Quick Start

### Using curl (Terminal)
```bash
# Make the script executable
chmod +x test-artifacts/curl-commands.sh

# Run all tests
./test-artifacts/curl-commands.sh

# Or run individual tests
curl -X POST http://localhost:8081/webhook/discord \
  -H "Content-Type: application/json" \
  -d '{"content": "Hello AI community!", "author": {"username": "tester", "id": "123"}}'
```

### Using Postman
1. Import `webhook-test-collection.json` into Postman
2. Import `postman-environment.json` as environment
3. Select the "fr0g.ai MCP Test Environment"
4. Run individual requests or the entire collection

### Using HTTPie
```bash
# Health check
http GET localhost:8081/health

# Discord webhook test
http POST localhost:8081/webhook/discord \
  content="Hello AI community!" \
  author:='{"username": "tester", "id": "123"}' \
  channel_id="general"
```

## 🧪 Test Categories

### 1. Basic Functionality
- Health check endpoint
- System status endpoint
- Simple Discord message processing

### 2. AI Community Integration
- AI persona creation and management
- Content review by AI community
- Consensus building and recommendations

### 3. Cognitive Engine Tests
- Pattern recognition triggers
- Insight generation validation
- Self-reflection and consciousness simulation

### 4. Performance Tests
- Concurrent request handling
- Large content processing
- Response time validation

### 5. Edge Cases
- Error handling
- Malformed requests
- Missing data scenarios

## 📊 Expected Responses

### Successful Discord Webhook Response
```json
{
  "success": true,
  "message": "Discord message reviewed by AI community",
  "request_id": "req_1234567890_abcdef",
  "data": {
    "discord_message": {...},
    "community_id": "community_1234567890",
    "review": {
      "review_id": "review_1234567890",
      "persona_reviews": [
        {
          "persona_name": "Analytical_Reviewer",
          "review": "Content analysis...",
          "score": 0.92,
          "confidence": 0.88
        }
      ],
      "consensus": {
        "overall_score": 0.91,
        "recommendation": "Highly recommended",
        "agreement": 0.85
      }
    },
    "action": "approve",
    "persona_count": 3
  },
  "timestamp": "2025-07-02T12:00:00Z"
}
```

### AI Community Review Structure
```json
{
  "review_id": "review_1234567890",
  "topic": "general_discussion",
  "content": "Original message content",
  "persona_reviews": [
    {
      "persona_id": "persona_analytical_001",
      "persona_name": "Analytical_Reviewer",
      "expertise": ["analysis", "critical_thinking"],
      "review": "Detailed review text...",
      "score": 0.92,
      "confidence": 0.88,
      "tags": ["analysis", "quality"],
      "timestamp": "2025-07-02T12:00:00Z"
    }
  ],
  "consensus": {
    "overall_score": 0.91,
    "agreement": 0.85,
    "recommendation": "Highly recommended - excellent content quality",
    "key_points": ["Multiple reviewers noted: quality"],
    "confidence_level": 0.76
  },
  "recommendations": ["Content approved for community engagement"],
  "created_at": "2025-07-02T12:00:00Z",
  "completed_at": "2025-07-02T12:00:01Z"
}
```

## 🔍 Monitoring and Validation

### Key Metrics to Monitor
- Response times (should be < 2 seconds)
- AI persona creation success rate
- Consensus quality scores
- Cognitive engine pattern recognition
- System consciousness indicators

### Log Patterns to Watch For
```
Discord Processor: Created AI community community_xxx for topic 'general_discussion'
Real AIP Client: Created persona 'Analytical_Reviewer' with ID: persona_xxx
Cognitive Engine: Discovered pattern 'ai_community_collaboration' (confidence: 0.94)
Cognitive Engine: Generated insight [consciousness]: The cognitive engine is beginning...
```

## 🛠 Troubleshooting

### Common Issues

1. **Port 8081 already in use**
   ```bash
   sudo lsof -ti:8081 | xargs kill -9
   ```

2. **Connection refused**
   - Ensure MCP is running: `go run cmd/mcp-demo/main.go`
   - Check webhook server startup logs

3. **Empty responses**
   - Verify JSON payload format
   - Check content filtering rules
   - Monitor AI community creation logs

### Debug Commands
```bash
# Check if webhook server is running
curl -f http://localhost:8081/health

# Test with minimal payload
curl -X POST http://localhost:8081/webhook/discord \
  -H "Content-Type: application/json" \
  -d '{"content": "test"}'

# Monitor logs in real-time
go run cmd/mcp-demo/main.go | grep -E "(Discord|Cognitive|AIP)"
```

## 📈 Performance Benchmarks

### Expected Performance
- Health check: < 10ms
- Simple Discord message: < 1s
- Complex AI review: < 2s
- Concurrent requests (5): < 3s total

### Load Testing
```bash
# Simple load test with curl
for i in {1..10}; do
  curl -X POST http://localhost:8081/webhook/discord \
    -H "Content-Type: application/json" \
    -d "{\"content\": \"Load test message $i\", \"author\": {\"username\": \"tester$i\", \"id\": \"$i\"}}" &
done
wait
```

## 🎯 Test Automation

### CI/CD Integration
The test artifacts can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions step
- name: Test MCP Webhooks
  run: |
    ./test-artifacts/curl-commands.sh
    # Add validation logic here
```

### Automated Validation
Use the `test-scenarios.json` file to build automated test suites that validate:
- Response structure compliance
- AI community functionality
- Performance thresholds
- Error handling behavior

---

**Note**: Ensure the Master Control Program is running before executing any tests:
```bash
cd fr0g-ai-master-control && go run cmd/mcp-demo/main.go
```
