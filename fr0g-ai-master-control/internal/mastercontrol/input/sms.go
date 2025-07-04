package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// SMSProcessor processes incoming SMS/text messages for threat analysis
type SMSProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *SMSConfig
}

// SMSConfig holds SMS processor configuration
type SMSConfig struct {
	Provider          string        `yaml:"provider"`           // "google_voice", "twilio", "webhook"
	APIKey            string        `yaml:"api_key"`
	APISecret         string        `yaml:"api_secret"`
	WebhookURL        string        `yaml:"webhook_url"`
	PhoneNumber       string        `yaml:"phone_number"`
	CommunityTopic    string        `yaml:"community_topic"`
	PersonaCount      int           `yaml:"persona_count"`
	ReviewTimeout     time.Duration `yaml:"review_timeout"`
	RequiredConsensus float64       `yaml:"required_consensus"`
	EnableSpamFilter  bool          `yaml:"enable_spam_filter"`
	BlockedNumbers    []string      `yaml:"blocked_numbers"`
}

// SMSMessage represents an incoming SMS message
type SMSMessage struct {
	ID          string            `json:"id"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	Body        string            `json:"body"`
	Timestamp   time.Time         `json:"timestamp"`
	Provider    string            `json:"provider"`
	MessageType string            `json:"message_type"` // "sms", "mms"
	Media       []MediaAttachment `json:"media,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

// MediaAttachment represents media in MMS messages
type MediaAttachment struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Filename    string `json:"filename"`
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
	return fmt.Sprintf("SMS/Text Message Threat Vector Interceptor via %s - Mobile communication intelligence gathering for AI community review on topic: %s", 
		s.config.Provider, s.config.CommunityTopic)
}

// ProcessWebhook processes an SMS webhook (for webhook-based SMS analysis)
func (s *SMSProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("SMS Processor: Processing SMS threat vector webhook")
	
	// Parse SMS message from request body
	sms, err := s.parseSMSMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMS message: %w", err)
	}
	
	// Check if number is blocked
	if s.isNumberBlocked(sms.From) {
		log.Printf("SMS Processor: Blocked number %s, rejecting message", sms.From)
		return &WebhookResponse{
			Success:   true,
			Message:   "SMS from blocked number rejected",
			RequestID: request.ID,
			Data: map[string]interface{}{
				"action":      "blocked",
				"reason":      "blocked_number",
				"from_number": sms.From,
			},
			Timestamp: time.Now(),
		}, nil
	}
	
	// Analyze SMS for threats using AI community
	threatLevel, consensus, err := s.analyzeSMSThreats(ctx, sms)
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
			"from_number":  sms.From,
			"message_type": sms.MessageType,
		},
		Timestamp: time.Now(),
	}
	
	log.Printf("SMS Processor: SMS analyzed - From: %s, Threat Level: %s, Consensus: %.2f", 
		sms.From, threatLevel, consensus)
	
	return response, nil
}

// analyzeSMSThreats analyzes SMS for threats using AI community
func (s *SMSProcessor) analyzeSMSThreats(ctx context.Context, sms *SMSMessage) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
SMS/Text Message Threat Analysis Request:
From: %s
To: %s
Message: %s
Type: %s
Media Attachments: %d
Timestamp: %s
Provider: %s

Please analyze this SMS message for potential threats including:
- Phishing attempts (fake bank alerts, verification codes)
- Smishing (SMS phishing)
- Malware links
- Social engineering attempts
- Spam/unwanted marketing
- Scam attempts (lottery, romance, tech support)
- Suspicious shortened URLs
- Requests for personal information
`, sms.From, sms.To, sms.Body, sms.MessageType, len(sms.Media), sms.Timestamp.Format(time.RFC3339), sms.Provider)
	
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
		
		if consensus >= 0.9 {
			threatLevel = "critical"
		} else if consensus >= 0.8 {
			threatLevel = "high"
		} else if consensus >= 0.6 {
			threatLevel = "medium"
		} else if consensus >= 0.4 {
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
	
	sms := &SMSMessage{
		From:        getStringFromMap(bodyMap, "from"),
		To:          getStringFromMap(bodyMap, "to"),
		Body:        getStringFromMap(bodyMap, "body"),
		Provider:    getStringFromMap(bodyMap, "provider"),
		MessageType: getStringFromMap(bodyMap, "message_type"),
		Timestamp:   time.Now(),
		Metadata:    make(map[string]string),
		Media:       []MediaAttachment{},
	}
	
	// Parse metadata
	if metadataData, ok := bodyMap["metadata"].(map[string]interface{}); ok {
		for key, value := range metadataData {
			if valueStr, ok := value.(string); ok {
				sms.Metadata[key] = valueStr
			}
		}
	}
	
	// Parse media attachments for MMS
	if mediaData, ok := bodyMap["media"].([]interface{}); ok {
		for _, media := range mediaData {
			if mediaMap, ok := media.(map[string]interface{}); ok {
				attachment := MediaAttachment{
					URL:         getStringFromMap(mediaMap, "url"),
					ContentType: getStringFromMap(mediaMap, "content_type"),
					Filename:    getStringFromMap(mediaMap, "filename"),
				}
				if size, ok := mediaMap["size"].(float64); ok {
					attachment.Size = int64(size)
				}
				sms.Media = append(sms.Media, attachment)
			}
		}
	}
	
	return sms, nil
}

// isNumberBlocked checks if a phone number is blocked
func (s *SMSProcessor) isNumberBlocked(number string) bool {
	for _, blocked := range s.config.BlockedNumbers {
		if strings.Contains(number, blocked) {
			return true
		}
	}
	return false
}
