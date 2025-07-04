package api

import (
	"strings"
	"testing"
)

func TestValidateChatCompletionRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *ChatCompletionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []Message{
					{Role: "user", Content: "Hello, world!"},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "request cannot be nil",
		},
		{
			name: "missing model",
			req: &ChatCompletionRequest{
				Messages: []Message{
					{Role: "user", Content: "Hello, world!"},
				},
			},
			wantErr: true,
			errMsg:  "model is required",
		},
		{
			name: "no messages",
			req: &ChatCompletionRequest{
				Model:    "gpt-3.5-turbo",
				Messages: []Message{},
			},
			wantErr: true,
			errMsg:  "at least one message is required",
		},
		{
			name: "too many messages",
			req: &ChatCompletionRequest{
				Model:    "gpt-3.5-turbo",
				Messages: make([]Message, 101), // Over limit
			},
			wantErr: true,
			errMsg:  "too many messages",
		},
		{
			name: "invalid temperature",
			req: &ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []Message{
					{Role: "user", Content: "Hello, world!"},
				},
				Temperature: func() *float64 { f := 3.0; return &f }(), // Over limit
			},
			wantErr: true,
			errMsg:  "temperature must be between 0 and 2",
		},
		{
			name: "invalid max tokens",
			req: &ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []Message{
					{Role: "user", Content: "Hello, world!"},
				},
				MaxTokens: func() *int32 { i := int32(-1); return &i }(), // Negative
			},
			wantErr: true,
			errMsg:  "max_tokens must be between 1 and 32000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateChatCompletionRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChatCompletionRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateChatCompletionRequest() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		content string
		wantErr bool
		errMsg  string
	}{
		{"valid message", "user", "Hello, world!", false, ""},
		{"empty role", "", "Hello, world!", true, "role is required"},
		{"empty content", "user", "", true, "content is required"},
		{"whitespace only content", "user", "   ", true, "content cannot be empty or whitespace only"},
		{"invalid role", "invalid", "Hello, world!", true, "invalid role"},
		{"content too long", "user", strings.Repeat("a", 32001), true, "content too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMessage(tt.role, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateMessage() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateModel(t *testing.T) {
	tests := []struct {
		name    string
		model   string
		wantErr bool
		errMsg  string
	}{
		{"valid supported model", "gpt-3.5-turbo", false, ""},
		{"valid custom model", "custom-model-v1", false, ""},
		{"empty model", "", true, "model cannot be empty"},
		{"invalid characters", "model with spaces", true, "invalid model name format"},
		{"special characters", "model@#$", true, "invalid model name format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModel(tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateModel() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidatePersonaPrompt(t *testing.T) {
	tests := []struct {
		name   string
		prompt *string
		wantErr bool
		errMsg string
	}{
		{"nil prompt", nil, false, ""},
		{"valid prompt", func() *string { s := "You are a helpful assistant"; return &s }(), false, ""},
		{"empty prompt", func() *string { s := ""; return &s }(), true, "persona prompt cannot be empty"},
		{"whitespace only", func() *string { s := "   "; return &s }(), true, "persona prompt cannot be empty"},
		{"too long prompt", func() *string { s := strings.Repeat("a", 8001); return &s }(), true, "persona prompt too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePersonaPrompt(tt.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePersonaPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidatePersonaPrompt() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateRequestSize(t *testing.T) {
	tests := []struct {
		name    string
		req     *ChatCompletionRequest
		wantErr bool
	}{
		{
			name: "small request",
			req: &ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "large request",
			req: &ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []Message{
					{Role: "user", Content: strings.Repeat("a", 100*1024)}, // 100KB content
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequestSize(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequestSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConversationFlow(t *testing.T) {
	tests := []struct {
		name     string
		messages []Message
		wantErr  bool
	}{
		{
			name: "valid conversation",
			messages: []Message{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
			},
			wantErr: false,
		},
		{
			name:     "empty messages",
			messages: []Message{},
			wantErr:  true,
		},
		{
			name: "single user message",
			messages: []Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr: false,
		},
		{
			name: "system message first",
			messages: []Message{
				{Role: "system", Content: "You are a helpful assistant"},
				{Role: "user", Content: "Hello"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConversationFlow(tt.messages)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConversationFlow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"user role", "user", true},
		{"assistant role", "assistant", true},
		{"system role", "system", true},
		{"function role", "function", true},
		{"invalid role", "invalid", false},
		{"empty role", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidRole(tt.role); got != tt.want {
				t.Errorf("isValidRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
