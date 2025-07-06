package config

import (
	"strings"
	"testing"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				HTTP: sharedconfig.HTTPConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				GRPC: sharedconfig.GRPCConfig{
					Port:             "9090",
					MaxRecvMsgSize:   4 * 1024 * 1024,
					MaxSendMsgSize:   4 * 1024 * 1024,
					EnableReflection: true,
				},
				Storage: sharedconfig.StorageConfig{
					Type: "memory",
				},
				Client: ClientConfig{
					Type: "local",
				},
				Security: sharedconfig.SecurityConfig{
					EnableAuth: false,
				},
			},
			wantErr: false,
		},
		{
			name: "missing HTTP port",
			config: &Config{
				HTTP: sharedconfig.HTTPConfig{
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				GRPC: sharedconfig.GRPCConfig{
					Port:             "9090",
					MaxRecvMsgSize:   4 * 1024 * 1024,
					MaxSendMsgSize:   4 * 1024 * 1024,
					EnableReflection: true,
				},
				Storage: sharedconfig.StorageConfig{Type: "memory"},
				Client:  ClientConfig{Type: "local"},
			},
			wantErr: true,
			errMsg:  "invalid port format",
		},
		{
			name: "port conflict",
			config: &Config{
				HTTP: sharedconfig.HTTPConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				GRPC: sharedconfig.GRPCConfig{
					Port:             "8080", // Same as HTTP
					MaxRecvMsgSize:   4 * 1024 * 1024,
					MaxSendMsgSize:   4 * 1024 * 1024,
					EnableReflection: true,
				},
				Storage: sharedconfig.StorageConfig{Type: "memory"},
				Client:  ClientConfig{Type: "local"},
			},
			wantErr: true,
			errMsg:  "HTTP and gRPC ports cannot be the same",
		},
		{
			name: "invalid storage type",
			config: &Config{
				HTTP: sharedconfig.HTTPConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
				GRPC: sharedconfig.GRPCConfig{
					Port:             "9090",
					MaxRecvMsgSize:   4 * 1024 * 1024,
					MaxSendMsgSize:   4 * 1024 * 1024,
					EnableReflection: true,
				},
				Storage: sharedconfig.StorageConfig{Type: "invalid"},
				Client:  ClientConfig{Type: "local"},
			},
			wantErr: true,
			errMsg:  "invalid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Config.Validate() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	cfg := &Config{
		HTTP: sharedconfig.HTTPConfig{
			Port: "8080",
			Host: "localhost",
		},
		GRPC: sharedconfig.GRPCConfig{
			Port: "9090",
			Host: "localhost",
		},
		Storage: sharedconfig.StorageConfig{
			Type:    "file",
			DataDir: "/tmp/data",
		},
		Client: ClientConfig{
			Type:      "local",
			ServerURL: "http://localhost:8080",
		},
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"http.port", "8080"},
		{"http.host", "localhost"},
		{"grpc.port", "9090"},
		{"grpc.host", "localhost"},
		{"storage.type", "file"},
		{"storage.data_dir", "/tmp/data"},
		{"client.type", "local"},
		{"client.server_url", "http://localhost:8080"},
		{"unknown.key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := cfg.GetString(tt.key)
			if result != tt.expected {
				t.Errorf("GetString(%s) = %s, want %s", tt.key, result, tt.expected)
			}
		})
	}
}
