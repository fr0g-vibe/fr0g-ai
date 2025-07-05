package discord

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Processor handles Discord threat detection and analysis
type Processor struct {
	config          *sharedconfig.DiscordConfig
	threatPatterns  map[string]*regexp.Regexp
	suspiciousWords []string
	userHistory     map[string]*UserInfo
	messageHistory  []DiscordMessage
	mu              sync.RWMutex
	isRunning       bool
	stopChan        chan struct{}
}

// DiscordMessage represents a Discord message
type DiscordMessage struct {
	ID          string                 `json:"id"`
	GuildID     string                 `json:"guild_id"`
	ChannelID   string                 `json:"channel_id"`
	UserID      string                 `json:"user_id"`
	Username    string                 `json:"username"`
	Content     string                 `json:"content"`
	MessageType string                 `json:"message_type"` // "message", "embed", "attachment"
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
	BotScore        float64   `json:"bot_score"`
	ScamScore       float64   `json:"scam_score"`
	SocialEngScore  float64   `json:"social_eng_score"`
	Indicators      []string  `json:"indicators"`
	Recommendations []string  `json:"recommendations"`
	ProcessedAt     time.Time `json:"processed_at"`
}

// UserInfo tracks information about Discord users
type UserInfo struct {
	UserID         string    `json:"user_id"`
	Username       string    `json:"username"`
	FirstSeen      time.Time `json:"first_seen"`
	LastSeen       time.Time `json:"last_seen"`
	MessageCount   int       `json:"message_count"`
	ThreatCount    int       `json:"threat_count"`
	IsBot          bool      `json:"is_bot"`
	IsBlacklisted  bool      `json:"is_blacklisted"`
	IsWhitelisted  bool      `json:"is_whitelisted"`
	Reputation     float64   `json:"reputation"` // 0.0-1.0, higher is better
	RecentMessages []string  `json:"recent_messages"`
}

// NewProcessor creates a new Discord processor instance
func NewProcessor(cfg *sharedconfig.DiscordConfig) *Processor {
	p := &Processor{
		config:         cfg,
		threatPatterns: make(map[string]*regexp.Regexp),
		userHistory:    make(map[string]*UserInfo),
		messageHistory: make([]DiscordMessage, 0),
		stopChan:       make(chan struct{}),
	}

	p.initializeThreatPatterns()
	p.initializeSuspiciousWords()

	return p
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "discord"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.config.Enabled
}

// initializeThreatPatterns sets up regex patterns for threat detection
func (p *Processor) initializeThreatPatterns() {
	patterns := map[string]string{
		"malware_url":      `(?i)(bit\.ly|tinyurl|t\.co|goo\.gl|short\.link)/[a-zA-Z0-9]+`,
		"phishing_url":     `(?i)(discord\.gg|discordapp\.com).*[^a-zA-Z0-9\-_]`,
		"crypto_scam":      `(?i)(bitcoin|btc|ethereum|eth|crypto|nft|airdrop|giveaway)`,
		"nitro_scam":       `(?i)(free nitro|discord nitro|nitro gift|steam gift)`,
		"bot_pattern":      `(?i)^(bot|service|auto)_?[a-z0-9]*$`,
		"spam_pattern":     `(?i)(free|win|prize|money|cash|earn|work from home)`,
		"malware_keywords": `(?i)(download|install|exe|zip|rar|torrent)`,
		"social_eng":       `(?i)(urgent|immediate|click here|verify now|suspended)`,
		"invite_spam":      `(?i)discord\.gg/[a-zA-Z0-9]+`,
		"dm_scam":          `(?i)(dm me|private message|check dms)`,
	}

	for name, pattern := range patterns {
		if compiled, err := regexp.Compile(pattern); err == nil {
			p.threatPatterns[name] = compiled
		} else {
			log.Printf("Failed to compile Discord threat pattern %s: %v", name, err)
		}
	}
}

// initializeSuspiciousWords sets up suspicious word detection
func (p *Processor) initializeSuspiciousWords() {
	p.suspiciousWords = []string{
		"hack", "crack", "warez", "keygen", "serial",
		"ddos", "dos", "raid", "spam", "bot",
		"trojan", "virus", "malware", "exploit",
		"phish", "scam", "fraud", "steal", "password",
		"nitro", "giveaway", "airdrop", "crypto", "nft",
		"bitcoin", "ethereum", "wallet", "seed phrase",
	}
}

// Start begins Discord message processing
func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.isRunning {
		p.mu.Unlock()
		return fmt.Errorf("Discord processor is already running")
	}
	p.isRunning = true
	p.mu.Unlock()

	log.Printf("Starting Discord processor with config: %+v", p.config)

	// Start bot connection if bot token is provided
	if p.config.BotToken != "" {
		go p.startBotConnection(ctx)
	}

	// Start webhook server if enabled
	if p.config.WebhookEnabled {
		go p.startWebhookServer(ctx)
	}

	// Start message processing goroutine
	go p.processMessages(ctx)

	log.Println("Discord processor started successfully")
	return nil
}

