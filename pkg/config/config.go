package config

import (
	"fmt"
	"strings"
	"time"
)

// BaseConfig provides common configuration functionality
type BaseConfig interface {
	Validate() ValidationErrors
}

// CommonConfig holds shared configuration fields
type CommonConfig struct {
	Environment string        `yaml:"environment" json:"environment"`
	LogLevel    string        `yaml:"log_level" json:"log_level"`
	Debug       bool          `yaml:"debug" json:"debug"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout"`
}

// ServerConfig holds common server configuration
type ServerConfig struct {
	Host            string        `yaml:"host" json:"host"`
	HTTPPort        int           `yaml:"http_port" json:"http_port"`
	GRPCPort        int           `yaml:"grpc_port" json:"grpc_port"`
	EnableTLS       bool          `yaml:"enable_tls" json:"enable_tls"`
	CertFile        string        `yaml:"cert_file" json:"cert_file"`
	KeyFile         string        `yaml:"key_file" json:"key_file"`
	ReadTimeout     time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" json:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout"`
}

// SecurityConfig holds common security configuration
type SecurityConfig struct {
	EnableAuth       bool     `yaml:"enable_auth" json:"enable_auth"`
	APIKey           string   `yaml:"api_key" json:"api_key"`
	AllowedAPIKeys   []string `yaml:"allowed_api_keys" json:"allowed_api_keys"`
	EnableCORS       bool     `yaml:"enable_cors" json:"enable_cors"`
	AllowedOrigins   []string `yaml:"allowed_origins" json:"allowed_origins"`
	RateLimitRPM     int      `yaml:"rate_limit_rpm" json:"rate_limit_rpm"`
	RequireAPIKey    bool     `yaml:"require_api_key" json:"require_api_key"`
	EnableReflection bool     `yaml:"enable_reflection" json:"enable_reflection"`
}

// StorageConfig holds common storage configuration
type StorageConfig struct {
	Type    string `yaml:"type" json:"type"`
	DataDir string `yaml:"data_dir" json:"data_dir"`
}

// MonitoringConfig holds common monitoring configuration
type MonitoringConfig struct {
	EnableMetrics       bool `yaml:"enable_metrics" json:"enable_metrics"`
	MetricsPort         int  `yaml:"metrics_port" json:"metrics_port"`
	HealthCheckInterval int  `yaml:"health_check_interval" json:"health_check_interval"`
	EnableTracing       bool `yaml:"enable_tracing" json:"enable_tracing"`
}

// OpenWebUIConfig holds OpenWebUI client configuration
type OpenWebUIConfig struct {
	BaseURL string `yaml:"base_url" json:"base_url"`
	APIKey  string `yaml:"api_key" json:"api_key"`
	Timeout int    `yaml:"timeout" json:"timeout"` // timeout in seconds
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level" json:"level"`
	Format     string `yaml:"format" json:"format"`
	Output     string `yaml:"output" json:"output"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
}

// SMSConfig holds SMS processing configuration
type SMSConfig struct {
	Enabled               bool     `yaml:"enabled" json:"enabled"`
	GoogleVoiceEnabled    bool     `yaml:"google_voice_enabled" json:"google_voice_enabled"`
	WebhookEnabled        bool     `yaml:"webhook_enabled" json:"webhook_enabled"`
	WebhookPort           int      `yaml:"webhook_port" json:"webhook_port"`
	ProcessingInterval    int      `yaml:"processing_interval" json:"processing_interval"`
	MaxHistorySize        int      `yaml:"max_history_size" json:"max_history_size"`
	ThreatThreshold       float64  `yaml:"threat_threshold" json:"threat_threshold"`
	ResponseEnabled       bool     `yaml:"response_enabled" json:"response_enabled"`
	ResponseTemplates     []string `yaml:"response_templates" json:"response_templates"`
}

// VoiceConfig holds voice processing configuration
type VoiceConfig struct {
	Enabled              bool    `yaml:"enabled" json:"enabled"`
	SpeechToTextEnabled  bool    `yaml:"speech_to_text_enabled" json:"speech_to_text_enabled"`
	CallRecordingEnabled bool    `yaml:"call_recording_enabled" json:"call_recording_enabled"`
	MonitoringInterval   int     `yaml:"monitoring_interval" json:"monitoring_interval"`
	MaxHistorySize       int     `yaml:"max_history_size" json:"max_history_size"`
	ThreatThreshold      float64 `yaml:"threat_threshold" json:"threat_threshold"`
	ResponseEnabled      bool    `yaml:"response_enabled" json:"response_enabled"`
	VoiceAPIKey          string  `yaml:"voice_api_key" json:"voice_api_key"`
}

// IRCConfig holds IRC processing configuration
type IRCConfig struct {
	Enabled           bool     `yaml:"enabled" json:"enabled"`
	Servers           []string `yaml:"servers" json:"servers"`
	Channels          []string `yaml:"channels" json:"channels"`
	Nickname          string   `yaml:"nickname" json:"nickname"`
	Username          string   `yaml:"username" json:"username"`
	Realname          string   `yaml:"realname" json:"realname"`
	Password          string   `yaml:"password" json:"password"`
	TLS               bool     `yaml:"tls" json:"tls"`
	TLSInsecure       bool     `yaml:"tls_insecure" json:"tls_insecure"`
	ReconnectInterval int      `yaml:"reconnect_interval" json:"reconnect_interval"`
	MaxHistorySize    int      `yaml:"max_history_size" json:"max_history_size"`
	ResponseEnabled   bool     `yaml:"response_enabled" json:"response_enabled"`
}

// ESMTPConfig holds ESMTP processing configuration
type ESMTPConfig struct {
	Enabled         bool     `yaml:"enabled" json:"enabled"`
	ListenAddress   string   `yaml:"listen_address" json:"listen_address"`
	Port            int      `yaml:"port" json:"port"`
	TLS             bool     `yaml:"tls" json:"tls"`
	CertFile        string   `yaml:"cert_file" json:"cert_file"`
	KeyFile         string   `yaml:"key_file" json:"key_file"`
	MaxMessageSize  int64    `yaml:"max_message_size" json:"max_message_size"`
	AllowedDomains  []string `yaml:"allowed_domains" json:"allowed_domains"`
	ResponseEnabled bool     `yaml:"response_enabled" json:"response_enabled"`
	SMTPRelay       string   `yaml:"smtp_relay" json:"smtp_relay"`
}

// DiscordConfig holds Discord processing configuration
type DiscordConfig struct {
	Enabled         bool     `yaml:"enabled" json:"enabled"`
	BotToken        string   `yaml:"bot_token" json:"bot_token"`
	WebhookEnabled  bool     `yaml:"webhook_enabled" json:"webhook_enabled"`
	WebhookPort     int      `yaml:"webhook_port" json:"webhook_port"`
	GuildIDs        []string `yaml:"guild_ids" json:"guild_ids"`
	ChannelIDs      []string `yaml:"channel_ids" json:"channel_ids"`
	ResponseEnabled bool     `yaml:"response_enabled" json:"response_enabled"`
	MaxHistorySize  int      `yaml:"max_history_size" json:"max_history_size"`
}

// WebhookConfig holds webhook processing configuration
type WebhookConfig struct {
	Enabled        bool          `yaml:"enabled" json:"enabled"`
	Port           int           `yaml:"port" json:"port"`
	Host           string        `yaml:"host" json:"host"`
	ReadTimeout    time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout" json:"write_timeout"`
	MaxRequestSize int64         `yaml:"max_request_size" json:"max_request_size"`
	EnableLogging  bool          `yaml:"enable_logging" json:"enable_logging"`
	AllowedOrigins []string      `yaml:"allowed_origins" json:"allowed_origins"`
}

