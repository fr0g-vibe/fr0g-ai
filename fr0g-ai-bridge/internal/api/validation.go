package api

import (
	"fmt"
	"strings"
)

// isValidRole checks if the role is one of the allowed values
func isValidRole(role string) bool {
	validRoles := []string{"system", "user", "assistant", "function"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// ValidationError represents a validation error
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

// ChatCompletionRequest represents a chat completion request for validation
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int32    `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ValidateChatCompletionRequest validates a chat completion request
func ValidateChatCompletionRequest(req *ChatCompletionRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	
	if err := ValidateModel(req.Model); err != nil {
		return err
	}
	
	if len(req.Messages) == 0 {
		return fmt.Errorf("at least one message is required")
	}
	
	// Validate each message
	for i, msg := range req.Messages {
		if err := ValidateMessage(msg.Role, msg.Content); err != nil {
			return fmt.Errorf("message[%d]: %v", i, err)
		}
	}
	
	// Validate optional parameters
	if req.Temperature != nil && (*req.Temperature < 0 || *req.Temperature > 2) {
		return fmt.Errorf("temperature must be between 0 and 2")
	}
	
	if req.MaxTokens != nil && *req.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens must be positive")
	}
	
	return nil
}

// ValidateMessage validates a single chat message
func ValidateMessage(role, content string) error {
	if role == "" {
		return fmt.Errorf("role is required")
	}
	if content == "" {
		return fmt.Errorf("content is required")
	}
	if len(content) > 10000 {
		return fmt.Errorf("content too long (max 10000 characters)")
	}
	if !isValidRole(role) {
		return fmt.Errorf("invalid role: %s", role)
	}
	return nil
}

// ValidateModel checks if the model name is valid
func ValidateModel(model string) error {
	if model == "" {
		return fmt.Errorf("model cannot be empty")
	}
	
	validModels := []string{"gpt-3.5-turbo", "gpt-4", "claude-3", "llama-2"}
	for _, validModel := range validModels {
		if model == validModel {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported model: %s", model)
}

// ValidatePersonaPrompt validates the persona prompt if provided
func ValidatePersonaPrompt(prompt *string) error {
	if prompt != nil {
		if len(*prompt) > 5000 {
			return fmt.Errorf("persona prompt too long (max 5000 characters)")
		}
		if strings.TrimSpace(*prompt) == "" {
			return fmt.Errorf("persona prompt cannot be empty or whitespace only")
		}
	}
	return nil
}
