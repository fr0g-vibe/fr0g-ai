package config

import (
	"os"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// ServiceRegistryConfig holds service registry configuration
type ServiceRegistryConfig struct {
	Enabled        bool              `yaml:"enabled"`
	URL            string            `yaml:"url"`
	ServiceName    string            `yaml:"service_name"`
	ServiceID      string            `yaml:"service_id"`
	Tags           []string          `yaml:"tags"`
	Meta           map[string]string `yaml:"meta"`
	HealthInterval time.Duration     `yaml:"health_interval"`
}

// Config represents the AIP service configuration
type Config struct {
	HTTP            sharedconfig.HTTPConfig     `yaml:"http"`
	GRPC            sharedconfig.GRPCConfig     `yaml:"grpc"`
	Storage         sharedconfig.StorageConfig  `yaml:"storage"`
	Security        sharedconfig.SecurityConfig `yaml:"security"`
	Validation      ValidationConfig            `yaml:"validation"`
	Client          ClientConfig                `yaml:"client"`
	ServiceRegistry ServiceRegistryConfig       `yaml:"service_registry"`
}

// ValidationConfig represents validation-specific configuration
type ValidationConfig struct {
	StrictMode   bool `yaml:"strict_mode"`
	EnableStrict bool `yaml:"enable_strict"`
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
			EnableReflection: getBoolEnv("GRPC_ENABLE_REFLECTION", false),
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
		ServiceRegistry: ServiceRegistryConfig{
			Enabled:        getBoolEnv("SERVICE_REGISTRY_ENABLED", false),
			URL:            getEnvOrDefault("SERVICE_REGISTRY_URL", "http://localhost:8500"),
			ServiceName:    getEnvOrDefault("SERVICE_NAME", "fr0g-ai-aip"),
			ServiceID:      getEnvOrDefault("SERVICE_ID", "fr0g-ai-aip-1"),
			Tags:           []string{"ai", "personas", "identities"},
			Meta:           map[string]string{"version": "1.0.0"},
			HealthInterval: 30 * time.Second,
		},
	}

	// Load from file if specified
	if configPath != "" {
		if err := loader.LoadFromFile(cfg); err != nil {
			return nil, err
		}
	}

	// Override service registry settings from environment
	if enableRegistry := os.Getenv("SERVICE_REGISTRY_ENABLED"); enableRegistry == "true" {
		cfg.ServiceRegistry.Enabled = true
	}

	if registryURL := os.Getenv("SERVICE_REGISTRY_URL"); registryURL != "" {
		cfg.ServiceRegistry.URL = registryURL
	}

	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		cfg.ServiceRegistry.ServiceName = serviceName
	}

	if serviceID := os.Getenv("SERVICE_ID"); serviceID != "" {
		cfg.ServiceRegistry.ServiceID = serviceID
	}

	if healthIntervalStr := os.Getenv("SERVICE_REGISTRY_HEALTH_INTERVAL"); healthIntervalStr != "" {
		if duration, err := time.ParseDuration(healthIntervalStr); err == nil {
			cfg.ServiceRegistry.HealthInterval = duration
		}
	}

	// Add port information to service registry metadata
	cfg.ServiceRegistry.Meta["http_port"] = cfg.HTTP.Port
	cfg.ServiceRegistry.Meta["grpc_port"] = cfg.GRPC.Port

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
	
	// Validate Service Registry configuration
	if c.ServiceRegistry.Enabled {
		if err := sharedconfig.ValidateURL(c.ServiceRegistry.URL, "service_registry.url"); err != nil {
			errors = append(errors, *err)
		}

		if err := sharedconfig.ValidateRequired(c.ServiceRegistry.ServiceName, "service_registry.service_name"); err != nil {
			errors = append(errors, *err)
		}

		if err := sharedconfig.ValidateRequired(c.ServiceRegistry.ServiceID, "service_registry.service_id"); err != nil {
			errors = append(errors, *err)
		}

		if c.ServiceRegistry.HealthInterval <= 0 {
			errors = append(errors, sharedconfig.ValidationError{
				Field:   "service_registry.health_interval",
				Message: "health interval must be positive",
			})
		}
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
