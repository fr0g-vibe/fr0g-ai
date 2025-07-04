package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	sharedconfig "pkg/config"
)




func (c *Config) validateHTTPConfig() []sharedconfig.ValidationError {
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