// Stop stops Discord message processing
func (p *Processor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return fmt.Errorf("Discord processor is not running")
	}

	close(p.stopChan)
	p.isRunning = false

	log.Println("Discord processor stopped")
	return nil
}

// ProcessMessage analyzes a Discord message for threats
func (p *Processor) ProcessMessage(msg DiscordMessage) (*DiscordMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Update user tracking
	p.updateUserInfo(msg.UserID, msg.Username)

	// Perform threat analysis
	analysis := p.analyzeThreat(msg)
	msg.Analysis = analysis
	msg.ThreatLevel = p.calculateThreatLevel(analysis)

	// Store message in history
	if p.config.MaxHistorySize > 0 {
		p.messageHistory = append(p.messageHistory, msg)
		if len(p.messageHistory) > p.config.MaxHistorySize {
			p.messageHistory = p.messageHistory[1:]
		}
	}

	return &msg, nil
}

// analyzeThreat performs comprehensive threat analysis on Discord message
func (p *Processor) analyzeThreat(msg DiscordMessage) *ThreatAnalysis {
	analysis := &ThreatAnalysis{
		ThreatTypes:     make([]string, 0),
		Indicators:      make([]string, 0),
		Recommendations: make([]string, 0),
		ProcessedAt:     time.Now(),
	}

	content := strings.ToLower(msg.Content)

	// Check threat patterns
	for patternName, pattern := range p.threatPatterns {
		if pattern.MatchString(msg.Content) {
			analysis.ThreatTypes = append(analysis.ThreatTypes, patternName)
			analysis.Indicators = append(analysis.Indicators, fmt.Sprintf("Pattern match: %s", patternName))
		}
	}

	// Calculate various threat scores
	analysis.SpamScore = p.calculateSpamScore(content)
	analysis.PhishingScore = p.calculatePhishingScore(content)
	analysis.MalwareScore = p.calculateMalwareScore(content)
	analysis.BotScore = p.calculateBotScore(msg)
	analysis.ScamScore = p.calculateScamScore(content)
	analysis.SocialEngScore = p.calculateSocialEngScore(content)

	// Calculate overall confidence
	analysis.Confidence = p.calculateOverallConfidence(analysis)

	// Generate recommendations
	analysis.Recommendations = p.generateRecommendations(analysis)

	return analysis
}

