package sms

// Config holds SMS processor configuration
type Config struct {
	Enabled              bool   `yaml:"enabled" json:"enabled"`
	GoogleVoiceEnabled   bool   `yaml:"google_voice_enabled" json:"google_voice_enabled"`
	GoogleVoiceAPIKey    string `yaml:"google_voice_api_key" json:"google_voice_api_key"`
	WebhookEnabled       bool   `yaml:"webhook_enabled" json:"webhook_enabled"`
	WebhookPort          int    `yaml:"webhook_port" json:"webhook_port"`
	WebhookPath          string `yaml:"webhook_path" json:"webhook_path"`
	ProcessingInterval   int    `yaml:"processing_interval" json:"processing_interval"`
	MaxHistorySize       int    `yaml:"max_history_size" json:"max_history_size"`
	ThreatThreshold      float64 `yaml:"threat_threshold" json:"threat_threshold"`
	SpamFilterEnabled    bool   `yaml:"spam_filter_enabled" json:"spam_filter_enabled"`
	BlacklistEnabled     bool   `yaml:"blacklist_enabled" json:"blacklist_enabled"`
	WhitelistEnabled     bool   `yaml:"whitelist_enabled" json:"whitelist_enabled"`
}

// DefaultConfig returns default SMS processor configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:              true,
		GoogleVoiceEnabled:   false,
		WebhookEnabled:       true,
		WebhookPort:          8082,
		WebhookPath:          "/sms/webhook",
		ProcessingInterval:   30,
		MaxHistorySize:       1000,
		ThreatThreshold:      0.5,
		SpamFilterEnabled:    true,
		BlacklistEnabled:     true,
		WhitelistEnabled:     true,
	}
}
