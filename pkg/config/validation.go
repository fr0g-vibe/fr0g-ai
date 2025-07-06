package config

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ves ValidationErrors) Error() string {
	var messages []string
	for _, ve := range ves {
		messages = append(messages, ve.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (ves ValidationErrors) HasErrors() bool {
	return len(ves) > 0
}

// Validator interface for configuration validation
type Validator interface {
	Validate() ValidationErrors
}

// ValidatePort validates a port number
func ValidatePort(port interface{}, fieldName string) *ValidationError {
	var portNum int

	switch p := port.(type) {
	case int:
		portNum = p
	case string:
		var err error
		portNum, err = strconv.Atoi(p)
		if err != nil {
			return &ValidationError{
				Field:   fieldName,
				Message: "invalid port format",
			}
		}
	default:
		return &ValidationError{
			Field:   fieldName,
			Message: "port must be int or string",
		}
	}

	if portNum <= 0 || portNum > 65535 {
		return &ValidationError{
			Field:   fieldName,
			Message: "port must be between 1 and 65535",
		}
	}

	return nil
}

// ValidateTimeout validates a timeout duration
func ValidateTimeout(timeout time.Duration, fieldName string) *ValidationError {
	if timeout <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "timeout must be positive",
		}
	}
	if timeout > 24*time.Hour {
		return &ValidationError{
			Field:   fieldName,
			Message: "timeout cannot exceed 24 hours",
		}
	}
	return nil
}

// ValidateRequired validates that a field is not empty
func ValidateRequired(value string, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "field is required",
		}
	}
	return nil
}

// ValidateAPIKey validates API key strength
func ValidateAPIKey(apiKey string, fieldName string) *ValidationError {
	if len(apiKey) < 16 {
		return &ValidationError{
			Field:   fieldName,
			Message: "API key must be at least 16 characters long",
		}
	}
	return nil
}

// ValidateModel validates model name format
func ValidateModel(model string) *ValidationError {
	if model == "" {
		return &ValidationError{
			Field:   "model",
			Message: "model cannot be empty",
		}
	}

	validModelPattern := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	if !validModelPattern.MatchString(model) {
		return &ValidationError{
			Field:   "model",
			Message: fmt.Sprintf("invalid model name format: %s", model),
		}
	}

	return nil
}

// ValidateNetworkAddress validates a network address
func ValidateNetworkAddress(address string, fieldName string) *ValidationError {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("invalid address format: %v", err),
		}
	}

	if host == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "host cannot be empty",
		}
	}

	if portErr := ValidatePort(port, fieldName+".port"); portErr != nil {
		return portErr
	}

	return nil
}

// ValidateRole checks if the role is one of the allowed values
func ValidateRole(role string) *ValidationError {
	validRoles := []string{"system", "user", "assistant", "function"}
	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}
	return &ValidationError{
		Field:   "role",
		Message: fmt.Sprintf("invalid role: %s", role),
	}
}

// Contains checks if a slice contains an item
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidateEnum validates that a value is in a list of allowed values
func ValidateEnum(value string, allowedValues []string, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "value is required",
		}
	}

	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}

	return &ValidationError{
		Field:   fieldName,
		Message: fmt.Sprintf("invalid value, must be one of: %s", strings.Join(allowedValues, ", ")),
	}
}

// ValidateLength validates string length
func ValidateLength(value string, min, max int, fieldName string) *ValidationError {
	length := len(value)
	if length < min {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at least %d characters", min),
		}
	}
	if length > max {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at most %d characters", max),
		}
	}
	return nil
}

// ValidateDirectoryPath validates that a directory path is valid
func ValidateDirectoryPath(path string, fieldName string) *ValidationError {
	if path == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "directory path is required",
		}
	}

	// Check if path exists or can be created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Try to create the directory
		if err := os.MkdirAll(path, 0755); err != nil {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("cannot create directory: %v", err),
			}
		}
	}

	return nil
}

