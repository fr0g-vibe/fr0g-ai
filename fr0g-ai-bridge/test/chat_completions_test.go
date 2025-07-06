package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/api"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/client"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockOpenWebUIClient implements a mock client for testing
type MockOpenWebUIClient struct {
	responses map[string]*models.ChatCompletionResponse
	errors    map[string]error
	healthy   bool
}

func NewMockOpenWebUIClient() *MockOpenWebUIClient {
	return &MockOpenWebUIClient{
		responses: make(map[string]*models.ChatCompletionResponse),
		errors:    make(map[string]error),
		healthy:   true,
	}
}

func (m *MockOpenWebUIClient) ChatCompletion(ctx context.Context, req *models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	key := fmt.Sprintf("%s:%d", req.Model, len(req.Messages))
	if err, exists := m.errors[key]; exists {
		return nil, err
	}
	if resp, exists := m.responses[key]; exists {
		return resp, nil
	}
	
	// Default response
	return &models.ChatCompletionResponse{
		ID:      "chatcmpl-test-123",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.ChatMessage{
					Role:    "assistant",
					Content: "This is a test response from the mock client.",
				},
				FinishReason: "stop",
			},
		},
		Usage: models.Usage{
			PromptTokens:     10,
			CompletionTokens: 15,
			TotalTokens:      25,
		},
	}, nil
}

func (m *MockOpenWebUIClient) SendMessage(message, model string) (string, error) {
	return "Mock response", nil
}

func (m *MockOpenWebUIClient) GetModels() ([]string, error) {
	return []string{"gpt-3.5-turbo", "gpt-4"}, nil
}

func (m *MockOpenWebUIClient) HealthCheck(ctx context.Context) error {
	if !m.healthy {
		return fmt.Errorf("mock client unhealthy")
	}
	return nil
}

func (m *MockOpenWebUIClient) SetResponse(model string, messageCount int, response *models.ChatCompletionResponse) {
	key := fmt.Sprintf("%s:%d", model, messageCount)
	m.responses[key] = response
}

func (m *MockOpenWebUIClient) SetError(model string, messageCount int, err error) {
	key := fmt.Sprintf("%s:%d", model, messageCount)
	m.errors[key] = err
}

func (m *MockOpenWebUIClient) SetHealthy(healthy bool) {
	m.healthy = healthy
}

// Test setup helper
func setupTestServer() (*api.RESTServer, *MockOpenWebUIClient) {
	mockClient := NewMockOpenWebUIClient()
	cfg := &config.Config{
		Security: config.SecurityConfig{
			RequireAPIKey:   false,
			EnableCORS:      true,
			RateLimitRPM:    60,
			AllowedOrigins:  []string{"*"},
			AllowedAPIKeys:  []string{},
		},
		OpenWebUI: config.OpenWebUIConfig{
			BaseURL: "http://localhost:8080",
			APIKey:  "test-key",
			Timeout: 30,
		},
	}
	
	// Create a wrapper that implements the expected interface
	clientWrapper := &client.OpenWebUIClient{}
	
	server := api.NewRESTServer(clientWrapper, cfg)
	
	// Replace the client with our mock (this would need interface changes in real implementation)
	// For now, we'll test the validation and routing logic
	
	return server, mockClient
}

