# FR0G-AI-IO TODO

## HIGH PRIORITY - External API Integration
- [ ] CRITICAL: Google Voice API client implementation
- [ ] CRITICAL: Complete gRPC bidirectional communication with master-control
- [ ] CRITICAL: SMS output with 99% delivery rate target
- [ ] MISSING: Error handling and retry mechanisms
- [ ] MISSING: External API rate limiting and throttling

## COMPLETED
- [x] 5 input processors operational (SMS, Voice, IRC, ESMTP, Discord)
- [x] Advanced output command review system
- [x] Basic gRPC communication structure
- [x] Threat vector interception framework

## IN PROGRESS
- [ ] Google Voice API integration
- [ ] Master-control gRPC bidirectional streams
- [ ] Production-ready SMS delivery
- [ ] External API error handling

## TECHNICAL DEBT
- [ ] Input processor optimization
- [ ] Real-time threat detection enhancement
- [ ] Multi-channel output coordination
- [ ] Communication pattern analysis
package googlevoice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	sharedconfig "pkg/config"
)

type Client struct {
	httpClient   *http.Client
	apiKey       string
	baseURL      string
	retryConfig  RetryConfig
}

type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

type SMSRequest struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
}

type SMSResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

type DeliveryStatus struct {
	MessageID     string    `json:"message_id"`
	Status        string    `json:"status"` // "sent", "delivered", "failed"
	DeliveredAt   time.Time `json:"delivered_at,omitempty"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	AttemptCount  int       `json:"attempt_count"`
}

func NewClient(config *Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:  config.APIKey,
		baseURL: config.BaseURL,
		retryConfig: RetryConfig{
			MaxRetries:    config.MaxRetries,
			InitialDelay:  time.Duration(config.InitialDelayMs) * time.Millisecond,
			MaxDelay:      time.Duration(config.MaxDelayMs) * time.Millisecond,
			BackoffFactor: config.BackoffFactor,
		},
	}
}

func (c *Client) SendSMS(ctx context.Context, req *SMSRequest) (*SMSResponse, error) {
	var lastErr error
	delay := c.retryConfig.InitialDelay

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.sendSMSAttempt(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if !c.shouldRetry(err) {
			break
		}

		// Exponential backoff
		delay = time.Duration(float64(delay) * c.retryConfig.BackoffFactor)
		if delay > c.retryConfig.MaxDelay {
			delay = c.retryConfig.MaxDelay
		}
	}

	return nil, fmt.Errorf("SMS send failed after %d attempts: %w", c.retryConfig.MaxRetries+1, lastErr)
}

func (c *Client) sendSMSAttempt(ctx context.Context, req *SMSRequest) (*SMSResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/sms/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", httpResp.StatusCode, string(body))
	}

	var smsResp SMSResponse
	if err := json.Unmarshal(body, &smsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &smsResp, nil
}

func (c *Client) GetDeliveryStatus(ctx context.Context, messageID string) (*DeliveryStatus, error) {
	url := fmt.Sprintf("%s/sms/status/%s", c.baseURL, messageID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var status DeliveryStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &status, nil
}

func (c *Client) shouldRetry(err error) bool {
	// Implement retry logic based on error type
	// For now, retry on network errors and 5xx status codes
	return true // Simplified for now
}
package googlevoice

import (
	sharedconfig "pkg/config"
)

type Config struct {
	APIKey         string  `yaml:"api_key" env:"GOOGLE_VOICE_API_KEY"`
	BaseURL        string  `yaml:"base_url" env:"GOOGLE_VOICE_BASE_URL"`
	MaxRetries     int     `yaml:"max_retries" env:"GOOGLE_VOICE_MAX_RETRIES"`
	InitialDelayMs int     `yaml:"initial_delay_ms" env:"GOOGLE_VOICE_INITIAL_DELAY_MS"`
	MaxDelayMs     int     `yaml:"max_delay_ms" env:"GOOGLE_VOICE_MAX_DELAY_MS"`
	BackoffFactor  float64 `yaml:"backoff_factor" env:"GOOGLE_VOICE_BACKOFF_FACTOR"`
	FromNumber     string  `yaml:"from_number" env:"GOOGLE_VOICE_FROM_NUMBER"`
}

func (c *Config) Validate() sharedconfig.ValidationErrors {
	var errors sharedconfig.ValidationErrors

	if err := sharedconfig.ValidateRequired(c.APIKey, "google_voice.api_key"); err != nil {
		errors = append(errors, *err)
	}

	if err := sharedconfig.ValidateRequired(c.BaseURL, "google_voice.base_url"); err != nil {
		errors = append(errors, *err)
	}

	if err := sharedconfig.ValidateRequired(c.FromNumber, "google_voice.from_number"); err != nil {
		errors = append(errors, *err)
	}

	if c.MaxRetries < 0 || c.MaxRetries > 10 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "google_voice.max_retries",
			Message: "must be between 0 and 10",
		})
	}

	if c.InitialDelayMs < 100 || c.InitialDelayMs > 10000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "google_voice.initial_delay_ms",
			Message: "must be between 100 and 10000",
		})
	}

	if c.MaxDelayMs < c.InitialDelayMs || c.MaxDelayMs > 60000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "google_voice.max_delay_ms",
			Message: "must be greater than initial_delay_ms and less than 60000",
		})
	}

	if c.BackoffFactor < 1.0 || c.BackoffFactor > 5.0 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "google_voice.backoff_factor",
			Message: "must be between 1.0 and 5.0",
		})
	}

	return errors
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL:        "https://api.googlevoice.com/v1",
		MaxRetries:     3,
		InitialDelayMs: 1000,
		MaxDelayMs:     30000,
		BackoffFactor:  2.0,
	}
}
package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	sharedconfig "pkg/config"
)

type Client struct {
	conn           *grpc.ClientConn
	config         *Config
	reconnectMutex sync.Mutex
	isConnected    bool
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

type Config struct {
	Host               string        `yaml:"host" env:"GRPC_HOST"`
	Port               int           `yaml:"port" env:"GRPC_PORT"`
	MaxRetries         int           `yaml:"max_retries" env:"GRPC_MAX_RETRIES"`
	ReconnectDelayMs   int           `yaml:"reconnect_delay_ms" env:"GRPC_RECONNECT_DELAY_MS"`
	HeartbeatIntervalMs int          `yaml:"heartbeat_interval_ms" env:"GRPC_HEARTBEAT_INTERVAL_MS"`
	TimeoutMs          int           `yaml:"timeout_ms" env:"GRPC_TIMEOUT_MS"`
}

type StreamHandler interface {
	HandleInputEvent(event *InputEvent) (*InputEventResponse, error)
	HandleOutputCommand(command *OutputCommand) error
}

type InputEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Priority  int                    `json:"priority"`
}

type InputEventResponse struct {
	EventID     string                 `json:"event_id"`
	Processed   bool                   `json:"processed"`
	Actions     []OutputAction         `json:"actions"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProcessedAt time.Time              `json:"processed_at"`
}

