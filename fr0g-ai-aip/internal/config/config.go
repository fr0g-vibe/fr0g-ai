package config

import (
	"os"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Config represents the AIP service configuration
type Config struct {
	HTTP       sharedconfig.HTTPConfig     `yaml:"http"`
	GRPC       sharedconfig.GRPCConfig     `yaml:"grpc"`
	Storage    sharedconfig.StorageConfig  `yaml:"storage"`
	Security   sharedconfig.SecurityConfig `yaml:"security"`
	Validation ValidationConfig            `yaml:"validation"`
}

// ValidationConfig represents validation-specific configuration
type ValidationConfig struct {
	StrictMode bool `yaml:"strict_mode"`
}

// LoadConfig loads the configuration from environment variables and files
func LoadConfig(configPath string) (*Config, error) {
	// Create default config
	cfg := &Config{
		HTTP: sharedconfig.HTTPConfig{
			Port:         getEnvOrDefault("HTTP_PORT", "8080"),
			Host:         getEnvOrDefault("HTTP_HOST", "0.0.0.0"),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			EnableTLS:    false,
		},
		GRPC: sharedconfig.GRPCConfig{
			Port: getEnvOrDefault("GRPC_PORT", "9090"),
			Host: getEnvOrDefault("GRPC_HOST", "0.0.0.0"),
		},
		Storage: sharedconfig.StorageConfig{
			Type:    getEnvOrDefault("STORAGE_TYPE", "file"),
			DataDir: getEnvOrDefault("DATA_DIR", "./data"),
		},
		Security: sharedconfig.SecurityConfig{
			EnableAuth:       false,
			EnableCORS:       true,
			AllowedOrigins:   []string{"*"},
			RateLimitRPM:     60,
			RequireAPIKey:    false,
			EnableReflection: true,
		},
		Validation: ValidationConfig{
			StrictMode: false,
		},
	}
	
	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	var errors sharedconfig.ValidationErrors
	
	// Validate HTTP config
	if err := sharedconfig.ValidatePort(c.HTTP.Port, "http.port"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate gRPC config
	if err := sharedconfig.ValidatePort(c.GRPC.Port, "grpc.port"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate storage config
	if err := sharedconfig.ValidateRequired(c.Storage.Type, "storage.type"); err != nil {
		errors = append(errors, *err)
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// Helper function to get environment variable or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
