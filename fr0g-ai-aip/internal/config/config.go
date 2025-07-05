package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// GetString gets a string value from environment or returns default
func (c *Config) GetString(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}

// Config holds all application configuration
type Config struct {
	// Server configuration using shared types
	HTTP sharedconfig.HTTPConfig `yaml:"http"`
	GRPC sharedconfig.GRPCConfig `yaml:"grpc"`
	
	// Storage configuration using shared types
	Storage sharedconfig.StorageConfig `yaml:"storage"`
	
	// Security configuration using shared types
	Security sharedconfig.SecurityConfig `yaml:"security"`
	
	// Logging configuration using shared types
	Logging sharedconfig.LoggingConfig `yaml:"logging"`
	
	// AIP-specific configuration
	Client ClientConfig `yaml:"client"`
	Validation ValidationConfig `yaml:"validation"`
}

type ClientConfig struct {
	Type      string `yaml:"type"`       // local, rest, grpc
	ServerURL string `yaml:"server_url"`
	Timeout   time.Duration `yaml:"timeout"`
}

type ValidationConfig struct {
	StrictMode bool `yaml:"strict_mode"`
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	config := &Config{
		HTTP: sharedconfig.HTTPConfig{
			Port:            getEnv("FR0G_HTTP_PORT", "8080"),
			Host:            getEnv("FR0G_HTTP_HOST", "0.0.0.0"),
			ReadTimeout:     getDurationEnv("FR0G_HTTP_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getDurationEnv("FR0G_HTTP_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getDurationEnv("FR0G_HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			EnableTLS:       getBoolEnv("FR0G_HTTP_ENABLE_TLS", false),
			CertFile:        getEnv("FR0G_HTTP_CERT_FILE", ""),
			KeyFile:         getEnv("FR0G_HTTP_KEY_FILE", ""),
		},
		GRPC: sharedconfig.GRPCConfig{
			Port:              getEnv("FR0G_GRPC_PORT", "9090"),
			Host:              getEnv("FR0G_GRPC_HOST", "0.0.0.0"),
			MaxRecvMsgSize:    getIntEnv("FR0G_GRPC_MAX_RECV_MSG_SIZE", 4*1024*1024), // 4MB
			MaxSendMsgSize:    getIntEnv("FR0G_GRPC_MAX_SEND_MSG_SIZE", 4*1024*1024), // 4MB
			ConnectionTimeout: getDurationEnv("FR0G_GRPC_CONNECTION_TIMEOUT", 5*time.Second),
			EnableTLS:         getBoolEnv("FR0G_GRPC_ENABLE_TLS", false),
			CertFile:          getEnv("FR0G_GRPC_CERT_FILE", ""),
			KeyFile:           getEnv("FR0G_GRPC_KEY_FILE", ""),
		},
		Storage: sharedconfig.StorageConfig{
			Type:    getEnv("FR0G_STORAGE_TYPE", "file"),
			DataDir: getEnv("FR0G_DATA_DIR", "./data"),
		},
		Security: sharedconfig.SecurityConfig{
			EnableAuth:       getBoolEnv("FR0G_ENABLE_AUTH", false),
			APIKey:           getEnv("FR0G_API_KEY", ""),
			EnableCORS:       getBoolEnv("FR0G_ENABLE_CORS", true),
			AllowedOrigins:   []string{"*"},
			RateLimitRPM:     getIntEnv("FR0G_RATE_LIMIT_RPM", 60),
			RequireAPIKey:    getBoolEnv("FR0G_REQUIRE_API_KEY", false),
			EnableReflection: getBoolEnv("FR0G_ENABLE_REFLECTION", true),
		},
		Logging: sharedconfig.LoggingConfig{
			Level:  getEnv("FR0G_LOG_LEVEL", "info"),
			Format: getEnv("FR0G_LOG_FORMAT", "text"),
		},
		Client: ClientConfig{
			Type:      getEnv("FR0G_CLIENT_TYPE", "grpc"),
			ServerURL: getEnv("FR0G_SERVER_URL", "localhost:9090"),
			Timeout:   getDurationEnv("FR0G_CLIENT_TIMEOUT", 30*time.Second),
		},
		Validation: ValidationConfig{
			StrictMode: getBoolEnv("FR0G_VALIDATION_STRICT", false),
		},
	}
	
	// Expand relative paths
	if !filepath.IsAbs(config.Storage.DataDir) {
		if abs, err := filepath.Abs(config.Storage.DataDir); err == nil {
			config.Storage.DataDir = abs
		}
	}
	
	return config
}

// LoadConfig loads configuration using shared config loader
func LoadConfig(configPath string) (*Config, error) {
	loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
		ConfigPath: configPath,
		EnvPrefix:  "FR0G_AIP",
	})
	
	// Load environment files
	if err := loader.LoadEnvFiles(); err != nil {
		// Non-fatal, just log warning
	}
	
	// Create default config
	cfg := &Config{
		HTTP: sharedconfig.HTTPConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
		GRPC: sharedconfig.GRPCConfig{
			Port: "9090",
			Host: "0.0.0.0",
		},
		Storage: sharedconfig.StorageConfig{
			Type:    "memory",
			DataDir: "./data",
		},
		Security: sharedconfig.SecurityConfig{
			EnableCORS:       true,
			AllowedOrigins:   []string{"*"},
			RateLimitRPM:     60,
			RequireAPIKey:    false,
			EnableReflection: true,
		},
		Logging: sharedconfig.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Client: ClientConfig{
			Type:      "grpc",
			ServerURL: "localhost:9090",
			Timeout:   30 * time.Second,
		},
		Validation: ValidationConfig{
			StrictMode: false,
		},
	}
	
	// Load from file
	if err := loader.LoadFromFile(cfg); err != nil {
		return nil, err
	}
	
	// Override with environment variables
	cfg.HTTP.Port = loader.GetEnvString("HTTP_PORT", cfg.HTTP.Port)
	cfg.HTTP.Host = loader.GetEnvString("HTTP_HOST", cfg.HTTP.Host)
	cfg.GRPC.Port = loader.GetEnvString("GRPC_PORT", cfg.GRPC.Port)
	cfg.GRPC.Host = loader.GetEnvString("GRPC_HOST", cfg.GRPC.Host)
	cfg.Logging.Level = loader.GetEnvString("LOG_LEVEL", cfg.Logging.Level)
	cfg.Logging.Format = loader.GetEnvString("LOG_FORMAT", cfg.Logging.Format)
	
	return cfg, nil
}

// Validate validates the entire configuration using shared validation
func (c *Config) Validate() error {
	var errors []sharedconfig.ValidationError
	
	// Validate HTTP config
	if c.HTTP.Port == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "http.port",
			Message: "HTTP port is required",
		})
	}
	
	// Validate gRPC config
	if c.GRPC.Port == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "grpc.port",
			Message: "gRPC port is required",
		})
	}
	
	// Validate other components using shared validation
	errors = append(errors, c.Security.Validate()...)
	errors = append(errors, c.Storage.Validate()...)
	
	// Validate client type
	validClientTypes := []string{"local", "rest", "grpc"}
	if !contains(validClientTypes, c.Client.Type) {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "client.type",
			Message: "invalid client type, must be one of: local, rest, grpc",
		})
	}
	
	// Validate timeout
	if c.Client.Timeout <= 0 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "client.timeout",
			Message: "timeout must be positive",
		})
	}
	
	if len(errors) > 0 {
		return sharedconfig.ValidationErrors(errors)
	}
	
	return nil
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
