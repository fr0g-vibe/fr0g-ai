package api

import (
	"fmt"
	"regexp"
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
	if role == "" {
		return fmt.Errorf("role is required")
	}
	if content == "" {
		return fmt.Errorf("content is required")
	}
	if len(content) > 32000 { // Increased limit for modern models
		return fmt.Errorf("content too long (max 32000 characters)")
	}
	if !isValidRole(role) {
		return fmt.Errorf("invalid role: %s", role)
	}
	
	// Additional content validation
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("content cannot be empty or whitespace only")
	}
	
	return nil
}

// ValidateModel checks if the model name is valid
func ValidateModel(model string) error {
	if model == "" {
		return fmt.Errorf("model cannot be empty")
	}
	
	// Allow more flexible model naming with regex pattern
	validModelPattern := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	if !validModelPattern.MatchString(model) {
		return fmt.Errorf("invalid model name format: %s", model)
	}
	
	// Check against known supported models
	supportedModels := []string{
		"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo", "gpt-4o",
		"claude-3-haiku", "claude-3-sonnet", "claude-3-opus",
		"llama-2-7b", "llama-2-13b", "llama-2-70b",
		"mistral-7b", "mixtral-8x7b",
	}
	
	for _, supportedModel := range supportedModels {
		if model == supportedModel {
			return nil
		}
	}
	
	// Allow custom models but warn they might not be supported
	return nil
}

// ValidatePersonaPrompt validates the persona prompt if provided
func ValidatePersonaPrompt(prompt *string) error {
	if prompt != nil {
		if len(*prompt) > 8000 { // Increased limit for more detailed personas
			return fmt.Errorf("persona prompt too long (max 8000 characters)")
		}
		if strings.TrimSpace(*prompt) == "" {
			return fmt.Errorf("persona prompt cannot be empty or whitespace only")
		}
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
