package config

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors []ValidationError

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
