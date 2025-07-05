package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// ServiceRegistryConfig holds service registry configuration
type ServiceRegistryConfig struct {
	Enabled     bool          `yaml:"enabled"`
	URL         string        `yaml:"url"`
	ServiceName string        `yaml:"service_name"`
	ServiceID   string        `yaml:"service_id"`
	Tags        []string      `yaml:"tags"`
	Meta        map[string]string `yaml:"meta"`
	HealthInterval time.Duration `yaml:"health_interval"`
}

// Config holds the application configuration
type Config struct {
	Server          sharedconfig.ServerConfig     `yaml:"server"`
	OpenWebUI       sharedconfig.OpenWebUIConfig  `yaml:"openwebui"`
	ServiceRegistry ServiceRegistryConfig         `yaml:"service_registry"`
	Logging         sharedconfig.LoggingConfig    `yaml:"logging"`
	Security        sharedconfig.SecurityConfig   `yaml:"security"`
	Monitoring      sharedconfig.MonitoringConfig `yaml:"monitoring"`
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
		ServiceRegistry: ServiceRegistryConfig{
			Enabled:        false, // Disabled by default
			URL:            "http://localhost:8500",
			ServiceName:    "fr0g-ai-bridge",
			ServiceID:      "fr0g-ai-bridge-1",
			Tags:           []string{"ai", "bridge", "openwebui"},
			Meta:           map[string]string{"version": "1.0.0"},
			HealthInterval: 30 * time.Second,
		},
		Logging: sharedconfig.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Security: sharedconfig.SecurityConfig{
			EnableCORS:       true,
			AllowedOrigins:   []string{"*"},
			RateLimitRPM:     60,
			RequireAPIKey:    false,
			EnableReflection: true,
		},
		Monitoring: sharedconfig.MonitoringConfig{
			EnableMetrics:       true,
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

	// Service Registry environment variables
	if enableRegistry := os.Getenv("SERVICE_REGISTRY_ENABLED"); enableRegistry == "true" {
		config.ServiceRegistry.Enabled = true
	}

	if registryURL := os.Getenv("SERVICE_REGISTRY_URL"); registryURL != "" {
		config.ServiceRegistry.URL = registryURL
	}

	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		config.ServiceRegistry.ServiceName = serviceName
	}

	if serviceID := os.Getenv("SERVICE_ID"); serviceID != "" {
		config.ServiceRegistry.ServiceID = serviceID
	}

	if healthIntervalStr := os.Getenv("SERVICE_REGISTRY_HEALTH_INTERVAL"); healthIntervalStr != "" {
		if duration, err := time.ParseDuration(healthIntervalStr); err == nil {
			config.ServiceRegistry.HealthInterval = duration
		}
	}

	return config, nil
}