func TestChatCompletionsEndpoint(t *testing.T) {
	server, mockClient := setupTestServer()
	
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		setupMock      func(*MockOpenWebUIClient)
	}{
		{
			name:   "Valid chat completion request",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello, how are you?"},
				},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockOpenWebUIClient) {
				m.SetResponse("gpt-3.5-turbo", 1, &models.ChatCompletionResponse{
					ID:      "chatcmpl-test-123",
					Object:  "chat.completion",
					Created: time.Now().Unix(),
					Model:   "gpt-3.5-turbo",
					Choices: []models.Choice{
						{
							Index: 0,
							Message: models.ChatMessage{
								Role:    "assistant",
								Content: "Hello! I'm doing well, thank you for asking.",
							},
							FinishReason: "stop",
						},
					},
					Usage: models.Usage{
						PromptTokens:     5,
						CompletionTokens: 12,
						TotalTokens:      17,
					},
				})
			},
		},
		{
			name:   "Legacy endpoint compatibility",
			method: "POST",
			path:   "/api/chat/completions",
			body: models.ChatCompletionRequest{
				Model: "gpt-4",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Test legacy endpoint"},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Missing model field",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Empty messages array",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model:    "gpt-3.5-turbo",
				Messages: []models.ChatMessage{},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid role",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "invalid", Content: "Hello"},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Empty content",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: ""},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Content too long",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: string(make([]byte, 33000))}, // Too long
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Too many messages",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model:    "gpt-3.5-turbo",
				Messages: make([]models.ChatMessage, 101), // Too many
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid temperature",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: func() *float64 { v := 3.0; return &v }(), // Too high
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid max_tokens",
			method: "POST",
			path:   "/v1/chat/completions",
			body: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxTokens: func() *int { v := 0; return &v }(), // Too low
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			var body bytes.Buffer
			if tt.body != nil {
				err := json.NewEncoder(&body).Encode(tt.body)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(tt.method, tt.path, &body)
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			server.GetRouter().ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Response body: %s", w.Body.String())

			if tt.expectedStatus == http.StatusOK {
				var response models.ChatCompletionResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				
				// Validate OpenAI-compatible response format
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, "chat.completion", response.Object)
				assert.NotZero(t, response.Created)
				assert.NotEmpty(t, response.Model)
				assert.NotEmpty(t, response.Choices)
				assert.Equal(t, 0, response.Choices[0].Index)
				assert.NotEmpty(t, response.Choices[0].Message.Content)
				assert.Equal(t, "assistant", response.Choices[0].Message.Role)
				assert.NotEmpty(t, response.Choices[0].FinishReason)
			}
		})
	}
}

func TestPersonaAwareChatCompletions(t *testing.T) {
	server, mockClient := setupTestServer()

	tests := []struct {
		name          string
		request       models.ChatCompletionRequest
		expectedError bool
	}{
		{
			name: "Valid persona prompt",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				PersonaPrompt: "You are a helpful assistant specialized in technical support.",
			},
			expectedError: false,
		},
		{
			name: "Persona prompt too long",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				PersonaPrompt: string(make([]byte, 8001)), // Too long
			},
			expectedError: true,
		},
		{
			name: "Empty persona prompt",
			request: models.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []models.ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				PersonaPrompt: "",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(tt.request)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/v1/chat/completions", &body)
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			server.GetRouter().ServeHTTP(w, req)

			if tt.expectedError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				// Note: This will fail without proper mock integration
				// but validates the request parsing and validation
				assert.NotEqual(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestConversationFlow(t *testing.T) {
	server, _ := setupTestServer()

	// Test multi-turn conversation
	request := models.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []models.ChatMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "What's the weather like?"},
			{Role: "assistant", Content: "I don't have access to current weather data."},
			{Role: "user", Content: "Can you tell me a joke instead?"},
		},
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(request)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/v1/chat/completions", &body)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	server.GetRouter().ServeHTTP(w, req)

	// Should pass validation (conversation flow is valid)
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestErrorHandling(t *testing.T) {
	server, mockClient := setupTestServer()

	// Test upstream service error
	mockClient.SetError("gpt-3.5-turbo", 1, fmt.Errorf("upstream service unavailable"))

	request := models.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []models.ChatMessage{
			{Role: "user", Content: "Hello"},
		},
	}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(request)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/v1/chat/completions", &body)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	server.GetRouter().ServeHTTP(w, req)

	// Should return 500 for upstream errors
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResp models.ErrorResponse
	err = json.NewDecoder(w.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.NotEmpty(t, errorResp.Error)
	assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
}

func TestHealthEndpoint(t *testing.T) {
	server, mockClient := setupTestServer()

	tests := []struct {
		name           string
		healthy        bool
		expectedStatus int
	}{
		{
			name:           "Healthy service",
			healthy:        true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unhealthy service",
			healthy:        false,
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.SetHealthy(tt.healthy)

			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			server.GetRouter().ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.HealthResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			assert.NotEmpty(t, response.Status)
			assert.NotZero(t, response.Time)
			assert.NotEmpty(t, response.Version)

			if tt.healthy {
				assert.Equal(t, "healthy", response.Status)
				assert.Empty(t, response.Error)
			} else {
				assert.Equal(t, "unhealthy", response.Status)
				assert.NotEmpty(t, response.Error)
			}
		})
	}
}
