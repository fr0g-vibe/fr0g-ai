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
	var errors []ValidationError
	
	// Validate port
	if err := sharedconfig.ValidateRequired(c.HTTP.Port, "http.port"); err != nil {
		errors = append(errors, *err)
	} else if err := sharedconfig.ValidatePort(c.HTTP.Port, "http.port"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate timeouts using shared validation
	if err := sharedconfig.ValidateTimeout(c.HTTP.ReadTimeout, "http.read_timeout"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidateTimeout(c.HTTP.WriteTimeout, "http.write_timeout"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidateTimeout(c.HTTP.ShutdownTimeout, "http.shutdown_timeout"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate TLS config
	if c.HTTP.EnableTLS {
		if err := sharedconfig.ValidateRequired(c.HTTP.CertFile, "http.cert_file"); err != nil {
			errors = append(errors, *err)
		}
		if err := sharedconfig.ValidateRequired(c.HTTP.KeyFile, "http.key_file"); err != nil {
			errors = append(errors, *err)
		}
	}
	
	return errors
}

func (c *Config) validateGRPCConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate port
	if err := sharedconfig.ValidateRequired(c.GRPC.Port, "grpc.port"); err != nil {
		errors = append(errors, *err)
	} else if err := sharedconfig.ValidatePort(c.GRPC.Port, "grpc.port"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate message sizes
	if err := sharedconfig.ValidatePositive(c.GRPC.MaxRecvMsgSize, "grpc.max_recv_msg_size"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidatePositive(c.GRPC.MaxSendMsgSize, "grpc.max_send_msg_size"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate connection timeout
	if err := sharedconfig.ValidateTimeout(c.GRPC.ConnectionTimeout, "grpc.connection_timeout"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate TLS config
	if c.GRPC.EnableTLS {
		if err := sharedconfig.ValidateRequired(c.GRPC.CertFile, "grpc.cert_file"); err != nil {
			errors = append(errors, *err)
		}
		if err := sharedconfig.ValidateRequired(c.GRPC.KeyFile, "grpc.key_file"); err != nil {
			errors = append(errors, *err)
		}
	}
	
	return errors
}

func (c *Config) validateStorageConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate storage type
	validTypes := []string{"memory", "file"}
	if !sharedconfig.Contains(validTypes, c.Storage.Type) {
		errors = append(errors, ValidationError{
			Field:   "storage.type",
			Message: fmt.Sprintf("invalid storage type, must be one of: %s", strings.Join(validTypes, ", ")),
		})
	}
	
	// Validate file storage specific config
	if c.Storage.Type == "file" {
		if err := sharedconfig.ValidateRequired(c.Storage.DataDir, "storage.data_dir"); err != nil {
			errors = append(errors, *err)
		}
	}
	
	return errors
}

func (c *Config) validateClientConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate client type
	validTypes := []string{"local", "rest", "grpc"}
	if c.Client.Type != "" && !sharedconfig.Contains(validTypes, c.Client.Type) {
		errors = append(errors, ValidationError{
			Field:   "client.type",
			Message: fmt.Sprintf("invalid client type, must be one of: %s", strings.Join(validTypes, ", ")),
		})
	}
	
	// Validate timeout
	if err := sharedconfig.ValidateTimeout(c.Client.Timeout, "client.timeout"); err != nil {
		errors = append(errors, *err)
	}
	
	return errors
}

func (c *Config) validateSecurityConfig() []ValidationError {
	var errors []ValidationError
	
	// Validate API key if auth is enabled
	if c.Security.EnableAuth {
		if err := sharedconfig.ValidateRequired(c.Security.APIKey, "security.api_key"); err != nil {
			errors = append(errors, *err)
		}
	}
	
	// Validate API key strength
	if c.Security.APIKey != "" {
		if err := sharedconfig.ValidateAPIKey(c.Security.APIKey, "security.api_key"); err != nil {
			errors = append(errors, *err)
		}
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

// ValidateNetworkAddress validates a network address using shared validation
func ValidateNetworkAddress(address string) error {
	if err := sharedconfig.ValidateNetworkAddress(address, "address"); err != nil {
		return fmt.Errorf(err.Message)
	}
	return nil
}