type OutputAction struct {
	Type     string                 `json:"type"`
	Target   string                 `json:"target"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

type OutputCommand struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Target   string                 `json:"target"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
	Priority int                    `json:"priority"`
}

func NewClient(config *Config) *Client {
	return &Client{
		config:   config,
		stopChan: make(chan struct{}),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	c.reconnectMutex.Lock()
	defer c.reconnectMutex.Unlock()

	if c.isConnected {
		return nil
	}

	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	
	conn, err := grpc.DialContext(ctx, address, 
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Duration(c.config.TimeoutMs)*time.Millisecond),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	c.conn = conn
	c.isConnected = true
	
	log.Printf("Connected to master-control gRPC server at %s", address)
	return nil
}

func (c *Client) StartBidirectionalStream(ctx context.Context, handler StreamHandler) error {
	if !c.isConnected {
		return fmt.Errorf("client not connected")
	}

	c.wg.Add(1)
	go c.maintainConnection(ctx, handler)

	return nil
}

func (c *Client) maintainConnection(ctx context.Context, handler StreamHandler) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopChan:
			return
		default:
			if err := c.runBidirectionalStream(ctx, handler); err != nil {
				log.Printf("Bidirectional stream error: %v", err)
				
				// Attempt reconnection
				if err := c.reconnect(ctx); err != nil {
					log.Printf("Reconnection failed: %v", err)
					time.Sleep(time.Duration(c.config.ReconnectDelayMs) * time.Millisecond)
				}
			}
		}
	}
}

