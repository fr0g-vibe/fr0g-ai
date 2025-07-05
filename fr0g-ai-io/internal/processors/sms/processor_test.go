package sms

import (
	"context"
	"testing"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/config"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.SMSConfig{
		Enabled:            true,
		ProcessingInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.GetType() != "sms" {
		t.Errorf("Expected processor type 'sms', got '%s'", processor.GetType())
	}

	if !processor.IsEnabled() {
		t.Error("Expected processor to be enabled")
	}
}

func TestProcessMessage(t *testing.T) {
	cfg := &config.SMSConfig{
		Enabled:            true,
		ProcessingInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	// Test normal message
	normalMsg := SMSMessage{
		ID:        "test-1",
		From:      "+1234567890",
		To:        "+0987654321",
		Body:      "Hello, how are you?",
		Timestamp: time.Now(),
	}

	result, err := processor.ProcessMessage(normalMsg)
	if err != nil {
		t.Fatalf("Error processing normal message: %v", err)
	}

	if result.ThreatLevel != ThreatLevelNone {
		t.Errorf("Expected threat level None for normal message, got %s", result.ThreatLevel.String())
	}

	// Test spam message
	spamMsg := SMSMessage{
		ID:        "test-2",
		From:      "+1234567890",
		To:        "+0987654321",
		Body:      "URGENT! You've won a FREE prize! Click here NOW to claim your lottery winnings!",
		Timestamp: time.Now(),
	}

	result, err = processor.ProcessMessage(spamMsg)
	if err != nil {
		t.Fatalf("Error processing spam message: %v", err)
	}

	if result.ThreatLevel == ThreatLevelNone {
		t.Error("Expected threat level higher than None for spam message")
	}

	if result.Analysis.SpamScore == 0 {
		t.Error("Expected spam score > 0 for spam message")
	}
}

func TestThreatAnalysis(t *testing.T) {
	cfg := &config.SMSConfig{
		Enabled:            true,
		ProcessingInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	// Test phishing message
	phishingMsg := SMSMessage{
		ID:        "test-3",
		From:      "+1234567890",
		To:        "+0987654321",
		Body:      "Your account has been suspended. Please verify your password immediately at bit.ly/verify123",
		Timestamp: time.Now(),
	}

	result, err := processor.ProcessMessage(phishingMsg)
	if err != nil {
		t.Fatalf("Error processing phishing message: %v", err)
	}

	if result.Analysis.PhishingScore == 0 {
		t.Error("Expected phishing score > 0 for phishing message")
	}

	if len(result.Analysis.ThreatTypes) == 0 {
		t.Error("Expected threat types to be detected")
	}

	if len(result.Analysis.Recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
}

func TestProcessorLifecycle(t *testing.T) {
	cfg := &config.SMSConfig{
		Enabled:            true,
		ProcessingInterval: 1, // Short interval for testing
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test start
	err := processor.Start(ctx)
	if err != nil {
		t.Fatalf("Error starting processor: %v", err)
	}

	// Test double start (should fail)
	err = processor.Start(ctx)
	if err == nil {
		t.Error("Expected error when starting already running processor")
	}

	// Test stop
	err = processor.Stop()
	if err != nil {
		t.Fatalf("Error stopping processor: %v", err)
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
		ProcessingInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	// Process some messages
	messages := []SMSMessage{
		{
			ID:        "test-1",
			From:      "+1234567890",
			Body:      "Normal message",
			Timestamp: time.Now(),
		},
		{
			ID:        "test-2",
			From:      "+1234567890",
			Body:      "URGENT! FREE MONEY! Click here NOW!",
			Timestamp: time.Now(),
		},
	}

	for _, msg := range messages {
		_, err := processor.ProcessMessage(msg)
		if err != nil {
			t.Fatalf("Error processing message: %v", err)
		}
	}

	stats := processor.GetStats()

	if stats["total_messages"].(int) != 2 {
		t.Errorf("Expected 2 total messages, got %v", stats["total_messages"])
	}

	if stats["unique_numbers"].(int) != 1 {
		t.Errorf("Expected 1 unique number, got %v", stats["unique_numbers"])
	}

	if stats["is_running"].(bool) != false {
		t.Error("Expected processor to not be running")
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

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.level.String())
		}
	}
}
