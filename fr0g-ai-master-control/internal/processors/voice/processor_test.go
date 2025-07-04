package voice

import (
	"context"
	"strings"
	"testing"
	"time"

	"fr0g-ai-master-control/internal/config"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 30,
		MaxHistorySize:     100,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.config != cfg {
		t.Error("Expected processor config to match input config")
	}

	if len(processor.threatPatterns) == 0 {
		t.Error("Expected threat patterns to be initialized")
	}

	if len(processor.suspiciousWords) == 0 {
		t.Error("Expected suspicious words to be initialized")
	}
}

func TestProcessCall(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:         true,
		MaxHistorySize:  100,
		ThreatThreshold: 0.5,
	}

	processor := NewProcessor(cfg)

	tests := []struct {
		name            string
		call            VoiceCall
		expectedThreat  ThreatLevel
		minConfidence   float64
	}{
		{
			name: "Normal call",
			call: VoiceCall{
				ID:         "call1",
				CallerID:   "+1234567890",
				StartTime:  time.Now().Add(-5 * time.Minute),
				EndTime:    time.Now(),
				Duration:   5 * time.Minute,
				Transcript: "Hello, I'm calling to check on my appointment tomorrow.",
			},
			expectedThreat: ThreatLevelNone,
			minConfidence:  0.0,
		},
		{
			name: "IRS scam call",
			call: VoiceCall{
				ID:         "call2",
				CallerID:   "+1111111111",
				StartTime:  time.Now().Add(-2 * time.Minute),
				EndTime:    time.Now(),
				Duration:   2 * time.Minute,
				Transcript: "This is the IRS. Your tax refund has been suspended. You must call back immediately or face arrest and legal action.",
			},
			expectedThreat: ThreatLevelLow,
			minConfidence:  0.2,
		},
		{
			name: "Tech support scam",
			call: VoiceCall{
				ID:         "call3",
				CallerID:   "+2222222222",
				StartTime:  time.Now().Add(-3 * time.Minute),
				EndTime:    time.Now(),
				Duration:   3 * time.Minute,
				Transcript: "This is Microsoft technical support. Your computer has been infected with a virus. We need immediate access to fix this urgent security breach.",
			},
			expectedThreat: ThreatLevelLow,
			minConfidence:  0.3,
		},
		{
			name: "Robocall",
			call: VoiceCall{
				ID:         "call4",
				CallerID:   "+3333333333",
				StartTime:  time.Now().Add(-15 * time.Second),
				EndTime:    time.Now(),
				Duration:   15 * time.Second,
				Transcript: "This is not a sales call. Press 1 to speak with a representative or press 9 to be removed from our list.",
			},
			expectedThreat: ThreatLevelLow,
			minConfidence:  0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.ProcessCall(tt.call)
			if err != nil {
				t.Fatalf("ProcessCall failed: %v", err)
			}

			if result.ThreatLevel < tt.expectedThreat {
				t.Errorf("Expected threat level >= %s, got %s", 
					tt.expectedThreat.String(), result.ThreatLevel.String())
			}

			if result.Analysis == nil {
				t.Fatal("Expected analysis to be present")
			}

			if result.Analysis.Confidence < tt.minConfidence {
				t.Errorf("Expected confidence >= %.2f, got %.2f", 
					tt.minConfidence, result.Analysis.Confidence)
			}
		})
	}
}

func TestCalculateScamScore(t *testing.T) {
	cfg := &config.VoiceConfig{Enabled: true}
	processor := NewProcessor(cfg)

	tests := []struct {
		name     string
		transcript string
		minScore float64
		maxScore float64
	}{
		{
			name:     "Normal conversation",
			transcript: "hello how are you today can we schedule a meeting",
			minScore: 0.0,
			maxScore: 0.2,
		},
		{
			name:     "Suspicious words",
			transcript: "urgent immediate verify confirm suspended account",
			minScore: 0.4,
			maxScore: 1.0,
		},
		{
			name:     "High threat content",
			transcript: "urgent arrest warrant legal action immediate payment required",
			minScore: 0.6,
			maxScore: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := processor.calculateScamScore(tt.transcript)
			
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Expected scam score between %.2f and %.2f, got %.2f", 
					tt.minScore, tt.maxScore, score)
			}
		})
	}
}

