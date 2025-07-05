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
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
}

// Config represents the main configuration structure used by applications
type Config struct {
	HTTP     HTTPConfig     `yaml:"http" json:"http"`
	GRPC     GRPCConfig     `yaml:"grpc" json:"grpc"`
	Storage  StorageConfig  `yaml:"storage" json:"storage"`
	Security SecurityConfig `yaml:"security" json:"security"`
	Logging  LoggingConfig  `yaml:"logging" json:"logging"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port            string        `yaml:"port" json:"port"`
	Host            string        `yaml:"host" json:"host"`
	ReadTimeout     time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" json:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout"`
	EnableTLS       bool          `yaml:"enable_tls" json:"enable_tls"`
	CertFile        string        `yaml:"cert_file" json:"cert_file"`
	KeyFile         string        `yaml:"key_file" json:"key_file"`
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Port              string        `yaml:"port" json:"port"`
	Host              string        `yaml:"host" json:"host"`
	MaxRecvMsgSize    int           `yaml:"max_recv_msg_size" json:"max_recv_msg_size"`
	MaxSendMsgSize    int           `yaml:"max_send_msg_size" json:"max_send_msg_size"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" json:"connection_timeout"`
	EnableTLS         bool          `yaml:"enable_tls" json:"enable_tls"`
	CertFile          string        `yaml:"cert_file" json:"cert_file"`
	KeyFile           string        `yaml:"key_file" json:"key_file"`
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
