package irc

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Processor handles IRC threat detection and analysis
type Processor struct {
	config          *sharedconfig.IRCConfig
	connections     map[string]*IRCConnection
	threatPatterns  map[string]*regexp.Regexp
	suspiciousWords []string
	userHistory     map[string]*UserInfo
	messageHistory  []IRCMessage
	mu              sync.RWMutex
	isRunning       bool
	stopChan        chan struct{}
}

// IRCConnection represents a connection to an IRC server
type IRCConnection struct {
	Server    string
	Conn      net.Conn
	Connected bool
	LastPing  time.Time
	Channels  []string
	mu        sync.RWMutex
}

// IRCMessage represents an IRC message
type IRCMessage struct {
	ID          string                 `json:"id"`
	Server      string                 `json:"server"`
	Channel     string                 `json:"channel"`
	Nick        string                 `json:"nick"`
	User        string                 `json:"user"`
	Host        string                 `json:"host"`
	Message     string                 `json:"message"`
	MessageType string                 `json:"message_type"` // PRIVMSG, NOTICE, JOIN, PART, etc.
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
	FloodScore      float64   `json:"flood_score"`
	SocialEngScore  float64   `json:"social_eng_score"`
	Indicators      []string  `json:"indicators"`
	Recommendations []string  `json:"recommendations"`
	ProcessedAt     time.Time `json:"processed_at"`
}

// UserInfo tracks information about IRC users
type UserInfo struct {
	Nick           string    `json:"nick"`
	User           string    `json:"user"`
	Host           string    `json:"host"`
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

// NewProcessor creates a new IRC processor instance
func NewProcessor(cfg *sharedconfig.IRCConfig) *Processor {
	p := &Processor{
		config:         cfg,
		connections:    make(map[string]*IRCConnection),
		threatPatterns: make(map[string]*regexp.Regexp),
		userHistory:    make(map[string]*UserInfo),
		messageHistory: make([]IRCMessage, 0),
		stopChan:       make(chan struct{}),
	}

	p.initializeThreatPatterns()
	p.initializeSuspiciousWords()

	return p
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "irc"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.config.Enabled
}

// initializeThreatPatterns sets up regex patterns for threat detection
func (p *Processor) initializeThreatPatterns() {
	patterns := map[string]string{
		"malware_url":      `(?i)(bit\.ly|tinyurl|t\.co|goo\.gl|short\.link)/[a-zA-Z0-9]+`,
		"phishing_url":     `(?i)(login|verify|account|security).*\.(tk|ml|ga|cf|pw)`,
		"bot_pattern":      `(?i)^(bot|service|auto)_?[a-z0-9]*$`,
		"flood_pattern":    `(.)\1{10,}|(.{1,3})\2{5,}`, // Repeated characters or patterns
		"spam_pattern":     `(?i)(free|win|prize|money|cash|earn|work from home)`,
		"malware_keywords": `(?i)(download|install|exe|zip|rar|torrent)`,
		"social_eng":       `(?i)(urgent|immediate|click here|verify now|suspended)`,
		"dcc_exploit":      `(?i)DCC (SEND|CHAT|GET)`,
		"ctcp_flood":       `\x01.*\x01`,
	}

	for name, pattern := range patterns {
		if compiled, err := regexp.Compile(pattern); err == nil {
			p.threatPatterns[name] = compiled
		} else {
			log.Printf("Failed to compile IRC threat pattern %s: %v", name, err)
		}
	}
}

// initializeSuspiciousWords sets up suspicious word detection
func (p *Processor) initializeSuspiciousWords() {
	p.suspiciousWords = []string{
		"hack", "crack", "warez", "keygen", "serial",
		"ddos", "dos", "flood", "spam", "bot",
		"trojan", "virus", "malware", "exploit",
		"phish", "scam", "fraud", "steal", "password",
		"credit card", "ssn", "social security",
	}
}

// Start begins IRC message processing
func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.isRunning {
		p.mu.Unlock()
		return fmt.Errorf("IRC processor is already running")
	}
	p.isRunning = true
	p.mu.Unlock()

	log.Printf("Starting IRC processor with config: %+v", p.config)

	// Connect to IRC servers
	for _, server := range p.config.Servers {
		go p.connectToServer(ctx, server)
	}

	// Start message processing goroutine
	go p.processMessages(ctx)

	log.Println("IRC processor started successfully")
	return nil
}

