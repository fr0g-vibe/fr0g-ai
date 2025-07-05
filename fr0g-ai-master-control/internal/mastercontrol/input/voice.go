package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// VoiceProcessor processes voice/audio messages for threat analysis
type VoiceProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *VoiceConfig
}

// VoiceConfig holds voice processor configuration
type VoiceConfig struct {
	Provider             string        `yaml:"provider"` // "twilio", "aws_transcribe", "google_speech", etc.
	APIKey               string        `yaml:"api_key"`
	APISecret            string        `yaml:"api_secret"`
	WebhookURL           string        `yaml:"webhook_url"`
	CommunityTopic       string        `yaml:"community_topic"`
	PersonaCount         int           `yaml:"persona_count"`
	ReviewTimeout        time.Duration `yaml:"review_timeout"`
	RequiredConsensus    float64       `yaml:"required_consensus"`
	TrustedNumbers       []string      `yaml:"trusted_numbers"`
	BlockedNumbers       []string      `yaml:"blocked_numbers"`
	MaxRecordingDuration time.Duration `yaml:"max_recording_duration"`
	SupportedFormats     []string      `yaml:"supported_formats"`
	TranscriptionEnabled bool          `yaml:"transcription_enabled"`
	SentimentAnalysis    bool          `yaml:"sentiment_analysis"`
	VoiceprintAnalysis   bool          `yaml:"voiceprint_analysis"`
	AudioStoragePath     string        `yaml:"audio_storage_path"`
}

// VoiceMessage represents a voice/audio message
type VoiceMessage struct {
	ID                string            `json:"id"`
	From              string            `json:"from"`
	To                string            `json:"to"`
	CallSID           string            `json:"call_sid"`
	RecordingURL      string            `json:"recording_url"`
	RecordingDuration float64           `json:"recording_duration"`
	Transcription     string            `json:"transcription,omitempty"`
	Confidence        float64           `json:"confidence,omitempty"`
	Language          string            `json:"language,omitempty"`
	AudioFormat       string            `json:"audio_format"`
	FileSize          int64             `json:"file_size"`
	Direction         string            `json:"direction"` // "inbound", "outbound"
	Status            string            `json:"status"`
	Timestamp         time.Time         `json:"timestamp"`
	Country           string            `json:"country,omitempty"`
	Carrier           string            `json:"carrier,omitempty"`
	VoiceAnalysis     *VoiceAnalysis    `json:"voice_analysis,omitempty"`
	Metadata          map[string]string `json:"metadata"`
}

// VoiceAnalysis represents voice analysis results
type VoiceAnalysis struct {
	SentimentScore  float64            `json:"sentiment_score"`
	EmotionScores   map[string]float64 `json:"emotion_scores"`
	StressLevel     float64            `json:"stress_level"`
	SpeechRate      float64            `json:"speech_rate"`
	VoiceprintID    string             `json:"voiceprint_id,omitempty"`
	SpeakerGender   string             `json:"speaker_gender,omitempty"`
	EstimatedAge    int                `json:"estimated_age,omitempty"`
	AccentRegion    string             `json:"accent_region,omitempty"`
	BackgroundNoise float64            `json:"background_noise"`
	AudioQuality    float64            `json:"audio_quality"`
}

// NewVoiceProcessor creates a new voice processor
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
	return fmt.Sprintf("Voice/Audio Threat Vector Interceptor via %s - Voice intelligence gathering for AI community review on topic: %s",
		v.config.Provider, v.config.CommunityTopic)
}

