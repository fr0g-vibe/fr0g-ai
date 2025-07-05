package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// ESMTPProcessor processes ESMTP email messages for threat analysis
type ESMTPProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *ESMTPConfig
}

// ESMTPConfig holds ESMTP processor configuration
type ESMTPConfig struct {
	Server            string        `yaml:"server"`
	Port              int           `yaml:"port"`
	UseSSL            bool          `yaml:"use_ssl"`
	Username          string        `yaml:"username"`
	Password          string        `yaml:"password"`
	Domain            string        `yaml:"domain"`
	CommunityTopic    string        `yaml:"community_topic"`
	PersonaCount      int           `yaml:"persona_count"`
	ReviewTimeout     time.Duration `yaml:"review_timeout"`
	RequiredConsensus float64       `yaml:"required_consensus"`
	TrustedDomains    []string      `yaml:"trusted_domains"`
	BlockedDomains    []string      `yaml:"blocked_domains"`
	MaxAttachmentSize int64         `yaml:"max_attachment_size"`
	ScanAttachments   bool          `yaml:"scan_attachments"`
	QuarantinePath    string        `yaml:"quarantine_path"`
}

// ESMTPMessage represents an ESMTP email message
type ESMTPMessage struct {
	ID          string            `json:"id"`
	From        string            `json:"from"`
	To          []string          `json:"to"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	HTMLBody    string            `json:"html_body,omitempty"`
	Headers     map[string]string `json:"headers"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	MessageID   string            `json:"message_id"`
	InReplyTo   string            `json:"in_reply_to,omitempty"`
	References  []string          `json:"references,omitempty"`
	Priority    string            `json:"priority,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Hash        string `json:"hash,omitempty"`
	Quarantined bool   `json:"quarantined"`
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
	return fmt.Sprintf("ESMTP Email Threat Vector Interceptor on %s:%d - Email intelligence gathering for AI community review on topic: %s",
		e.config.Server, e.config.Port, e.config.CommunityTopic)
}

// ProcessWebhook processes an ESMTP message webhook
func (e *ESMTPProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("ESMTP Processor: Processing email threat vector webhook")

	// Parse email message from request body
	emailMsg, err := e.parseESMTPMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ESMTP message: %w", err)
	}

	// Check if domain is blocked
	if e.isDomainBlocked(emailMsg.From) {
		log.Printf("ESMTP Processor: Blocked domain %s, quarantining email", emailMsg.From)
		return &WebhookResponse{
			Success:   true,
			Message:   "Email from blocked domain quarantined",
			RequestID: request.ID,
			Data: map[string]interface{}{
				"action":  "quarantined",
				"reason":  "blocked_domain",
				"from":    emailMsg.From,
				"subject": emailMsg.Subject,
			},
			Timestamp: time.Now(),
		}, nil
	}

	// Check if domain is trusted (lower threat threshold)
	isTrusted := e.isDomainTrusted(emailMsg.From)

	// Analyze email for threats using AI community
	threatLevel, consensus, err := e.analyzeEmailThreats(ctx, emailMsg, isTrusted)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze email threats: %w", err)
	}

	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "Email threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"threat_level":     threatLevel,
			"consensus":        consensus,
			"review_id":        fmt.Sprintf("email_review_%d", time.Now().UnixNano()),
			"from":             emailMsg.From,
			"subject":          emailMsg.Subject,
			"attachment_count": len(emailMsg.Attachments),
			"is_trusted":       isTrusted,
		},
		Timestamp: time.Now(),
	}

	log.Printf("ESMTP Processor: Email analyzed - From: %s, Subject: %s, Threat Level: %s, Consensus: %.2f",
		emailMsg.From, emailMsg.Subject, threatLevel, consensus)

	return response, nil
}

// analyzeEmailThreats analyzes email for threats using AI community
func (e *ESMTPProcessor) analyzeEmailThreats(ctx context.Context, emailMsg *ESMTPMessage, isTrusted bool) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
Email Threat Analysis Request:
From: %s
To: %s
Subject: %s
Message ID: %s
Is Trusted Domain: %v
Timestamp: %s
Attachment Count: %d

Email Body:
%s

`, emailMsg.From, strings.Join(emailMsg.To, ", "), emailMsg.Subject,
		emailMsg.MessageID, isTrusted, emailMsg.Timestamp.Format(time.RFC3339),
		len(emailMsg.Attachments), emailMsg.Body)

	// Add attachment information
	if len(emailMsg.Attachments) > 0 {
		analysisContent += "\nAttachments:\n"
		for _, attachment := range emailMsg.Attachments {
			analysisContent += fmt.Sprintf("- %s (%s, %d bytes)\n",
				attachment.Filename, attachment.ContentType, attachment.Size)
		}
	}

	analysisContent += `
Please analyze this email for potential threats including:
- Phishing attempts and credential harvesting
- Malware attachments and malicious links
- Business Email Compromise (BEC) attacks
- Social engineering and pretexting
- Spam and unwanted commercial email
- Email spoofing and domain impersonation
- Ransomware delivery mechanisms
- Data exfiltration attempts
- CEO fraud and invoice scams
- Cryptocurrency and financial scams
- Identity theft and personal information harvesting
- Suspicious file attachments and executables
`

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

		// Adjust thresholds for trusted domains
		thresholds := map[string]float64{
			"critical": 0.9,
			"high":     0.8,
			"medium":   0.6,
			"low":      0.4,
		}

		if isTrusted {
			// Higher thresholds for trusted domains
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

// parseESMTPMessage parses an email message from the request body
func (e *ESMTPProcessor) parseESMTPMessage(body interface{}) (*ESMTPMessage, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}

	emailMsg := &ESMTPMessage{
		ID:        getStringFromMap(bodyMap, "id"),
		From:      getStringFromMap(bodyMap, "from"),
		Subject:   getStringFromMap(bodyMap, "subject"),
		Body:      getStringFromMap(bodyMap, "body"),
		HTMLBody:  getStringFromMap(bodyMap, "html_body"),
		MessageID: getStringFromMap(bodyMap, "message_id"),
		InReplyTo: getStringFromMap(bodyMap, "in_reply_to"),
		Priority:  getStringFromMap(bodyMap, "priority"),
		Timestamp: time.Now(),
		Headers:   make(map[string]string),
		Metadata:  make(map[string]string),
	}

	// Parse To, CC, BCC arrays
	if toData, ok := bodyMap["to"].([]interface{}); ok {
		for _, to := range toData {
			if toStr, ok := to.(string); ok {
				emailMsg.To = append(emailMsg.To, toStr)
			}
		}
	}

	if ccData, ok := bodyMap["cc"].([]interface{}); ok {
		for _, cc := range ccData {
			if ccStr, ok := cc.(string); ok {
				emailMsg.CC = append(emailMsg.CC, ccStr)
			}
		}
	}

	if bccData, ok := bodyMap["bcc"].([]interface{}); ok {
		for _, bcc := range bccData {
			if bccStr, ok := bcc.(string); ok {
				emailMsg.BCC = append(emailMsg.BCC, bccStr)
			}
		}
	}

	// Parse references array
	if referencesData, ok := bodyMap["references"].([]interface{}); ok {
		for _, ref := range referencesData {
			if refStr, ok := ref.(string); ok {
				emailMsg.References = append(emailMsg.References, refStr)
			}
		}
	}

	// Parse headers
	if headersData, ok := bodyMap["headers"].(map[string]interface{}); ok {
		for key, value := range headersData {
			if valueStr, ok := value.(string); ok {
				emailMsg.Headers[key] = valueStr
			}
		}
	}

	// Parse attachments
	if attachmentsData, ok := bodyMap["attachments"].([]interface{}); ok {
		for _, attachmentData := range attachmentsData {
			if attachmentMap, ok := attachmentData.(map[string]interface{}); ok {
				attachment := EmailAttachment{
					Filename:    getStringFromMap(attachmentMap, "filename"),
					ContentType: getStringFromMap(attachmentMap, "content_type"),
					Hash:        getStringFromMap(attachmentMap, "hash"),
					Quarantined: getBoolFromMap(attachmentMap, "quarantined"),
				}

				if size, ok := attachmentMap["size"].(float64); ok {
					attachment.Size = int64(size)
				}

				emailMsg.Attachments = append(emailMsg.Attachments, attachment)
			}
		}
	}

	return emailMsg, nil
}

// isDomainBlocked checks if a domain is in the blocked list
func (e *ESMTPProcessor) isDomainBlocked(email string) bool {
	domain := e.extractDomain(email)
	for _, blocked := range e.config.BlockedDomains {
		if strings.EqualFold(domain, blocked) {
			return true
		}
	}
	return false
}

// isDomainTrusted checks if a domain is in the trusted list
func (e *ESMTPProcessor) isDomainTrusted(email string) bool {
	domain := e.extractDomain(email)
	for _, trusted := range e.config.TrustedDomains {
		if strings.EqualFold(domain, trusted) {
			return true
		}
	}
	return false
}

// extractDomain extracts domain from email address
func (e *ESMTPProcessor) extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
