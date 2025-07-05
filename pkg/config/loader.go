package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)


// LoaderOptions configures how configuration is loaded
type LoaderOptions struct {
	ConfigPath   string
	EnvPrefix    string
	EnvFilePaths []string
}

// Loader handles configuration loading from multiple sources
type Loader struct {
	options LoaderOptions
}

// NewLoader creates a new configuration loader
func NewLoader(options LoaderOptions) *Loader {
	if len(options.EnvFilePaths) == 0 {
		options.EnvFilePaths = []string{".env", "../.env", "../../.env"}
	}
	return &Loader{options: options}
}

// LoadEnvFile loads environment variables from a .env file
func (l *Loader) LoadEnvFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // File doesn't exist, skip silently
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}

	return nil
}

// LoadFromFile loads configuration from YAML file
func (l *Loader) LoadFromFile(config interface{}) error {
	if l.options.ConfigPath == "" {
		return nil
	}

	if _, err := os.Stat(l.options.ConfigPath); os.IsNotExist(err) {
		return nil // File doesn't exist, skip silently
	}

	data, err := os.ReadFile(l.options.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// LoadEnvFiles loads all configured .env files
func (l *Loader) LoadEnvFiles() error {
	for _, envPath := range l.options.EnvFilePaths {
		if err := l.LoadEnvFile(envPath); err != nil {
			fmt.Printf("Warning: failed to load %s: %v\n", envPath, err)
		}
	}
	return nil
}

// GetEnvString gets string from environment with optional prefix
func (l *Loader) GetEnvString(key, defaultValue string) string {
	if l.options.EnvPrefix != "" {
		key = l.options.EnvPrefix + "_" + key
	}
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


// GetEnvInt gets integer from environment with optional prefix
func (l *Loader) GetEnvInt(key string, defaultValue int) int {
	if l.options.EnvPrefix != "" {
		key = l.options.EnvPrefix + "_" + key
	}
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// GetEnvBool gets boolean from environment with optional prefix
func (l *Loader) GetEnvBool(key string, defaultValue bool) bool {
	if l.options.EnvPrefix != "" {
		key = l.options.EnvPrefix + "_" + key
	}
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

// GetEnvDuration gets duration from environment with optional prefix
func (l *Loader) GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if l.options.EnvPrefix != "" {
		key = l.options.EnvPrefix + "_" + key
	}
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
