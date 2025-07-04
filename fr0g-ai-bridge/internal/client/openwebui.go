package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fr0g-ai-bridge/internal/models"
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
		fmt.Printf("DEBUG: Using API key for OpenWebUI: %s...\n", c.apiKey[:min(10, len(c.apiKey))])
	} else {
		fmt.Println("DEBUG: No API key configured for OpenWebUI")
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
		fmt.Printf("DEBUG: OpenWebUI API error - Status: %d, Response: %s\n", resp.StatusCode, string(respBody))
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
		fmt.Printf("DEBUG: HealthCheck using API key: %s...\n", c.apiKey[:min(10, len(c.apiKey))])
	} else {
		fmt.Println("DEBUG: HealthCheck - No API key configured")
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("DEBUG: HealthCheck failed - Status: %d, Response: %s\n", resp.StatusCode, string(body))
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}

// Legacy methods for backward compatibility with REST API
// SendMessage sends a simple message and returns the response
func (c *OpenWebUIClient) SendMessage(message, model string) (string, error) {
	req := &models.ChatCompletionRequest{
		Model: model,
		Messages: []models.ChatMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	resp, err := c.ChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

// Model represents an AI model for legacy compatibility
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// GetModels retrieves available models from OpenWebUI
func (c *OpenWebUIClient) GetModels() ([]Model, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		fmt.Printf("DEBUG: GetModels using API key: %s...\n", c.apiKey[:min(10, len(c.apiKey))])
	} else {
		fmt.Println("DEBUG: GetModels - No API key configured")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("DEBUG: GetModels failed - Status: %d, Response: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var modelsResp struct {
		Object string  `json:"object"`
		Data   []Model `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return modelsResp.Data, nil
}
