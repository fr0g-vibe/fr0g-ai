package irc

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"fr0g-ai-master-control/internal/config"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled:              true,
		Servers:              []string{"irc.example.com:6667"},
		Channels:             []string{"#test"},
		Nickname:             "testbot",
		Username:             "testuser",
		Realname:             "Test Bot",
		ReconnectInterval:    30,
		MaxHistorySize:       1000,
		ThreatThreshold:      0.5,
		MonitorPrivateMsg:    true,
		MonitorChannelMsg:    true,
		BotDetectionEnabled:  true,
		FloodProtection:      true,
		MaxMessagesPerMinute: 10,
	}

	processor := NewProcessor(cfg)

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.config != cfg {
		t.Error("Expected config to be set correctly")
	}

	if len(processor.threatPatterns) == 0 {
		t.Error("Expected threat patterns to be initialized")
	}

	if len(processor.suspiciousWords) == 0 {
		t.Error("Expected suspicious words to be initialized")
	}
}

func TestProcessMessage(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled:         true,
		MaxHistorySize:  100,
		ThreatThreshold: 0.5,
	}

	processor := NewProcessor(cfg)

	tests := []struct {
		name            string
		message         IRCMessage
		expectedThreat  ThreatLevel
		minConfidence   float64
	}{
		{
			name: "Normal message",
			message: IRCMessage{
				ID:          "test1",
				Server:      "irc.example.com",
				Channel:     "#test",
				Nick:        "user1",
				User:        "user1",
				Host:        "example.com",
				Message:     "Hello everyone!",
				MessageType: "PRIVMSG",
				Timestamp:   time.Now(),
			},
			expectedThreat: ThreatLevelNone,
			minConfidence:  0.0,
		},
		{
			name: "Spam message",
			message: IRCMessage{
				ID:          "test2",
				Server:      "irc.example.com",
				Channel:     "#test",
				Nick:        "spammer",
				User:        "spam",
				Host:        "spam.com",
				Message:     "FREE MONEY!!! WIN BIG PRIZES!!! CLICK HERE NOW!!!",
				MessageType: "PRIVMSG",
				Timestamp:   time.Now(),
			},
			expectedThreat: ThreatLevelMedium,
			minConfidence:  0.3,
		},
		{
			name: "Phishing message",
			message: IRCMessage{
				ID:          "test3",
				Server:      "irc.example.com",
				Channel:     "#test",
				Nick:        "phisher",
				User:        "bad",
				Host:        "evil.com",
				Message:     "URGENT: Your account has been suspended! Click here to verify: http://fake-bank.tk/login",
				MessageType: "PRIVMSG",
				Timestamp:   time.Now(),
			},
			expectedThreat: ThreatLevelHigh,
			minConfidence:  0.5,
		},
		{
			name: "Bot-like message",
			message: IRCMessage{
				ID:          "test4",
				Server:      "irc.example.com",
				Channel:     "#test",
				Nick:        "bot_user123",
				User:        "bot",
				Host:        "automated.com",
				Message:     "[AUTO] System status: OK",
				MessageType: "PRIVMSG",
				Timestamp:   time.Now(),
			},
			expectedThreat: ThreatLevelLow,
			minConfidence:  0.2,
		},
		{
			name: "Malware message",
			message: IRCMessage{
				ID:          "test5",
				Server:      "irc.example.com",
				Channel:     "#test",
				Nick:        "hacker",
				User:        "evil",
				Host:        "malware.com",
				Message:     "Download this crack.exe for free software! Install now!",
				MessageType: "PRIVMSG",
				Timestamp:   time.Now(),
			},
			expectedThreat: ThreatLevelHigh,
			minConfidence:  0.4,
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

			if len(result.Analysis.Recommendations) == 0 && result.ThreatLevel > ThreatLevelNone {
				t.Error("Expected recommendations for threat message")
			}
		})
	}
}

