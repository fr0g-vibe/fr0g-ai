package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// GoogleVoiceAPI represents the Google Voice API client
type GoogleVoiceAPI struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// SMSRequest represents a request to send an SMS
type SMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
	From    string `json:"from,omitempty"`
}

// SMSResponse represents the response from SMS API
type SMSResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

// Processor handles SMS output operations
type Processor struct {
	config     *sharedconfig.Config
	googleAPI  *GoogleVoiceAPI
	isEnabled  bool
}

// NewProcessor creates a new SMS output processor
func NewProcessor(cfg *sharedconfig.Config) *Processor {
	// Initialize Google Voice API client
	googleAPI := &GoogleVoiceAPI{
		apiKey:  cfg.GetString("GOOGLE_VOICE_API_KEY"),
		baseURL: "https://www.googleapis.com/voice/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return &Processor{
		config:    cfg,
		googleAPI: googleAPI,
		isEnabled: cfg.GetString("GOOGLE_VOICE_API_KEY") != "",
	}
}

// Process executes an SMS output command
func (p *Processor) Process(command *outputs.OutputCommand) (*outputs.OutputResult, error) {
	if !p.isEnabled {
		return &outputs.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: "SMS processor disabled - missing Google Voice API key",
			Metadata:     map[string]string{"processor": "sms"},
			CompletedAt:  time.Now(),
		}, nil
	}

	// Validate command
	if command.Target == "" {
		return &outputs.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: "SMS target phone number is required",
			Metadata:     map[string]string{"processor": "sms"},
			CompletedAt:  time.Now(),
		}, nil
	}

	if command.Content == "" {
		return &outputs.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: "SMS message content is required",
			Metadata:     map[string]string{"processor": "sms"},
			CompletedAt:  time.Now(),
		}, nil
	}

	// Send SMS via Google Voice API
	messageID, err := p.sendSMS(command.Target, command.Content, command.Metadata)
	if err != nil {
		return &outputs.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to send SMS: %v", err),
			Metadata:     map[string]string{"processor": "sms"},
			CompletedAt:  time.Now(),
		}, nil
	}

	return &outputs.OutputResult{
		CommandID:   command.ID,
		Success:     true,
		Metadata: map[string]string{
			"processor":  "sms",
			"message_id": messageID,
			"target":     command.Target,
			"api":        "google_voice",
		},
		CompletedAt: time.Now(),
	}, nil
}

// sendSMS sends an SMS message via Google Voice API
func (p *Processor) sendSMS(phoneNumber, message string, metadata map[string]string) (string, error) {
	// Prepare SMS request
	smsReq := SMSRequest{
		To:      phoneNumber,
		Message: message,
		From:    metadata["from_number"], // Optional sender number
	}

	// Marshal request to JSON
	reqBody, err := json.Marshal(smsReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/messages", p.googleAPI.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.googleAPI.apiKey))
	req.Header.Set("User-Agent", "fr0g-ai-io/1.0")

	// Execute request with retry logic
	var resp *http.Response
	var lastErr error
	
	for attempt := 0; attempt < 3; attempt++ {
		resp, lastErr = p.googleAPI.httpClient.Do(req)
		if lastErr == nil && resp.StatusCode < 500 {
			break
		}
		
		if resp != nil {
			resp.Body.Close()
		}
		
		// Exponential backoff
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}

	if lastErr != nil {
		return "", fmt.Errorf("HTTP request failed after retries: %w", lastErr)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var smsResp SMSResponse
	if err := json.Unmarshal(respBody, &smsResp); err != nil {
		return "", fmt.Errorf("failed to parse SMS response: %w", err)
	}

	if smsResp.Error != "" {
		return "", fmt.Errorf("SMS API error: %s", smsResp.Error)
	}

	return smsResp.MessageID, nil
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "sms"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.isEnabled
}

// GetStatus returns the current status of the SMS processor
func (p *Processor) GetStatus() map[string]interface{} {
	status := map[string]interface{}{
		"type":    "sms",
		"enabled": p.isEnabled,
		"api":     "google_voice",
	}

	if !p.isEnabled {
		status["reason"] = "missing Google Voice API key"
	}

	return status
}
