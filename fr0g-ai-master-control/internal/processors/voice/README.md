# Voice Threat Detection Processor

The Voice processor is a comprehensive threat detection system for analyzing voice calls and identifying potential security threats including scam calls, phishing attempts, robocalls, and social engineering attacks.

## Features

### Threat Detection
- **Scam Detection**: Pattern-based analysis for common scam indicators
- **Phishing Detection**: Information request and verification scam identification
- **Social Engineering**: Authority impersonation and fear tactic detection
- **Robocall Detection**: Automated call pattern recognition
- **Emotional Manipulation**: Pressure tactic and urgency indicator analysis

### Integration Options
- **Speech-to-Text**: Real-time transcript generation from audio streams
- **Call Recording**: Audio capture and storage with compliance features
- **Telephony Integration**: Support for Twilio, Asterisk, and other providers
- **Webhook Server**: HTTP endpoint for receiving call events from external services

### Advanced Analytics
- **Caller Tracking**: Reputation scoring and behavioral analysis
- **Call History**: Configurable call retention with threat correlation
- **Speech Pattern Analysis**: Script reading, repetition, and density detection
- **Threat Scoring**: Multi-factor confidence calculation with severity levels

## Configuration

```yaml
voice:
  enabled: true
  speech_to_text_enabled: false
  speech_to_text_api_key: "your-api-key"
  call_recording_enabled: false
  recording_path: "/tmp/voice_recordings"
  monitoring_interval: 60
  max_history_size: 500
  threat_threshold: 0.5
  telephony_provider: "twilio"
  telephony_api_key: "your-telephony-key"
  webhook_enabled: true
  webhook_port: 8083
  webhook_path: "/voice/webhook"
```

## Usage

### Basic Usage

```go
import "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/processors/voice"

// Create processor with configuration
cfg := &config.VoiceConfig{
    Enabled: true,
    MonitoringInterval: 60,
    MaxHistorySize: 500,
}

processor := voice.NewProcessor(cfg)

// Start processing
ctx := context.Background()
err := processor.Start(ctx)
if err != nil {
    log.Fatal(err)
}

// Process a call
call := voice.VoiceCall{
    ID: "call123",
    CallerID: "+1234567890",
    StartTime: time.Now().Add(-5 * time.Minute),
    EndTime: time.Now(),
    Duration: 5 * time.Minute,
    Transcript: "This is the IRS calling about your tax refund...",
}

result, err := processor.ProcessCall(call)
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
    ThreatTypes        []string  // Detected threat categories
    Confidence         float64   // Overall confidence (0.0-1.0)
    ScamScore          float64   // General scam likelihood
    PhishingScore      float64   // Information theft likelihood
    SocialEngScore     float64   // Social engineering likelihood
    RobocallScore      float64   // Automated call likelihood
    EmotionalManipScore float64  // Emotional manipulation likelihood
    Indicators         []string  // Specific threat indicators
    Recommendations    []string  // Security recommendations
    SpeechPatterns     []string  // Detected speech patterns
    ProcessedAt        time.Time // Analysis timestamp
}
```

## Threat Levels

- **None**: No threats detected (confidence < 0.2)
- **Low**: Minor suspicious indicators (0.2 ≤ confidence < 0.4)
- **Medium**: Moderate threat indicators (0.4 ≤ confidence < 0.6)
- **High**: Strong threat indicators (0.6 ≤ confidence < 0.8)
- **Critical**: Severe threat detected (confidence ≥ 0.8)

## Detection Patterns

### Scam Indicators
- IRS/tax-related threats and refund claims
- Tech support impersonation (Microsoft, Windows)
- Bank fraud and account suspension notices
- Social Security and Medicare scams
- Lottery and prize winner notifications

### Phishing Indicators
- Requests for personal information (SSN, bank details)
- Account verification and confirmation requests
- Password and PIN number solicitation
- Credit card and financial information requests

### Social Engineering Indicators
- Authority impersonation (IRS, FBI, police, banks)
- Fear tactics (arrest, warrant, legal action)
- Urgency and pressure tactics
- Emergency and crisis scenarios

### Robocall Indicators
- Very short call duration (< 30 seconds)
- "Press 1" or "Press 9" instructions
- "This is not a sales call" disclaimers
- Automated message announcements

### Speech Pattern Analysis
- Script reading detection (formal punctuation patterns)
- Repetitive word usage
- High word density (fast talking)
- Automated speech characteristics

## Caller Tracking

The processor maintains reputation scores for callers based on:
- Call frequency and duration patterns
- Historical threat detection
- Blacklist/whitelist status
- Average call length analysis
- Behavioral pattern recognition

## Statistics and Monitoring

```go
stats := processor.GetStats()
// Returns:
// - total_calls: Total processed calls
// - unique_callers: Unique caller IDs seen
// - threat_counts: Breakdown by threat level
// - average_call_duration: Mean call length
// - total_call_duration: Cumulative call time
// - is_running: Processor status
// - speech_to_text_enabled: STT service status
// - call_recording_enabled: Recording service status
```

## Integration Examples

### Speech-to-Text Integration
```go
cfg.SpeechToTextEnabled = true
cfg.SpeechToTextAPIKey = "your-api-key"
```

### Call Recording
```go
cfg.CallRecordingEnabled = true
cfg.RecordingPath = "/secure/recordings"
```

### Telephony Provider
```go
cfg.TelephonyProvider = "twilio"
cfg.TelephonyAPIKey = "your-twilio-key"
```

### Webhook Server
```go
cfg.WebhookEnabled = true
cfg.WebhookPort = 8083
cfg.WebhookPath = "/voice/webhook"
```

## Security Considerations

- Call recordings must comply with local recording laws
- API keys should be stored securely and rotated regularly
- Webhook endpoints should use HTTPS and authentication
- Personal information in transcripts should be handled per privacy regulations
- Threat detection results should be logged securely
- Caller data should comply with telecommunications privacy requirements

## Testing

Run the test suite:
```bash
cd fr0g-ai-master-control
go test ./internal/processors/voice/...
```

## Future Enhancements

- Real-time audio stream analysis
- Machine learning-based voice pattern recognition
- Integration with external threat intelligence feeds
- Advanced natural language processing for transcripts
- Real-time threat correlation across multiple call vectors
- Automated call blocking and response capabilities
- Voice biometric analysis for caller identification
- Integration with law enforcement reporting systems

## Supported Telephony Providers

- **Twilio**: Full API integration for call events and recording
- **Asterisk**: PBX integration for enterprise environments
- **FreePBX**: Open-source PBX system integration
- **3CX**: Business phone system integration
- **Generic SIP**: Standard SIP protocol support

## Compliance Features

- **Call Recording Laws**: Configurable compliance with local recording regulations
- **Data Retention**: Automatic cleanup of old recordings and transcripts
- **Privacy Protection**: Redaction of sensitive information in logs
- **Audit Trails**: Comprehensive logging of all threat detection activities
- **Reporting**: Automated threat reports for law enforcement cooperation