func (c *Client) runBidirectionalStream(ctx context.Context, handler StreamHandler) error {
	// This would implement the actual bidirectional streaming
	// For now, implementing a simplified version
	
	ticker := time.NewTicker(time.Duration(c.config.HeartbeatIntervalMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.stopChan:
			return nil
		case <-ticker.C:
			// Send heartbeat or process pending messages
			if err := c.sendHeartbeat(ctx); err != nil {
				return fmt.Errorf("heartbeat failed: %w", err)
			}
		}
	}
}

func (c *Client) sendHeartbeat(ctx context.Context) error {
	// Implement heartbeat logic
	return nil
}

func (c *Client) reconnect(ctx context.Context) error {
	c.reconnectMutex.Lock()
	defer c.reconnectMutex.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.isConnected = false
	}

	return c.Connect(ctx)
}

func (c *Client) SendInputEvent(ctx context.Context, event *InputEvent) (*InputEventResponse, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("client not connected")
	}

	// Implement actual gRPC call
	// For now, return a mock response
	return &InputEventResponse{
		EventID:     event.ID,
		Processed:   true,
		ProcessedAt: time.Now(),
	}, nil
}

func (c *Client) Close() error {
	close(c.stopChan)
	c.wg.Wait()

	c.reconnectMutex.Lock()
	defer c.reconnectMutex.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.isConnected = false
		return err
	}

	return nil
}

func (c *Config) Validate() sharedconfig.ValidationErrors {
	var errors sharedconfig.ValidationErrors

	if err := sharedconfig.ValidateRequired(c.Host, "grpc.host"); err != nil {
		errors = append(errors, *err)
	}

	if err := sharedconfig.ValidatePort(c.Port, "grpc.port"); err != nil {
		errors = append(errors, *err)
	}

	if c.MaxRetries < 0 || c.MaxRetries > 10 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "grpc.max_retries",
			Message: "must be between 0 and 10",
		})
	}

	if c.ReconnectDelayMs < 100 || c.ReconnectDelayMs > 60000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "grpc.reconnect_delay_ms",
			Message: "must be between 100 and 60000",
		})
	}

	if c.HeartbeatIntervalMs < 1000 || c.HeartbeatIntervalMs > 300000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "grpc.heartbeat_interval_ms",
			Message: "must be between 1000 and 300000",
		})
	}

	if c.TimeoutMs < 1000 || c.TimeoutMs > 60000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "grpc.timeout_ms",
			Message: "must be between 1000 and 60000",
		})
	}

	return errors
}

func DefaultGRPCConfig() *Config {
	return &Config{
		Host:                "localhost",
		Port:                9092,
		MaxRetries:          3,
		ReconnectDelayMs:    5000,
		HeartbeatIntervalMs: 30000,
		TimeoutMs:           10000,
	}
}
package sms

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"fr0g-ai-io/internal/external/googlevoice"
	sharedconfig "pkg/config"
)

type Service struct {
	client          *googlevoice.Client
	config          *Config
	deliveryTracker *DeliveryTracker
	rateLimiter     *RateLimiter
}

type Config struct {
	GoogleVoice     *googlevoice.Config `yaml:"google_voice"`
	MaxConcurrent   int                 `yaml:"max_concurrent" env:"SMS_MAX_CONCURRENT"`
	RateLimit       int                 `yaml:"rate_limit" env:"SMS_RATE_LIMIT"`
	DeliveryTimeout int                 `yaml:"delivery_timeout_ms" env:"SMS_DELIVERY_TIMEOUT_MS"`
}

