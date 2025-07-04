package config

import (
	"fmt"
	"net"
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
