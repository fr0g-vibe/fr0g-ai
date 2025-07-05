package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents the MCP server configuration
type Config struct {
	// Server configuration
	HTTPPort int    `yaml:"http_port" env:"MCP_HTTP_PORT"`
	Host     string `yaml:"host" env:"MCP_HOST"`
	
	// Logging configuration
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL"`
	LogFormat string `yaml:"log_format" env:"LOG_FORMAT"`
	
	// MCP-specific configuration
	LearningEnabled         bool          `yaml:"learning_enabled" env:"MCP_LEARNING_ENABLED"`
	SystemConsciousness     bool          `yaml:"system_consciousness" env:"MCP_SYSTEM_CONSCIOUSNESS"`
	EmergentCapabilities    bool          `yaml:"emergent_capabilities" env:"MCP_EMERGENT_CAPABILITIES"`
	MaxConcurrentWorkflows  int           `yaml:"max_concurrent_workflows" env:"MCP_MAX_CONCURRENT_WORKFLOWS"`
	HealthCheckInterval     time.Duration `yaml:"health_check_interval" env:"MCP_HEALTH_CHECK_INTERVAL"`
	
	// Input processor configuration
	ESMTPEnabled bool   `yaml:"esmtp_enabled" env:"ESMTP_ENABLED"`
	ESMTPPort    int    `yaml:"esmtp_port" env:"ESMTP_PORT"`
	ESMTPHost    string `yaml:"esmtp_host" env:"ESMTP_HOST"`
	
	DiscordEnabled bool `yaml:"discord_enabled" env:"DISCORD_ENABLED"`
	SMSEnabled     bool `yaml:"sms_enabled" env:"SMS_ENABLED"`
	VoiceEnabled   bool `yaml:"voice_enabled" env:"VOICE_ENABLED"`
	IRCEnabled     bool `yaml:"irc_enabled" env:"IRC_ENABLED"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		// Default values
		HTTPPort:                8081,
		Host:                    "0.0.0.0",
		LogLevel:                "info",
		LogFormat:               "json",
		LearningEnabled:         true,
		SystemConsciousness:     true,
		EmergentCapabilities:    true,
		MaxConcurrentWorkflows:  10,
		HealthCheckInterval:     30 * time.Second,
		ESMTPEnabled:            true,
		ESMTPPort:               2525,
		ESMTPHost:               "0.0.0.0",
		DiscordEnabled:          true,
		SMSEnabled:              true,
		VoiceEnabled:            true,
		IRCEnabled:              true,
	}
	
	// Load from environment variables
	if port := os.Getenv("MCP_HTTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.HTTPPort = p
		}
	}
	
	if host := os.Getenv("MCP_HOST"); host != "" {
		config.Host = host
	}
	
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.LogLevel = level
	}
	
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.LogFormat = format
	}
	
	// Boolean environment variables
	config.LearningEnabled = getBoolEnv("MCP_LEARNING_ENABLED", config.LearningEnabled)
	config.SystemConsciousness = getBoolEnv("MCP_SYSTEM_CONSCIOUSNESS", config.SystemConsciousness)
	config.EmergentCapabilities = getBoolEnv("MCP_EMERGENT_CAPABILITIES", config.EmergentCapabilities)
	config.ESMTPEnabled = getBoolEnv("ESMTP_ENABLED", config.ESMTPEnabled)
	config.DiscordEnabled = getBoolEnv("DISCORD_ENABLED", config.DiscordEnabled)
	config.SMSEnabled = getBoolEnv("SMS_ENABLED", config.SMSEnabled)
	config.VoiceEnabled = getBoolEnv("VOICE_ENABLED", config.VoiceEnabled)
	config.IRCEnabled = getBoolEnv("IRC_ENABLED", config.IRCEnabled)
	
	// Integer environment variables
	if workflows := os.Getenv("MCP_MAX_CONCURRENT_WORKFLOWS"); workflows != "" {
		if w, err := strconv.Atoi(workflows); err == nil {
			config.MaxConcurrentWorkflows = w
		}
	}
	
	if esmtpPort := os.Getenv("ESMTP_PORT"); esmtpPort != "" {
		if p, err := strconv.Atoi(esmtpPort); err == nil {
			config.ESMTPPort = p
		}
	}
	
	if esmtpHost := os.Getenv("ESMTP_HOST"); esmtpHost != "" {
		config.ESMTPHost = esmtpHost
	}
	
	// Duration environment variables
	if interval := os.Getenv("MCP_HEALTH_CHECK_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			config.HealthCheckInterval = d
		}
	}
	
	return config, nil
}

// getBoolEnv gets a boolean environment variable with a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}
