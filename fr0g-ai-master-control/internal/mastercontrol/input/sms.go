package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// SMSProcessor processes SMS messages for threat analysis
type SMSProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *SMSConfig
}

// SMSConfig holds SMS processor configuration
type SMSConfig struct {
	Provider             string        `yaml:"provider"` // "twilio", "aws_sns", "nexmo", etc.
	AccountSID           string        `yaml:"account_sid"`
	AuthToken            string        `yaml:"auth_token"`
	PhoneNumber          string        `yaml:"phone_number"`
	WebhookURL           string        `yaml:"webhook_url"`
	CommunityTopic       string        `yaml:"community_topic"`
	PersonaCount         int           `yaml:"persona_count"`
	ReviewTimeout        time.Duration `yaml:"review_timeout"`
	RequiredConsensus    float64       `yaml:"required_consensus"`
	TrustedNumbers       []string      `yaml:"trusted_numbers"`
	BlockedNumbers       []string      `yaml:"blocked_numbers"`
	RateLimitWindow      time.Duration `yaml:"rate_limit_window"`
	MaxMessagesPerWindow int           `yaml:"max_messages_per_window"`
	EnableMMS            bool          `yaml:"enable_mms"`
	MaxMediaSize         int64         `yaml:"max_media_size"`
}

// SMSMessage represents an SMS message
type SMSMessage struct {
	ID          string            `json:"id"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	Body        string            `json:"body"`
	MediaURLs   []string          `json:"media_urls,omitempty"`
	MessageSID  string            `json:"message_sid"`
	Status      string            `json:"status"`
	Direction   string            `json:"direction"` // "inbound", "outbound"
	Timestamp   time.Time         `json:"timestamp"`
	Country     string            `json:"country,omitempty"`
	Region      string            `json:"region,omitempty"`
	Carrier     string            `json:"carrier,omitempty"`
	MessageType string            `json:"message_type"` // "sms", "mms"
	Metadata    map[string]string `json:"metadata"`
}

// NewSMSProcessor creates a new SMS processor
func NewSMSProcessor(config *SMSConfig, aiClient AIPersonaCommunityClient) (*SMSProcessor, error) {
	return &SMSProcessor{
		aiClient: aiClient,
		config:   config,
	}, nil
}

// GetTag returns the processor tag
func (s *SMSProcessor) GetTag() string {
	return "sms"
}

// GetDescription returns the processor description
func (s *SMSProcessor) GetDescription() string {
	return fmt.Sprintf("SMS Threat Vector Interceptor via %s - SMS intelligence gathering for AI community review on topic: %s",
		s.config.Provider, s.config.CommunityTopic)
}

// ProcessWebhook processes an SMS message webhook
func (s *SMSProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("SMS Processor: Processing SMS threat vector webhook")

	// Parse SMS message from request body
	smsMsg, err := s.parseSMSMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMS message: %w", err)
	}

	// Check if number is blocked
	if s.isNumberBlocked(smsMsg.From) {
		log.Printf("SMS Processor: Blocked number %s, dropping message", smsMsg.From)
		return &WebhookResponse{
			Success:   true,
			Message:   "SMS from blocked number dropped",
			RequestID: request.ID,
			Data: map[string]interface{}{
				"action": "blocked",
				"reason": "blocked_number",
				"from":   smsMsg.From,
			},
			Timestamp: time.Now(),
		}, nil
	}

	// Check if number is trusted (lower threat threshold)
	isTrusted := s.isNumberTrusted(smsMsg.From)

	// Analyze SMS for threats using AI community
	threatLevel, consensus, err := s.analyzeSMSThreats(ctx, smsMsg, isTrusted)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze SMS threats: %w", err)
	}

	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "SMS threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"threat_level": threatLevel,
			"consensus":    consensus,
			"review_id":    fmt.Sprintf("sms_review_%d", time.Now().UnixNano()),
			"from":         smsMsg.From,
			"message_type": smsMsg.MessageType,
			"has_media":    len(smsMsg.MediaURLs) > 0,
			"is_trusted":   isTrusted,
		},
		Timestamp: time.Now(),
	}

	log.Printf("SMS Processor: Message analyzed - From: %s, Type: %s, Threat Level: %s, Consensus: %.2f",
		smsMsg.From, smsMsg.MessageType, threatLevel, consensus)

	return response, nil
}

// analyzeSMSThreats analyzes SMS for threats using AI community
func (s *SMSProcessor) analyzeSMSThreats(ctx context.Context, smsMsg *SMSMessage, isTrusted bool) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
SMS Threat Analysis Request:
From: %s
To: %s
Message Type: %s
Direction: %s
Status: %s
Is Trusted Number: %v
Country: %s
Carrier: %s
Timestamp: %s
Media URLs Count: %d

Message Body:
%s

`, smsMsg.From, smsMsg.To, smsMsg.MessageType, smsMsg.Direction,
		smsMsg.Status, isTrusted, smsMsg.Country, smsMsg.Carrier,
		smsMsg.Timestamp.Format(time.RFC3339), len(smsMsg.MediaURLs), smsMsg.Body)

	// Add media URL information
	if len(smsMsg.MediaURLs) > 0 {
		analysisContent += "\nMedia URLs:\n"
		for _, url := range smsMsg.MediaURLs {
			analysisContent += fmt.Sprintf("- %s\n", url)
		}
	}

	analysisContent += `
Please analyze this SMS message for potential threats including:
- Phishing attempts and credential harvesting
- Smishing (SMS phishing) attacks
- Malicious links and URL shorteners
- Premium rate SMS scams
- Social engineering and pretexting
- Identity theft attempts
- Two-factor authentication bypass attempts
- Fake delivery notifications
- Cryptocurrency and investment scams
- Romance and dating scams
- Tech support scams
- Banking and financial fraud
- Malware distribution via MMS
- Spam and unwanted commercial messages
`

	// Create AI community for threat analysis
	community, err := s.aiClient.CreateCommunity(ctx, s.config.CommunityTopic, s.config.PersonaCount)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to create AI community: %w", err)
	}

	// Submit for AI community review
	review, err := s.aiClient.SubmitForReview(ctx, community.ID, analysisContent)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to submit for review: %w", err)
	}

	// Determine threat level based on consensus
	threatLevel := "unknown"
	consensus := 0.0

	if review.Consensus != nil {
		consensus = review.Consensus.OverallScore

		// Adjust thresholds for trusted numbers
		thresholds := map[string]float64{
			"critical": 0.9,
			"high":     0.8,
			"medium":   0.6,
			"low":      0.4,
		}

		if isTrusted {
			// Higher thresholds for trusted numbers
			thresholds["critical"] = 0.95
			thresholds["high"] = 0.85
			thresholds["medium"] = 0.7
			thresholds["low"] = 0.5
		}

		if consensus >= thresholds["critical"] {
			threatLevel = "critical"
		} else if consensus >= thresholds["high"] {
			threatLevel = "high"
		} else if consensus >= thresholds["medium"] {
			threatLevel = "medium"
		} else if consensus >= thresholds["low"] {
			threatLevel = "low"
		} else {
			threatLevel = "minimal"
		}
	}

	return threatLevel, consensus, nil
}

