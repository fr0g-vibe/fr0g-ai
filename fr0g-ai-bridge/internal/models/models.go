package models

import (
	"fmt"
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

// Validate validates the health response format
func (hr *HealthResponse) Validate() error {
	if hr.Status == "" {
		return fmt.Errorf("status field is required")
	}
	
	validStatuses := []string{"healthy", "unhealthy", "degraded"}
	isValidStatus := false
	for _, validStatus := range validStatuses {
		if hr.Status == validStatus {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return fmt.Errorf("status must be one of: %v", validStatuses)
	}
	
	if hr.Time.IsZero() {
		return fmt.Errorf("time field is required")
	}
	
	if hr.Version == "" {
		return fmt.Errorf("version field is required")
	}
	
	return nil
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
