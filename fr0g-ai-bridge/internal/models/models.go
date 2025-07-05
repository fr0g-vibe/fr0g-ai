package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ChatMessage represents a single chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model         string        `json:"model"`
	Messages      []ChatMessage `json:"messages"`
	PersonaPrompt string        `json:"persona_prompt,omitempty"`
	Temperature   *float64      `json:"temperature,omitempty"`
	MaxTokens     *int          `json:"max_tokens,omitempty"`
	Stream        *bool         `json:"stream,omitempty"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse represents a chat completion response
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string                 `json:"status"`
	Time    time.Time              `json:"time"`
	Version string                 `json:"version"`
	Error   string                 `json:"error,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Validate validates the health response format comprehensively
func (hr *HealthResponse) Validate() error {
	var errors []string
	
	// Validate required status field
	if hr.Status == "" {
		errors = append(errors, "status field is required")
	} else {
		validStatuses := []string{"healthy", "unhealthy", "degraded"}
		isValidStatus := false
		for _, validStatus := range validStatuses {
			if hr.Status == validStatus {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			errors = append(errors, fmt.Sprintf("status must be one of: %v, got: %s", validStatuses, hr.Status))
		}
	}
	
	// Validate required time field
	if hr.Time.IsZero() {
		errors = append(errors, "time field is required and cannot be zero")
	} else {
		// Ensure time is not in the future (with 1 minute tolerance)
		if hr.Time.After(time.Now().Add(time.Minute)) {
			errors = append(errors, "time cannot be in the future")
		}
		// Ensure time is not too old (24 hours)
		if hr.Time.Before(time.Now().Add(-24 * time.Hour)) {
			errors = append(errors, "time cannot be older than 24 hours")
		}
	}
	
	// Validate required version field
	if hr.Version == "" {
		errors = append(errors, "version field is required")
	} else {
		// Basic semantic version validation
		if len(hr.Version) > 50 {
			errors = append(errors, "version field too long (max 50 characters)")
		}
	}
	
	// Validate error field constraints
	if hr.Error != "" {
		if len(hr.Error) > 1000 {
			errors = append(errors, "error field too long (max 1000 characters)")
		}
		// If error is present, status should be unhealthy or degraded
		if hr.Status == "healthy" {
			errors = append(errors, "error field should not be present when status is healthy")
		}
	}
	
	// Validate details field
	if hr.Details != nil {
		if len(hr.Details) > 20 {
			errors = append(errors, "details field has too many entries (max 20)")
		}
		// Check for sensitive information in details
		for key, value := range hr.Details {
			if len(key) > 100 {
				errors = append(errors, fmt.Sprintf("details key '%s' too long (max 100 characters)", key))
			}
			if valueStr, ok := value.(string); ok && len(valueStr) > 500 {
				errors = append(errors, fmt.Sprintf("details value for key '%s' too long (max 500 characters)", key))
			}
		}
	}
	
	// Return combined errors
	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

// ValidateForStatusCode validates the health response and returns appropriate HTTP status code
func (hr *HealthResponse) ValidateForStatusCode() (int, error) {
	if err := hr.Validate(); err != nil {
		return http.StatusInternalServerError, err
	}
	
	switch hr.Status {
	case "healthy":
		return http.StatusOK, nil
	case "degraded":
		return http.StatusOK, nil // Still operational but with issues
	case "unhealthy":
		return http.StatusServiceUnavailable, nil
	default:
		return http.StatusInternalServerError, fmt.Errorf("unknown status: %s", hr.Status)
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
