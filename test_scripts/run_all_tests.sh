#!/bin/bash

# Master test script to run all webhook processor tests
# This script runs all individual test scripts in sequence

echo "=========================================="
echo "FR0G-AI Master Control Webhook Tests"
echo "=========================================="
echo

# Check if the server is running
echo "Checking if webhook server is running..."
if ! curl -s "http://localhost:8080/health" > /dev/null 2>&1; then
    echo "❌ Webhook server is not running on localhost:8080"
    echo ""
    echo "To start the webhook server:"
    echo "1. cd fr0g-ai-master-control"
    echo "2. go mod tidy"
    echo "3. go run cmd/webhook-server/main.go"
    echo ""
    echo "Then run this test script again."
    exit 1
fi
echo "✅ Webhook server is running"
echo

# Make all test scripts executable
chmod +x test_scripts/*.sh

# Run IRC tests
echo "🔄 Running IRC webhook tests..."
./test_scripts/test_irc_webhook.sh
echo

# Run SMS tests
echo "🔄 Running SMS webhook tests..."
./test_scripts/test_sms_webhook.sh
echo

# Run Voice tests
echo "🔄 Running Voice webhook tests..."
./test_scripts/test_voice_webhook.sh
echo

# Run ESMTP tests
echo "🔄 Running ESMTP webhook tests..."
./test_scripts/test_esmtp_webhook.sh
echo

echo "=========================================="
echo "All webhook processor tests completed!"
echo "=========================================="
echo
echo "Test Summary:"
echo "- IRC Processor: Chat threat analysis"
echo "- SMS Processor: SMS/MMS threat analysis"
echo "- Voice Processor: Voice call threat analysis"
echo "- ESMTP Processor: Email threat analysis"
echo
echo "Each processor tests various threat scenarios including:"
echo "- Phishing attempts"
echo "- Social engineering"
echo "- Malware distribution"
echo "- Scam attempts"
echo "- Trusted vs untrusted sources"
echo
echo "Review the JSON responses above to see threat analysis results."
