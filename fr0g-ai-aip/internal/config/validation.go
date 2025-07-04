package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	sharedconfig "pkg/config"
)

// ValidationError represents a validation error (alias to shared type)
type ValidationError = sharedconfig.ValidationError

// ValidationErrors represents multiple validation errors (alias to shared type)
type ValidationErrors = sharedconfig.ValidationErrors

// Config represents the application configuration structure
type Config struct {
	HTTP     HTTPConfig     `yaml:"http" json:"http"`
	GRPC     GRPCConfig     `yaml:"grpc" json:"grpc"`
	Storage  StorageConfig  `yaml:"storage" json:"storage"`
	Client   ClientConfig   `yaml:"client" json:"client"`
	Security SecurityConfig `yaml:"security" json:"security"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port            string        `yaml:"port" json:"port"`
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
	MaxRecvMsgSize    int           `yaml:"max_recv_msg_size" json:"max_recv_msg_size"`
	MaxSendMsgSize    int           `yaml:"max_send_msg_size" json:"max_send_msg_size"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" json:"connection_timeout"`
	EnableTLS         bool          `yaml:"enable_tls" json:"enable_tls"`
	CertFile          string        `yaml:"cert_file" json:"cert_file"`
	KeyFile           string        `yaml:"key_file" json:"key_file"`
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type    string `yaml:"type" json:"type"`
	DataDir string `yaml:"data_dir" json:"data_dir"`
}

// ClientConfig holds client configuration
type ClientConfig struct {
	Type    string        `yaml:"type" json:"type"`
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	EnableAuth bool   `yaml:"enable_auth" json:"enable_auth"`
	APIKey     string `yaml:"api_key" json:"api_key"`
}

// Validate validates the entire configuration
func (c *Config) Validate() ValidationErrors {
	var errors []ValidationError
	
	errors = append(errors, c.validateHTTPConfig()...)
	errors = append(errors, c.validateGRPCConfig()...)
	errors = append(errors, c.validateStorageConfig()...)
	errors = append(errors, c.validateClientConfig()...)
	errors = append(errors, c.validateSecurityConfig()...)
	errors = append(errors, c.validateCrossConfig()...)
	
	return ValidationErrors(errors)
}

func (c *Config) validateHTTPConfig() []ValidationError {
	var errors []sharedconfig.ValidationError
	
	// Validate port
	if c.HTTP.Port == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "http.port",
			Message: "port is required",
		})
	} else if err := sharedconfig.ValidatePort(c.HTTP.Port, "http.port"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate timeouts
	if c.HTTP.ReadTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "http.read_timeout",
			Message: "read timeout must be positive",
		})
	}
	
	if c.HTTP.WriteTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "http.write_timeout",
			Message: "write timeout must be positive",
		})
	}
	
	if c.HTTP.ShutdownTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "http.shutdown_timeout",
			Message: "shutdown timeout must be positive",
		})
	}
	
	// Validate TLS config
	if c.HTTP.EnableTLS {
		if c.HTTP.CertFile == "" {
			errors = append(errors, ValidationError{
				Field:   "http.cert_file",
				Message: "cert file is required when TLS is enabled",
			})
		}
		if c.HTTP.KeyFile == "" {
			errors = append(errors, ValidationError{
				Field:   "http.key_file",
				Message: "key file is required when TLS is enabled",
			})
		}
	}
	
	return errors
}

func (c *Config) validateGRPCConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate port
	if c.GRPC.Port == "" {
		errors = append(errors, ValidationError{
			Field:   "grpc.port",
			Message: "port is required",
		})
	} else if !isValidPort(c.GRPC.Port) {
		errors = append(errors, ValidationError{
			Field:   "grpc.port",
			Message: "invalid port number",
		})
	}
	
	// Validate message sizes
	if c.GRPC.MaxRecvMsgSize <= 0 {
		errors = append(errors, ValidationError{
			Field:   "grpc.max_recv_msg_size",
			Message: "max receive message size must be positive",
		})
	}
	
	if c.GRPC.MaxSendMsgSize <= 0 {
		errors = append(errors, ValidationError{
			Field:   "grpc.max_send_msg_size",
			Message: "max send message size must be positive",
		})
	}
	
	// Validate connection timeout
	if c.GRPC.ConnectionTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "grpc.connection_timeout",
			Message: "connection timeout must be positive",
		})
	}
	
	// Validate TLS config
	if c.GRPC.EnableTLS {
		if c.GRPC.CertFile == "" {
			errors = append(errors, ValidationError{
				Field:   "grpc.cert_file",
				Message: "cert file is required when TLS is enabled",
			})
		}
		if c.GRPC.KeyFile == "" {
			errors = append(errors, ValidationError{
				Field:   "grpc.key_file",
				Message: "key file is required when TLS is enabled",
			})
		}
	}
	
	return errors
}

func (c *Config) validateStorageConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate storage type
	validTypes := []string{"memory", "file"}
	if !contains(validTypes, c.Storage.Type) {
		errors = append(errors, ValidationError{
			Field:   "storage.type",
			Message: "invalid storage type",
		})
	}
	
	// Validate file storage specific config
	if c.Storage.Type == "file" && c.Storage.DataDir == "" {
		errors = append(errors, ValidationError{
			Field:   "storage.data_dir",
			Message: "data directory is required for file storage",
		})
	}
	
	return errors
}

func (c *Config) validateClientConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate client type
	validTypes := []string{"local", "rest", "grpc"}
	if c.Client.Type != "" && !contains(validTypes, c.Client.Type) {
		errors = append(errors, ValidationError{
			Field:   "client.type",
			Message: fmt.Sprintf("invalid client type, must be one of: %s", strings.Join(validTypes, ", ")),
		})
	}
	
	// Validate timeout
	if c.Client.Timeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "client.timeout",
			Message: "client timeout must be positive",
		})
	}
	
	return errors
}

func (c *Config) validateSecurityConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate API key if auth is enabled
	if c.Security.EnableAuth && c.Security.APIKey == "" {
		errors = append(errors, ValidationError{
			Field:   "security.api_key",
			Message: "API key is required when authentication is enabled",
		})
	}
	
	// Validate API key strength
	if c.Security.APIKey != "" && len(c.Security.APIKey) < 16 {
		errors = append(errors, ValidationError{
			Field:   "security.api_key",
			Message: "API key must be at least 16 characters long",
		})
	}
	
	return errors
}

func (c *Config) validateCrossConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate port conflicts
	if c.HTTP.Port == c.GRPC.Port {
		errors = append(errors, ValidationError{
			Field:   "ports",
			Message: "HTTP and gRPC ports cannot be the same",
		})
	}
	
	return errors
}

// Helper functions
func isValidPort(port string) bool {
	return sharedconfig.ValidatePort(port, "port") == nil
}

func contains(slice []string, item string) bool {
	return sharedconfig.Contains(slice, item)
}

// ValidateNetworkAddress validates a network address
func ValidateNetworkAddress(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid address format: %v", err)
	}
	
	// Validate host
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	
	// Validate port
	if !isValidPort(port) {
		return fmt.Errorf("invalid port: %s", port)
	}
	
	return nil
}

// ValidateTimeout validates a timeout duration
func ValidateTimeout(timeout time.Duration, name string) error {
	if timeout <= 0 {
		return fmt.Errorf("%s timeout must be positive", name)
	}
	if timeout > 24*time.Hour {
		return fmt.Errorf("%s timeout cannot exceed 24 hours", name)
	}
	return nil
}
