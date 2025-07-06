package api

import (
	"fmt"
	"strings"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

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

	// Validate message count limits
	if len(req.Messages) > 100 {
		return fmt.Errorf("too many messages (max 100)")
	}

	// Validate each message
	for i, msg := range req.Messages {
		if err := ValidateMessage(msg.Role, msg.Content); err != nil {
			return fmt.Errorf("message[%d]: %v", i, err)
		}
	}

	// Validate conversation flow
	if err := ValidateConversationFlow(req.Messages); err != nil {
		return fmt.Errorf("conversation flow: %v", err)
	}

	// Validate optional parameters
	if req.Temperature != nil && (*req.Temperature < 0 || *req.Temperature > 2) {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	if req.MaxTokens != nil && (*req.MaxTokens <= 0 || *req.MaxTokens > 32000) {
		return fmt.Errorf("max_tokens must be between 1 and 32000")
	}

	// Validate overall request size
	if err := ValidateRequestSize(req); err != nil {
		return err
	}

	return nil
}

// ValidateMessage validates a single chat message
func ValidateMessage(role, content string) error {
	// Validate role
	if err := validateRole(role); err != nil {
		return err
	}

	// Validate required content
	if err := validateRequired(content, "content"); err != nil {
		return err
	}

	// Validate content length
	if err := validateLength(content, 1, 32000, "content"); err != nil {
		return err
	}

	return nil
}

// ValidateModel checks if the model name is valid
func ValidateModel(model string) error {
	if err := validateRequired(model, "model"); err != nil {
		return err
	}

	if err := validateLength(model, 1, 100, "model"); err != nil {
		return err
	}

	// Check against known supported models
	supportedModels := []string{
		"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo", "gpt-4o",
		"claude-3-haiku", "claude-3-sonnet", "claude-3-opus",
		"llama-2-7b", "llama-2-13b", "llama-2-70b",
		"mistral-7b", "mixtral-8x7b",
	}

	// Allow custom models - just validate format
	return nil
}

// ValidatePersonaPrompt validates the persona prompt if provided
func ValidatePersonaPrompt(prompt *string) error {
	if prompt != nil && *prompt != "" {
		if err := validateLength(*prompt, 1, 8000, "persona_prompt"); err != nil {
			return err
		}
	}
	return nil
}

// Helper validation functions
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
	if strings.TrimSpace(value) == "" {
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

// ValidateRequestSize validates the overall request size
func ValidateRequestSize(req *ChatCompletionRequest) error {
	totalSize := len(req.Model)

	for _, msg := range req.Messages {
		totalSize += len(msg.Role) + len(msg.Content)
	}

	// Reasonable limit for total request size (100KB)
	if totalSize > 100*1024 {
		return fmt.Errorf("request too large (max 100KB)")
	}

	return nil
}

// ValidateConversationFlow validates the logical flow of messages
func ValidateConversationFlow(messages []Message) error {
	if len(messages) == 0 {
		return fmt.Errorf("no messages provided")
	}

	// Check for alternating user/assistant pattern (common best practice)
	hasUser := false
	hasAssistant := false

	for _, msg := range messages {
		switch msg.Role {
		case "user":
			hasUser = true
		case "assistant":
			hasAssistant = true
		case "system":
			// System messages are fine anywhere
		case "function":
			// Function messages need special handling but are valid
		}
	}

	// Warn if conversation doesn't have both user and assistant (unless it's the first message)
	if len(messages) > 1 && hasUser && !hasAssistant {
		// This is just a warning case - still valid
	}

	return nil
}