type DeliveryTracker struct {
	mutex     sync.RWMutex
	messages  map[string]*MessageStatus
	callbacks map[string]func(*MessageStatus)
}

type MessageStatus struct {
	ID           string
	Status       string
	SentAt       time.Time
	DeliveredAt  *time.Time
	AttemptCount int
	LastError    string
}

type RateLimiter struct {
	tokens    chan struct{}
	refillRate time.Duration
	stopChan  chan struct{}
}

type SMSCommand struct {
	ID       string            `json:"id"`
	To       string            `json:"to"`
	Content  string            `json:"content"`
	Priority int               `json:"priority"`
	Metadata map[string]string `json:"metadata"`
}

type SMSResult struct {
	CommandID    string    `json:"command_id"`
	Success      bool      `json:"success"`
	MessageID    string    `json:"message_id,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	SentAt       time.Time `json:"sent_at"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
}

func NewService(config *Config) (*Service, error) {
	if err := config.Validate(); len(err) > 0 {
		return nil, fmt.Errorf("invalid SMS service configuration: %v", err)
	}

	client := googlevoice.NewClient(config.GoogleVoice)
	
	tracker := &DeliveryTracker{
		messages:  make(map[string]*MessageStatus),
		callbacks: make(map[string]func(*MessageStatus)),
	}

	rateLimiter := &RateLimiter{
		tokens:     make(chan struct{}, config.RateLimit),
		refillRate: time.Second / time.Duration(config.RateLimit),
		stopChan:   make(chan struct{}),
	}

	// Fill initial tokens
	for i := 0; i < config.RateLimit; i++ {
		rateLimiter.tokens <- struct{}{}
	}

	service := &Service{
		client:          client,
		config:          config,
		deliveryTracker: tracker,
		rateLimiter:     rateLimiter,
	}

	// Start rate limiter refill
	go service.refillTokens()

	return service, nil
}

func (s *Service) SendSMS(ctx context.Context, cmd *SMSCommand) (*SMSResult, error) {
	// Rate limiting
	select {
	case <-s.rateLimiter.tokens:
		// Token acquired
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Track message
	status := &MessageStatus{
		ID:           cmd.ID,
		Status:       "sending",
		SentAt:       time.Now(),
		AttemptCount: 1,
	}
	s.deliveryTracker.trackMessage(status)

	// Send via Google Voice API
	req := &googlevoice.SMSRequest{
		To:      cmd.To,
		From:    s.config.GoogleVoice.FromNumber,
		Message: cmd.Content,
	}

	resp, err := s.client.SendSMS(ctx, req)
	if err != nil {
		status.Status = "failed"
		status.LastError = err.Error()
		s.deliveryTracker.updateMessage(status)
		
		return &SMSResult{
			CommandID:    cmd.ID,
			Success:      false,
			ErrorMessage: err.Error(),
			SentAt:       status.SentAt,
		}, err
	}

	status.Status = "sent"
	status.LastError = ""
	s.deliveryTracker.updateMessage(status)

	// Start delivery tracking
	go s.trackDelivery(ctx, resp.MessageID, cmd.ID)

	return &SMSResult{
		CommandID: cmd.ID,
		Success:   true,
		MessageID: resp.MessageID,
		SentAt:    status.SentAt,
	}, nil
}

func (s *Service) trackDelivery(ctx context.Context, messageID, commandID string) {
	timeout := time.Duration(s.config.DeliveryTimeout) * time.Millisecond
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timeoutTimer.C:
			// Delivery timeout
			s.deliveryTracker.mutex.Lock()
			if status, exists := s.deliveryTracker.messages[commandID]; exists {
				if status.Status == "sent" {
					status.Status = "timeout"
					status.LastError = "delivery confirmation timeout"
				}
			}
			s.deliveryTracker.mutex.Unlock()
			return
		case <-ticker.C:
			// Check delivery status
			deliveryStatus, err := s.client.GetDeliveryStatus(ctx, messageID)
			if err != nil {
				log.Printf("Failed to get delivery status for message %s: %v", messageID, err)
				continue
			}

			s.deliveryTracker.mutex.Lock()
			if status, exists := s.deliveryTracker.messages[commandID]; exists {
				status.Status = deliveryStatus.Status
				if deliveryStatus.Status == "delivered" && !deliveryStatus.DeliveredAt.IsZero() {
					status.DeliveredAt = &deliveryStatus.DeliveredAt
				}
				if deliveryStatus.ErrorMessage != "" {
					status.LastError = deliveryStatus.ErrorMessage
				}
				status.AttemptCount = deliveryStatus.AttemptCount
			}
			s.deliveryTracker.mutex.Unlock()

			// Stop tracking if delivered or failed
			if deliveryStatus.Status == "delivered" || deliveryStatus.Status == "failed" {
				return
			}
		}
	}
}

