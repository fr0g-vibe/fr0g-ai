package esmtp

import (
	"fmt"
	"log"
	"net"
	"net/mail"
	"regexp"
	"strings"
	"sync"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/types"
)

// Processor handles ESMTP threat detection and analysis
type Processor struct {
	config         *sharedconfig.ESMTPConfig
	threatPatterns map[string]*regexp.Regexp
	spamKeywords   []string
	emailHistory   []EmailMessage
	senderInfo     map[string]*SenderInfo
	mu             sync.RWMutex
	isRunning      bool
	stopChan       chan struct{}
	smtpServer     *SMTPServer
}

// EmailMessage represents an incoming email message
type EmailMessage struct {
	ID          string                 `json:"id"`
	From        string                 `json:"from"`
	To          []string               `json:"to"`
	Subject     string                 `json:"subject"`
	Body        string                 `json:"body"`
	Headers     map[string]string      `json:"headers"`
	Timestamp   time.Time              `json:"timestamp"`
	ThreatLevel ThreatLevel            `json:"threat_level"`
	Analysis    *ThreatAnalysis        `json:"analysis,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ThreatLevel represents the severity of detected threats
type ThreatLevel int

const (
	ThreatLevelNone ThreatLevel = iota
	ThreatLevelLow
	ThreatLevelMedium
	ThreatLevelHigh
	ThreatLevelCritical
)

func (t ThreatLevel) String() string {
	switch t {
	case ThreatLevelNone:
		return "none"
	case ThreatLevelLow:
		return "low"
	case ThreatLevelMedium:
		return "medium"
	case ThreatLevelHigh:
		return "high"
	case ThreatLevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// ThreatAnalysis contains detailed threat analysis results
type ThreatAnalysis struct {
	ThreatTypes     []string  `json:"threat_types"`
	Confidence      float64   `json:"confidence"`
	SpamScore       float64   `json:"spam_score"`
	PhishingScore   float64   `json:"phishing_score"`
	MalwareScore    float64   `json:"malware_score"`
	SpoofingScore   float64   `json:"spoofing_score"`
	Indicators      []string  `json:"indicators"`
	Recommendations []string  `json:"recommendations"`
	ProcessedAt     time.Time `json:"processed_at"`
}

// SenderInfo tracks information about email senders
type SenderInfo struct {
	Email         string    `json:"email"`
	Domain        string    `json:"domain"`
	FirstSeen     time.Time `json:"first_seen"`
	LastSeen      time.Time `json:"last_seen"`
	MessageCount  int       `json:"message_count"`
	ThreatCount   int       `json:"threat_count"`
	IsBlacklisted bool      `json:"is_blacklisted"`
	IsWhitelisted bool      `json:"is_whitelisted"`
	Reputation    float64   `json:"reputation"` // 0.0-1.0, higher is better
}

// SMTPServer represents the SMTP server component
type SMTPServer struct {
	listener  net.Listener
	config    *sharedconfig.ESMTPConfig
	processor *Processor
}

// NewProcessor creates a new ESMTP processor instance
func NewProcessor(cfg *sharedconfig.ESMTPConfig) *Processor {
	p := &Processor{
		config:         cfg,
		threatPatterns: make(map[string]*regexp.Regexp),
		emailHistory:   make([]EmailMessage, 0),
		senderInfo:     make(map[string]*SenderInfo),
		stopChan:       make(chan struct{}),
	}

	p.initializeThreatPatterns()
	p.initializeSpamKeywords()

	return p
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "esmtp"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.config.Enabled
}



// Process processes an input event (converts to EmailMessage and back)
func (p *Processor) Process(event *types.InputEvent) (*types.InputEventResponse, error) {
	// Convert InputEvent to EmailMessage
	emailMsg, err := p.convertToEmailMessage(event)
	if err != nil {
		return &types.InputEventResponse{
			EventID:     event.ID,
			Processed:   false,
			Actions:     []types.OutputAction{},
			Metadata:    map[string]interface{}{"error": err.Error()},
			ProcessedAt: time.Now(),
		}, nil
	}

	// Process the email message
	processedMsg, err := p.ProcessMessage(*emailMsg)
	if err != nil {
		return &types.InputEventResponse{
			EventID:     event.ID,
			Processed:   false,
			Actions:     []types.OutputAction{},
			Metadata:    map[string]interface{}{"error": err.Error()},
			ProcessedAt: time.Now(),
		}, nil
	}

	// Convert back to InputEventResponse
	response := &types.InputEventResponse{
		EventID:     event.ID,
		Processed:   true,
		Actions:     p.generateActions(processedMsg),
		Metadata:    p.convertMetadata(processedMsg),
		ProcessedAt: time.Now(),
	}

	return response, nil
}

// initializeThreatPatterns sets up regex patterns for threat detection
func (p *Processor) initializeThreatPatterns() {
	patterns := map[string]string{
		"phishing_url":         `(?i)(bit\.ly|tinyurl|t\.co|goo\.gl|short\.link)/[a-zA-Z0-9]+`,
		"suspicious_attachment": `\.(exe|scr|bat|com|pif|vbs|js|jar|zip|rar)$`,
		"urgent_action":        `(?i)(urgent|immediate|act now|limited time|expires|deadline)`,
		"financial_scam":       `(?i)(wire transfer|bitcoin|cryptocurrency|investment|lottery|inheritance)`,
		"credential_theft":     `(?i)(verify|confirm|update|suspended|locked|security alert)`,
		"spoofed_sender":       `(?i)(noreply|no-reply|donotreply|admin|support|security)@`,
		"malware_keywords":     `(?i)(download|install|update|click here|verify account)`,
		"social_eng":           `(?i)(ceo|manager|boss|urgent request|confidential)`,
	}

	for name, pattern := range patterns {
		if compiled, err := regexp.Compile(pattern); err == nil {
			p.threatPatterns[name] = compiled
		} else {
			log.Printf("Failed to compile threat pattern %s: %v", name, err)
		}
	}
}

// initializeSpamKeywords sets up common spam keywords
func (p *Processor) initializeSpamKeywords() {
	p.spamKeywords = []string{
		"free", "winner", "congratulations", "prize", "lottery",
		"urgent", "immediate", "act now", "limited time", "expires",
		"click here", "call now", "reply now", "unsubscribe",
		"viagra", "cialis", "weight loss", "debt relief",
		"work from home", "make money", "earn cash", "get rich",
		"no obligation", "risk free", "guaranteed", "100% free",
		"nigerian prince", "inheritance", "beneficiary",
	}
}

// ProcessMessage analyzes an email message for threats
func (p *Processor) ProcessMessage(msg EmailMessage) (*EmailMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Parse sender information
	senderEmail := p.extractSenderEmail(msg.From)
	p.updateSenderInfo(senderEmail)

	// Perform threat analysis
	analysis := p.analyzeThreat(msg)
	msg.Analysis = analysis
	msg.ThreatLevel = p.calculateThreatLevel(analysis)

	// Store message in history (use default limit)
	maxHistory := 1000
	
	p.emailHistory = append(p.emailHistory, msg)
	if len(p.emailHistory) > maxHistory {
		p.emailHistory = p.emailHistory[1:]
	}

	log.Printf("Processed email from %s: threat_level=%s, confidence=%.2f",
		senderEmail, msg.ThreatLevel.String(), analysis.Confidence)

	return &msg, nil
}

// analyzeThreat performs comprehensive threat analysis on email content
func (p *Processor) analyzeThreat(msg EmailMessage) *ThreatAnalysis {
	analysis := &ThreatAnalysis{
		ThreatTypes:     make([]string, 0),
		Indicators:      make([]string, 0),
		Recommendations: make([]string, 0),
		ProcessedAt:     time.Now(),
	}

	// Analyze subject and body
	content := strings.ToLower(msg.Subject + " " + msg.Body)

	// Check threat patterns
	for patternName, pattern := range p.threatPatterns {
		if pattern.MatchString(msg.Subject + " " + msg.Body) {
			analysis.ThreatTypes = append(analysis.ThreatTypes, patternName)
			analysis.Indicators = append(analysis.Indicators, fmt.Sprintf("Pattern match: %s", patternName))
		}
	}

	// Calculate various threat scores
	analysis.SpamScore = p.calculateSpamScore(content)
	analysis.PhishingScore = p.calculatePhishingScore(msg)
	analysis.MalwareScore = p.calculateMalwareScore(content)
	analysis.SpoofingScore = p.calculateSpoofingScore(msg)

	// Calculate overall confidence
	analysis.Confidence = p.calculateOverallConfidence(analysis)

	// Generate recommendations
	analysis.Recommendations = p.generateRecommendations(analysis)

	return analysis
}

// calculateSpamScore calculates spam likelihood
func (p *Processor) calculateSpamScore(content string) float64 {
	score := 0.0
	keywordCount := 0

	for _, keyword := range p.spamKeywords {
		if strings.Contains(content, keyword) {
			keywordCount++
			score += 0.2 // Increased sensitivity
		}
	}

	// Bonus for multiple keywords
	if keywordCount > 2 { // Reduced threshold
		score += float64(keywordCount-2) * 0.1 // Increased bonus
	}

	// Check for URL shorteners (high phishing indicator)
	if strings.Contains(content, "bit.ly") || strings.Contains(content, "tinyurl") {
		score += 0.6
	}

	// Check for account suspension language
	if strings.Contains(content, "suspended") || strings.Contains(content, "verify") {
		score += 0.4
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculatePhishingScore calculates phishing likelihood
func (p *Processor) calculatePhishingScore(msg EmailMessage) float64 {
	score := 0.0
	content := strings.ToLower(msg.Subject + " " + msg.Body)

	// Check for suspicious URLs
	if p.threatPatterns["phishing_url"].MatchString(msg.Body) {
		score += 0.4
	}

	// Check for credential theft indicators
	if strings.Contains(content, "verify") || strings.Contains(content, "confirm") {
		score += 0.3
	}

	// Check for urgency indicators
	if strings.Contains(content, "urgent") || strings.Contains(content, "immediate") {
		score += 0.2
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateMalwareScore calculates malware likelihood
func (p *Processor) calculateMalwareScore(content string) float64 {
	score := 0.0

	// Check for download/install requests
	if strings.Contains(content, "download") || strings.Contains(content, "install") {
		score += 0.3
	}

	// Check for attachment indicators
	if strings.Contains(content, "attachment") || strings.Contains(content, "file") {
		score += 0.2
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateSpoofingScore calculates spoofing likelihood
func (p *Processor) calculateSpoofingScore(msg EmailMessage) float64 {
	score := 0.0

	// Check sender reputation
	senderEmail := p.extractSenderEmail(msg.From)
	if info, exists := p.senderInfo[senderEmail]; exists {
		if info.Reputation < 0.3 {
			score += 0.4
		}
	}

	// Check for spoofed sender patterns
	if p.threatPatterns["spoofed_sender"].MatchString(msg.From) {
		score += 0.3
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateOverallConfidence calculates overall threat confidence
func (p *Processor) calculateOverallConfidence(analysis *ThreatAnalysis) float64 {
	scores := []float64{
		analysis.SpamScore,
		analysis.PhishingScore,
		analysis.MalwareScore,
		analysis.SpoofingScore,
	}

	total := 0.0
	for _, score := range scores {
		total += score
	}

	confidence := total / float64(len(scores))
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// calculateThreatLevel determines threat level based on analysis
func (p *Processor) calculateThreatLevel(analysis *ThreatAnalysis) ThreatLevel {
	if analysis.Confidence >= 0.7 {
		return ThreatLevelCritical
	} else if analysis.Confidence >= 0.5 {
		return ThreatLevelHigh
	} else if analysis.Confidence >= 0.3 {
		return ThreatLevelMedium
	} else if analysis.Confidence >= 0.15 {
		return ThreatLevelLow
	}
	return ThreatLevelNone
}

// generateRecommendations generates security recommendations
func (p *Processor) generateRecommendations(analysis *ThreatAnalysis) []string {
	recommendations := make([]string, 0)

	if analysis.SpamScore > 0.5 {
		recommendations = append(recommendations, "Mark as spam and block sender")
	}

	if analysis.PhishingScore > 0.5 {
		recommendations = append(recommendations, "Do not click any links or provide credentials")
	}

	if analysis.MalwareScore > 0.5 {
		recommendations = append(recommendations, "Do not open attachments or download files")
	}

	if analysis.SpoofingScore > 0.5 {
		recommendations = append(recommendations, "Verify sender identity through alternative channels")
	}

	if analysis.Confidence > 0.7 {
		recommendations = append(recommendations, "Quarantine message and report to security team")
	}

	return recommendations
}

// extractSenderEmail extracts email address from sender field
func (p *Processor) extractSenderEmail(from string) string {
	if addr, err := mail.ParseAddress(from); err == nil {
		return addr.Address
	}
	return from
}

// updateSenderInfo updates tracking information for email senders
func (p *Processor) updateSenderInfo(email string) {
	info, exists := p.senderInfo[email]
	if !exists {
		domain := ""
		if parts := strings.Split(email, "@"); len(parts) == 2 {
			domain = parts[1]
		}

		info = &SenderInfo{
			Email:        email,
			Domain:       domain,
			FirstSeen:    time.Now(),
			LastSeen:     time.Now(),
			MessageCount: 0,
			ThreatCount:  0,
			Reputation:   0.5, // Neutral starting reputation
		}
		p.senderInfo[email] = info
	}

	info.LastSeen = time.Now()
	info.MessageCount++
}

// convertToEmailMessage converts InputEvent to EmailMessage
func (p *Processor) convertToEmailMessage(event *types.InputEvent) (*EmailMessage, error) {
	emailMsg := &EmailMessage{
		ID:        event.ID,
		Timestamp: event.Timestamp,
		Metadata:  event.Metadata,
	}

	// Extract email-specific fields from metadata
	if from, ok := event.Metadata["from"].(string); ok {
		emailMsg.From = from
	} else {
		emailMsg.From = event.Source
	}

	if to, ok := event.Metadata["to"].([]string); ok {
		emailMsg.To = to
	} else if to, ok := event.Metadata["to"].(string); ok {
		emailMsg.To = []string{to}
	}

	if subject, ok := event.Metadata["subject"].(string); ok {
		emailMsg.Subject = subject
	}

	if headers, ok := event.Metadata["headers"].(map[string]string); ok {
		emailMsg.Headers = headers
	}

	emailMsg.Body = event.Content

	return emailMsg, nil
}

// generateActions generates output actions based on threat analysis
func (p *Processor) generateActions(msg *EmailMessage) []types.OutputAction {
	actions := []types.OutputAction{}

	if msg.Analysis == nil {
		return actions
	}

	// Generate alert action for high-threat emails
	if msg.ThreatLevel >= ThreatLevelHigh {
		action := types.OutputAction{
			ID:       fmt.Sprintf("alert-%s", msg.ID),
			Type:     "alert",
			Target:   "security-team@company.com", // This should be configurable
			Content:  p.formatThreatAlert(msg),
			Priority: int(msg.ThreatLevel),
			Metadata: map[string]interface{}{
				"threat_level": msg.ThreatLevel.String(),
				"confidence":   msg.Analysis.Confidence,
				"source_email": msg.From,
			},
			CreatedAt: time.Now(),
		}
		actions = append(actions, action)
	}

	// Generate quarantine action for critical threats
	if msg.ThreatLevel == ThreatLevelCritical {
		action := types.OutputAction{
			ID:       fmt.Sprintf("quarantine-%s", msg.ID),
			Type:     "quarantine",
			Target:   "email-system",
			Content:  fmt.Sprintf("Quarantine email %s from %s", msg.ID, msg.From),
			Priority: 10,
			Metadata: map[string]interface{}{
				"action":       "quarantine",
				"email_id":     msg.ID,
				"sender":       msg.From,
				"threat_types": msg.Analysis.ThreatTypes,
			},
			CreatedAt: time.Now(),
		}
		actions = append(actions, action)
	}

	return actions
}

// formatThreatAlert formats a threat alert message
func (p *Processor) formatThreatAlert(msg *EmailMessage) string {
	alert := fmt.Sprintf("THREAT DETECTED: Email from %s\n", msg.From)
	alert += fmt.Sprintf("Subject: %s\n", msg.Subject)
	alert += fmt.Sprintf("Threat Level: %s\n", msg.ThreatLevel.String())
	alert += fmt.Sprintf("Confidence: %.2f\n", msg.Analysis.Confidence)
	alert += fmt.Sprintf("Threat Types: %v\n", msg.Analysis.ThreatTypes)
	alert += fmt.Sprintf("Recommendations: %v\n", msg.Analysis.Recommendations)
	return alert
}

// convertMetadata converts EmailMessage analysis to metadata
func (p *Processor) convertMetadata(msg *EmailMessage) map[string]interface{} {
	metadata := make(map[string]interface{})

	if msg.Analysis != nil {
		metadata["threat_analysis"] = msg.Analysis
		metadata["threat_level"] = msg.ThreatLevel.String()
		metadata["confidence"] = msg.Analysis.Confidence
	}

	metadata["sender"] = msg.From
	metadata["subject"] = msg.Subject
	metadata["processed_at"] = time.Now()

	return metadata
}

// GetStats returns processor statistics
func (p *Processor) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	threatCounts := make(map[string]int)
	for _, msg := range p.emailHistory {
		threatCounts[msg.ThreatLevel.String()]++
	}

	return map[string]interface{}{
		"total_messages":      len(p.emailHistory),
		"unique_senders":      len(p.senderInfo),
		"threat_counts":       threatCounts,
		"is_running":          p.isRunning,
		"smtp_server_enabled": p.config.Enabled,
	}
}
