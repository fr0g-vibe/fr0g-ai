package voice

// Config holds voice processor configuration
type Config struct {
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

// DefaultConfig returns default voice processor configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:              true,
		SpeechToTextEnabled:  false,
		CallRecordingEnabled: false,
		RecordingPath:        "/tmp/voice_recordings",
		MonitoringInterval:   60,
		MaxHistorySize:       500,
		ThreatThreshold:      0.5,
		TelephonyProvider:    "twilio",
		WebhookEnabled:       true,
		WebhookPort:          8083,
		WebhookPath:          "/voice/webhook",
	}
}
