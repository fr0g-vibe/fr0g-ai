package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// ESMTPProcessor processes ESMTP emails for threat analysis
type ESMTPProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *ESMTPConfig
}

// ESMTPConfig holds ESMTP processor configuration
type ESMTPConfig struct {
	Host              string        `yaml:"host"`
	Port              int           `yaml:"port"`
	TLSPort           int           `yaml:"tls_port"`
	Hostname          string        `yaml:"hostname"`
	MaxMessageSize    int64         `yaml:"max_message_size"`
	Timeout           time.Duration `yaml:"timeout"`
	EnableTLS         bool          `yaml:"enable_tls"`
	CertFile          string        `yaml:"cert_file"`
	KeyFile           string        `yaml:"key_file"`
	CommunityTopic    string        `yaml:"community_topic"`
	PersonaCount      int           `yaml:"persona_count"`
	ReviewTimeout     time.Duration `yaml:"review_timeout"`
	RequiredConsensus float64       `yaml:"required_consensus"`
}

// EmailMessage represents an email message
type EmailMessage struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	Attachments []string          `json:"attachments"`
	Timestamp   time.Time         `json:"timestamp"`
}

// NewESMTPProcessor creates a new ESMTP processor
func NewESMTPProcessor(config *ESMTPConfig, aiClient AIPersonaCommunityClient) (*ESMTPProcessor, error) {
	return &ESMTPProcessor{
		aiClient: aiClient,
		config:   config,
	}, nil
}

// GetTag returns the processor tag
func (e *ESMTPProcessor) GetTag() string {
	return "esmtp"
}

// GetDescription returns the processor description
func (e *ESMTPProcessor) GetDescription() string {
	return fmt.Sprintf("ESMTP Threat Vector Interceptor on %s:%d - Email intelligence gathering for AI community review on topic: %s", 
		e.config.Host, e.config.Port, e.config.CommunityTopic)
}

// ProcessWebhook processes an ESMTP webhook (for webhook-based email analysis)
func (e *ESMTPProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("ESMTP Processor: Processing email threat vector webhook")
	
	// Parse email message from request body
	email, err := e.parseEmailMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email message: %w", err)
	}
	
	// Analyze email for threats using AI community
	threatLevel, consensus, err := e.analyzeEmailThreats(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze email threats: %w", err)
	}
	
	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "Email threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"threat_level": threatLevel,
			"consensus":    consensus,
			"review_id":    fmt.Sprintf("review_%d", time.Now().UnixNano()),
		},
		Timestamp: time.Now(),
	}
	
	log.Printf("ESMTP Processor: Email analyzed - Threat Level: %s, Consensus: %.2f", threatLevel, consensus)
	
	return response, nil
}

// analyzeEmailThreats analyzes email for threats using AI community
func (e *ESMTPProcessor) analyzeEmailThreats(ctx context.Context, email *EmailMessage) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
Email Threat Analysis Request:
From: %s
To: %s
Subject: %s
Body: %s
Headers: %v
Attachments: %v

Please analyze this email for potential threats including:
- Phishing attempts
- Malware delivery
- Social engineering
- Spam/unwanted content
- Suspicious links or attachments
`, email.From, strings.Join(email.To, ", "), email.Subject, email.Body, email.Headers, email.Attachments)
	
	// Create AI community for threat analysis
	community, err := e.aiClient.CreateCommunity(ctx, e.config.CommunityTopic, e.config.PersonaCount)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to create AI community: %w", err)
	}
	
	// Submit for AI community review
	review, err := e.aiClient.SubmitForReview(ctx, community.ID, analysisContent)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to submit for review: %w", err)
	}
	
	// Determine threat level based on consensus
	threatLevel := "unknown"
	consensus := 0.0
	
	if review.Consensus != nil {
		consensus = review.Consensus.OverallScore
		
		if consensus >= 0.8 {
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

// parseEmailMessage parses an email message from the request body
func (e *ESMTPProcessor) parseEmailMessage(body interface{}) (*EmailMessage, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}
	
	email := &EmailMessage{
		From:        getStringFromMap(bodyMap, "from"),
		Subject:     getStringFromMap(bodyMap, "subject"),
		Body:        getStringFromMap(bodyMap, "body"),
		Timestamp:   time.Now(),
		Headers:     make(map[string]string),
		Attachments: []string{},
	}
	
	// Parse To field
	if toData, ok := bodyMap["to"].([]interface{}); ok {
		for _, to := range toData {
			if toStr, ok := to.(string); ok {
				email.To = append(email.To, toStr)
			}
		}
	}
	
	// Parse Headers
	if headersData, ok := bodyMap["headers"].(map[string]interface{}); ok {
		for key, value := range headersData {
			if valueStr, ok := value.(string); ok {
				email.Headers[key] = valueStr
			}
		}
	}
	
	// Parse Attachments
	if attachmentsData, ok := bodyMap["attachments"].([]interface{}); ok {
		for _, attachment := range attachmentsData {
			if attachmentStr, ok := attachment.(string); ok {
				email.Attachments = append(email.Attachments, attachmentStr)
			}
		}
	}
	
	return email, nil
}
