package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"fr0g-ai-bridge/internal/client"
	"fr0g-ai-bridge/internal/config"
	"fr0g-ai-bridge/internal/models"
)

// RESTServer handles REST API requests
type RESTServer struct {
	client     *client.OpenWebUIClient
	router     *mux.Router
	config     *config.Config
	rateLimiter *RateLimiter
}

// RateLimiter implements a simple rate limiting mechanism
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    requestsPerMinute,
		window:   time.Minute,
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	
	// Clean old requests
	if requests, exists := rl.requests[ip]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.window {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[ip] = validRequests
	}

	// Check if under limit
	if len(rl.requests[ip]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// NewRESTServer creates a new REST server
func NewRESTServer(openWebUIClient *client.OpenWebUIClient, cfg *config.Config) *RESTServer {
	server := &RESTServer{
		client:      openWebUIClient,
		router:      mux.NewRouter(),
		config:      cfg,
		rateLimiter: NewRateLimiter(cfg.Security.RateLimitRPM),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the REST API routes
func (s *RESTServer) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Chat completion endpoint (OpenAI compatible)
	s.router.HandleFunc("/api/chat/completions", s.handleChatCompletion).Methods("POST")

	// Legacy simple chat endpoint
	s.router.HandleFunc("/api/v1/chat", s.handleSimpleChat).Methods("POST")
	
	// Models endpoint
	s.router.HandleFunc("/api/v1/models", s.handleModels).Methods("GET")

	// Add middleware
	s.router.Use(s.securityHeadersMiddleware)
	s.router.Use(s.requestSizeLimitMiddleware)
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.rateLimitMiddleware)
	if s.config.Security.RequireAPIKey {
		s.router.Use(s.apiKeyMiddleware)
	}
	if s.config.Security.EnableCORS {
		s.router.Use(s.corsMiddleware)
	}
}

// GetRouter returns the configured router
func (s *RESTServer) GetRouter() *mux.Router {
	return s.router
}

// handleHealth handles health check requests
func (s *RESTServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check OpenWebUI health
	err := s.client.HealthCheck(ctx)
	
	response := models.HealthResponse{
		Time:    time.Now(),
		Version: "1.0.0",
	}

	if err != nil {
		response.Status = "unhealthy"
		response.Error = err.Error()
		response.Details = map[string]interface{}{
			"openwebui_url": s.config.OpenWebUI.BaseURL,
			"has_api_key":   s.config.OpenWebUI.APIKey != "",
			"api_key_prefix": func() string {
				if s.config.OpenWebUI.APIKey != "" && len(s.config.OpenWebUI.APIKey) > 8 {
					return s.config.OpenWebUI.APIKey[:8] + "..."
				}
				return "none"
			}(),
			"timeout_seconds": s.config.OpenWebUI.Timeout,
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Printf("Health check failed: %v", err)
	} else {
		response.Status = "healthy"
		response.Details = map[string]interface{}{
			"openwebui_url": s.config.OpenWebUI.BaseURL,
			"authenticated": true,
		}
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleChatCompletion handles OpenAI-compatible chat completion requests
func (s *RESTServer) handleChatCompletion(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := s.validateChatCompletionRequest(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Forward to OpenWebUI
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	resp, err := s.client.ChatCompletion(ctx, &req)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to process chat completion", err)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// SimpleChatRequest represents a simple chat request
type SimpleChatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model,omitempty"`
}

// SimpleChatResponse represents a simple chat response
type SimpleChatResponse struct {
	Response  string    `json:"response"`
	Model     string    `json:"model"`
	Timestamp time.Time `json:"timestamp"`
}

// handleSimpleChat handles simple chat requests (legacy endpoint)
func (s *RESTServer) handleSimpleChat(w http.ResponseWriter, r *http.Request) {
	var req SimpleChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if req.Message == "" {
		s.writeError(w, http.StatusBadRequest, "Message is required", nil)
		return
	}

	// Validate message length and content
	if len(req.Message) > 10000 {
		s.writeError(w, http.StatusBadRequest, "Message too long (max 10000 characters)", nil)
		return
	}

	// Sanitize input
	req.Message = strings.TrimSpace(req.Message)
	if req.Model != "" {
		req.Model = strings.TrimSpace(req.Model)
	}

	// Bridge request to OpenWebUI using legacy method
	response, err := s.client.SendMessage(req.Message, req.Model)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to process request", err)
		return
	}

	chatResp := SimpleChatResponse{
		Response:  response,
		Model:     req.Model,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatResp)
}

// handleModels handles requests for available AI models
func (s *RESTServer) handleModels(w http.ResponseWriter, r *http.Request) {
	models, err := s.client.GetModels()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to fetch models", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"models":    models,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// validateChatCompletionRequest validates the chat completion request
func (s *RESTServer) validateChatCompletionRequest(req *models.ChatCompletionRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	if len(req.Messages) == 0 {
		return fmt.Errorf("messages are required")
	}
	if len(req.Messages) > 100 {
		return fmt.Errorf("too many messages (max 100)")
	}
	for i, msg := range req.Messages {
		if msg.Role == "" {
			return fmt.Errorf("message %d: role is required", i)
		}
		if msg.Content == "" {
			return fmt.Errorf("message %d: content is required", i)
		}
		if len(msg.Content) > 10000 {
			return fmt.Errorf("message %d: content too long (max 10000 characters)", i)
		}
		// Sanitize role to prevent injection
		if !isValidRole(msg.Role) {
			return fmt.Errorf("message %d: invalid role", i)
		}
	}
	// Validate persona prompt length
	if len(req.PersonaPrompt) > 5000 {
		return fmt.Errorf("persona prompt too long (max 5000 characters)")
	}
	// Validate optional parameters
	if req.Temperature != nil && (*req.Temperature < 0 || *req.Temperature > 2) {
		return fmt.Errorf("temperature must be between 0 and 2")
	}
	if req.MaxTokens != nil && (*req.MaxTokens < 1 || *req.MaxTokens > 4096) {
		return fmt.Errorf("max_tokens must be between 1 and 4096")
	}
	return nil
}


// writeError writes an error response
func (s *RESTServer) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
	log.Printf("API Error: %s - %v", message, err)
	
	errorResp := models.ErrorResponse{
		Error:   message,
		Code:    statusCode,
	}
	
	if err != nil {
		errorResp.Message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}

// corsMiddleware adds CORS headers
func (s *RESTServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := "*"
		if len(s.config.Security.AllowedOrigins) > 0 && s.config.Security.AllowedOrigins[0] != "*" {
			origin = s.config.Security.AllowedOrigins[0] // Simplified for now
		}
		
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware implements rate limiting
func (s *RESTServer) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client IP
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		if !s.rateLimiter.Allow(clientIP) {
			s.writeError(w, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs all requests
func (s *RESTServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v", 
			r.Method, r.RequestURI, r.RemoteAddr, wrapped.statusCode, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// securityHeadersMiddleware adds security headers to all responses
func (s *RESTServer) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Enforce HTTPS (when behind a proxy)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Prevent referrer leakage
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'none'; object-src 'none'")
		
		next.ServeHTTP(w, r)
	})
}

// requestSizeLimitMiddleware limits the size of request bodies
func (s *RESTServer) requestSizeLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limit request body size to 1MB
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
		next.ServeHTTP(w, r)
	})
}

// apiKeyMiddleware validates API keys
func (s *RESTServer) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip API key check for health endpoint
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract API key from Authorization header or X-API-Key header
		var apiKey string
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			apiKey = r.Header.Get("X-API-Key")
		}

		if apiKey == "" {
			s.writeError(w, http.StatusUnauthorized, "API key required", nil)
			return
		}

		// Validate API key
		if !s.isValidAPIKey(apiKey) {
			s.writeError(w, http.StatusUnauthorized, "Invalid API key", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isValidAPIKey checks if the provided API key is valid
func (s *RESTServer) isValidAPIKey(apiKey string) bool {
	if len(s.config.Security.AllowedAPIKeys) == 0 {
		return true // No API keys configured, allow all
	}
	
	for _, validKey := range s.config.Security.AllowedAPIKeys {
		if apiKey == validKey {
			return true
		}
	}
	return false
}
