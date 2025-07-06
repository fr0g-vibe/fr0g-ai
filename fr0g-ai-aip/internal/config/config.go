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
	Client     ClientConfig                `yaml:"client"`
}

// ValidationConfig represents validation-specific configuration
type ValidationConfig struct {
	StrictMode bool `yaml:"strict_mode"`
}

// ClientConfig represents client-specific configuration
type ClientConfig struct {
	Type      string `yaml:"type"`
	ServerURL string `yaml:"server_url"`
}

// LoadConfig loads the configuration from environment variables and files
func LoadConfig(configPath string) (*Config, error) {
	// Use the centralized loader
	loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
		ConfigPath: configPath,
		EnvPrefix:  "FR0G",
	})

	// Create config with direct field assignment
	cfg := &Config{
		HTTP: sharedconfig.HTTPConfig{
			Port:         getEnvOrDefault("HTTP_PORT", "8080"),
			Host:         getEnvOrDefault("HTTP_HOST", "0.0.0.0"),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			EnableTLS:    false,
		},
		GRPC: sharedconfig.GRPCConfig{
			Port:             getEnvOrDefault("GRPC_PORT", "9090"),
			Host:             getEnvOrDefault("GRPC_HOST", "0.0.0.0"),
			EnableReflection: getBoolEnv("GRPC_ENABLE_REFLECTION", true),
		},
		Storage: sharedconfig.StorageConfig{
			Type:    getEnvOrDefault("FR0G_STORAGE_TYPE", "file"),
			DataDir: getEnvOrDefault("FR0G_DATA_DIR", "./data"),
		},
		Security: sharedconfig.SecurityConfig{
			EnableAuth:       getBoolEnv("ENABLE_AUTH", false),
			EnableCORS:       getBoolEnv("ENABLE_CORS", true),
			AllowedOrigins:   []string{"*"},
			RateLimitRPM:     60,
			RequireAPIKey:    getBoolEnv("REQUIRE_API_KEY", false),
			EnableReflection: getBoolEnv("GRPC_ENABLE_REFLECTION", true),
		},
		Validation: ValidationConfig{
			StrictMode: getBoolEnv("VALIDATION_STRICT_MODE", false),
		},
		Client: ClientConfig{
			Type:      getEnvOrDefault("FR0G_CLIENT_TYPE", "local"),
			ServerURL: getEnvOrDefault("FR0G_SERVER_URL", "http://localhost:8080"),
		},
	}

	// Load from file if specified
	if configPath != "" {
		if err := loader.LoadFromFile(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// Load loads the configuration with default path
func Load() *Config {
	cfg, _ := LoadConfig("")
	return cfg
}

// Validate validates the configuration using centralized validation
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
	
	// Check for port conflicts
	if c.HTTP.Port == c.GRPC.Port {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "ports",
			Message: "HTTP and gRPC ports cannot be the same",
		})
	}
	
	// Validate storage config
	if err := sharedconfig.ValidateRequired(c.Storage.Type, "storage.type"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate storage type is supported
	if err := sharedconfig.ValidateEnum(c.Storage.Type, []string{"file", "memory"}, "storage.type"); err != nil {
		errors = append(errors, *err)
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// GetString implements the interface for getting string configuration values
func (c *Config) GetString(key string) string {
	switch key {
	case "http.port":
		return c.HTTP.Port
	case "http.host":
		return c.HTTP.Host
	case "grpc.port":
		return c.GRPC.Port
	case "grpc.host":
		return c.GRPC.Host
	case "storage.type":
		return c.Storage.Type
	case "storage.data_dir":
		return c.Storage.DataDir
	case "client.type":
		return c.Client.Type
	case "client.server_url":
		return c.Client.ServerURL
	default:
		return ""
	}
}

// Helper function to get environment variable or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper function to get boolean environment variable or default
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