func TestCalculatePhishingScore(t *testing.T) {
	cfg := &config.VoiceConfig{Enabled: true}
	processor := NewProcessor(cfg)

	tests := []struct {
		name     string
		transcript string
		minScore float64
		maxScore float64
	}{
		{
			name:     "Normal conversation",
			transcript: "hello how are you doing today",
			minScore: 0.0,
			maxScore: 0.2,
		},
		{
			name:     "Information request",
			transcript: "we need to verify your social security number and bank account",
			minScore: 0.5,
			maxScore: 1.0,
		},
		{
			name:     "Account suspension",
			transcript: "your account has been suspended please confirm your details",
			minScore: 0.4,
			maxScore: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := processor.calculatePhishingScore(tt.transcript)
			
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Expected phishing score between %.2f and %.2f, got %.2f", 
					tt.minScore, tt.maxScore, score)
			}
		})
	}
}

func TestCalculateRobocallScore(t *testing.T) {
	cfg := &config.VoiceConfig{Enabled: true}
	processor := NewProcessor(cfg)

	tests := []struct {
		name     string
		call     VoiceCall
		minScore float64
		maxScore float64
	}{
		{
			name: "Normal length call",
			call: VoiceCall{
				Duration:   5 * time.Minute,
				Transcript: "Hello, this is a normal conversation",
			},
			minScore: 0.0,
			maxScore: 0.3,
		},
		{
			name: "Short robocall",
			call: VoiceCall{
				Duration:   15 * time.Second,
				Transcript: "Press 1 to speak with a representative",
			},
			minScore: 0.6,
			maxScore: 1.0,
		},
		{
			name: "Automated message",
			call: VoiceCall{
				Duration:   45 * time.Second,
				Transcript: "This is an automated message. Press 9 to be removed.",
			},
			minScore: 0.3,
			maxScore: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := processor.calculateRobocallScore(tt.call)
			
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Expected robocall score between %.2f and %.2f, got %.2f", 
					tt.minScore, tt.maxScore, score)
			}
		})
	}
}

func TestCalculateSocialEngScore(t *testing.T) {
	cfg := &config.VoiceConfig{Enabled: true}
	processor := NewProcessor(cfg)

	tests := []struct {
		name     string
		transcript string
		minScore float64
		maxScore float64
	}{
		{
			name:     "Normal conversation",
			transcript: "hello how can I help you today",
			minScore: 0.0,
			maxScore: 0.2,
		},
		{
			name:     "Authority claim",
			transcript: "this is the irs calling about your tax situation",
			minScore: 0.1,
			maxScore: 0.5,
		},
		{
			name:     "Fear tactics",
			transcript: "police will arrest you if you don't pay this fine immediately",
			minScore: 0.4,
			maxScore: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := processor.calculateSocialEngScore(tt.transcript)
			
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Expected social engineering score between %.2f and %.2f, got %.2f", 
					tt.minScore, tt.maxScore, score)
			}
		})
	}
}

func TestAnalyzeSpeechPatterns(t *testing.T) {
	cfg := &config.VoiceConfig{Enabled: true}
	processor := NewProcessor(cfg)

	tests := []struct {
		name       string
		transcript string
		expectPattern bool
		patternType   string
	}{
		{
			name:       "Normal speech",
			transcript: "Hello how are you doing today",
			expectPattern: false,
		},
		{
			name:       "Script reading",
			transcript: "Hello. This is a formal message. Please listen carefully. We have important information. Thank you.",
			expectPattern: true,
			patternType:   "formal_script_reading",
		},
		{
			name:       "High word density",
			transcript: strings.Repeat("word ", 250),
			expectPattern: true,
			patternType:   "high_word_density",
		},
		{
			name:       "Repetitive words",
			transcript: "urgent urgent urgent urgent this is urgent please urgent",
			expectPattern: true,
			patternType:   "repetitive_word_urgent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns := processor.analyzeSpeechPatterns(tt.transcript)
			
			if tt.expectPattern {
				if len(patterns) == 0 {
					t.Errorf("Expected to find speech patterns, got none")
				}
				if tt.patternType != "" {
					found := false
					for _, pattern := range patterns {
						if strings.Contains(pattern, tt.patternType) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected to find pattern type %s, got %v", tt.patternType, patterns)
					}
				}
			} else {
				if len(patterns) > 0 {
					t.Errorf("Expected no patterns for normal speech, got %v", patterns)
				}
			}
		})
	}
}

