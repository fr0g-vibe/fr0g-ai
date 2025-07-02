package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fr0g-vibe/fr0g-ai-bridge/internal/models"
)

// OpenWebUIClient handles communication with OpenWebUI API
type OpenWebUIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewOpenWebUIClient creates a new OpenWebUI client
func NewOpenWebUIClient(baseURL, apiKey string, timeout time.Duration) *OpenWebUIClient {
	return &OpenWebUIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ChatCompletion sends a chat completion request to OpenWebUI
func (c *OpenWebUIClient) ChatCompletion(ctx context.Context, req *models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	// Prepare the request for OpenWebUI
	openWebUIReq := c.prepareOpenWebUIRequest(req)

	// Marshal the request
	reqBody, err := json.Marshal(openWebUIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenWebUI API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var chatResp models.ChatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &chatResp, nil
}

// prepareOpenWebUIRequest converts our request format to OpenWebUI format
func (c *OpenWebUIClient) prepareOpenWebUIRequest(req *models.ChatCompletionRequest) *models.ChatCompletionRequest {
	// Create a copy of the request
	openWebUIReq := *req

	// If persona prompt is provided, prepend it as a system message
	if req.PersonaPrompt != "" {
		systemMessage := models.ChatMessage{
			Role:    "system",
			Content: req.PersonaPrompt,
		}

		// Check if there's already a system message
		hasSystemMessage := false
		for i, msg := range openWebUIReq.Messages {
			if msg.Role == "system" {
				// Prepend persona prompt to existing system message
				openWebUIReq.Messages[i].Content = req.PersonaPrompt + "\n\n" + msg.Content
				hasSystemMessage = true
				break
			}
		}

		// If no system message exists, add one at the beginning
		if !hasSystemMessage {
			openWebUIReq.Messages = append([]models.ChatMessage{systemMessage}, openWebUIReq.Messages...)
		}
	}

	// Clear persona prompt as it's not part of OpenWebUI API
	openWebUIReq.PersonaPrompt = ""

	return &openWebUIReq
}

// HealthCheck performs a health check against OpenWebUI
func (c *OpenWebUIClient) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenWebUIClient handles communication with OpenWebUI API
type OpenWebUIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewOpenWebUIClient creates a new OpenWebUI client
func NewOpenWebUIClient(baseURL, apiKey string, timeout time.Duration) *OpenWebUIClient {
	return &OpenWebUIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ChatMessage represents a chat message for OpenWebUI
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat request to OpenWebUI
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ChatResponse represents a response from OpenWebUI
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// Model represents an AI model
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelsResponse represents the response from the models endpoint
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// SendMessage sends a message to OpenWebUI and returns the response
func (c *OpenWebUIClient) SendMessage(message, model string) (string, error) {
	if model == "" {
		model = "gpt-3.5-turbo" // Default model
	}

	chatReq := ChatRequest{
		Model: model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
		Stream: false,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// GetModels retrieves available models from OpenWebUI
func (c *OpenWebUIClient) GetModels() ([]Model, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return modelsResp.Data, nil
}