func TestThreatPatterns(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled: true,
	}

	processor := NewProcessor(cfg)

	tests := []struct {
		name     string
		message  string
		patterns []string
	}{
		{
			name:     "Malware URL",
			message:  "Check out this link: http://bit.ly/malware123",
			patterns: []string{"malware_url"},
		},
		{
			name:     "Phishing URL",
			message:  "Verify your account at http://login-bank.tk/secure",
			patterns: []string{"phishing_url"},
		},
		{
			name:     "Bot pattern",
			message:  "Message from bot_service_auto",
			patterns: []string{"bot_pattern"},
		},
		{
			name:     "Flood pattern",
			message:  "AAAAAAAAAAAAAAAAAAA",
			patterns: []string{"flood_pattern"},
		},
		{
			name:     "DCC exploit",
			message:  "DCC SEND malware.exe 192.168.1.1 1234 12345",
			patterns: []string{"dcc_exploit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := IRCMessage{
				Message: tt.message,
			}

			analysis := processor.analyzeThreat(msg)

			for _, expectedPattern := range tt.patterns {
				found := false
				for _, threatType := range analysis.ThreatTypes {
					if threatType == expectedPattern {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected pattern %s to be detected in message: %s", 
						expectedPattern, tt.message)
				}
			}
		})
	}
}

func TestUserTracking(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled:        true,
		MaxHistorySize: 100,
	}

	processor := NewProcessor(cfg)

	// Send multiple messages from same user
	for i := 0; i < 5; i++ {
		msg := IRCMessage{
			ID:          fmt.Sprintf("test%d", i),
			Nick:        "testuser",
			User:        "test",
			Host:        "example.com",
			Message:     fmt.Sprintf("Message %d", i),
			MessageType: "PRIVMSG",
			Timestamp:   time.Now(),
		}

		_, err := processor.ProcessMessage(msg)
		if err != nil {
			t.Fatalf("ProcessMessage failed: %v", err)
		}
	}

	// Check user tracking
	key := "testuser!test@example.com"
	userInfo, exists := processor.userHistory[key]
	if !exists {
		t.Fatal("Expected user to be tracked")
	}

	if userInfo.MessageCount != 5 {
		t.Errorf("Expected message count 5, got %d", userInfo.MessageCount)
	}

	if userInfo.Nick != "testuser" {
		t.Errorf("Expected nick 'testuser', got '%s'", userInfo.Nick)
	}
}

func TestSpamScoreCalculation(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled: true,
	}

	processor := NewProcessor(cfg)

	tests := []struct {
		name          string
		message       string
		minSpamScore  float64
	}{
		{
			name:         "Normal message",
			message:      "Hello, how are you today?",
			minSpamScore: 0.0,
		},
		{
			name:         "Spam keywords",
			message:      "free money win prize cash",
			minSpamScore: 0.5,
		},
		{
			name:         "Excessive punctuation",
			message:      "Buy now!!!! Limited time!!!!",
			minSpamScore: 0.3,
		},
		{
			name:         "All caps",
			message:      "URGENT OFFER FOR YOU",
			minSpamScore: 0.4,
		},
		{
			name:         "Combined spam indicators",
			message:      "FREE MONEY!!! WIN BIG PRIZES!!! WORK FROM HOME!!!",
			minSpamScore: 0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := processor.calculateSpamScore(strings.ToLower(tt.message))
			if score < tt.minSpamScore {
				t.Errorf("Expected spam score >= %.2f, got %.2f for message: %s", 
					tt.minSpamScore, score, tt.message)
			}
		})
	}
}

func TestBotDetection(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled:             true,
		BotDetectionEnabled: true,
	}

	processor := NewProcessor(cfg)

	tests := []struct {
		name       string
		nick       string
		message    string
		minBotScore float64
	}{
		{
			name:        "Normal user",
			nick:        "john_doe",
			message:     "Hello everyone!",
			minBotScore: 0.0,
		},
		{
			name:        "Bot-like nick",
			nick:        "bot_service",
			message:     "Status update",
			minBotScore: 0.5,
		},
		{
			name:        "Automated message format",
			nick:        "normaluser",
			message:     "[AUTO] System notification",
			minBotScore: 0.2,
		},
		{
			name:        "Bot nick with automated message",
			nick:        "auto_bot123",
			message:     "[STATUS] All systems operational",
			minBotScore: 0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := IRCMessage{
				Nick:    tt.nick,
				Message: tt.message,
			}

			score := processor.calculateBotScore(msg)
			if score < tt.minBotScore {
				t.Errorf("Expected bot score >= %.2f, got %.2f for nick: %s, message: %s", 
					tt.minBotScore, score, tt.nick, tt.message)
			}
		})
	}
}

