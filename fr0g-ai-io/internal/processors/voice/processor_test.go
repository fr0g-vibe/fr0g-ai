package voice

import (
	"context"
	"testing"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/config"
)

func TestNewVoiceProcessor(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.GetType() != "voice" {
		t.Errorf("Expected processor type 'voice', got '%s'", processor.GetType())
	}

	if !processor.IsEnabled() {
		t.Error("Expected processor to be enabled")
	}
}

func TestProcessCall(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	// Test normal call
	normalCall := VoiceCall{
		ID:          "test-1",
		CallerID:    "+1234567890",
		RecipientID: "+0987654321",
		StartTime:   time.Now().Add(-5 * time.Minute),
		EndTime:     time.Now(),
		Duration:    5 * time.Minute,
		Transcript:  "Hello, I'm calling to check on your account status.",
	}

	result, err := processor.ProcessCall(normalCall)
	if err != nil {
		t.Fatalf("Error processing normal call: %v", err)
	}

	if result.ThreatLevel != ThreatLevelNone {
		t.Errorf("Expected threat level None for normal call, got %s", result.ThreatLevel.String())
	}

	// Test scam call
	scamCall := VoiceCall{
		ID:          "test-2",
		CallerID:    "+1234567890",
		RecipientID: "+0987654321",
		StartTime:   time.Now().Add(-2 * time.Minute),
		EndTime:     time.Now(),
		Duration:    2 * time.Minute,
		Transcript:  "This is the IRS. Your account has been suspended due to suspicious activity. You must verify your social security number immediately or face arrest.",
	}

	result, err = processor.ProcessCall(scamCall)
	if err != nil {
		t.Fatalf("Error processing scam call: %v", err)
	}

	if result.ThreatLevel == ThreatLevelNone {
		t.Error("Expected threat level higher than None for scam call")
	}

	if result.Analysis.ScamScore == 0 {
		t.Error("Expected scam score > 0 for scam call")
	}

	if result.Analysis.SocialEngScore == 0 {
		t.Error("Expected social engineering score > 0 for scam call")
	}
}

func TestThreatAnalysis(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	// Test robocall
	robocall := VoiceCall{
		ID:          "test-3",
		CallerID:    "+1234567890",
		RecipientID: "+0987654321",
		StartTime:   time.Now().Add(-15 * time.Second),
		EndTime:     time.Now(),
		Duration:    15 * time.Second,
		Transcript:  "This is not a sales call. Press 1 to speak with a representative about your car warranty.",
	}

	result, err := processor.ProcessCall(robocall)
	if err != nil {
		t.Fatalf("Error processing robocall: %v", err)
	}

	if result.Analysis.RobocallScore == 0 {
		t.Error("Expected robocall score > 0 for robocall")
	}

	if len(result.Analysis.ThreatTypes) == 0 {
		t.Error("Expected threat types to be detected")
	}

	if len(result.Analysis.Recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
}

func TestVoiceProcessorLifecycle(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 1, // Short interval for testing
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

func TestVoiceGetStats(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:               true,
		MonitoringInterval:    30,
		MaxHistorySize:        1000,
		ThreatThreshold:       0.5,
		SpeechToTextEnabled:   true,
		CallRecordingEnabled:  false,
	}

	processor := NewProcessor(cfg)

	// Process some calls
	calls := []VoiceCall{
		{
			ID:          "test-1",
			CallerID:    "+1234567890",
			Duration:    2 * time.Minute,
			Transcript:  "Normal business call",
		},
		{
			ID:          "test-2",
			CallerID:    "+1234567890",
			Duration:    30 * time.Second,
			Transcript:  "This is the IRS calling about your tax refund. You must act immediately!",
		},
	}

	for _, call := range calls {
		_, err := processor.ProcessCall(call)
		if err != nil {
			t.Fatalf("Error processing call: %v", err)
		}
	}

	stats := processor.GetStats()

	if stats["total_calls"].(int) != 2 {
		t.Errorf("Expected 2 total calls, got %v", stats["total_calls"])
	}

	if stats["unique_callers"].(int) != 1 {
		t.Errorf("Expected 1 unique caller, got %v", stats["unique_callers"])
	}

	if stats["is_running"].(bool) != false {
		t.Error("Expected processor to not be running")
	}

	if stats["speech_to_text_enabled"].(bool) != true {
		t.Error("Expected speech-to-text to be enabled")
	}

	if stats["call_recording_enabled"].(bool) != false {
		t.Error("Expected call recording to be disabled")
	}
}

func TestVoiceThreatLevelString(t *testing.T) {
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

func TestSpeechPatternAnalysis(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	// Test script reading detection
	scriptCall := VoiceCall{
		ID:         "test-script",
		CallerID:   "+1234567890",
		Duration:   3 * time.Minute,
		Transcript: "Hello. This is a formal message. We are calling regarding your account. Please listen carefully. Your account status requires immediate attention. Thank you for your time.",
	}

	result, err := processor.ProcessCall(scriptCall)
	if err != nil {
		t.Fatalf("Error processing script call: %v", err)
	}

	if len(result.Analysis.SpeechPatterns) == 0 {
		t.Error("Expected speech patterns to be detected")
	}
}

func TestEmotionalManipulationDetection(t *testing.T) {
	cfg := &config.VoiceConfig{
		Enabled:            true,
		MonitoringInterval: 30,
		MaxHistorySize:     1000,
		ThreatThreshold:    0.5,
	}

	processor := NewProcessor(cfg)

	// Test emotional manipulation
	manipCall := VoiceCall{
		ID:         "test-manip",
		CallerID:   "+1234567890",
		Duration:   2 * time.Minute,
		Transcript: "This is an emergency! Your family is in danger! You must act now to protect your loved ones! Time is running out!",
	}

	result, err := processor.ProcessCall(manipCall)
	if err != nil {
		t.Fatalf("Error processing manipulation call: %v", err)
	}

	if result.Analysis.EmotionalManipScore == 0 {
		t.Error("Expected emotional manipulation score > 0")
	}

	if result.ThreatLevel == ThreatLevelNone {
		t.Error("Expected threat level higher than None for emotional manipulation")
	}
}
