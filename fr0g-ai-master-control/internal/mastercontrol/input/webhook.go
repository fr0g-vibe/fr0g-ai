package input

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// WebhookManager handles incoming webhook requests and routes them for processing
type WebhookManager struct {
	server     *http.Server
	router     *mux.Router
	processors map[string]WebhookProcessor
	config     *WebhookConfig
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// WebhookConfig holds webhook manager configuration
type WebhookConfig struct {
	Port            int           `yaml:"port"`
	Host            string        `yaml:"host"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	MaxRequestSize  int64         `yaml:"max_request_size"`
	EnableLogging   bool          `yaml:"enable_logging"`
	AllowedOrigins  []string      `yaml:"allowed_origins"`
}

// WebhookRequest represents an incoming webhook request
type WebhookRequest struct {
	ID          string                 `json:"id"`
	Source      string                 `json:"source"`
	Tag         string                 `json:"tag"`
	Timestamp   time.Time              `json:"timestamp"`
	Headers     map[string]string      `json:"headers"`
	Body        interface{}            `json:"body"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty"`
}

// WebhookResponse represents the response to a webhook
type WebhookResponse struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// WebhookProcessor defines the interface for processing webhooks
type WebhookProcessor interface {
	ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error)
	GetTag() string
	GetDescription() string
}

// NewWebhookManager creates a new webhook manager
func NewWebhookManager(config *WebhookConfig) *WebhookManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	wm := &WebhookManager{
		router:     mux.NewRouter(),
		processors: make(map[string]WebhookProcessor),
		config:     config,
		ctx:        ctx,
		cancel:     cancel,
	}
	
	wm.setupRoutes()
	wm.setupServer()
	
	return wm
}

// Start begins webhook manager operation
func (wm *WebhookManager) Start() error {
	log.Printf("Webhook Manager: Starting webhook server on %s:%d", wm.config.Host, wm.config.Port)
	
	go func() {
		if err := wm.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Webhook Manager: Server error: %v", err)
		}
	}()
	
	log.Println("Webhook Manager: Webhook server started successfully")
	return nil
}

// Stop gracefully stops the webhook manager
func (wm *WebhookManager) Stop() error {
	log.Println("Webhook Manager: Stopping webhook server...")
	
	wm.cancel()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := wm.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("webhook server shutdown error: %w", err)
	}
	
	log.Println("Webhook Manager: Webhook server stopped")
	return nil
}

// RegisterProcessor registers a webhook processor for a specific tag
func (wm *WebhookManager) RegisterProcessor(processor WebhookProcessor) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	
	tag := processor.GetTag()
	if _, exists := wm.processors[tag]; exists {
		return fmt.Errorf("processor for tag '%s' already registered", tag)
	}
	
	wm.processors[tag] = processor
	log.Printf("Webhook Manager: Registered processor for tag '%s': %s", tag, processor.GetDescription())
	
	return nil
}

// UnregisterProcessor removes a webhook processor
func (wm *WebhookManager) UnregisterProcessor(tag string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	
	delete(wm.processors, tag)
	log.Printf("Webhook Manager: Unregistered processor for tag '%s'", tag)
}

// GetRegisteredProcessors returns all registered processors
func (wm *WebhookManager) GetRegisteredProcessors() map[string]string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	
	processors := make(map[string]string)
	for tag, processor := range wm.processors {
		processors[tag] = processor.GetDescription()
	}
	
	return processors
}

// setupRoutes configures the HTTP routes
func (wm *WebhookManager) setupRoutes() {
	// Generic webhook endpoint with tag parameter
	wm.router.HandleFunc("/webhook/{tag}", wm.handleWebhook).Methods("POST")
	
	// Health check endpoint
	wm.router.HandleFunc("/health", wm.handleHealth).Methods("GET")
	
	// Status endpoint
	wm.router.HandleFunc("/status", wm.handleStatus).Methods("GET")
	
	// Apply middleware
	wm.router.Use(wm.loggingMiddleware)
	wm.router.Use(wm.corsMiddleware)
	wm.router.Use(wm.requestSizeMiddleware)
}

// setupServer configures the HTTP server
func (wm *WebhookManager) setupServer() {
	wm.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", wm.config.Host, wm.config.Port),
		Handler:      wm.router,
		ReadTimeout:  wm.config.ReadTimeout,
		WriteTimeout: wm.config.WriteTimeout,
	}
}

// handleWebhook processes incoming webhook requests
func (wm *WebhookManager) handleWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]
	
	// Create webhook request
	webhookReq := &WebhookRequest{
		ID:        generateRequestID(),
		Source:    r.RemoteAddr,
		Tag:       tag,
		Timestamp: time.Now(),
		Headers:   extractHeaders(r),
		Metadata:  make(map[string]interface{}),
	}
	
	// Parse request body
	var body interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		wm.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON body", webhookReq.ID)
		return
	}
	webhookReq.Body = body
	
	// Find processor for tag
	wm.mu.RLock()
	processor, exists := wm.processors[tag]
	wm.mu.RUnlock()
	
	if !exists {
		wm.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("No processor found for tag '%s'", tag), webhookReq.ID)
		return
	}
	
	// Process webhook
	ctx, cancel := context.WithTimeout(wm.ctx, 30*time.Second)
	defer cancel()
	
	response, err := processor.ProcessWebhook(ctx, webhookReq)
	if err != nil {
		log.Printf("Webhook Manager: Error processing webhook for tag '%s': %v", tag, err)
		wm.writeErrorResponse(w, http.StatusInternalServerError, "Processing error", webhookReq.ID)
		return
	}
	
	// Write successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	
	if wm.config.EnableLogging {
		log.Printf("Webhook Manager: Successfully processed webhook for tag '%s', request ID: %s", tag, webhookReq.ID)
	}
}

// handleHealth provides health check endpoint
func (wm *WebhookManager) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":     "healthy",
		"timestamp":  time.Now(),
		"processors": len(wm.processors),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

// handleStatus provides status information
func (wm *WebhookManager) handleStatus(w http.ResponseWriter, r *http.Request) {
	processors := wm.GetRegisteredProcessors()
	
	status := map[string]interface{}{
		"webhook_manager": "running",
		"processors":      processors,
		"timestamp":       time.Now(),
		"config": map[string]interface{}{
			"port":             wm.config.Port,
			"host":             wm.config.Host,
			"max_request_size": wm.config.MaxRequestSize,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// writeErrorResponse writes an error response
func (wm *WebhookManager) writeErrorResponse(w http.ResponseWriter, statusCode int, message, requestID string) {
	response := &WebhookResponse{
		Success:   false,
		Message:   message,
		RequestID: requestID,
		Timestamp: time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Middleware functions

func (wm *WebhookManager) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wm.config.EnableLogging {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			log.Printf("Webhook Manager: %s %s - %v", r.Method, r.URL.Path, duration)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (wm *WebhookManager) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := "*"
		if len(wm.config.AllowedOrigins) > 0 && wm.config.AllowedOrigins[0] != "*" {
			origin = wm.config.AllowedOrigins[0]
		}
		
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (wm *WebhookManager) requestSizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, wm.config.MaxRequestSize)
		next.ServeHTTP(w, r)
	})
}

// Utility functions

func generateRequestID() string {
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), randomString(8))
}

func extractHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for name, values := range r.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}
	return headers
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
