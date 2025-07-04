package input

import (
	"context"
	"fmt"
	"log"
	"time"
)

// VoiceProcessor processes incoming voice calls for threat analysis
type VoiceProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *VoiceConfig
}

// VoiceConfig holds voice call processor configuration
type VoiceConfig struct {
	Provider              string        `yaml:"provider"`                // "google_voice", "twilio", "webrtc"
	APIKey                string        `yaml:"api_key"`
	APISecret             string        `yaml:"api_secret"`
	WebhookURL            string        `yaml:"webhook_url"`
	PhoneNumber           string        `yaml:"phone_number"`
	SpeechToTextEnabled   bool          `yaml:"speech_to_text_enabled"`
	GoogleCloudProjectID  string        `yaml:"google_cloud_project_id"`
	GoogleCloudKeyFile    string        `yaml:"google_cloud_key_file"`
	CommunityTopic        string        `yaml:"community_topic"`
	PersonaCount          int           `yaml:"persona_count"`
	ReviewTimeout         time.Duration `yaml:"review_timeout"`
	RequiredConsensus     float64       `yaml:"required_consensus"`
	AutoRecording         bool          `yaml:"auto_recording"`
	MaxCallDuration       time.Duration `yaml:"max_call_duration"`
	BlockedNumbers        []string      `yaml:"blocked_numbers"`
	SuspiciousKeywords    []string      `yaml:"suspicious_keywords"`
}

// VoiceCall represents an incoming voice call
type VoiceCall struct {
	ID               string            `json:"id"`
	From             string            `json:"from"`
	To               string            `json:"to"`
	StartTime        time.Time         `json:"start_time"`
	EndTime          *time.Time        `json:"end_time,omitempty"`
	Duration         time.Duration     `json:"duration"`
	Status           string            `json:"status"` // "ringing", "answered", "completed", "failed"
	Provider         string            `json:"provider"`
	RecordingURL     string            `json:"recording_url,omitempty"`
	TranscriptText   string            `json:"transcript_text,omitempty"`
	TranscriptStatus string            `json:"transcript_status,omitempty"`
	CallerID         *CallerIDInfo     `json:"caller_id,omitempty"`
	Metadata         map[string]string `json:"metadata"`
}

// CallerIDInfo represents caller identification information
type CallerIDInfo struct {
	Name         string `json:"name"`
	Number       string `json:"number"`
	Location     string `json:"location"`
	CarrierName  string `json:"carrier_name"`
	LineType     string `json:"line_type"` // "mobile", "landline", "voip"
	IsSpam       bool   `json:"is_spam"`
	IsRoboCall   bool   `json:"is_robocall"`
	TrustScore   float64 `json:"trust_score"`
}

// NewVoiceProcessor creates a new voice call processor
func NewVoiceProcessor(config *VoiceConfig, aiClient AIPersonaCommunityClient) (*VoiceProcessor, error) {
	return &VoiceProcessor{
		aiClient: aiClient,
		config:   config,
	}, nil
}

// GetTag returns the processor tag
func (v *VoiceProcessor) GetTag() string {
	return "voice"
}

// GetDescription returns the processor description
func (v *VoiceProcessor) GetDescription() string {
	return fmt.Sprintf("Voice Call Threat Vector Interceptor via %s - Telephony intelligence gathering with speech-to-text analysis for AI community review on topic: %s", 
		v.config.Provider, v.config.CommunityTopic)
}

// ProcessWebhook processes a voice call webhook
func (v *VoiceProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("Voice Processor: Processing voice call threat vector webhook")
	
	// Parse voice call from request body
	call, err := v.parseVoiceCall(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse voice call: %w", err)
	}
	
	// Check if number is blocked
	if v.isNumberBlocked(call.From) {
		log.Printf("Voice Processor: Blocked number %s, rejecting call", call.From)
		return &WebhookResponse{
			Success:   true,
			Message:   "Call from blocked number rejected",
			RequestID: request.ID,
			Data: map[string]interface{}{
				"action":      "blocked",
				"reason":      "blocked_number",
				"from_number": call.From,
				"call_id":     call.ID,
			},
			Timestamp: time.Now(),
		}, nil
	}
	
	// Analyze call for threats using AI community
	threatLevel, consensus, err := v.analyzeCallThreats(ctx, call)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze call threats: %w", err)
	}
	
	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "Voice call threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"threat_level":    threatLevel,
			"consensus":       consensus,
			"review_id":       fmt.Sprintf("voice_review_%d", time.Now().UnixNano()),
			"call_id":         call.ID,
			"from_number":     call.From,
			"has_transcript":  call.TranscriptText != "",
			"call_duration":   call.Duration.Seconds(),
		},
		Timestamp: time.Now(),
	}
	
	log.Printf("Voice Processor: Call analyzed - From: %s, Duration: %.1fs, Threat Level: %s, Consensus: %.2f", 
		call.From, call.Duration.Seconds(), threatLevel, consensus)
	
	return response, nil
}

// analyzeCallThreats analyzes voice call for threats using AI community
func (v *VoiceProcessor) analyzeCallThreats(ctx context.Context, call *VoiceCall) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
Voice Call Threat Analysis Request:
From: %s
To: %s
Duration: %.1f seconds
Status: %s
Provider: %s
Start Time: %s

