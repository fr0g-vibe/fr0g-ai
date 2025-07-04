package config

import (
	sharedconfig "pkg/config"
)

// ValidationError represents a validation error
type ValidationError = sharedconfig.ValidationError

// ValidationErrors represents a collection of validation errors
type ValidationErrors = sharedconfig.ValidationErrors

// SMSConfig holds SMS processor configuration
type SMSConfig struct {
	Enabled              bool    `yaml:"enabled" json:"enabled"`
	GoogleVoiceEnabled   bool    `yaml:"google_voice_enabled" json:"google_voice_enabled"`
	GoogleVoiceAPIKey    string  `yaml:"google_voice_api_key" json:"google_voice_api_key"`
	WebhookEnabled       bool    `yaml:"webhook_enabled" json:"webhook_enabled"`
	WebhookPort          int     `yaml:"webhook_port" json:"webhook_port"`
	WebhookPath          string  `yaml:"webhook_path" json:"webhook_path"`
	ProcessingInterval   int     `yaml:"processing_interval" json:"processing_interval"`
	MaxHistorySize       int     `yaml:"max_history_size" json:"max_history_size"`
	ThreatThreshold      float64 `yaml:"threat_threshold" json:"threat_threshold"`
	SpamFilterEnabled    bool    `yaml:"spam_filter_enabled" json:"spam_filter_enabled"`
	BlacklistEnabled     bool    `yaml:"blacklist_enabled" json:"blacklist_enabled"`
	WhitelistEnabled     bool    `yaml:"whitelist_enabled" json:"whitelist_enabled"`
}

// VoiceConfig holds voice processor configuration
type VoiceConfig struct {
	Enabled              bool    `yaml:"enabled" json:"enabled"`
	SpeechToTextEnabled  bool    `yaml:"speech_to_text_enabled" json:"speech_to_text_enabled"`
	SpeechToTextAPIKey   string  `yaml:"speech_to_text_api_key" json:"speech_to_text_api_key"`
	CallRecordingEnabled bool    `yaml:"call_recording_enabled" json:"call_recording_enabled"`
	RecordingPath        string  `yaml:"recording_path" json:"recording_path"`
	MonitoringInterval   int     `yaml:"monitoring_interval" json:"monitoring_interval"`
	MaxHistorySize       int     `yaml:"max_history_size" json:"max_history_size"`
	ThreatThreshold      float64 `yaml:"threat_threshold" json:"threat_threshold"`
	TelephonyProvider    string  `yaml:"telephony_provider" json:"telephony_provider"`
	TelephonyAPIKey      string  `yaml:"telephony_api_key" json:"telephony_api_key"`
	WebhookEnabled       bool    `yaml:"webhook_enabled" json:"webhook_enabled"`
	WebhookPort          int     `yaml:"webhook_port" json:"webhook_port"`
	WebhookPath          string  `yaml:"webhook_path" json:"webhook_path"`
}

// IRCConfig holds IRC processor configuration
type IRCConfig struct {
	Enabled              bool     `yaml:"enabled" json:"enabled"`
	Servers              []string `yaml:"servers" json:"servers"`
	Channels             []string `yaml:"channels" json:"channels"`
	Nickname             string   `yaml:"nickname" json:"nickname"`
	Username             string   `yaml:"username" json:"username"`
	Realname             string   `yaml:"realname" json:"realname"`
	Password             string   `yaml:"password" json:"password"`
	TLS                  bool     `yaml:"tls" json:"tls"`
	TLSInsecure          bool     `yaml:"tls_insecure" json:"tls_insecure"`
	ReconnectInterval    int      `yaml:"reconnect_interval" json:"reconnect_interval"`
	MaxHistorySize       int      `yaml:"max_history_size" json:"max_history_size"`
	ThreatThreshold      float64  `yaml:"threat_threshold" json:"threat_threshold"`
	MonitorPrivateMsg    bool     `yaml:"monitor_private_msg" json:"monitor_private_msg"`
	MonitorChannelMsg    bool     `yaml:"monitor_channel_msg" json:"monitor_channel_msg"`
	LogChannelActivity   bool     `yaml:"log_channel_activity" json:"log_channel_activity"`
	BotDetectionEnabled  bool     `yaml:"bot_detection_enabled" json:"bot_detection_enabled"`
	FloodProtection      bool     `yaml:"flood_protection" json:"flood_protection"`
	MaxMessagesPerMinute int      `yaml:"max_messages_per_minute" json:"max_messages_per_minute"`
}