// QueueConfig holds message queue configuration
type QueueConfig struct {
	Type           string        `yaml:"type" json:"type"` // "memory", "redis", "rabbitmq"
	Address        string        `yaml:"address" json:"address"`
	MaxSize        int           `yaml:"max_size" json:"max_size"`
	RetryAttempts  int           `yaml:"retry_attempts" json:"retry_attempts"`
	RetryDelay     time.Duration `yaml:"retry_delay" json:"retry_delay"`
	PersistEnabled bool          `yaml:"persist_enabled" json:"persist_enabled"`
}

// MasterControlConfig holds master control communication configuration
type MasterControlConfig struct {
	Address    string        `yaml:"address" json:"address"`
	Port       int           `yaml:"port" json:"port"`
	TLS        bool          `yaml:"tls" json:"tls"`
	Timeout    time.Duration `yaml:"timeout" json:"timeout"`
	RetryCount int           `yaml:"retry_count" json:"retry_count"`
	RetryDelay time.Duration `yaml:"retry_delay" json:"retry_delay"`
}

// Config represents the main configuration structure used by applications
type Config struct {
	HTTP          HTTPConfig          `yaml:"http" json:"http"`
	GRPC          GRPCConfig          `yaml:"grpc" json:"grpc"`
	Storage       StorageConfig       `yaml:"storage" json:"storage"`
	Security      SecurityConfig      `yaml:"security" json:"security"`
	Logging       LoggingConfig       `yaml:"logging" json:"logging"`
	SMS           SMSConfig           `yaml:"sms" json:"sms"`
	Voice         VoiceConfig         `yaml:"voice" json:"voice"`
	IRC           IRCConfig           `yaml:"irc" json:"irc"`
	ESMTP         ESMTPConfig         `yaml:"esmtp" json:"esmtp"`
	Discord       DiscordConfig       `yaml:"discord" json:"discord"`
	Webhook       WebhookConfig       `yaml:"webhook" json:"webhook"`
	Queue         QueueConfig         `yaml:"queue" json:"queue"`
	MasterControl MasterControlConfig `yaml:"master_control" json:"master_control"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Address         string        `yaml:"address" json:"address"`
	Port            string        `yaml:"port" json:"port"`
	Host            string        `yaml:"host" json:"host"`
	ReadTimeout     time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" json:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout"`
	EnableTLS       bool          `yaml:"enable_tls" json:"enable_tls"`
	CertFile        string        `yaml:"cert_file" json:"cert_file"`
	KeyFile         string        `yaml:"key_file" json:"key_file"`
	EnableCORS      bool          `yaml:"enable_cors" json:"enable_cors"`
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Address           string        `yaml:"address" json:"address"`
	Port              string        `yaml:"port" json:"port"`
	Host              string        `yaml:"host" json:"host"`
	MaxRecvMsgSize    int           `yaml:"max_recv_msg_size" json:"max_recv_msg_size"`
	MaxSendMsgSize    int           `yaml:"max_send_msg_size" json:"max_send_msg_size"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" json:"connection_timeout"`
	EnableTLS         bool          `yaml:"enable_tls" json:"enable_tls"`
	TLS               bool          `yaml:"tls" json:"tls"`
	CertFile          string        `yaml:"cert_file" json:"cert_file"`
	KeyFile           string        `yaml:"key_file" json:"key_file"`
	EnableReflection  bool          `yaml:"enable_reflection" json:"enable_reflection"`
}