Caller ID Information:
- Name: %s
- Location: %s
- Carrier: %s
- Line Type: %s
- Spam Indicator: %v
- Robocall Indicator: %v
- Trust Score: %.2f

`, call.From, call.To, call.Duration.Seconds(), call.Status, call.Provider, 
	call.StartTime.Format(time.RFC3339),
	v.getCallerIDField(call.CallerID, "name"),
	v.getCallerIDField(call.CallerID, "location"),
	v.getCallerIDField(call.CallerID, "carrier"),
	v.getCallerIDField(call.CallerID, "line_type"),
	call.CallerID != nil && call.CallerID.IsSpam,
	call.CallerID != nil && call.CallerID.IsRoboCall,
	v.getCallerIDTrustScore(call.CallerID))
	
	// Add transcript if available
	if call.TranscriptText != "" {
		analysisContent += fmt.Sprintf(`
Call Transcript:
%s

Transcript Status: %s
`, call.TranscriptText, call.TranscriptStatus)
	}
	
	analysisContent += `
Please analyze this voice call for potential threats including:
- Robocalls and automated spam
- Phishing attempts (fake IRS, bank, tech support)
- Social engineering attacks
- Vishing (voice phishing)
- Scam attempts (lottery, romance, investment)
- Identity theft attempts
- Suspicious caller behavior patterns
- Fraudulent caller ID spoofing
- Telemarketing violations
- Threatening or harassing calls
`
	
	// Create AI community for threat analysis
	community, err := v.aiClient.CreateCommunity(ctx, v.config.CommunityTopic, v.config.PersonaCount)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to create AI community: %w", err)
	}
	
	// Submit for AI community review
	review, err := v.aiClient.SubmitForReview(ctx, community.ID, analysisContent)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to submit for review: %w", err)
	}
	
	// Determine threat level based on consensus
	threatLevel := "unknown"
	consensus := 0.0
	
	if review.Consensus != nil {
		consensus = review.Consensus.OverallScore
		
		// Voice calls have different threat thresholds
		if consensus >= 0.9 {
			threatLevel = "critical"
		} else if consensus >= 0.75 {
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

// parseVoiceCall parses a voice call from the request body
func (v *VoiceProcessor) parseVoiceCall(body interface{}) (*VoiceCall, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}
	
	call := &VoiceCall{
		ID:               getStringFromMap(bodyMap, "id"),
		From:             getStringFromMap(bodyMap, "from"),
		To:               getStringFromMap(bodyMap, "to"),
		Status:           getStringFromMap(bodyMap, "status"),
		Provider:         getStringFromMap(bodyMap, "provider"),
		RecordingURL:     getStringFromMap(bodyMap, "recording_url"),
		TranscriptText:   getStringFromMap(bodyMap, "transcript_text"),
		TranscriptStatus: getStringFromMap(bodyMap, "transcript_status"),
		StartTime:        time.Now(),
		Metadata:         make(map[string]string),
	}
	
	// Parse duration
	if durationStr := getStringFromMap(bodyMap, "duration"); durationStr != "" {
		if duration, err := time.ParseDuration(durationStr); err == nil {
			call.Duration = duration
		}
	}
	
	// Parse caller ID information
	if callerIDData, ok := bodyMap["caller_id"].(map[string]interface{}); ok {
		call.CallerID = &CallerIDInfo{
			Name:        getStringFromMap(callerIDData, "name"),
			Number:      getStringFromMap(callerIDData, "number"),
			Location:    getStringFromMap(callerIDData, "location"),
			CarrierName: getStringFromMap(callerIDData, "carrier_name"),
			LineType:    getStringFromMap(callerIDData, "line_type"),
			IsSpam:      getBoolFromMap(callerIDData, "is_spam"),
			IsRoboCall:  getBoolFromMap(callerIDData, "is_robocall"),
		}
		
		if trustScore, ok := callerIDData["trust_score"].(float64); ok {
			call.CallerID.TrustScore = trustScore
		}
	}
	
	// Parse metadata
	if metadataData, ok := bodyMap["metadata"].(map[string]interface{}); ok {
		for key, value := range metadataData {
			if valueStr, ok := value.(string); ok {
				call.Metadata[key] = valueStr
			}
		}
	}
	
	return call, nil
}

// isNumberBlocked checks if a phone number is blocked
func (v *VoiceProcessor) isNumberBlocked(number string) bool {
	for _, blocked := range v.config.BlockedNumbers {
		if number == blocked {
			return true
		}
	}
	return false
}

// Helper functions for caller ID
func (v *VoiceProcessor) getCallerIDField(callerID *CallerIDInfo, field string) string {
	if callerID == nil {
		return "unknown"
	}
	
	switch field {
	case "name":
		return callerID.Name
	case "location":
		return callerID.Location
	case "carrier":
		return callerID.CarrierName
	case "line_type":
		return callerID.LineType
	default:
		return "unknown"
	}
}

func (v *VoiceProcessor) getCallerIDTrustScore(callerID *CallerIDInfo) float64 {
	if callerID == nil {
		return 0.0
	}
	return callerID.TrustScore
}