// ProcessWebhook processes a voice message webhook
func (v *VoiceProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("Voice Processor: Processing voice threat vector webhook")

	// Parse voice message from request body
	voiceMsg, err := v.parseVoiceMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse voice message: %w", err)
	}

	// Check if number is blocked
	if v.isNumberBlocked(voiceMsg.From) {
		log.Printf("Voice Processor: Blocked number %s, dropping recording", voiceMsg.From)
		return &WebhookResponse{
			Success:   true,
			Message:   "Voice call from blocked number dropped",
			RequestID: request.ID,
			Data: map[string]interface{}{
				"action": "blocked",
				"reason": "blocked_number",
				"from":   voiceMsg.From,
			},
			Timestamp: time.Now(),
		}, nil
	}

	// Check if number is trusted (lower threat threshold)
	isTrusted := v.isNumberTrusted(voiceMsg.From)

	// Analyze voice message for threats using AI community
	threatLevel, consensus, err := v.analyzeVoiceThreats(ctx, voiceMsg, isTrusted)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze voice threats: %w", err)
	}

	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "Voice threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"threat_level":             threatLevel,
			"consensus":                consensus,
			"review_id":                fmt.Sprintf("voice_review_%d", time.Now().UnixNano()),
			"from":                     voiceMsg.From,
			"duration":                 voiceMsg.RecordingDuration,
			"has_transcription":        voiceMsg.Transcription != "",
			"transcription_confidence": voiceMsg.Confidence,
			"is_trusted":               isTrusted,
		},
		Timestamp: time.Now(),
	}

	log.Printf("Voice Processor: Recording analyzed - From: %s, Duration: %.2fs, Threat Level: %s, Consensus: %.2f",
		voiceMsg.From, voiceMsg.RecordingDuration, threatLevel, consensus)

	return response, nil
}

// analyzeVoiceThreats analyzes voice message for threats using AI community
func (v *VoiceProcessor) analyzeVoiceThreats(ctx context.Context, voiceMsg *VoiceMessage, isTrusted bool) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
Voice/Audio Threat Analysis Request:
From: %s
To: %s
Call SID: %s
Direction: %s
Duration: %.2f seconds
Audio Format: %s
File Size: %d bytes
Is Trusted Number: %v
Country: %s
Carrier: %s
Timestamp: %s

`, voiceMsg.From, voiceMsg.To, voiceMsg.CallSID, voiceMsg.Direction,
		voiceMsg.RecordingDuration, voiceMsg.AudioFormat, voiceMsg.FileSize,
		isTrusted, voiceMsg.Country, voiceMsg.Carrier,
		voiceMsg.Timestamp.Format(time.RFC3339))

	// Add transcription if available
	if voiceMsg.Transcription != "" {
		analysisContent += fmt.Sprintf(`
Transcription (Confidence: %.2f):
%s

`, voiceMsg.Confidence, voiceMsg.Transcription)
	}

	// Add voice analysis if available
	if voiceMsg.VoiceAnalysis != nil {
		analysisContent += fmt.Sprintf(`
Voice Analysis:
- Sentiment Score: %.2f
- Stress Level: %.2f
- Speech Rate: %.2f words/minute
- Speaker Gender: %s
- Estimated Age: %d
- Accent Region: %s
- Background Noise: %.2f
- Audio Quality: %.2f

`, voiceMsg.VoiceAnalysis.SentimentScore, voiceMsg.VoiceAnalysis.StressLevel,
			voiceMsg.VoiceAnalysis.SpeechRate, voiceMsg.VoiceAnalysis.SpeakerGender,
			voiceMsg.VoiceAnalysis.EstimatedAge, voiceMsg.VoiceAnalysis.AccentRegion,
			voiceMsg.VoiceAnalysis.BackgroundNoise, voiceMsg.VoiceAnalysis.AudioQuality)

		if len(voiceMsg.VoiceAnalysis.EmotionScores) > 0 {
			analysisContent += "Emotion Scores:\n"
			for emotion, score := range voiceMsg.VoiceAnalysis.EmotionScores {
				analysisContent += fmt.Sprintf("- %s: %.2f\n", emotion, score)
			}
		}
	}

	analysisContent += `
Please analyze this voice/audio message for potential threats including:
- Vishing (voice phishing) attacks
- Social engineering and pretexting
- Robocalls and automated scams
- Tech support scams
- IRS and government impersonation
- Banking and financial fraud calls
- Romance and relationship scams
- Charity and donation scams
- Warranty and insurance scams
- Debt collection scams
- Medical and health scams
- Voice deepfakes and impersonation
- Threatening or harassing calls
- Identity theft attempts
- Cryptocurrency and investment scams
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