func TestGetStats(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled:    true,
		Servers:    []string{"irc1.example.com", "irc2.example.com"},
		Channels:   []string{"#test1", "#test2"},
	}

	processor := NewProcessor(cfg)

	// Add some test data
	processor.messageHistory = []IRCMessage{
		{ThreatLevel: ThreatLevelNone},
		{ThreatLevel: ThreatLevelLow},
		{ThreatLevel: ThreatLevelHigh},
	}

	processor.userHistory = map[string]*UserInfo{
		"user1!test@example.com": {},
		"user2!test@example.com": {},
	}

	stats := processor.GetStats()

	if stats["total_messages"] != 3 {
		t.Errorf("Expected total_messages 3, got %v", stats["total_messages"])
	}

	if stats["unique_users"] != 2 {
		t.Errorf("Expected unique_users 2, got %v", stats["unique_users"])
	}

	if stats["total_servers"] != 2 {
		t.Errorf("Expected total_servers 2, got %v", stats["total_servers"])
	}

	if stats["monitored_channels"] != 2 {
		t.Errorf("Expected monitored_channels 2, got %v", stats["monitored_channels"])
	}

	threatCounts, ok := stats["threat_counts"].(map[string]int)
	if !ok {
		t.Fatal("Expected threat_counts to be map[string]int")
	}

	if threatCounts["none"] != 1 {
		t.Errorf("Expected 1 'none' threat, got %d", threatCounts["none"])
	}

	if threatCounts["low"] != 1 {
		t.Errorf("Expected 1 'low' threat, got %d", threatCounts["low"])
	}

	if threatCounts["high"] != 1 {
		t.Errorf("Expected 1 'high' threat, got %d", threatCounts["high"])
	}
}

func TestStartStop(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled:           true,
		Servers:           []string{}, // Empty to avoid actual connections
		ReconnectInterval: 1,
	}

	processor := NewProcessor(cfg)

	// Test starting
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := processor.Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if !processor.isRunning {
		t.Error("Expected processor to be running")
	}

	// Test starting again (should fail)
	err = processor.Start(ctx)
	if err == nil {
		t.Error("Expected error when starting already running processor")
	}

	// Test stopping
	err = processor.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if processor.isRunning {
		t.Error("Expected processor to be stopped")
	}

	// Test stopping again (should fail)
	err = processor.Stop()
	if err == nil {
		t.Error("Expected error when stopping already stopped processor")
	}
}

func TestParseIRCMessage(t *testing.T) {
	cfg := &config.IRCConfig{
		Enabled: true,
	}

	processor := NewProcessor(cfg)

	tests := []struct {
		name     string
		line     string
		expected *IRCMessage
	}{
		{
			name: "PRIVMSG",
			line: ":nick!user@host.com PRIVMSG #channel :Hello world",
			expected: &IRCMessage{
				Server:      "irc.example.com",
				Channel:     "#channel",
				Nick:        "nick",
				User:        "user",
				Host:        "host.com",
				Message:     "Hello world",
				MessageType: "PRIVMSG",
			},
		},
		{
			name: "NOTICE",
			line: ":nick!user@host.com NOTICE target :This is a notice",
			expected: &IRCMessage{
				Server:      "irc.example.com",
				Channel:     "target",
				Nick:        "nick",
				User:        "user",
				Host:        "host.com",
				Message:     "This is a notice",
				MessageType: "NOTICE",
			},
		},
		{
			name:     "Invalid message",
			line:     "INVALID LINE",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.parseIRCMessage("irc.example.com", tt.line)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil result, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected result, got nil")
			}

			if result.Server != tt.expected.Server {
				t.Errorf("Expected server %s, got %s", tt.expected.Server, result.Server)
			}

			if result.Channel != tt.expected.Channel {
				t.Errorf("Expected channel %s, got %s", tt.expected.Channel, result.Channel)
			}

			if result.Nick != tt.expected.Nick {
				t.Errorf("Expected nick %s, got %s", tt.expected.Nick, result.Nick)
			}

			if result.User != tt.expected.User {
				t.Errorf("Expected user %s, got %s", tt.expected.User, result.User)
			}

			if result.Host != tt.expected.Host {
				t.Errorf("Expected host %s, got %s", tt.expected.Host, result.Host)
			}

			if result.Message != tt.expected.Message {
				t.Errorf("Expected message %s, got %s", tt.expected.Message, result.Message)
			}

			if result.MessageType != tt.expected.MessageType {
				t.Errorf("Expected message type %s, got %s", tt.expected.MessageType, result.MessageType)
			}
		})
	}
}