// parseSMSMessage parses an SMS message from the request body
func (s *SMSProcessor) parseSMSMessage(body interface{}) (*SMSMessage, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}

	smsMsg := &SMSMessage{
		ID:          getStringFromMap(bodyMap, "id"),
		From:        getStringFromMap(bodyMap, "from"),
		To:          getStringFromMap(bodyMap, "to"),
		Body:        getStringFromMap(bodyMap, "body"),
		MessageSID:  getStringFromMap(bodyMap, "message_sid"),
		Status:      getStringFromMap(bodyMap, "status"),
		Direction:   getStringFromMap(bodyMap, "direction"),
		Country:     getStringFromMap(bodyMap, "country"),
		Region:      getStringFromMap(bodyMap, "region"),
		Carrier:     getStringFromMap(bodyMap, "carrier"),
		MessageType: getStringFromMap(bodyMap, "message_type"),
		Timestamp:   time.Now(),
		Metadata:    make(map[string]string),
	}

	// Parse media URLs
	if mediaData, ok := bodyMap["media_urls"].([]interface{}); ok {
		for _, media := range mediaData {
			if mediaStr, ok := media.(string); ok {
				smsMsg.MediaURLs = append(smsMsg.MediaURLs, mediaStr)
			}
		}
	}

	// Parse metadata
	if metadataData, ok := bodyMap["metadata"].(map[string]interface{}); ok {
		for key, value := range metadataData {
			if valueStr, ok := value.(string); ok {
				smsMsg.Metadata[key] = valueStr
			}
		}
	}

	return smsMsg, nil
}

// isNumberBlocked checks if a phone number is in the blocked list
func (s *SMSProcessor) isNumberBlocked(phoneNumber string) bool {
	normalizedNumber := s.normalizePhoneNumber(phoneNumber)
	for _, blocked := range s.config.BlockedNumbers {
		if s.normalizePhoneNumber(blocked) == normalizedNumber {
			return true
		}
	}
	return false
}

// isNumberTrusted checks if a phone number is in the trusted list
func (s *SMSProcessor) isNumberTrusted(phoneNumber string) bool {
	normalizedNumber := s.normalizePhoneNumber(phoneNumber)
	for _, trusted := range s.config.TrustedNumbers {
		if s.normalizePhoneNumber(trusted) == normalizedNumber {
			return true
		}
	}
	return false
}

// normalizePhoneNumber normalizes phone number format for comparison
func (s *SMSProcessor) normalizePhoneNumber(phoneNumber string) string {
	// Remove all non-digit characters
	normalized := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phoneNumber)

	// Add country code if missing (assuming US +1)
	if len(normalized) == 10 {
		normalized = "1" + normalized
	}

	return normalized
}