// Validate validates the CommonConfig
func (c *CommonConfig) Validate() ValidationErrors {
	var errors []ValidationError

	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if c.LogLevel != "" && !Contains(validLogLevels, c.LogLevel) {
		errors = append(errors, ValidationError{
			Field:   "log_level",
			Message: fmt.Sprintf("invalid log level, must be one of: %s", strings.Join(validLogLevels, ", ")),
		})
	}

	if c.Timeout > 0 {
		if err := ValidateTimeout(c.Timeout, "timeout"); err != nil {
			errors = append(errors, *err)
		}
	}

	return ValidationErrors(errors)
}

// Validate validates the ServerConfig
func (c *ServerConfig) Validate() ValidationErrors {
	var errors []ValidationError

	if err := ValidatePort(c.HTTPPort, "http_port"); err != nil {
		errors = append(errors, *err)
	}

	if err := ValidatePort(c.GRPCPort, "grpc_port"); err != nil {
		errors = append(errors, *err)
	}

	if c.HTTPPort == c.GRPCPort {
		errors = append(errors, ValidationError{
			Field:   "ports",
			Message: "HTTP and gRPC ports cannot be the same",
		})
	}

	if c.EnableTLS {
		if err := ValidateRequired(c.CertFile, "cert_file"); err != nil {
			errors = append(errors, *err)
		}
		if err := ValidateRequired(c.KeyFile, "key_file"); err != nil {
			errors = append(errors, *err)
		}
	}

	if c.ReadTimeout > 0 {
		if err := ValidateTimeout(c.ReadTimeout, "read_timeout"); err != nil {
			errors = append(errors, *err)
		}
	}

	if c.WriteTimeout > 0 {
		if err := ValidateTimeout(c.WriteTimeout, "write_timeout"); err != nil {
			errors = append(errors, *err)
		}
	}

	if c.ShutdownTimeout > 0 {
		if err := ValidateTimeout(c.ShutdownTimeout, "shutdown_timeout"); err != nil {
			errors = append(errors, *err)
		}
	}

	return ValidationErrors(errors)
}

// Validate validates the SecurityConfig
func (c *SecurityConfig) Validate() ValidationErrors {
	var errors []ValidationError

	if c.EnableAuth || c.RequireAPIKey {
		if err := ValidateRequired(c.APIKey, "api_key"); err != nil {
			errors = append(errors, *err)
		} else if err := ValidateAPIKey(c.APIKey, "api_key"); err != nil {
			errors = append(errors, *err)
		}
	}

	if c.RateLimitRPM < 0 {
		errors = append(errors, ValidationError{
			Field:   "rate_limit_rpm",
			Message: "rate limit must be non-negative",
		})
	}

	return ValidationErrors(errors)
}

