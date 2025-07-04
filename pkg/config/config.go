package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// BaseConfig provides common configuration functionality
type BaseConfig interface {
	Validate() error
	LoadFromFile(path string) error
	LoadFromEnv() error
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
