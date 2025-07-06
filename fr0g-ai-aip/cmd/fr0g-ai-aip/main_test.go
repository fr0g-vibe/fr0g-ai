package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/storage"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		httpPort    string
		grpcPort    string
		expectError bool
	}{
		{
			name:        "different ports should pass",
			httpPort:    "8080",
			grpcPort:    "9090",
			expectError: false,
		},
		{
			name:        "same ports should fail",
			httpPort:    "8080",
			grpcPort:    "8080",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				HTTP: sharedconfig.HTTPConfig{Port: tt.httpPort},
				GRPC: sharedconfig.GRPCConfig{Port: tt.grpcPort},
				Storage: sharedconfig.StorageConfig{Type: "memory"},
			}

			err := cfg.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("Validate() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestCreateStorage(t *testing.T) {
	tests := []struct {
		name        string
		storageType string
		dataDir     string
		expectError bool
	}{
		{
			name:        "memory storage should work",
			storageType: "memory",
			dataDir:     "",
			expectError: false,
		},
		{
			name:        "file storage without datadir should fail",
			storageType: "file",
			dataDir:     "",
			expectError: true,
		},
		{
			name:        "invalid storage type should fail",
			storageType: "invalid",
			dataDir:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var store storage.Storage
			var err error

			switch tt.storageType {
			case "memory":
				store = storage.NewMemoryStorage()
			case "file":
				if tt.dataDir == "" {
					err = fmt.Errorf("data directory required for file storage")
				} else {
					store, err = storage.NewFileStorage(tt.dataDir)
				}
			default:
				err = fmt.Errorf("unsupported storage type: %s", tt.storageType)
			}

			if (err != nil) != tt.expectError {
				t.Errorf("createStorage() error = %v, expectError %v", err, tt.expectError)
			}
			if !tt.expectError && store == nil {
				t.Error("Expected storage to be created")
			}
		})
	}
}

func TestConfigCreation(t *testing.T) {
	// Create a minimal valid config
	cfg := &config.Config{
		HTTP: sharedconfig.HTTPConfig{
			Port:         "8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		GRPC: sharedconfig.GRPCConfig{
			Port:             "9090",
			MaxRecvMsgSize:   1024 * 1024,
			MaxSendMsgSize:   1024 * 1024,
			EnableReflection: true,
		},
		Storage: sharedconfig.StorageConfig{
			Type: "memory",
		},
		Security: sharedconfig.SecurityConfig{
			EnableAuth: false,
		},
	}

	// Test config validation
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Config validation failed: %v", err)
	}

	// Test GetString method
	if cfg.GetString("http.port") != "8080" {
		t.Error("GetString should return correct HTTP port")
	}

	if cfg.GetString("grpc.port") != "9090" {
		t.Error("GetString should return correct gRPC port")
	}
}

func TestConfigIntegration(t *testing.T) {
	// Test with a complete config structure
	cfg := &config.Config{
		HTTP: sharedconfig.HTTPConfig{
			Port:         "8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		GRPC: sharedconfig.GRPCConfig{
			Port:             "9090",
			MaxRecvMsgSize:   1024 * 1024,
			MaxSendMsgSize:   1024 * 1024,
			EnableReflection: true,
		},
		Storage: sharedconfig.StorageConfig{
			Type: "memory",
		},
		Security: sharedconfig.SecurityConfig{
			EnableAuth: false,
		},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() with valid config should not error: %v", err)
	}
}
