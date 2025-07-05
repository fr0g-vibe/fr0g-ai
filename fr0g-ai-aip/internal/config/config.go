package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// GetString gets a string value from environment or returns default
func (c *Config) GetString(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}

// Config holds all application configuration
type Config struct {
	// Server configuration
	HTTP HTTPConfig `yaml:"http"`
	GRPC GRPCConfig `yaml:"grpc"`
	
	// Storage configuration
	Storage StorageConfig `yaml:"storage"`
	
	// Client configuration
	Client ClientConfig `yaml:"client"`
	
	// Security configuration
	Security SecurityConfig `yaml:"security"`
	
	// Logging configuration
	Logging LoggingConfig `yaml:"logging"`
}

type HTTPConfig struct {
	Port            string        `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	EnableTLS       bool          `yaml:"enable_tls"`
	CertFile        string        `yaml:"cert_file"`
	KeyFile         string        `yaml:"key_file"`
}

type GRPCConfig struct {
	Port            string        `yaml:"port"`
	MaxRecvMsgSize  int           `yaml:"max_recv_msg_size"`
	MaxSendMsgSize  int           `yaml:"max_send_msg_size"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
	EnableTLS       bool          `yaml:"enable_tls"`
	CertFile        string        `yaml:"cert_file"`
	KeyFile         string        `yaml:"key_file"`
}

type StorageConfig struct {
	Type    string `yaml:"type"` // memory, file
	DataDir string `yaml:"data_dir"`
}

type ClientConfig struct {
	Type      string `yaml:"type"`       // local, rest, grpc
	ServerURL string `yaml:"server_url"`
	Timeout   time.Duration `yaml:"timeout"`
}

type SecurityConfig struct {
	EnableAuth bool   `yaml:"enable_auth"`
	APIKey     string `yaml:"api_key"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"` // json, text
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	config := &Config{
		HTTP: HTTPConfig{
			Port:            getEnv("FR0G_HTTP_PORT", "8080"),
			ReadTimeout:     getDurationEnv("FR0G_HTTP_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getDurationEnv("FR0G_HTTP_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getDurationEnv("FR0G_HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			EnableTLS:       getBoolEnv("FR0G_HTTP_ENABLE_TLS", false),
			CertFile:        getEnv("FR0G_HTTP_CERT_FILE", ""),
			KeyFile:         getEnv("FR0G_HTTP_KEY_FILE", ""),
		},
		GRPC: GRPCConfig{
			Port:              getEnv("FR0G_GRPC_PORT", "9090"),
			MaxRecvMsgSize:    getIntEnv("FR0G_GRPC_MAX_RECV_MSG_SIZE", 4*1024*1024), // 4MB
			MaxSendMsgSize:    getIntEnv("FR0G_GRPC_MAX_SEND_MSG_SIZE", 4*1024*1024), // 4MB
			ConnectionTimeout: getDurationEnv("FR0G_GRPC_CONNECTION_TIMEOUT", 5*time.Second),
			EnableTLS:         getBoolEnv("FR0G_GRPC_ENABLE_TLS", false),
			CertFile:          getEnv("FR0G_GRPC_CERT_FILE", ""),
			KeyFile:           getEnv("FR0G_GRPC_KEY_FILE", ""),
		},
		Storage: StorageConfig{
			Type:    getEnv("FR0G_STORAGE_TYPE", "file"),
			DataDir: getEnv("FR0G_DATA_DIR", "./data"),
		},
		Client: ClientConfig{
			Type:      getEnv("FR0G_CLIENT_TYPE", "grpc"),
			ServerURL: getEnv("FR0G_SERVER_URL", "localhost:9090"),
			Timeout:   getDurationEnv("FR0G_CLIENT_TIMEOUT", 30*time.Second),
		},
		Security: SecurityConfig{
			EnableAuth: getBoolEnv("FR0G_ENABLE_AUTH", false),
			APIKey:     getEnv("FR0G_API_KEY", ""),
		},
		Logging: LoggingConfig{
			Level:  getEnv("FR0G_LOG_LEVEL", "info"),
			Format: getEnv("FR0G_LOG_FORMAT", "text"),
		},
	}
	
	// Expand relative paths
	if !filepath.IsAbs(config.Storage.DataDir) {
		if abs, err := filepath.Abs(config.Storage.DataDir); err == nil {
			config.Storage.DataDir = abs
		}
	}
	
	return config
}

// Validate validates the entire configuration using shared validation
func (c *Config) Validate() error {
	var errors sharedconfig.ValidationErrors
	
	// Validate HTTP port
	if err := sharedconfig.ValidatePort(c.HTTP.Port, "http.port"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate gRPC port
	if err := sharedconfig.ValidatePort(c.GRPC.Port, "grpc.port"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate storage type
	validStorageTypes := []string{"memory", "file"}
	if err := sharedconfig.ValidateEnum(c.Storage.Type, validStorageTypes, "storage.type"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate data directory if using file storage
	if c.Storage.Type == "file" {
		if err := sharedconfig.ValidateRequired(c.Storage.DataDir, "storage.data_dir"); err != nil {
			errors = append(errors, *err)
		}
		if err := sharedconfig.ValidateDirectoryPath(c.Storage.DataDir, "storage.data_dir"); err != nil {
			errors = append(errors, *err)
		}
	}
	
	// Validate client type
	validClientTypes := []string{"local", "rest", "grpc"}
	if err := sharedconfig.ValidateEnum(c.Client.Type, validClientTypes, "client.type"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate timeout
	if err := sharedconfig.ValidateTimeout(c.Client.Timeout, "client.timeout"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate API key if auth is enabled
	if c.Security.EnableAuth {
		if err := sharedconfig.ValidateRequired(c.Security.APIKey, "security.api_key"); err != nil {
			errors = append(errors, *err)
		}
		if err := sharedconfig.ValidateAPIKey(c.Security.APIKey, "security.api_key"); err != nil {
			errors = append(errors, *err)
		}
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