// ValidateURL validates a URL format
func ValidateURL(url string, fieldName string) *ValidationError {
	if url == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "URL cannot be empty",
		}
	}

	// Basic URL validation - starts with http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return &ValidationError{
			Field:   fieldName,
			Message: "URL must start with http:// or https://",
		}
	}

	return nil
}

// ValidateRange validates that a numeric value is within a specified range
func ValidateRange(value float64, min, max float64, fieldName string) *ValidationError {
	if value < min || value > max {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("value must be between %.2f and %.2f", min, max),
		}
	}
	return nil
}

// ValidatePositive validates that a numeric value is positive
func ValidatePositive(value int, fieldName string) *ValidationError {
	if value <= 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "value must be positive",
		}
	}
	return nil
}

// ValidateStringSliceNotEmpty validates that a string slice is not empty
func ValidateStringSliceNotEmpty(slice []string, fieldName string) *ValidationError {
	if len(slice) == 0 {
		return &ValidationError{
			Field:   fieldName,
			Message: "at least one item is required",
		}
	}
	return nil
}

// ValidateFilePath validates that a file path exists and is accessible
func ValidateFilePath(path string, fieldName string) *ValidationError {
	if path == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "file path is required",
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("file does not exist: %s", path),
		}
	}

	return nil
}

// ValidatePersonaName validates persona name format and length
func ValidatePersonaName(name string, fieldName string) *ValidationError {
	if err := ValidateRequired(name, fieldName); err != nil {
		return err
	}

	if err := ValidateLength(name, 1, 100, fieldName); err != nil {
		return err
	}

	// Check for valid characters (alphanumeric, spaces, hyphens, underscores)
	validNamePattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validNamePattern.MatchString(name) {
		return &ValidationError{
			Field:   fieldName,
			Message: "name can only contain letters, numbers, spaces, hyphens, and underscores",
		}
	}

	return nil
}

// ValidatePersonaTopic validates persona topic format
func ValidatePersonaTopic(topic string, fieldName string) *ValidationError {
	if err := ValidateRequired(topic, fieldName); err != nil {
		return err
	}

	if err := ValidateLength(topic, 1, 200, fieldName); err != nil {
		return err
	}

	return nil
}

// ValidatePersonaPrompt validates persona prompt content
func ValidatePersonaPrompt(prompt string, fieldName string) *ValidationError {
	if err := ValidateRequired(prompt, fieldName); err != nil {
		return err
	}

	if err := ValidateLength(prompt, 10, 5000, fieldName); err != nil {
		return err
	}

	return nil
}

// ValidatePersonaContext validates persona context map
func ValidatePersonaContext(context map[string]string, fieldName string) *ValidationError {
	if context == nil {
		return nil // Context is optional
	}

	for key, value := range context {
		if key == "" {
			return &ValidationError{
				Field:   fieldName,
				Message: "context keys cannot be empty",
			}
		}

		if len(key) > 50 {
			return &ValidationError{
				Field:   fieldName,
				Message: "context keys must be 50 characters or less",
			}
		}

		if len(value) > 500 {
			return &ValidationError{
				Field:   fieldName,
				Message: "context values must be 500 characters or less",
			}
		}
	}

	return nil
}

// ValidateIdentityName validates identity name format
func ValidateIdentityName(name string, fieldName string) *ValidationError {
	return ValidatePersonaName(name, fieldName) // Same rules as persona names
}

// ValidateIdentityBackground validates identity background content
func ValidateIdentityBackground(background string, fieldName string) *ValidationError {
	if background == "" {
		return nil // Background is optional
	}

	if err := ValidateLength(background, 0, 2000, fieldName); err != nil {
		return err
	}

	return nil
}

// ValidateLogLevel validates log level values
func ValidateLogLevel(level string, fieldName string) *ValidationError {
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	return ValidateEnum(level, validLevels, fieldName)
}

// ValidateLogFormat validates log format values
func ValidateLogFormat(format string, fieldName string) *ValidationError {
	validFormats := []string{"json", "text"}
	return ValidateEnum(format, validFormats, fieldName)
}