// Stop stops IRC message processing
func (p *Processor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return fmt.Errorf("IRC processor is not running")
	}

	close(p.stopChan)

	// Close all connections
	for _, conn := range p.connections {
		if conn.Connected && conn.Conn != nil {
			conn.Conn.Close()
		}
	}

	p.isRunning = false
	log.Println("IRC processor stopped")
	return nil
}

// connectToServer establishes connection to an IRC server
func (p *Processor) connectToServer(ctx context.Context, server string) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		default:
			if err := p.establishConnection(server); err != nil {
				log.Printf("Failed to connect to IRC server %s: %v", server, err)
				time.Sleep(time.Duration(p.config.ReconnectInterval) * time.Second)
				continue
			}

			// Handle messages from this connection
			p.handleConnection(ctx, server)

			// Reconnect after disconnect
			time.Sleep(time.Duration(p.config.ReconnectInterval) * time.Second)
		}
	}
}

// establishConnection creates connection to IRC server
func (p *Processor) establishConnection(server string) error {
	var conn net.Conn
	var err error

	if p.config.TLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: p.config.TLSInsecure,
		}
		conn, err = tls.Dial("tcp", server, tlsConfig)
	} else {
		conn, err = net.Dial("tcp", server)
	}

	if err != nil {
		return err
	}

	p.mu.Lock()
	p.connections[server] = &IRCConnection{
		Server:    server,
		Conn:      conn,
		Connected: true,
		LastPing:  time.Now(),
		Channels:  make([]string, 0),
	}
	p.mu.Unlock()

	// Send IRC registration
	p.sendIRCCommand(server, fmt.Sprintf("NICK %s", p.config.Nickname))
	p.sendIRCCommand(server, fmt.Sprintf("USER %s 0 * :%s", p.config.Username, p.config.Realname))

	if p.config.Password != "" {
		p.sendIRCCommand(server, fmt.Sprintf("PASS %s", p.config.Password))
	}

	// Join channels
	for _, channel := range p.config.Channels {
		p.sendIRCCommand(server, fmt.Sprintf("JOIN %s", channel))
	}

	return nil
}

// handleConnection processes messages from IRC connection
func (p *Processor) handleConnection(ctx context.Context, server string) {
	p.mu.RLock()
	conn := p.connections[server]
	p.mu.RUnlock()

	if conn == nil || !conn.Connected {
		return
	}

	buffer := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		default:
			conn.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			n, err := conn.Conn.Read(buffer)
			if err != nil {
				log.Printf("IRC connection error for %s: %v", server, err)
				conn.Connected = false
				return
			}

			data := string(buffer[:n])
			lines := strings.Split(strings.TrimSpace(data), "\n")

			for _, line := range lines {
				if line != "" {
					p.processIRCLine(server, line)
				}
			}
		}
	}
}

// processIRCLine processes a single IRC protocol line
func (p *Processor) processIRCLine(server, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Handle PING
	if strings.HasPrefix(line, "PING ") {
		pong := strings.Replace(line, "PING", "PONG", 1)
		p.sendIRCCommand(server, pong)
		return
	}

	// Parse IRC message
	msg := p.parseIRCMessage(server, line)
	if msg == nil {
		return
	}

	// Process message for threats
	if msg.MessageType == "PRIVMSG" || msg.MessageType == "NOTICE" {
		processedMsg, err := p.ProcessMessage(*msg)
		if err != nil {
			log.Printf("Error processing IRC message: %v", err)
			return
		}

		// Log high-threat messages
		if processedMsg.ThreatLevel >= ThreatLevelHigh {
			log.Printf("High-threat IRC message detected: server=%s, channel=%s, nick=%s, threat=%s",
				processedMsg.Server, processedMsg.Channel, processedMsg.Nick, processedMsg.ThreatLevel.String())
		}
	}
}

