# FR0G-AI Webhook Processor Test Scripts

This directory contains comprehensive test scripts for all webhook processors in the FR0G-AI Master Control system.

## Prerequisites

1. **Server Running**: Ensure the webhook server is running on `localhost:8080`
2. **Dependencies**: 
   - `curl` - for making HTTP requests
   - `jq` - for pretty-printing JSON responses
3. **Permissions**: Make scripts executable with `chmod +x test_scripts/*.sh`

## Test Scripts

### Individual Processor Tests

- **`test_irc_webhook.sh`** - Tests IRC chat message threat analysis
- **`test_sms_webhook.sh`** - Tests SMS/MMS message threat analysis  
- **`test_voice_webhook.sh`** - Tests voice call threat analysis
- **`test_esmtp_webhook.sh`** - Tests email threat analysis

### Master Test Script

- **`run_all_tests.sh`** - Runs all processor tests in sequence

## Running Tests

### Run All Tests
```bash
./test_scripts/run_all_tests.sh
```

### Run Individual Tests
```bash
./test_scripts/test_irc_webhook.sh
./test_scripts/test_sms_webhook.sh
./test_scripts/test_voice_webhook.sh
./test_scripts/test_esmtp_webhook.sh
```

## Test Scenarios

Each processor is tested with multiple threat scenarios:

### IRC Processor Tests
1. **Basic channel message** - Suspicious link sharing
2. **Private message phishing** - Credential harvesting attempt
3. **Trusted user message** - Lower threat threshold testing
4. **Bot spam** - Cryptocurrency pump scheme
5. **Malware distribution** - File sharing with malicious links

### SMS Processor Tests
1. **SMS phishing** - Fake bank security alert
2. **MMS with media** - Suspicious media attachments
3. **Cryptocurrency scam** - Bitcoin giveaway fraud
4. **Trusted number** - Legitimate bank notification
5. **Romance scam** - Military deployment fraud
6. **2FA bypass** - Social engineering for verification codes

### Voice Processor Tests
1. **Tech support scam** - Microsoft impersonation
2. **IRS impersonation** - Government authority fraud
3. **Voice deepfake** - AI-generated family emergency
4. **Trusted caller** - Legitimate bank notification
5. **Investment scam** - Cryptocurrency opportunity fraud

### ESMTP Processor Tests
1. **Bank phishing** - Account verification fraud
2. **BEC attack** - CEO impersonation for wire transfer
3. **Malware attachment** - Suspicious executable files
4. **Trusted domain** - Legitimate GitHub notification
5. **Crypto scam** - Fake Elon Musk giveaway
6. **Ransomware** - Legal notice with malicious PDF

## Expected Responses

Each test should return a JSON response with:

```json
{
  "success": true,
  "message": "Threat vector submitted for community review",
  "request_id": "test_id",
  "data": {
    "threat_level": "high|medium|low|minimal|critical",
    "consensus": 0.85,
    "review_id": "review_12345",
    "is_trusted": false,
    // ... processor-specific data
  },
  "timestamp": "2023-12-01T12:00:00Z"
}
```

## Threat Level Interpretation

- **Critical** (0.9+): Immediate threat requiring action
- **High** (0.8+): Significant threat, high confidence
- **Medium** (0.6+): Moderate threat, investigate further
- **Low** (0.4+): Minor threat, monitor
- **Minimal** (<0.4): Low risk, likely benign

## Trusted vs Untrusted Sources

The tests include both trusted and untrusted sources to verify:
- **Trusted sources** have higher threat thresholds
- **Blocked sources** are immediately rejected
- **Unknown sources** use standard thresholds

## Troubleshooting

### Server Not Running
```bash
# Check if server is running
curl -s http://localhost:8080/health

# If not running, start the webhook server first
```

### Missing Dependencies
```bash
# Install jq on Ubuntu/Debian
sudo apt-get install jq

# Install jq on macOS
brew install jq
```

### Permission Denied
```bash
# Make scripts executable
chmod +x test_scripts/*.sh
```

## Customization

To modify test scenarios:

1. **Edit webhook payloads** in the test scripts
2. **Add new test cases** by copying existing curl commands
3. **Modify threat scenarios** by changing message content
4. **Test different configurations** by updating processor settings

## Integration with CI/CD

These scripts can be integrated into CI/CD pipelines:

```bash
# Run tests and capture exit code
./test_scripts/run_all_tests.sh
if [ $? -eq 0 ]; then
    echo "All tests passed"
else
    echo "Tests failed"
    exit 1
fi
```
