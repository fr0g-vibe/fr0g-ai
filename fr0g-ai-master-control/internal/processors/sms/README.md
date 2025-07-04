# SMS Threat Detection Processor

The SMS processor is a comprehensive threat detection system for analyzing SMS messages and identifying potential security threats including spam, phishing, malware, and social engineering attacks.

## Features

### Threat Detection
- **Spam Detection**: Keyword-based filtering with confidence scoring
- **Phishing Detection**: URL analysis and social engineering pattern recognition
- **Malware Detection**: Suspicious download/install request identification
- **Social Engineering**: Account verification and urgency-based attack detection

### Integration Options
- **Google Voice API**: Direct integration with Google Voice for message monitoring
- **Webhook Server**: HTTP endpoint for receiving SMS messages from external services
- **Real-time Processing**: Continuous message analysis with configurable intervals

### Advanced Analytics
- **Phone Number Tracking**: Reputation scoring and behavioral analysis
- **Message History**: Configurable message retention with threat correlation
- **Threat Scoring**: Multi-factor confidence calculation with severity levels
- **Pattern Recognition**: Regex-based threat pattern matching

## Configuration

```yaml
sms:
  enabled: true
  google_voice_enabled: false
  google_voice_api_key: "your-api-key"
  webhook_enabled: true
  webhook_port: 8082
  webhook_path: "/sms/webhook"
  processing_interval: 30
  max_history_size: 1000
  threat_threshold: 0.5
  spam_filter_enabled: true
  blacklist_enabled: true
  whitelist_enabled: true
```

## Usage

### Basic Usage

```go
import "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/processors/sms"

// Create processor with configuration
cfg := &config.SMSConfig{
    Enabled: true,
    ProcessingInterval: 30,
    MaxHistorySize: 1000,
}

processor := sms.NewProcessor(cfg)

// Start processing
ctx := context.Background()
err := processor.Start(ctx)
if err != nil {
    log.Fatal(err)
}

// Process a message
message := sms.SMSMessage{
    ID: "msg123",
    From: "+1234567890",
    To: "+0987654321",
    Body: "Suspicious message content",
    Timestamp: time.Now(),
}

result, err := processor.ProcessMessage(message)
if err != nil {
    log.Printf("Processing failed: %v", err)
} else {
    log.Printf("Threat Level: %s, Confidence: %.2f", 
        result.ThreatLevel.String(), result.Analysis.Confidence)
}
```

### Threat Analysis Results

```go
type ThreatAnalysis struct {
    ThreatTypes     []string  // Detected threat categories
    Confidence      float64   // Overall confidence (0.0-1.0)
    SpamScore       float64   // Spam likelihood
    PhishingScore   float64   // Phishing likelihood
    MalwareScore    float64   // Malware likelihood
    SocialEngScore  float64   // Social engineering likelihood
    Indicators      []string  // Specific threat indicators
    Recommendations []string  // Security recommendations
    ProcessedAt     time.Time // Analysis timestamp
}
```

## Threat Levels

- **None**: No threats detected (confidence < 0.2)
- **Low**: Minor suspicious indicators (0.2 ≤ confidence < 0.4)
- **Medium**: Moderate threat indicators (0.4 ≤ confidence < 0.6)
- **High**: Strong threat indicators (0.6 ≤ confidence < 0.8)
- **Critical**: Severe threat detected (confidence ≥ 0.8)

## Detection Patterns

### Spam Indicators
- Excessive use of promotional keywords (free, winner, prize)
- Urgency language (urgent, immediate, limited time)
- Excessive punctuation and capitalization
- Multiple spam keywords in single message

### Phishing Indicators
- Shortened URLs (bit.ly, tinyurl, etc.)
- Account verification requests
- Urgency combined with action requests
- Financial information solicitation

### Malware Indicators
- Download/install requests
- Software update notifications
- Suspicious link clicking instructions
- File attachment references

### Social Engineering Indicators
- Security alerts and account suspension notices
- Identity verification requests
- Information update requirements
- Locked account notifications

## Phone Number Tracking

The processor maintains reputation scores for phone numbers based on:
- Message frequency and patterns
- Historical threat detection
- Blacklist/whitelist status
- Behavioral analysis over time

## Statistics and Monitoring

```go
stats := processor.GetStats()
// Returns:
// - total_messages: Total processed messages
// - unique_numbers: Unique phone numbers seen
// - threat_counts: Breakdown by threat level
// - is_running: Processor status
// - google_voice_enabled: Integration status
// - webhook_enabled: Webhook server status
```

## Integration Examples

### Google Voice Integration
```go
cfg.GoogleVoiceEnabled = true
cfg.GoogleVoiceAPIKey = "your-api-key"
```

### Webhook Server
```go
cfg.WebhookEnabled = true
cfg.WebhookPort = 8082
cfg.WebhookPath = "/sms/webhook"
```

## Security Considerations

- API keys should be stored securely and rotated regularly
- Webhook endpoints should use HTTPS and authentication
- Message content should be handled according to privacy regulations
- Threat detection results should be logged securely
- Phone number data should comply with privacy requirements

## Testing

Run the test suite:
```bash
cd fr0g-ai-master-control
go test ./internal/processors/sms/...
```

## Future Enhancements

- Machine learning-based threat detection
- Integration with external threat intelligence feeds
- Advanced natural language processing
- Real-time threat correlation across multiple vectors
- Automated response and blocking capabilities
