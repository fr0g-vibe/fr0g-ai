package config_test

import (
	"os"
	"testing"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *sharedconfig.Config
		expectError bool
		errorCount  int
	}{
		{
			name: "valid configuration",
			config: &sharedconfig.Config{
				HTTP: sharedconfig.HTTPConfig{
					Port: 8080,
					Host: "localhost",
				},
				GRPC: sharedconfig.GRPCConfig{
					Port: 9090,
					Host: "localhost",
				},
				Storage: sharedconfig.StorageConfig{
					Type: "file",
					Path: "/tmp/test",
				},
			},
			expectError: false,
			errorCount:  0,
		},
		{
			name: "invalid HTTP port",
			config: &sharedconfig.Config{
				HTTP: sharedconfig.HTTPConfig{
					Port: 99999,
					Host: "localhost",
				},
				GRPC: sharedconfig.GRPCConfig{
					Port: 9090,
					Host: "localhost",
				},
				Storage: sharedconfig.StorageConfig{
					Type: "file",
					Path: "/tmp/test",
				},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "missing storage type",
			config: &sharedconfig.Config{
				HTTP: sharedconfig.HTTPConfig{
					Port: 8080,
					Host: "localhost",
				},
				GRPC: sharedconfig.GRPCConfig{
					Port: 9090,
					Host: "localhost",
				},
				Storage: sharedconfig.StorageConfig{
					Path: "/tmp/test",
				},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "port conflict",
			config: &sharedconfig.Config{
				HTTP: sharedconfig.HTTPConfig{
					Port: 8080,
					Host: "localhost",
				},
				GRPC: sharedconfig.GRPCConfig{
					Port: 8080, // Same as HTTP port
					Host: "localhost",
				},
				Storage: sharedconfig.StorageConfig{
					Type: "file",
					Path: "/tmp/test",
				},
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.config.Validate()
			
			if tt.expectError && len(errors) == 0 {
				t.Errorf("Expected validation errors but got none")
			}
			
			if !tt.expectError && len(errors) > 0 {
				t.Errorf("Expected no validation errors but got %d: %v", len(errors), errors)
			}
			
			if tt.errorCount > 0 && len(errors) != tt.errorCount {
				t.Errorf("Expected %d validation errors but got %d", tt.errorCount, len(errors))
			}
		})
	}
}

func TestConfigLoader(t *testing.T) {
	// Create temporary config file
	configContent := `
http:
  port: 8080
  host: "localhost"
grpc:
  port: 9090
  host: "localhost"
storage:
  type: "file"
  path: "/tmp/test"
`
	
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()
	
	// Test loading configuration
	loader := sharedconfig.NewLoader()
	config, err := loader.LoadFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Validate loaded configuration
	if config.HTTP.Port != 8080 {
		t.Errorf("Expected HTTP port 8080, got %d", config.HTTP.Port)
	}
	
	if config.GRPC.Port != 9090 {
		t.Errorf("Expected gRPC port 9090, got %d", config.GRPC.Port)
	}
	
	if config.Storage.Type != "file" {
		t.Errorf("Expected storage type 'file', got '%s'", config.Storage.Type)
	}
}

func TestEnvironmentVariableOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("HTTP_PORT", "8081")
	os.Setenv("GRPC_PORT", "9091")
	os.Setenv("STORAGE_TYPE", "database")
	defer func() {
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("GRPC_PORT")
		os.Unsetenv("STORAGE_TYPE")
	}()
	
	loader := sharedconfig.NewLoader()
	config, err := loader.LoadFromEnv()
	if err != nil {
		t.Fatalf("Failed to load config from env: %v", err)
	}
	
	if config.HTTP.Port != 8081 {
		t.Errorf("Expected HTTP port 8081 from env, got %d", config.HTTP.Port)
	}
	
	if config.GRPC.Port != 9091 {
		t.Errorf("Expected gRPC port 9091 from env, got %d", config.GRPC.Port)
	}
	
	if config.Storage.Type != "database" {
		t.Errorf("Expected storage type 'database' from env, got '%s'", config.Storage.Type)
	}
}

func TestValidationHelpers(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() *sharedconfig.ValidationError
		expectError bool
	}{
		{
			name: "valid port",
			testFunc: func() *sharedconfig.ValidationError {
				return sharedconfig.ValidatePort(8080, "test_port")
			},
			expectError: false,
		},
		{
			name: "invalid port - too high",
			testFunc: func() *sharedconfig.ValidationError {
				return sharedconfig.ValidatePort(99999, "test_port")
			},
			expectError: true,
		},
		{
			name: "invalid port - negative",
			testFunc: func() *sharedconfig.ValidationError {
				return sharedconfig.ValidatePort(-1, "test_port")
			},
			expectError: true,
		},
		{
			name: "valid timeout",
			testFunc: func() *sharedconfig.ValidationError {
				return sharedconfig.ValidateTimeout(30*time.Second, "test_timeout")
			},
			expectError: false,
		},
		{
			name: "invalid timeout - zero",
			testFunc: func() *sharedconfig.ValidationError {
				return sharedconfig.ValidateTimeout(0, "test_timeout")
			},
			expectError: true,
		},
		{
			name: "valid required field",
			testFunc: func() *sharedconfig.ValidationError {
				return sharedconfig.ValidateRequired("test_value", "test_field")
			},
			expectError: false,
		},
		{
			name: "invalid required field - empty",
			testFunc: func() *sharedconfig.ValidationError {
				return sharedconfig.ValidateRequired("", "test_field")
			},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			
			if tt.expectError && err == nil {
				t.Errorf("Expected validation error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

func BenchmarkConfigValidation(b *testing.B) {
	config := &sharedconfig.Config{
		HTTP: sharedconfig.HTTPConfig{
			Port: 8080,
			Host: "localhost",
		},
		GRPC: sharedconfig.GRPCConfig{
			Port: 9090,
			Host: "localhost",
		},
		Storage: sharedconfig.StorageConfig{
			Type: "file",
			Path: "/tmp/test",
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}

func BenchmarkConfigLoader(b *testing.B) {
	// Create temporary config file
	configContent := `
http:
  port: 8080
  host: "localhost"
grpc:
  port: 9090
  host: "localhost"
storage:
  type: "file"
  path: "/tmp/test"
`
	
	tmpFile, err := os.CreateTemp("", "bench-config-*.yaml")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(configContent); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()
	
	loader := sharedconfig.NewLoader()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loader.LoadFromFile(tmpFile.Name())
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}
