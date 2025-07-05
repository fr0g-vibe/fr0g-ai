package voice

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

// Processor handles voice call threat detection and analysis
type Processor struct {
	config          *config.VoiceConfig
	threatPatterns  map[string]*regexp.Regexp
	suspiciousWords []string
	callHistory     []VoiceCall
	callerInfo      map[string]*CallerInfo
	mu              sync.RWMutex
	isRunning       bool
	stopChan        chan struct{}
}

// VoiceCall represents a voice call with analysis
type VoiceCall struct {
	ID          string                 `json:"id"`
	CallerID    string                 `json:"caller_id"`
	RecipientID string                 `json:"recipient_id"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	Duration    time.Duration          `json:"duration"`
	Transcript  string                 `json:"transcript"`
	ThreatLevel ThreatLevel            `json:"threat_level"`
	Analysis    *ThreatAnalysis        `json:"analysis,omitempty"`
	AudioFile   string                 `json:"audio_file,omitempty"`
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

// ThreatAnalysis contains detailed voice threat analysis results
type ThreatAnalysis struct {
	ThreatTypes         []string  `json:"threat_types"`
	Confidence          float64   `json:"confidence"`
	ScamScore           float64   `json:"scam_score"`
	PhishingScore       float64   `json:"phishing_score"`
	SocialEngScore      float64   `json:"social_eng_score"`
	RobocallScore       float64   `json:"robocall_score"`
	EmotionalManipScore float64   `json:"emotional_manip_score"`
	Indicators          []string  `json:"indicators"`
	Recommendations     []string  `json:"recommendations"`
	SpeechPatterns      []string  `json:"speech_patterns"`
	ProcessedAt         time.Time `json:"processed_at"`
}

// CallerInfo tracks information about phone numbers
type CallerInfo struct {
	CallerID      string        `json:"caller_id"`
	FirstSeen     time.Time     `json:"first_seen"`
	LastSeen      time.Time     `json:"last_seen"`
	CallCount     int           `json:"call_count"`
	ThreatCount   int           `json:"threat_count"`
	IsBlacklisted bool          `json:"is_blacklisted"`
	IsWhitelisted bool          `json:"is_whitelisted"`
	Reputation    float64       `json:"reputation"` // 0.0-1.0, higher is better
	AvgCallLength time.Duration `json:"avg_call_length"`
}

// NewProcessor creates a new voice processor instance
func NewProcessor(cfg *sharedconfig.VoiceConfig) *Processor {
	p := &Processor{
		config:         cfg,
		threatPatterns: make(map[string]*regexp.Regexp),
		callHistory:    make([]VoiceCall, 0),
		callerInfo:     make(map[string]*CallerInfo),
		stopChan:       make(chan struct{}),
	}

	p.initializeThreatPatterns()
	p.initializeSuspiciousWords()

	return p
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "voice"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.config.Enabled
}

// initializeThreatPatterns sets up regex patterns for voice threat detection
func (p *Processor) initializeThreatPatterns() {
	patterns := map[string]string{
		"irs_scam":        `(?i)(irs|internal revenue|tax|refund|arrest|warrant)`,
		"tech_support":    `(?i)(microsoft|windows|computer|virus|infected|technical support)`,
		"bank_fraud":      `(?i)(bank|account|suspended|verify|security|fraud department)`,
		"social_security": `(?i)(social security|ssn|benefits|suspended|medicare)`,
		"credit_card":     `(?i)(credit card|debt|consolidation|lower interest|pre-approved)`,
		"lottery_scam":    `(?i)(lottery|winner|prize|congratulations|claim|sweepstakes)`,
		"charity_scam":    `(?i)(charity|donation|veterans|police|firefighters|urgent)`,
		"utility_scam":    `(?i)(electric|gas|water|utility|disconnect|payment|overdue)`,
		"insurance_scam":  `(?i)(insurance|policy|premium|coverage|expired|renewal)`,
		"investment_scam": `(?i)(investment|stocks|crypto|bitcoin|guaranteed return)`,
	}

	for name, pattern := range patterns {
		if compiled, err := regexp.Compile(pattern); err == nil {
			p.threatPatterns[name] = compiled
		} else {
			log.Printf("Failed to compile voice threat pattern %s: %v", name, err)
		}
	}
}

// initializeSuspiciousWords sets up suspicious words for voice analysis
func (p *Processor) initializeSuspiciousWords() {
	p.suspiciousWords = []string{
		"urgent", "immediate", "act now", "limited time", "expires today",
		"verify", "confirm", "update", "suspended", "locked", "frozen",
		"arrest", "warrant", "legal action", "court", "lawsuit", "fine",
		"refund", "owe", "debt", "payment", "overdue", "collection",
		"winner", "prize", "lottery", "sweepstakes", "congratulations",
		"free", "no cost", "guaranteed", "risk free", "limited offer",
		"social security", "medicare", "benefits", "disability",
		"bank account", "credit card", "routing number", "pin number",
		"wire transfer", "money order", "gift card", "prepaid card",
		"computer", "virus", "infected", "hacked", "security breach",
	}
}

// Start begins voice call processing
func (p *Processor) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.isRunning {
		p.mu.Unlock()
		return fmt.Errorf("voice processor is already running")
	}
	p.isRunning = true
	p.mu.Unlock()

	log.Printf("Starting voice processor with config: %+v", p.config)

	// Start call monitoring goroutine
	go p.monitorCalls(ctx)

	// Start speech-to-text service if configured
	if p.config.SpeechToTextEnabled {
		go p.startSpeechToTextService(ctx)
	}

	// Start call recording service if configured
	if p.config.CallRecordingEnabled {
		go p.startCallRecordingService(ctx)
	}

	log.Println("Voice processor started successfully")
	return nil
}

// Stop stops voice call processing
func (p *Processor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return fmt.Errorf("voice processor is not running")
	}

	close(p.stopChan)
	p.isRunning = false

	log.Println("Voice processor stopped")
	return nil
}

// ProcessCall analyzes a voice call for threats
func (p *Processor) ProcessCall(call VoiceCall) (*VoiceCall, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Update caller information
	p.updateCallerInfo(call.CallerID, call.Duration)

	// Perform threat analysis on transcript
	analysis := p.analyzeThreat(call)
	call.Analysis = analysis
	call.ThreatLevel = p.calculateThreatLevel(analysis)

	// Store call in history
	if p.config.MaxHistorySize > 0 {
		p.callHistory = append(p.callHistory, call)

		// Limit history size
		if len(p.callHistory) > p.config.MaxHistorySize {
			p.callHistory = p.callHistory[1:]
		}
	} else {
		// Always store at least some calls for testing
		p.callHistory = append(p.callHistory, call)
		if len(p.callHistory) > 1000 { // Default limit
			p.callHistory = p.callHistory[1:]
		}
	}

	log.Printf("Processed voice call from %s: threat_level=%s, confidence=%.2f, duration=%v",
		call.CallerID, call.ThreatLevel.String(), analysis.Confidence, call.Duration)

	return &call, nil
}

// analyzeThreat performs comprehensive threat analysis on voice call transcript
func (p *Processor) analyzeThreat(call VoiceCall) *ThreatAnalysis {
	analysis := &ThreatAnalysis{
		ThreatTypes:     make([]string, 0),
		Indicators:      make([]string, 0),
		Recommendations: make([]string, 0),
		SpeechPatterns:  make([]string, 0),
		ProcessedAt:     time.Now(),
	}

	transcript := strings.ToLower(call.Transcript)

	// Check threat patterns
	for patternName, pattern := range p.threatPatterns {
		if pattern.MatchString(call.Transcript) {
			analysis.ThreatTypes = append(analysis.ThreatTypes, patternName)
			analysis.Indicators = append(analysis.Indicators, fmt.Sprintf("Pattern match: %s", patternName))
		}
	}

	// Calculate individual threat scores
	analysis.ScamScore = p.calculateScamScore(transcript)
	analysis.PhishingScore = p.calculatePhishingScore(transcript)
	analysis.SocialEngScore = p.calculateSocialEngScore(transcript)
	analysis.RobocallScore = p.calculateRobocallScore(call)
	analysis.EmotionalManipScore = p.calculateEmotionalManipScore(transcript)

	// Calculate overall confidence
	analysis.Confidence = p.calculateOverallConfidence(analysis)

	// Analyze speech patterns
	analysis.SpeechPatterns = p.analyzeSpeechPatterns(transcript)

	// Generate recommendations
	analysis.Recommendations = p.generateRecommendations(analysis)

	return analysis
}

// calculateScamScore calculates general scam likelihood
func (p *Processor) calculateScamScore(transcript string) float64 {
	score := 0.0
	wordCount := 0

	for _, word := range p.suspiciousWords {
		if strings.Contains(transcript, word) {
			wordCount++
			score += 0.1
		}
	}

	// Bonus for multiple suspicious words
	if wordCount > 3 {
		score += float64(wordCount-3) * 0.05
	}

	// Check for urgency indicators
	urgencyWords := []string{"urgent", "immediate", "act now", "expires", "deadline"}
	for _, word := range urgencyWords {
		if strings.Contains(transcript, word) {
			score += 0.2
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculatePhishingScore calculates phishing likelihood
func (p *Processor) calculatePhishingScore(transcript string) float64 {
	score := 0.0

	// Check for information requests
	infoRequests := []string{"social security", "bank account", "credit card", "pin", "password"}
	for _, request := range infoRequests {
		if strings.Contains(transcript, request) {
			score += 0.3
		}
	}

	// Check for verification requests
	if strings.Contains(transcript, "verify") || strings.Contains(transcript, "confirm") {
		score += 0.2
	}

	// Check for account issues
	if strings.Contains(transcript, "suspended") || strings.Contains(transcript, "locked") {
		score += 0.3
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateSocialEngScore calculates social engineering likelihood
func (p *Processor) calculateSocialEngScore(transcript string) float64 {
	score := 0.0

	// Check for authority claims
	authorities := []string{"irs", "fbi", "police", "bank", "microsoft", "government"}
	for _, auth := range authorities {
		if strings.Contains(transcript, auth) {
			score += 0.2
		}
	}

	// Check for fear tactics
	fearWords := []string{"arrest", "warrant", "legal action", "court", "fine", "penalty"}
	for _, word := range fearWords {
		if strings.Contains(transcript, word) {
			score += 0.25
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateRobocallScore calculates robocall likelihood
func (p *Processor) calculateRobocallScore(call VoiceCall) float64 {
	score := 0.0

	// Very short calls are often robocalls
	if call.Duration < 30*time.Second {
		score += 0.3
	}

	// Check for robocall indicators in transcript
	robocallPhrases := []string{
		"this is not a sales call",
		"press 1 to",
		"press 9 to be removed",
		"recorded message",
		"automated message",
	}

	for _, phrase := range robocallPhrases {
		if strings.Contains(strings.ToLower(call.Transcript), phrase) {
			score += 0.4
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateEmotionalManipScore calculates emotional manipulation likelihood
func (p *Processor) calculateEmotionalManipScore(transcript string) float64 {
	score := 0.0

	// Check for emotional triggers
	emotionalWords := []string{
		"urgent", "emergency", "crisis", "help", "save", "protect",
		"family", "loved ones", "children", "grandchildren",
		"limited time", "last chance", "expires today",
	}

	for _, word := range emotionalWords {
		if strings.Contains(transcript, word) {
			score += 0.15
		}
	}

	// Check for pressure tactics
	pressureWords := []string{"act now", "don't wait", "call back immediately", "time sensitive"}
	for _, word := range pressureWords {
		if strings.Contains(transcript, word) {
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
		analysis.ScamScore,
		analysis.PhishingScore,
		analysis.SocialEngScore,
		analysis.RobocallScore,
		analysis.EmotionalManipScore,
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

// analyzeSpeechPatterns analyzes speech patterns for additional insights
func (p *Processor) analyzeSpeechPatterns(transcript string) []string {
	patterns := make([]string, 0)

	// Check for script reading indicators
	if strings.Count(transcript, ".") > len(transcript)/50 {
		patterns = append(patterns, "formal_script_reading")
	}

	// Check for repetitive phrases
	words := strings.Fields(transcript)
	if len(words) > 5 { // Reduced threshold from 10 to 5
		// Simple repetition check
		wordCount := make(map[string]int)
		for _, word := range words {
			cleanWord := strings.ToLower(strings.Trim(word, ".,!?"))
			if len(cleanWord) > 2 { // Only count words longer than 2 characters
				wordCount[cleanWord]++
			}
		}

		for word, count := range wordCount {
			if count > 2 && len(word) > 2 { // Reduced threshold from 3 to 2
				patterns = append(patterns, fmt.Sprintf("repetitive_word_%s", word))
			}
		}
	}

	// Check for fast talking (high word density)
	if len(words) > 200 {
		patterns = append(patterns, "high_word_density")
	}

	return patterns
}

// generateRecommendations generates security recommendations
func (p *Processor) generateRecommendations(analysis *ThreatAnalysis) []string {
	recommendations := make([]string, 0)

	if analysis.ScamScore > 0.5 {
		recommendations = append(recommendations, "Hang up immediately and block the number")
	}

	if analysis.PhishingScore > 0.5 {
		recommendations = append(recommendations, "Never provide personal information over unsolicited calls")
	}

	if analysis.SocialEngScore > 0.5 {
		recommendations = append(recommendations, "Verify caller identity through official channels")
	}

	if analysis.RobocallScore > 0.5 {
		recommendations = append(recommendations, "Report robocall to FTC and add number to block list")
	}

	if analysis.EmotionalManipScore > 0.5 {
		recommendations = append(recommendations, "Take time to think - legitimate organizations don't pressure for immediate action")
	}

	if analysis.Confidence > 0.7 {
		recommendations = append(recommendations, "Report to authorities and warn others about this scam")
	}

	return recommendations
}

// updateCallerInfo updates tracking information for callers
func (p *Processor) updateCallerInfo(callerID string, duration time.Duration) {
	info, exists := p.callerInfo[callerID]
	if !exists {
		info = &CallerInfo{
			CallerID:      callerID,
			FirstSeen:     time.Now(),
			LastSeen:      time.Now(),
			CallCount:     0,
			ThreatCount:   0,
			Reputation:    0.5, // Neutral starting reputation
			AvgCallLength: 0,
		}
		p.callerInfo[callerID] = info
	}

	info.LastSeen = time.Now()
	info.CallCount++

	// Update average call length
	if info.CallCount == 1 {
		info.AvgCallLength = duration
	} else {
		info.AvgCallLength = time.Duration((int64(info.AvgCallLength)*(int64(info.CallCount)-1) + int64(duration)) / int64(info.CallCount))
	}
}

// monitorCalls handles continuous call monitoring
func (p *Processor) monitorCalls(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(p.config.MonitoringInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			// Monitor for new calls
			p.checkForNewCalls()
		}
	}
}

// checkForNewCalls checks for new incoming calls
func (p *Processor) checkForNewCalls() {
	// Implementation would depend on telephony system integration
	// For now, this is a placeholder for future integration
	log.Println("Monitoring for new voice calls...")
}

// startSpeechToTextService starts speech-to-text processing
func (p *Processor) startSpeechToTextService(ctx context.Context) {
	log.Println("Starting speech-to-text service...")

	// TODO: Implement speech-to-text integration
	// This would involve:
	// 1. Integration with speech recognition APIs (Google, AWS, Azure)
	// 2. Real-time audio stream processing
	// 3. Transcript generation and formatting
	// 4. Integration with threat analysis pipeline

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-time.After(30 * time.Second):
			// Placeholder for periodic speech processing
			log.Println("Processing speech-to-text queue...")
		}
	}
}

// startCallRecordingService starts call recording service
func (p *Processor) startCallRecordingService(ctx context.Context) {
	log.Println("Starting call recording service...")

	// TODO: Implement call recording
	// This would involve:
	// 1. Audio capture from telephony system
	// 2. Audio file storage and management
	// 3. Compliance with recording laws
	// 4. Integration with analysis pipeline

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-time.After(60 * time.Second):
			// Placeholder for recording management
			log.Println("Managing call recordings...")
		}
	}
}

// GetStats returns processor statistics
func (p *Processor) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	threatCounts := make(map[string]int)
	totalDuration := time.Duration(0)

	for _, call := range p.callHistory {
		threatCounts[call.ThreatLevel.String()]++
		totalDuration += call.Duration
	}

	avgDuration := time.Duration(0)
	if len(p.callHistory) > 0 {
		avgDuration = totalDuration / time.Duration(len(p.callHistory))
	}

	return map[string]interface{}{
		"total_calls":            len(p.callHistory),
		"unique_callers":         len(p.callerInfo),
		"threat_counts":          threatCounts,
		"average_call_duration":  avgDuration.String(),
		"total_call_duration":    totalDuration.String(),
		"is_running":             p.isRunning,
		"speech_to_text_enabled": p.config.SpeechToTextEnabled,
		"call_recording_enabled": p.config.CallRecordingEnabled,
	}
}