// parseVoiceMessage parses a voice message from the request body
func (v *VoiceProcessor) parseVoiceMessage(body interface{}) (*VoiceMessage, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}

	voiceMsg := &VoiceMessage{
		ID:            getStringFromMap(bodyMap, "id"),
		From:          getStringFromMap(bodyMap, "from"),
		To:            getStringFromMap(bodyMap, "to"),
		CallSID:       getStringFromMap(bodyMap, "call_sid"),
		RecordingURL:  getStringFromMap(bodyMap, "recording_url"),
		Transcription: getStringFromMap(bodyMap, "transcription"),
		Language:      getStringFromMap(bodyMap, "language"),
		AudioFormat:   getStringFromMap(bodyMap, "audio_format"),
		Direction:     getStringFromMap(bodyMap, "direction"),
		Status:        getStringFromMap(bodyMap, "status"),
		Country:       getStringFromMap(bodyMap, "country"),
		Carrier:       getStringFromMap(bodyMap, "carrier"),
		Timestamp:     time.Now(),
		Metadata:      make(map[string]string),
	}

	// Parse numeric fields
	if duration, ok := bodyMap["recording_duration"].(float64); ok {
		voiceMsg.RecordingDuration = duration
	}

	if confidence, ok := bodyMap["confidence"].(float64); ok {
		voiceMsg.Confidence = confidence
	}

	if fileSize, ok := bodyMap["file_size"].(float64); ok {
		voiceMsg.FileSize = int64(fileSize)
	}

	// Parse voice analysis if available
	if analysisData, ok := bodyMap["voice_analysis"].(map[string]interface{}); ok {
		voiceMsg.VoiceAnalysis = &VoiceAnalysis{
			SpeakerGender: getStringFromMap(analysisData, "speaker_gender"),
			AccentRegion:  getStringFromMap(analysisData, "accent_region"),
			VoiceprintID:  getStringFromMap(analysisData, "voiceprint_id"),
		}

		if sentimentScore, ok := analysisData["sentiment_score"].(float64); ok {
			voiceMsg.VoiceAnalysis.SentimentScore = sentimentScore
		}

		if stressLevel, ok := analysisData["stress_level"].(float64); ok {
			voiceMsg.VoiceAnalysis.StressLevel = stressLevel
		}

		if speechRate, ok := analysisData["speech_rate"].(float64); ok {
			voiceMsg.VoiceAnalysis.SpeechRate = speechRate
		}

		if estimatedAge, ok := analysisData["estimated_age"].(float64); ok {
			voiceMsg.VoiceAnalysis.EstimatedAge = int(estimatedAge)
		}

		if backgroundNoise, ok := analysisData["background_noise"].(float64); ok {
			voiceMsg.VoiceAnalysis.BackgroundNoise = backgroundNoise
		}

		if audioQuality, ok := analysisData["audio_quality"].(float64); ok {
			voiceMsg.VoiceAnalysis.AudioQuality = audioQuality
		}

		// Parse emotion scores
		if emotionData, ok := analysisData["emotion_scores"].(map[string]interface{}); ok {
			voiceMsg.VoiceAnalysis.EmotionScores = make(map[string]float64)
			for emotion, score := range emotionData {
				if scoreFloat, ok := score.(float64); ok {
					voiceMsg.VoiceAnalysis.EmotionScores[emotion] = scoreFloat
				}
			}
		}
	}

	// Parse metadata
	if metadataData, ok := bodyMap["metadata"].(map[string]interface{}); ok {
		for key, value := range metadataData {
			if valueStr, ok := value.(string); ok {
				voiceMsg.Metadata[key] = valueStr
			}
		}
	}

	return voiceMsg, nil
}

// isNumberBlocked checks if a phone number is in the blocked list
func (v *VoiceProcessor) isNumberBlocked(phoneNumber string) bool {
	normalizedNumber := v.normalizePhoneNumber(phoneNumber)
	for _, blocked := range v.config.BlockedNumbers {
		if v.normalizePhoneNumber(blocked) == normalizedNumber {
			return true
		}
	}
	return false
}

// isNumberTrusted checks if a phone number is in the trusted list
func (v *VoiceProcessor) isNumberTrusted(phoneNumber string) bool {
	normalizedNumber := v.normalizePhoneNumber(phoneNumber)
	for _, trusted := range v.config.TrustedNumbers {
		if v.normalizePhoneNumber(trusted) == normalizedNumber {
			return true
		}
	}
	return false
}

// normalizePhoneNumber normalizes phone number format for comparison
func (v *VoiceProcessor) normalizePhoneNumber(phoneNumber string) string {
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
