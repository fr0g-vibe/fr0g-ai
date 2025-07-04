package config

import (
	"strings"
	"testing"
	"time"
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
				HTTP: HTTPConfig{
					Port:            "8080",
					ReadTimeout:     30 * time.Second,
					WriteTimeout:    30 * time.Second,
					ShutdownTimeout: 5 * time.Second,
				},
				GRPC: GRPCConfig{
					Port:              "9090",
					MaxRecvMsgSize:    4 * 1024 * 1024,
					MaxSendMsgSize:    4 * 1024 * 1024,
					ConnectionTimeout: 10 * time.Second,
				},
				Storage: StorageConfig{
					Type: "memory",
				},
				Client: ClientConfig{
					Type:    "local",
					Timeout: 30 * time.Second,
				},
				Security: SecurityConfig{
					EnableAuth: false,
				},
			},
			wantErr: false,
		},
		{
			name: "missing HTTP port",
			config: &Config{
				HTTP: HTTPConfig{
					ReadTimeout:     30 * time.Second,
					WriteTimeout:    30 * time.Second,
					ShutdownTimeout: 5 * time.Second,
				},
				GRPC: GRPCConfig{
					Port:              "9090",
					MaxRecvMsgSize:    4 * 1024 * 1024,
					MaxSendMsgSize:    4 * 1024 * 1024,
					ConnectionTimeout: 10 * time.Second,
				},
				Storage: StorageConfig{Type: "memory"},
				Client:  ClientConfig{Type: "local", Timeout: 30 * time.Second},
			},
			wantErr: true,
			errMsg:  "http.port: port is required",
		},
		{
			name: "port conflict",
			config: &Config{
				HTTP: HTTPConfig{
					Port:            "8080",
					ReadTimeout:     30 * time.Second,
					WriteTimeout:    30 * time.Second,
					ShutdownTimeout: 5 * time.Second,
				},
				GRPC: GRPCConfig{
					Port:              "8080", // Same as HTTP
					MaxRecvMsgSize:    4 * 1024 * 1024,
					MaxSendMsgSize:    4 * 1024 * 1024,
					ConnectionTimeout: 10 * time.Second,
				},
				Storage: StorageConfig{Type: "memory"},
				Client:  ClientConfig{Type: "local", Timeout: 30 * time.Second},
			},
			wantErr: true,
			errMsg:  "ports: HTTP and gRPC ports cannot be the same",
		},
		{
			name: "invalid storage type",
			config: &Config{
				HTTP: HTTPConfig{
					Port:            "8080",
					ReadTimeout:     30 * time.Second,
					WriteTimeout:    30 * time.Second,
					ShutdownTimeout: 5 * time.Second,
				},
				GRPC: GRPCConfig{
					Port:              "9090",
					MaxRecvMsgSize:    4 * 1024 * 1024,
					MaxSendMsgSize:    4 * 1024 * 1024,
					ConnectionTimeout: 10 * time.Second,
				},
				Storage: StorageConfig{Type: "invalid"},
				Client:  ClientConfig{Type: "local", Timeout: 30 * time.Second},
			},
			wantErr: true,
			errMsg:  "storage.type: invalid storage type",
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

func TestValidateNetworkAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{"valid address", "localhost:8080", false},
		{"valid IP address", "127.0.0.1:9090", false},
		{"invalid port", "localhost:99999", true},
		{"missing port", "localhost", true},
		{"empty host", ":8080", true},
		{"invalid format", "invalid-address", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNetworkAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNetworkAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		wantErr bool
	}{
		{"valid timeout", 30 * time.Second, false},
		{"zero timeout", 0, true},
		{"negative timeout", -5 * time.Second, true},
		{"too long timeout", 25 * time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeout(tt.timeout, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