func (s *Service) GetDeliveryStatus(commandID string) (*MessageStatus, bool) {
	s.deliveryTracker.mutex.RLock()
	defer s.deliveryTracker.mutex.RUnlock()
	
	status, exists := s.deliveryTracker.messages[commandID]
	return status, exists
}

func (s *Service) GetDeliveryRate() float64 {
	s.deliveryTracker.mutex.RLock()
	defer s.deliveryTracker.mutex.RUnlock()

	total := len(s.deliveryTracker.messages)
	if total == 0 {
		return 0.0
	}

	delivered := 0
	for _, status := range s.deliveryTracker.messages {
		if status.Status == "delivered" {
			delivered++
		}
	}

	return float64(delivered) / float64(total) * 100.0
}

func (s *Service) refillTokens() {
	ticker := time.NewTicker(s.rateLimiter.refillRate)
	defer ticker.Stop()

	for {
		select {
		case <-s.rateLimiter.stopChan:
			return
		case <-ticker.C:
			select {
			case s.rateLimiter.tokens <- struct{}{}:
				// Token added
			default:
				// Channel full, skip
			}
		}
	}
}

func (dt *DeliveryTracker) trackMessage(status *MessageStatus) {
	dt.mutex.Lock()
	defer dt.mutex.Unlock()
	dt.messages[status.ID] = status
}

func (dt *DeliveryTracker) updateMessage(status *MessageStatus) {
	dt.mutex.Lock()
	defer dt.mutex.Unlock()
	if existing, exists := dt.messages[status.ID]; exists {
		existing.Status = status.Status
		existing.LastError = status.LastError
		existing.AttemptCount = status.AttemptCount
		if status.DeliveredAt != nil {
			existing.DeliveredAt = status.DeliveredAt
		}
	}
}

func (s *Service) Close() error {
	close(s.rateLimiter.stopChan)
	return nil
}

func (c *Config) Validate() sharedconfig.ValidationErrors {
	var errors sharedconfig.ValidationErrors

	if c.GoogleVoice == nil {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "sms.google_voice",
			Message: "Google Voice configuration is required",
		})
	} else {
		if gvErrors := c.GoogleVoice.Validate(); len(gvErrors) > 0 {
			errors = append(errors, gvErrors...)
		}
	}

	if c.MaxConcurrent < 1 || c.MaxConcurrent > 100 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "sms.max_concurrent",
			Message: "must be between 1 and 100",
		})
	}

	if c.RateLimit < 1 || c.RateLimit > 1000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "sms.rate_limit",
			Message: "must be between 1 and 1000",
		})
	}

	if c.DeliveryTimeout < 5000 || c.DeliveryTimeout > 300000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "sms.delivery_timeout_ms",
			Message: "must be between 5000 and 300000",
		})
	}

	return errors
}

func DefaultSMSConfig() *Config {
	return &Config{
		GoogleVoice:     googlevoice.DefaultConfig(),
		MaxConcurrent:   10,
		RateLimit:       60, // 60 SMS per second
		DeliveryTimeout: 60000, // 60 seconds
	}
}
