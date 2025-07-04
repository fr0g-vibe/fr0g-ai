package sms

import (
	"context"
	"testing"
	"time"

	"../../../config"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.SMSConfig{
		Enabled:            true,
		ProcessingInterval: 30,
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

	if len(processor.spamKeywords) == 0 {
		t.Error("Expected spam keywords to be initialized")
	}
}

func TestProcessMessage(t *testing.T) {
	cfg := &config.SMSConfig{
		Enabled:         true,
		MaxHistorySize:  100,
		ThreatThreshold: 0.5,
	}

	processor := NewProcessor(cfg)

	tests := []struct {
		name            string
		message         SMSMessage
		expectedThreat  ThreatLevel
		minConfidence   float64
	}{
		{
			name: "Normal message",
			message: SMSMessage{
				ID:        "test1",
				From:      "+1234567890",
				To:        "+0987654321",
				Body:      "Hello, how are you?",
				Timestamp: time.Now(),
			},
			expectedThreat: ThreatLevelNone,
			minConfidence:  0.0,
		},
		{
			name: "Spam message",
			message: SMSMessage{
				ID:        "test2",
				From:      "+1111111111",
				To:        "+0987654321",
				Body:      "CONGRATULATIONS! You've won a FREE prize! Click here to claim now! Limited time offer!",
				Timestamp: time.Now(),
			},
			expectedThreat: ThreatLevelMedium,
			minConfidence:  0.3,
		},
		{
			name: "Phishing message",
			message: SMSMessage{
				ID:        "test3",
				From:      "+2222222222",
				To:        "+0987654321",
				Body:      "URGENT: Your account has been suspended. Verify your identity immediately: bit.ly/verify123",
				Timestamp: time.Now(),
			},
			expectedThreat: ThreatLevelHigh,
			minConfidence:  0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.ProcessMessage(tt.message)
			if err != nil {
				t.Fatalf("ProcessMessage failed: %v", err)
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

func TestCalculateSpamScore(t *testing.T) {
	cfg := &config.SMSConfig{Enabled: true}
	processor := NewProcessor(cfg)

	tests := []struct {
		name     string
		body     string
		minScore float64
		maxScore float64
	}{
		{
			name:     "Normal text",
			body:     "hello how are you today",
			minScore: 0.0,
			maxScore: 0.2,
		},
		{
			name:     "Spam keywords",
			body:     "free prize winner congratulations urgent",
			minScore: 0.4,
			maxScore: 1.0,
		},
		{
			name:     "Excessive caps and punctuation",
			body:     "FREE MONEY NOW!!! CLICK HERE!!!",
			minScore: 0.5,
			maxScore: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := processor.calculateSpamScore(tt.body)
			
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Expected spam score between %.2f and %.2f, got %.2f", 
					tt.minScore, tt.maxScore, score)
			}
		})
	}
}

func TestStartStop(t *testing.T) {
	cfg := &config.SMSConfig{
		Enabled:            true,
		ProcessingInterval: 1,
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
	cfg := &config.SMSConfig{
		Enabled:            true,
		GoogleVoiceEnabled: true,
		WebhookEnabled:     true,
	}

	processor := NewProcessor(cfg)

	// Process some test messages
	testMessages := []SMSMessage{
		{ID: "1", From: "+1111111111", Body: "Normal message", Timestamp: time.Now()},
		{ID: "2", From: "+2222222222", Body: "FREE PRIZE!!!", Timestamp: time.Now()},
	}

	for _, msg := range testMessages {
		_, err := processor.ProcessMessage(msg)
		if err != nil {
			t.Fatalf("Failed to process message: %v", err)
		}
	}

	stats := processor.GetStats()

	expectedKeys := []string{
		"total_messages", "unique_numbers", "threat_counts", 
		"is_running", "google_voice_enabled", "webhook_enabled",
	}

	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected stats key %s to exist", key)
		}
	}

	if stats["total_messages"].(int) != len(testMessages) {
		t.Errorf("Expected total_messages to be %d, got %v", 
			len(testMessages), stats["total_messages"])
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