func TestStartStop(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 1,
	}

	processor := NewProcessor(cfg)
	ctx := context.Background()

	// Test start
	err := processor.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start processor: %v", err)
	}

	// Test double start (should fail)
	err = processor.Start(ctx)
	if err == nil {
		t.Error("Expected error when starting already running processor")
	}

	// Test stop
	err = processor.Stop()
	if err != nil {
		t.Fatalf("Failed to stop processor: %v", err)
	}

	// Test double stop (should fail)
	err = processor.Stop()
	if err == nil {
		t.Error("Expected error when stopping already stopped processor")
	}
}

func TestGetStats(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:              true,
		SpeechToTextEnabled:  true,
		CallRecordingEnabled: true,
	}

	processor := NewProcessor(cfg)

	// Process some test calls
	testCalls := []VoiceCall{
		{
			ID:         "1",
			CallerID:   "+1111111111",
			Duration:   2 * time.Minute,
			Transcript: "Normal call",
		},
		{
			ID:         "2",
			CallerID:   "+2222222222",
			Duration:   30 * time.Second,
			Transcript: "This is the IRS calling about urgent tax matter",
		},
	}

	for _, call := range testCalls {
		_, err := processor.ProcessCall(call)
		if err != nil {
			t.Fatalf("Failed to process call: %v", err)
		}
	}

	stats := processor.GetStats()

	expectedKeys := []string{
		"total_calls", "unique_callers", "threat_counts", 
		"average_call_duration", "total_call_duration",
		"is_running", "speech_to_text_enabled", "call_recording_enabled",
	}

	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected stats key %s to exist", key)
		}
	}

	if stats["total_calls"].(int) != len(testCalls) {
		t.Errorf("Expected total_calls to be %d, got %v", 
			len(testCalls), stats["total_calls"])
	}
}

func TestThreatLevelString(t *testing.T) {
	tests := []struct {
		level    ThreatLevel
		expected string
	}{
		{ThreatLevelNone, "none"},
		{ThreatLevelLow, "low"},
		{ThreatLevelMedium, "medium"},
		{ThreatLevelHigh, "high"},
		{ThreatLevelCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.level.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.level.String())
			}
		})
	}
}

func TestUpdateCallerInfo(t *testing.T) {
	cfg := &config.VoiceConfig{Enabled: true}
	processor := NewProcessor(cfg)

	callerID := "+1234567890"
	duration1 := 2 * time.Minute
	duration2 := 4 * time.Minute

	// First call
	processor.updateCallerInfo(callerID, duration1)
	
	info, exists := processor.callerInfo[callerID]
	if !exists {
		t.Fatal("Expected caller info to be created")
	}

	if info.CallCount != 1 {
		t.Errorf("Expected call count to be 1, got %d", info.CallCount)
	}

	if info.AvgCallLength != duration1 {
		t.Errorf("Expected avg call length to be %v, got %v", duration1, info.AvgCallLength)
	}

	// Second call
	processor.updateCallerInfo(callerID, duration2)
	
	if info.CallCount != 2 {
		t.Errorf("Expected call count to be 2, got %d", info.CallCount)
	}

	expectedAvg := (duration1 + duration2) / 2
	if info.AvgCallLength != expectedAvg {
		t.Errorf("Expected avg call length to be %v, got %v", expectedAvg, info.AvgCallLength)
	}
}