// parseIRCMessage parses IRC protocol message
func (p *Processor) parseIRCMessage(server, line string) *IRCMessage {
	// Basic IRC message parsing
	// Format: :nick!user@host COMMAND target :message

	if !strings.HasPrefix(line, ":") {
		return nil
	}

	parts := strings.SplitN(line[1:], " ", 4)
	if len(parts) < 3 {
		return nil
	}

	// Parse nick!user@host
	hostmask := parts[0]
	nickParts := strings.SplitN(hostmask, "!", 2)
	nick := nickParts[0]

	var user, host string
	if len(nickParts) > 1 {
		userHost := strings.SplitN(nickParts[1], "@", 2)
		user = userHost[0]
		if len(userHost) > 1 {
			host = userHost[1]
		}
	}

	command := parts[1]
	target := parts[2]

	var message string
	if len(parts) > 3 && strings.HasPrefix(parts[3], ":") {
		message = parts[3][1:]
	}

	return &IRCMessage{
		ID:          fmt.Sprintf("%s-%d", server, time.Now().UnixNano()),
		Server:      server,
		Channel:     target,
		Nick:        nick,
		User:        user,
		Host:        host,
		Message:     message,
		MessageType: command,
		Timestamp:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}
}

// ProcessMessage analyzes an IRC message for threats
func (p *Processor) ProcessMessage(msg IRCMessage) (*IRCMessage, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Update user tracking
	p.updateUserInfo(msg.Nick, msg.User, msg.Host)

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

// analyzeThreat performs comprehensive threat analysis on IRC message
func (p *Processor) analyzeThreat(msg IRCMessage) *ThreatAnalysis {
	analysis := &ThreatAnalysis{
		ThreatTypes:     make([]string, 0),
		Indicators:      make([]string, 0),
		Recommendations: make([]string, 0),
		ProcessedAt:     time.Now(),
	}

	message := strings.ToLower(msg.Message)

	// Check threat patterns
	for patternName, pattern := range p.threatPatterns {
		if pattern.MatchString(msg.Message) {
			analysis.ThreatTypes = append(analysis.ThreatTypes, patternName)
			analysis.Indicators = append(analysis.Indicators, fmt.Sprintf("Pattern match: %s", patternName))
		}
	}

	// Calculate various threat scores
	analysis.SpamScore = p.calculateSpamScore(message)
	analysis.PhishingScore = p.calculatePhishingScore(message)
	analysis.MalwareScore = p.calculateMalwareScore(message)
	analysis.BotScore = p.calculateBotScore(msg)
	analysis.FloodScore = p.calculateFloodScore(msg)
	analysis.SocialEngScore = p.calculateSocialEngScore(message)

	// Calculate overall confidence
	analysis.Confidence = p.calculateOverallConfidence(analysis)

	// Generate recommendations
	analysis.Recommendations = p.generateRecommendations(analysis)

	return analysis
}

// calculateSpamScore calculates spam likelihood
func (p *Processor) calculateSpamScore(message string) float64 {
	score := 0.0

	for _, word := range p.suspiciousWords {
		if strings.Contains(message, word) {
			score += 0.2
		}
	}

	// Check for excessive punctuation
	if strings.Count(message, "!") > 3 {
		score += 0.3
	}

	// Check for excessive caps
	capsCount := 0
	for _, r := range message {
		if r >= 'A' && r <= 'Z' {
			capsCount++
		}
	}
	if len(message) > 0 && float64(capsCount)/float64(len(message)) > 0.5 {
		score += 0.4
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculatePhishingScore calculates phishing likelihood
func (p *Processor) calculatePhishingScore(message string) float64 {
	score := 0.0

	phishingKeywords := []string{
		"login", "verify", "account", "password", "suspended",
		"click here", "urgent", "immediate", "security alert",
	}

	for _, keyword := range phishingKeywords {
		if strings.Contains(message, keyword) {
			score += 0.25
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateMalwareScore calculates malware likelihood
func (p *Processor) calculateMalwareScore(message string) float64 {
	score := 0.0

	malwareKeywords := []string{
		"download", "install", "exe", "zip", "rar",
		"crack", "keygen", "warez", "torrent",
	}

	for _, keyword := range malwareKeywords {
		if strings.Contains(message, keyword) {
			score += 0.3
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateBotScore calculates bot likelihood
func (p *Processor) calculateBotScore(msg IRCMessage) float64 {
	score := 0.0

	// Check nick pattern
	if p.threatPatterns["bot_pattern"].MatchString(msg.Nick) {
		score += 0.5
	}

	// Check for automated patterns
	if strings.Contains(msg.Message, "[") && strings.Contains(msg.Message, "]") {
		score += 0.2
	}

	// Check user info (avoid deadlock by not acquiring lock again)
	key := fmt.Sprintf("%s!%s@%s", msg.Nick, msg.User, msg.Host)
	userInfo := p.userHistory[key]

	if userInfo != nil {
		// High message frequency indicates bot
		if userInfo.MessageCount > 100 && time.Since(userInfo.FirstSeen) < time.Hour {
			score += 0.4
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateFloodScore calculates flood likelihood
func (p *Processor) calculateFloodScore(msg IRCMessage) float64 {
	score := 0.0

	// Check for repeated characters using simple pattern
	if len(msg.Message) > 10 {
		// Count consecutive repeated characters
		maxRepeats := 0
		currentRepeats := 1
		for i := 1; i < len(msg.Message); i++ {
			if msg.Message[i] == msg.Message[i-1] {
				currentRepeats++
			} else {
				if currentRepeats > maxRepeats {
					maxRepeats = currentRepeats
				}
				currentRepeats = 1
			}
		}
		if maxRepeats > 10 {
			score += 0.6
		}
	}

	// Check message frequency (avoid deadlock by not acquiring lock again)
	key := fmt.Sprintf("%s!%s@%s", msg.Nick, msg.User, msg.Host)
	userInfo := p.userHistory[key]

	if userInfo != nil && len(userInfo.RecentMessages) > 0 {
		// Check for identical messages
		for _, recentMsg := range userInfo.RecentMessages {
			if recentMsg == msg.Message {
				score += 0.4
				break
			}
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateSocialEngScore calculates social engineering likelihood
func (p *Processor) calculateSocialEngScore(message string) float64 {
	score := 0.0

	socialEngKeywords := []string{
		"urgent", "immediate", "act now", "limited time",
		"verify", "confirm", "suspended", "locked",
	}

	for _, keyword := range socialEngKeywords {
		if strings.Contains(message, keyword) {
			score += 0.2
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
		analysis.FloodScore,
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
		recommendations = append(recommendations, "Consider ignoring or blocking user")
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

	if analysis.FloodScore > 0.5 {
		recommendations = append(recommendations, "Potential flood attack - consider rate limiting")
	}

	if analysis.Confidence > 0.7 {
		recommendations = append(recommendations, "High threat detected - consider reporting to channel operators")
	}

	return recommendations
}

// updateUserInfo updates tracking information for users
func (p *Processor) updateUserInfo(nick, user, host string) {
	key := fmt.Sprintf("%s!%s@%s", nick, user, host)

	info, exists := p.userHistory[key]
	if !exists {
		info = &UserInfo{
			Nick:           nick,
			User:           user,
			Host:           host,
			FirstSeen:      time.Now(),
			LastSeen:       time.Now(),
			MessageCount:   0,
			ThreatCount:    0,
			Reputation:     0.5, // Neutral starting reputation
			RecentMessages: make([]string, 0),
		}
		p.userHistory[key] = info
	}

	info.LastSeen = time.Now()
	info.MessageCount++

	// Keep only recent messages for flood detection
	if len(info.RecentMessages) >= 10 {
		info.RecentMessages = info.RecentMessages[1:]
	}
}

// sendIRCCommand sends a command to IRC server
func (p *Processor) sendIRCCommand(server, command string) error {
	p.mu.RLock()
	conn := p.connections[server]
	p.mu.RUnlock()

	if conn == nil || !conn.Connected {
		return fmt.Errorf("not connected to server %s", server)
	}

	_, err := conn.Conn.Write([]byte(command + "\r\n"))
	return err
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
	for key, user := range p.userHistory {
		if user.LastSeen.Before(cutoff) {
			delete(p.userHistory, key)
		}
	}

	// Check connection health
	for server, conn := range p.connections {
		if conn.Connected && time.Since(conn.LastPing) > 5*time.Minute {
			log.Printf("IRC connection to %s appears stale, will reconnect", server)
			conn.Connected = false
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

	connectedServers := 0
	for _, conn := range p.connections {
		if conn.Connected {
			connectedServers++
		}
	}

	return map[string]interface{}{
		"total_messages":     len(p.messageHistory),
		"unique_users":       len(p.userHistory),
		"threat_counts":      threatCounts,
		"is_running":         p.isRunning,
		"connected_servers":  connectedServers,
		"total_servers":      len(p.config.Servers),
		"monitored_channels": len(p.config.Channels),
	}
}
