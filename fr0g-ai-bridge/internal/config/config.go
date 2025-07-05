package config

import (
	"fmt"
	"os"
	"strconv"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Config holds the application configuration
type Config struct {
	Server     sharedconfig.ServerConfig     `yaml:"server"`
	OpenWebUI  sharedconfig.OpenWebUIConfig  `yaml:"openwebui"`
	Logging    sharedconfig.LoggingConfig    `yaml:"logging"`
	Security   sharedconfig.SecurityConfig   `yaml:"security"`
	Monitoring sharedconfig.MonitoringConfig `yaml:"monitoring"`
}

// Validate validates the configuration using shared validation
func (c *Config) Validate() sharedconfig.ValidationErrors {
	var errors []sharedconfig.ValidationError
	
	// Validate shared configurations
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Security.Validate()...)
	errors = append(errors, c.Monitoring.Validate()...)
	
	// Validate OpenWebUI configuration
	if err := sharedconfig.ValidateURL(c.OpenWebUI.BaseURL, "openwebui.base_url"); err != nil {
		errors = append(errors, *err)
	}
	
	if c.OpenWebUI.Timeout <= 0 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "openwebui.timeout",
			Message: "timeout must be positive",
		})
	}
	
	return sharedconfig.ValidationErrors(errors)
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Create loader with standard options
	loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
		ConfigPath: configPath,
		EnvPrefix:  "",
		EnvFilePaths: []string{
			".env",
			"../.env",
			"../../.env", // For when running from fr0g-ai-bridge subdirectory
		},
	})
	
	// Load environment files
	if err := loader.LoadEnvFiles(); err != nil {
		fmt.Printf("Warning: failed to load env files: %v\n", err)
	}

	config := &Config{
		Server: sharedconfig.ServerConfig{
			HTTPPort: 8080,
			GRPCPort: 9090,
			Host:     "0.0.0.0",
		},
		OpenWebUI: sharedconfig.OpenWebUIConfig{
			BaseURL: "http://localhost:3000",
			Timeout: 30,
		},
		Logging: sharedconfig.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Security: sharedconfig.SecurityConfig{
			EnableCORS:           true,
			AllowedOrigins:       []string{"*"},
			RateLimitRPM:         60,
			RequireAPIKey:        false,
			EnableReflection:     true,
		},
		Monitoring: sharedconfig.MonitoringConfig{
			EnableMetrics:        true,
			MetricsPort:         8082,
			HealthCheckInterval: 30,
			EnableTracing:       false,
		},
	}

	// Load from file using standardized loader
	if err := loader.LoadFromFile(config); err != nil {
		return nil, err
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
		fmt.Printf("OpenWebUI API Key loaded: %s\n", apiKey)
	} else {
		fmt.Println("WARNING: No OPENWEBUI_API_KEY found in environment")
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