// calculateSpamScore calculates spam likelihood
func (p *Processor) calculateSpamScore(content string) float64 {
	score := 0.0

	for _, word := range p.suspiciousWords {
		if strings.Contains(content, word) {
			score += 0.15
		}
	}

	// Check for excessive mentions
	if strings.Count(content, "@") > 3 {
		score += 0.3
	}

	// Check for excessive emojis
	emojiCount := strings.Count(content, ":")
	if emojiCount > 10 {
		score += 0.2
	}

	// Check for excessive caps
	capsCount := 0
	for _, r := range content {
		if r >= 'A' && r <= 'Z' {
			capsCount++
		}
	}
	if len(content) > 0 && float64(capsCount)/float64(len(content)) > 0.5 {
		score += 0.3
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculatePhishingScore calculates phishing likelihood
func (p *Processor) calculatePhishingScore(content string) float64 {
	score := 0.0

	phishingKeywords := []string{
		"verify", "account", "password", "suspended",
		"click here", "urgent", "immediate", "security alert",
		"login", "authenticate", "confirm identity",
	}

	for _, keyword := range phishingKeywords {
		if strings.Contains(content, keyword) {
			score += 0.2
		}
	}

	// Check for suspicious Discord invites
	if strings.Contains(content, "discord.gg") && strings.Contains(content, "free") {
		score += 0.4
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateMalwareScore calculates malware likelihood
func (p *Processor) calculateMalwareScore(content string) float64 {
	score := 0.0

	malwareKeywords := []string{
		"download", "install", "exe", "zip", "rar",
		"crack", "keygen", "warez", "torrent", "hack tool",
	}

	for _, keyword := range malwareKeywords {
		if strings.Contains(content, keyword) {
			score += 0.25
		}
	}

	// Check for file attachments in suspicious contexts
	if strings.Contains(content, "attachment") && strings.Contains(content, "free") {
		score += 0.3
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateBotScore calculates bot likelihood
func (p *Processor) calculateBotScore(msg DiscordMessage) float64 {
	score := 0.0

	// Check username pattern
	if p.threatPatterns["bot_pattern"].MatchString(msg.Username) {
		score += 0.4
	}

	// Check for automated patterns
	if strings.Contains(msg.Content, "[") && strings.Contains(msg.Content, "]") {
		score += 0.2
	}

	// Check user info
	userInfo := p.userHistory[msg.UserID]
	if userInfo != nil {
		// High message frequency indicates bot
		if userInfo.MessageCount > 50 && time.Since(userInfo.FirstSeen) < time.Hour {
			score += 0.3
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateScamScore calculates scam likelihood
func (p *Processor) calculateScamScore(content string) float64 {
	score := 0.0

	scamKeywords := []string{
		"free nitro", "giveaway", "airdrop", "crypto", "nft",
		"bitcoin", "ethereum", "wallet", "seed phrase",
		"investment", "trading", "profit", "guaranteed",
	}

	for _, keyword := range scamKeywords {
		if strings.Contains(content, keyword) {
			score += 0.2
		}
	}

	// Check for common Discord scam patterns
	if strings.Contains(content, "nitro") && strings.Contains(content, "free") {
		score += 0.5
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateSocialEngScore calculates social engineering likelihood
func (p *Processor) calculateSocialEngScore(content string) float64 {
	score := 0.0

	socialEngKeywords := []string{
		"urgent", "immediate", "act now", "limited time",
		"verify", "confirm", "suspended", "locked",
		"dm me", "private message", "check dms",
	}

	for _, keyword := range socialEngKeywords {
		if strings.Contains(content, keyword) {
			score += 0.15
		}
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
		analysis.BotScore,
		analysis.ScamScore,
		analysis.SocialEngScore,
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
	if analysis.Confidence >= 0.8 {
		return ThreatLevelCritical
	} else if analysis.Confidence >= 0.6 {
		return ThreatLevelHigh
	} else if analysis.Confidence >= 0.4 {
		return ThreatLevelMedium
	} else if analysis.Confidence >= 0.2 {
		return ThreatLevelLow
	}
	return ThreatLevelNone
}

// generateRecommendations generates security recommendations
func (p *Processor) generateRecommendations(analysis *ThreatAnalysis) []string {
	recommendations := make([]string, 0)

	if analysis.SpamScore > 0.5 {
		recommendations = append(recommendations, "Consider muting or banning user")
	}

	if analysis.PhishingScore > 0.5 {
		recommendations = append(recommendations, "Do not click any links or provide credentials")
	}

	if analysis.MalwareScore > 0.5 {
		recommendations = append(recommendations, "Do not download or execute any files")
	}

	if analysis.BotScore > 0.5 {
		recommendations = append(recommendations, "User may be automated - verify human interaction")
	}

	if analysis.ScamScore > 0.5 {
		recommendations = append(recommendations, "Potential scam detected - warn other users")
	}

	if analysis.Confidence > 0.7 {
		recommendations = append(recommendations, "High threat detected - consider reporting to Discord Trust & Safety")
	}

	return recommendations
}

// updateUserInfo updates tracking information for users
func (p *Processor) updateUserInfo(userID, username string) {
	info, exists := p.userHistory[userID]
	if !exists {
		info = &UserInfo{
			UserID:         userID,
			Username:       username,
			FirstSeen:      time.Now(),
			LastSeen:       time.Now(),
			MessageCount:   0,
			ThreatCount:    0,
			Reputation:     0.5, // Neutral starting reputation
			RecentMessages: make([]string, 0),
		}
		p.userHistory[userID] = info
	}

	info.LastSeen = time.Now()
	info.MessageCount++
	info.Username = username // Update username in case it changed
}

// startBotConnection starts Discord bot connection
func (p *Processor) startBotConnection(ctx context.Context) {
	log.Println("Starting Discord bot connection...")

	// TODO: Implement Discord bot API integration
	// This would involve:
	// 1. Authentication with Discord API using bot token
	// 2. Connecting to Discord Gateway
	// 3. Listening for message events
	// 4. Processing messages through threat analysis

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-time.After(30 * time.Second):
			// Placeholder for bot connection maintenance
			log.Println("Discord bot connection maintenance...")
		}
	}
}

// startWebhookServer starts webhook server for Discord events
func (p *Processor) startWebhookServer(ctx context.Context) {
	log.Printf("Starting Discord webhook server on port %d...", p.config.WebhookPort)

	// TODO: Implement webhook server
	// This would involve:
	// 1. HTTP server setup for Discord webhooks
	// 2. Webhook signature verification
	// 3. Event parsing and processing
	// 4. Integration with threat analysis pipeline

	select {
	case <-ctx.Done():
		return
	case <-p.stopChan:
		return
	}
}

// processMessages handles continuous message processing
func (p *Processor) processMessages(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			// Periodic maintenance tasks
			p.performMaintenance()
		}
	}
}

// performMaintenance performs periodic maintenance tasks
func (p *Processor) performMaintenance() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Clean up old user history
	cutoff := time.Now().Add(-24 * time.Hour)
	for userID, user := range p.userHistory {
		if user.LastSeen.Before(cutoff) {
			delete(p.userHistory, userID)
		}
	}
}

// GetStats returns processor statistics
func (p *Processor) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	threatCounts := make(map[string]int)
	for _, msg := range p.messageHistory {
		threatCounts[msg.ThreatLevel.String()]++
	}

	return map[string]interface{}{
		"total_messages":     len(p.messageHistory),
		"unique_users":       len(p.userHistory),
		"threat_counts":      threatCounts,
		"is_running":         p.isRunning,
		"bot_enabled":        p.config.BotToken != "",
		"webhook_enabled":    p.config.WebhookEnabled,
		"monitored_guilds":   len(p.config.GuildIDs),
		"monitored_channels": len(p.config.ChannelIDs),
	}
}
