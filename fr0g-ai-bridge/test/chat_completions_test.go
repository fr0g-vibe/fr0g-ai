package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/models"
)

// TestValidationFunctions tests the validation functions directly
func TestValidationFunctions(t *testing.T) {
	t.Run("ValidateRole", func(t *testing.T) {
		validRoles := []string{"user", "assistant", "system", "function"}
		for _, role := range validRoles {
			if err := validateRole(role); err != nil {
				t.Errorf("Expected role %s to be valid, got error: %v", role, err)
			}
		}
		
		invalidRoles := []string{"", "invalid", "bot", "human"}
		for _, role := range invalidRoles {
			if err := validateRole(role); err == nil {
				t.Errorf("Expected role %s to be invalid", role)
			}
		}
	})
	
	t.Run("ValidateRequired", func(t *testing.T) {
		if err := validateRequired("valid", "test"); err != nil {
			t.Errorf("Expected non-empty string to be valid, got: %v", err)
		}
		
		if err := validateRequired("", "test"); err == nil {
			t.Error("Expected empty string to be invalid")
		}
		
		if err := validateRequired("   ", "test"); err == nil {
			t.Error("Expected whitespace-only string to be invalid")
		}
	})
	
	t.Run("ValidateLength", func(t *testing.T) {
		if err := validateLength("hello", 1, 10, "test"); err != nil {
			t.Errorf("Expected valid length string to pass, got: %v", err)
		}
		
		if err := validateLength("", 1, 10, "test"); err == nil {
			t.Error("Expected too short string to fail")
		}
		
		if err := validateLength("this is way too long for the limit", 1, 10, "test"); err == nil {
			t.Error("Expected too long string to fail")
		}
	})
}

// Helper validation functions for testing
func validateRole(role string) error {
	if role == "" {
		return fmt.Errorf("role is required")
	}
	
	validRoles := []string{"user", "assistant", "system", "function"}
	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}
	return fmt.Errorf("invalid role: %s", role)
}

func validateRequired(value, fieldName string) error {
	if len(value) == 0 || len(value) == 0 {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func validateLength(value string, min, max int, fieldName string) error {
	length := len(value)
	if length < min {
		return fmt.Errorf("%s must be at least %d characters", fieldName, min)
	}
	if length > max {
		return fmt.Errorf("%s must be at most %d characters", fieldName, max)
	}
	return nil
}

func TestChatCompletionRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     models.ChatCompletionRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid request",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			expectError: false,
		},
		{
			name: "Missing model",
			request: models.ChatCompletionRequest{
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			expectError: true,
			errorMsg:    "model",
		},
		{
			name: "Empty messages",
			request: models.ChatCompletionRequest{
				Model:    "gpt-3.5-turbo",
				Messages: []models.ChatMessage{},
			},
			expectError: true,
			errorMsg:    "messages",
		},
		{
			name: "Invalid role",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "invalid", Content: "Hello"},
				},
			},
			expectError: true,
			errorMsg:    "role",
		},
		{
			name: "Empty content",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: ""},
				},
			},
			expectError: true,
			errorMsg:    "content",
		},
		{
			name: "Valid with persona prompt",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				PersonaPrompt: "You are a helpful assistant.",
			},
			expectError: false,
		},
		{
			name: "Valid with temperature",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: func() *float64 { v := 0.7; return &v }(),
			},
			expectError: false,
		},
		{
			name: "Invalid temperature too high",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: func() *float64 { v := 3.0; return &v }(),
			},
			expectError: true,
			errorMsg:    "temperature",
		},
		{
			name: "Valid max tokens",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxTokens: func() *int { v := 100; return &v }(),
			},
			expectError: false,
		},
		{
			name: "Invalid max tokens too low",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxTokens: func() *int { v := 0; return &v }(),
			},
			expectError: true,
			errorMsg:    "max_tokens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChatCompletionRequest(&tt.request)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', but got no error", tt.errorMsg)
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', but got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// validateChatCompletionRequest validates a chat completion request
func validateChatCompletionRequest(req *models.ChatCompletionRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	
	if len(req.Messages) == 0 {
		return fmt.Errorf("at least one message is required")
	}
	
	if len(req.Messages) > 100 {
		return fmt.Errorf("too many messages (max 100)")
	}
	
	for i, msg := range req.Messages {
		if err := validateRole(msg.Role); err != nil {
			return fmt.Errorf("message[%d]: %v", i, err)
		}
		if err := validateRequired(msg.Content, "content"); err != nil {
			return fmt.Errorf("message[%d]: %v", i, err)
		}
		if len(msg.Content) > 32000 {
			return fmt.Errorf("message[%d]: content too long (max 32000 characters)", i)
		}
	}
	
	if req.Temperature != nil && (*req.Temperature < 0 || *req.Temperature > 2) {
		return fmt.Errorf("temperature must be between 0 and 2")
	}
	
	if req.MaxTokens != nil && (*req.MaxTokens <= 0 || *req.MaxTokens > 32000) {
		return fmt.Errorf("max_tokens must be between 1 and 32000")
	}
	
	if len(req.PersonaPrompt) > 8000 {
		return fmt.Errorf("persona_prompt too long (max 8000 characters)")
	}
	
	return nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		func() bool {
			for i := 1; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())))
}

func TestHealthResponseValidation(t *testing.T) {
	tests := []struct {
		name        string
		response    models.HealthResponse
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid healthy response",
			response: models.HealthResponse{
				Status:  "healthy",
				Time:    time.Now(),
				Version: "1.0.0",
			},
			expectError: false,
		},
		{
			name: "Valid unhealthy response",
			response: models.HealthResponse{
				Status:  "unhealthy",
				Time:    time.Now(),
				Version: "1.0.0",
				Error:   "Service unavailable",
			},
			expectError: false,
		},
		{
			name: "Missing status",
			response: models.HealthResponse{
				Time:    time.Now(),
				Version: "1.0.0",
			},
			expectError: true,
			errorMsg:    "status",
		},
		{
			name: "Invalid status",
			response: models.HealthResponse{
				Status:  "invalid",
				Time:    time.Now(),
				Version: "1.0.0",
			},
			expectError: true,
			errorMsg:    "status",
		},
		{
			name: "Missing time",
			response: models.HealthResponse{
				Status:  "healthy",
				Version: "1.0.0",
			},
			expectError: true,
			errorMsg:    "time",
		},
		{
			name: "Missing version",
			response: models.HealthResponse{
				Status: "healthy",
				Time:   time.Now(),
			},
			expectError: true,
			errorMsg:    "version",
		},
		{
			name: "Error with healthy status",
			response: models.HealthResponse{
				Status:  "healthy",
				Time:    time.Now(),
				Version: "1.0.0",
				Error:   "Should not have error when healthy",
			},
			expectError: true,
			errorMsg:    "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.response.Validate()
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', but got no error", tt.errorMsg)
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', but got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestHealthResponseStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		response       models.HealthResponse
		expectedStatus int
	}{
		{
			name: "Healthy status",
			response: models.HealthResponse{
				Status:  "healthy",
				Time:    time.Now(),
				Version: "1.0.0",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Degraded status",
			response: models.HealthResponse{
				Status:  "degraded",
				Time:    time.Now(),
				Version: "1.0.0",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unhealthy status",
			response: models.HealthResponse{
				Status:  "unhealthy",
				Time:    time.Now(),
				Version: "1.0.0",
				Error:   "Service down",
			},
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, err := tt.response.ValidateForStatusCode()
			
			if err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
			
			if statusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, statusCode)
			}
		})
	}
}
