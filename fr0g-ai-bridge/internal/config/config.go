package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	OpenWebUI  OpenWebUIConfig  `yaml:"openwebui"`
	Logging    LoggingConfig    `yaml:"logging"`
	Security   SecurityConfig   `yaml:"security"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	HTTPPort int    `yaml:"http_port"`
	GRPCPort int    `yaml:"grpc_port"`
	Host     string `yaml:"host"`
}

// OpenWebUIConfig holds OpenWebUI API configuration
type OpenWebUIConfig struct {
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	Timeout int    `yaml:"timeout"` // timeout in seconds
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	EnableCORS           bool     `yaml:"enable_cors"`
	AllowedOrigins       []string `yaml:"allowed_origins"`
	RateLimitRPM         int      `yaml:"rate_limit_requests_per_minute"`
	RequireAPIKey        bool     `yaml:"require_api_key"`
	AllowedAPIKeys       []string `yaml:"allowed_api_keys"`
	EnableReflection     bool     `yaml:"enable_reflection"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	EnableMetrics         bool `yaml:"enable_metrics"`
	MetricsPort          int  `yaml:"metrics_port"`
	HealthCheckInterval  int  `yaml:"health_check_interval"`
	EnableTracing        bool `yaml:"enable_tracing"`
}

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // File doesn't exist, skip silently
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		// Only set if not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return nil
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Try to load .env file from current directory and parent directory
	envPaths := []string{
		".env",
		"../.env",
		"../../.env", // For when running from fr0g-ai-bridge subdirectory
	}
	
	for _, envPath := range envPaths {
		if err := loadEnvFile(envPath); err != nil {
			fmt.Printf("Warning: failed to load %s: %v\n", envPath, err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			HTTPPort: 8080,
			GRPCPort: 9090,
			Host:     "0.0.0.0",
		},
		OpenWebUI: OpenWebUIConfig{
			BaseURL: "http://localhost:3000",
			Timeout: 30,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Security: SecurityConfig{
			EnableCORS:           true,
			AllowedOrigins:       []string{"*"},
			RateLimitRPM:         60,
			RequireAPIKey:        false,
			EnableReflection:     true,
		},
		Monitoring: MonitoringConfig{
			EnableMetrics:        true,
			MetricsPort:         8082,
			HealthCheckInterval: 30,
			EnableTracing:       false,
		},
	}

	// Load from file if it exists
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}

			if err := yaml.Unmarshal(data, config); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		}
	}

	// Override with environment variables
	if port := os.Getenv("HTTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.HTTPPort = p
		}
	}

	if port := os.Getenv("GRPC_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.GRPCPort = p
		}
	}

	if host := os.Getenv("HOST"); host != "" {
		config.Server.Host = host
	}

	if baseURL := os.Getenv("OPENWEBUI_BASE_URL"); baseURL != "" {
		config.OpenWebUI.BaseURL = baseURL
	}

	if apiKey := os.Getenv("OPENWEBUI_API_KEY"); apiKey != "" {
		config.OpenWebUI.APIKey = apiKey
	}

	if timeout := os.Getenv("OPENWEBUI_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			config.OpenWebUI.Timeout = t
		}
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}

	// Security environment variables
	if enableCORS := os.Getenv("ENABLE_CORS"); enableCORS == "false" {
		config.Security.EnableCORS = false
	}

	if rateLimitStr := os.Getenv("RATE_LIMIT_RPM"); rateLimitStr != "" {
		if rateLimit, err := strconv.Atoi(rateLimitStr); err == nil {
			config.Security.RateLimitRPM = rateLimit
		}
	}

	if enableReflection := os.Getenv("ENABLE_REFLECTION"); enableReflection == "false" {
		config.Security.EnableReflection = false
	}

	// Monitoring environment variables
	if metricsPortStr := os.Getenv("METRICS_PORT"); metricsPortStr != "" {
		if metricsPort, err := strconv.Atoi(metricsPortStr); err == nil {
			config.Monitoring.MetricsPort = metricsPort
		}
	}

	if enableMetrics := os.Getenv("ENABLE_METRICS"); enableMetrics == "false" {
		config.Monitoring.EnableMetrics = false
	}

	return config, nil
}