// Validate validates the StorageConfig
func (c *StorageConfig) Validate() ValidationErrors {
	var errors []ValidationError

	validTypes := []string{"memory", "file", "redis", "postgres"}
	if !Contains(validTypes, c.Type) {
		errors = append(errors, ValidationError{
			Field:   "storage.type",
			Message: fmt.Sprintf("invalid storage type, must be one of: %s", strings.Join(validTypes, ", ")),
		})
	}

	if c.Type == "file" {
		if err := ValidateRequired(c.DataDir, "data_dir"); err != nil {
			errors = append(errors, *err)
		}
	}

	return ValidationErrors(errors)
}

// Validate validates the MonitoringConfig
func (c *MonitoringConfig) Validate() ValidationErrors {
	var errors []ValidationError

	if c.EnableMetrics {
		if err := ValidatePort(c.MetricsPort, "metrics_port"); err != nil {
			errors = append(errors, *err)
		}
	}

	if c.HealthCheckInterval < 0 {
		errors = append(errors, ValidationError{
			Field:   "health_check_interval",
			Message: "health check interval must be non-negative",
		})
	}

	return ValidationErrors(errors)
}

// Validate validates the main Config
func (c *Config) Validate() error {
	var allErrors []ValidationError

	// Validate HTTP config
	if c.HTTP.Port == "" {
		allErrors = append(allErrors, ValidationError{
			Field:   "http.port",
			Message: "HTTP port is required",
		})
	}

	// Validate gRPC config
	if c.GRPC.Port == "" {
		allErrors = append(allErrors, ValidationError{
			Field:   "grpc.port",
			Message: "gRPC port is required",
		})
	}

	// Validate other components
	allErrors = append(allErrors, c.Security.Validate()...)
	allErrors = append(allErrors, c.Storage.Validate()...)

	if len(allErrors) > 0 {
		return ValidationErrors(allErrors)
	}

	return nil
}

// GetDefaults returns default configuration values
func GetDefaults() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Address:      "localhost:8083",
			Port:         "8083",
			Host:         "localhost",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			EnableCORS:   true,
		},
		GRPC: GRPCConfig{
			Address: "localhost:9092",
			Port:    "9092",
			Host:    "localhost",
			TLS:     false,
		},
		SMS: SMSConfig{
			Enabled:            true,
			ProcessingInterval: 30,
			MaxHistorySize:     1000,
			ThreatThreshold:    0.5,
			ResponseEnabled:    false,
		},
		Voice: VoiceConfig{
			Enabled:            true,
			MonitoringInterval: 30,
			MaxHistorySize:     1000,
			ThreatThreshold:    0.5,
			ResponseEnabled:    false,
		},
		IRC: IRCConfig{
			Enabled:           true,
			Servers:           []string{"irc.libera.chat:6667"},
			Channels:          []string{"#fr0g-ai-test"},
			Nickname:          "fr0g-ai",
			Username:          "fr0g-ai",
			Realname:          "fr0g.ai Security Bot",
			ReconnectInterval: 30,
			MaxHistorySize:    1000,
			ResponseEnabled:   false,
		},
		ESMTP: ESMTPConfig{
			Enabled:        true,
			ListenAddress:  "0.0.0.0",
			Port:           2525,
			MaxMessageSize: 10 * 1024 * 1024, // 10MB
			ResponseEnabled: false,
		},
		Discord: DiscordConfig{
			Enabled:        true,
			WebhookEnabled: true,
			WebhookPort:    8084,
			MaxHistorySize: 1000,
			ResponseEnabled: false,
		},
		Webhook: WebhookConfig{
			Enabled:        true,
			Port:           8085,
			Host:           "0.0.0.0",
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxRequestSize: 1024 * 1024, // 1MB
			EnableLogging:  true,
		},
		Queue: QueueConfig{
			Type:          "memory",
			MaxSize:       10000,
			RetryAttempts: 3,
			RetryDelay:    5 * time.Second,
		},
		MasterControl: MasterControlConfig{
			Address:    "localhost:8081",
			Port:       8081,
			TLS:        false,
			Timeout:    30 * time.Second,
			RetryCount: 3,
			RetryDelay: 5 * time.Second,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		},
		Security: SecurityConfig{
			EnableAuth:     false,
			EnableCORS:     true,
			RateLimitRPM:   1000,
			RequireAPIKey:  false,
		},
		Storage: StorageConfig{
			Type:    "memory",
			DataDir: "/tmp/fr0g-ai-io",
		},
	}
}
