package email

import (
	"context"
	"testing"
	"time"
)

func TestNewEmailProcessor(t *testing.T) {
	// Test with nil config (should use defaults)
	processor := NewEmailProcessor(nil)
	if processor == nil {
		t.Fatal("Expected processor to be created with default config")
	}

	if processor.config.SMTPPort != 25 {
		t.Errorf("Expected default SMTP port 25, got %d", processor.config.SMTPPort)
	}

	if processor.config.SMTPSPort != 465 {
		t.Errorf("Expected default SMTPS port 465, got %d", processor.config.SMTPSPort)
	}

	// Test with custom config
	config := &EmailConfig{
		SMTPPort:          2525,
		SMTPSPort:         4465,
		Hostname:          "test.example.com",
		MaxMessageSize:    5 * 1024 * 1024,
		MaxConnections:    50,
		ConnectionTimeout: 15 * time.Second,
		EnableTLS:         false,
		ThreatThreshold:   0.8,
	}

	processor = NewEmailProcessor(config)
	if processor.config.SMTPPort != 2525 {
		t.Errorf("Expected custom SMTP port 2525, got %d", processor.config.SMTPPort)
	}
}

func TestEmailConfigValidation(t *testing.T) {
	config := &EmailConfig{
		SMTPPort:          25,
		SMTPSPort:         465,
		Hostname:          "localhost",
		MaxMessageSize:    1024,
		MaxConnections:    10,
		ConnectionTimeout: 30 * time.Second,
		ThreatThreshold:   0.7,
	}

	errors := config.Validate()
	if len(errors) != 0 {
		t.Errorf("Expected no validation errors for valid config, got %d errors", len(errors))
	}

	// Test invalid config
	invalidConfig := &EmailConfig{
		SMTPPort:          -1,
		SMTPSPort:         70000,
		Hostname:          "",
		MaxMessageSize:    -1,
		MaxConnections:    -1,
		ConnectionTimeout: -1,
		ThreatThreshold:   2.0,
	}

	errors = invalidConfig.Validate()
	if len(errors) == 0 {
		t.Error("Expected validation errors for invalid config")
	}
}

func TestThreatLevelCalculation(t *testing.T) {
	processor := NewEmailProcessor(nil)

	// Test high threat
	analysis := &EmailThreatAnalysis{
		SpamScore:     0.9,
		PhishingScore: 0.8,
		MalwareScore:  0.7,
	}

	level := processor.calculateThreatLevel(analysis)
	if level != ThreatLevelCritical {
		t.Errorf("Expected critical threat level, got %s", level)
	}

	// Test medium threat
	analysis = &EmailThreatAnalysis{
		SpamScore:     0.5,
		PhishingScore: 0.4,
		MalwareScore:  0.3,
	}

	level = processor.calculateThreatLevel(analysis)
	if level != ThreatLevelMedium {
		t.Errorf("Expected medium threat level, got %s", level)
	}

	// Test low threat
	analysis = &EmailThreatAnalysis{
		SpamScore:     0.1,
		PhishingScore: 0.1,
		MalwareScore:  0.1,
	}

	level = processor.calculateThreatLevel(analysis)
	if level != ThreatLevelLow {
		t.Errorf("Expected low threat level, got %s", level)
	}
}

func TestProcessorLifecycle(t *testing.T) {
	config := &EmailConfig{
		SMTPPort:          2525, // Use non-standard port for testing
		SMTPSPort:         4465,
		Hostname:          "localhost",
		MaxMessageSize:    1024 * 1024,
		MaxConnections:    10,
		ConnectionTimeout: 5 * time.Second,
		EnableTLS:         false, // Disable TLS for testing
	}

	processor := NewEmailProcessor(config)

	if processor.IsRunning() {
		t.Error("Processor should not be running initially")
	}

	// Note: We can't easily test Start() without proper network setup
	// This would require integration tests with actual network listeners
	
	ctx := context.Background()
	
	// Test that we can create the processor without errors
	if processor.GetType() != "email" {
		t.Errorf("Expected processor type 'email', got '%s'", processor.GetType())
	}

	// Test stop when not running
	err := processor.Stop(ctx)
	if err != nil {
		t.Errorf("Stop should not error when processor is not running: %v", err)
	}
}
