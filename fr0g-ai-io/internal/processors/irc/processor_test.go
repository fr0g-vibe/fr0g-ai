package irc

import (
	"context"
	"testing"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

func TestNewIRCProcessor(t *testing.T) {
	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{"irc.example.com:6667"},
		Channels:          []string{"#test"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 30,
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.GetType() != "irc" {
		t.Errorf("Expected processor type 'irc', got '%s'", processor.GetType())
	}

	if !processor.IsEnabled() {
		t.Error("Expected processor to be enabled")
	}
}

func TestProcessMessage(t *testing.T) {
	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{"irc.example.com:6667"},
		Channels:          []string{"#test"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 30,
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)

	// Test normal message
	normalMsg := IRCMessage{
		ID:          "test-1",
		Server:      "irc.example.com",
		Channel:     "#test",
		Nick:        "testuser",
		User:        "test",
		Host:        "example.com",
		Message:     "Hello everyone!",
		MessageType: "PRIVMSG",
		Timestamp:   time.Now(),
	}

	result, err := processor.ProcessMessage(normalMsg)
	if err != nil {
		t.Fatalf("Error processing normal message: %v", err)
	}

	if result.ThreatLevel != ThreatLevelNone {
		t.Errorf("Expected threat level None for normal message, got %s", result.ThreatLevel.String())
	}

	// Test spam message
	spamMsg := IRCMessage{
		ID:          "test-2",
		Server:      "irc.example.com",
		Channel:     "#test",
		Nick:        "spammer",
		User:        "spam",
		Host:        "spam.com",
		Message:     "FREE MONEY!!! WIN BIG PRIZES!!! CLICK HERE NOW!!!",
		MessageType: "PRIVMSG",
		Timestamp:   time.Now(),
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
	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{"irc.example.com:6667"},
		Channels:          []string{"#test"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 30,
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)

	// Test malware message
	malwareMsg := IRCMessage{
		ID:          "test-3",
		Server:      "irc.example.com",
		Channel:     "#test",
		Nick:        "hacker",
		User:        "hack",
		Host:        "malware.com",
		Message:     "Download this crack.exe file from bit.ly/hack123 for free warez!",
		MessageType: "PRIVMSG",
		Timestamp:   time.Now(),
	}

	result, err := processor.ProcessMessage(malwareMsg)
	if err != nil {
		t.Fatalf("Error processing malware message: %v", err)
	}

	if result.Analysis.MalwareScore == 0 {
		t.Error("Expected malware score > 0 for malware message")
	}

	if len(result.Analysis.ThreatTypes) == 0 {
		t.Error("Expected threat types to be detected")
	}

	if len(result.Analysis.Recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
}

func TestBotDetection(t *testing.T) {
	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{"irc.example.com:6667"},
		Channels:          []string{"#test"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 30,
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)

	// Test bot message
	botMsg := IRCMessage{
		ID:          "test-4",
		Server:      "irc.example.com",
		Channel:     "#test",
		Nick:        "bot_service123",
		User:        "bot",
		Host:        "automated.com",
		Message:     "[AUTO] System status: operational",
		MessageType: "PRIVMSG",
		Timestamp:   time.Now(),
	}

	result, err := processor.ProcessMessage(botMsg)
	if err != nil {
		t.Fatalf("Error processing bot message: %v", err)
	}

	if result.Analysis.BotScore == 0 {
		t.Error("Expected bot score > 0 for bot message")
	}
}

func TestFloodDetection(t *testing.T) {
	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{"irc.example.com:6667"},
		Channels:          []string{"#test"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 30,
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)

	// Test flood message
	floodMsg := IRCMessage{
		ID:          "test-5",
		Server:      "irc.example.com",
		Channel:     "#test",
		Nick:        "flooder",
		User:        "flood",
		Host:        "flood.com",
		Message:     "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		MessageType: "PRIVMSG",
		Timestamp:   time.Now(),
	}

	result, err := processor.ProcessMessage(floodMsg)
	if err != nil {
		t.Fatalf("Error processing flood message: %v", err)
	}

	if result.Analysis.FloodScore == 0 {
		t.Error("Expected flood score > 0 for flood message")
	}
}

func TestIRCProcessorLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping IRC processor lifecycle test in short mode")
	}

	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{}, // Empty servers for testing
		Channels:          []string{"#test"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 1, // Short interval for testing
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func TestIRCGetStats(t *testing.T) {
	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{"irc.example.com:6667"},
		Channels:          []string{"#test", "#security"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 30,
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)

	// Process some messages
	messages := []IRCMessage{
		{
			ID:          "test-1",
			Server:      "irc.example.com",
			Channel:     "#test",
			Nick:        "user1",
			Message:     "Normal message",
			MessageType: "PRIVMSG",
			Timestamp:   time.Now(),
		},
		{
			ID:          "test-2",
			Server:      "irc.example.com",
			Channel:     "#test",
			Nick:        "spammer",
			Message:     "FREE MONEY!!! CLICK HERE!!!",
			MessageType: "PRIVMSG",
			Timestamp:   time.Now(),
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

	if stats["unique_users"].(int) != 2 {
		t.Errorf("Expected 2 unique users, got %v", stats["unique_users"])
	}

	if stats["is_running"].(bool) != false {
		t.Error("Expected processor to not be running")
	}

	if stats["total_servers"].(int) != 1 {
		t.Errorf("Expected 1 total server, got %v", stats["total_servers"])
	}

	if stats["monitored_channels"].(int) != 2 {
		t.Errorf("Expected 2 monitored channels, got %v", stats["monitored_channels"])
	}
}

func TestIRCThreatLevelString(t *testing.T) {
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

func TestIRCMessageParsing(t *testing.T) {
	cfg := &sharedconfig.IRCConfig{
		Enabled:           true,
		Servers:           []string{"irc.example.com:6667"},
		Channels:          []string{"#test"},
		Nickname:          "fr0g-ai",
		Username:          "fr0g-ai",
		Realname:          "fr0g.ai Security Bot",
		ReconnectInterval: 30,
		MaxHistorySize:    1000,
	}

	processor := NewProcessor(cfg)

	// Test IRC message parsing
	testLine := ":nick!user@host.com PRIVMSG #channel :Hello world!"
	msg := processor.parseIRCMessage("irc.example.com", testLine)

	if msg == nil {
		t.Fatal("Expected parsed message, got nil")
	}

	if msg.Nick != "nick" {
		t.Errorf("Expected nick 'nick', got '%s'", msg.Nick)
	}

	if msg.User != "user" {
		t.Errorf("Expected user 'user', got '%s'", msg.User)
	}

	if msg.Host != "host.com" {
		t.Errorf("Expected host 'host.com', got '%s'", msg.Host)
	}

	if msg.Channel != "#channel" {
		t.Errorf("Expected channel '#channel', got '%s'", msg.Channel)
	}

	if msg.Message != "Hello world!" {
		t.Errorf("Expected message 'Hello world!', got '%s'", msg.Message)
	}

	if msg.MessageType != "PRIVMSG" {
		t.Errorf("Expected message type 'PRIVMSG', got '%s'", msg.MessageType)
	}
}
