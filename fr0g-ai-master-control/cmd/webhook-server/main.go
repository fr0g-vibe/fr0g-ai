package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"fr0g-ai-master-control/internal/mastercontrol/input"
)

// MockAIClient implements a mock AI client for testing
type MockAIClient struct{}

func (m *MockAIClient) CreateCommunity(ctx context.Context, topic string, personaCount int) (*input.Community, error) {
	return &input.Community{
		ID:        fmt.Sprintf("community_%d", time.Now().UnixNano()),
		Topic:     topic,
		Members:   []input.PersonaInfo{},
		CreatedAt: time.Now(),
		Status:    "active",
	}, nil
}

func (m *MockAIClient) SubmitForReview(ctx context.Context, communityID string, content string) (*input.CommunityReview, error) {
	// Mock threat analysis based on content keywords
	score := 0.3 // Default low threat
	
	// Simple keyword-based threat detection for demo
	threatKeywords := []string{
		"urgent", "click here", "verify", "suspended", "bitcoin", "giveaway",
		"microsoft", "technical support", "irs", "arrest", "wire transfer",
		"malware", "virus", "crack", "download", "free money", "investment",
		"cryptocurrency", "phishing", "scam", "fraud", "suspicious",
	}
	
	contentLower := strings.ToLower(content)
	for _, keyword := range threatKeywords {
		if strings.Contains(contentLower, keyword) {
			score += 0.15
		}
	}
	
	if score > 1.0 {
		score = 0.95
	}
	
	now := time.Now()
	return &input.CommunityReview{
		ReviewID:    fmt.Sprintf("review_%d", time.Now().UnixNano()),
		Topic:       "threat_analysis",
		Content:     content,
		Consensus: &input.Consensus{
			OverallScore:    score,
			Agreement:       0.85,
			Recommendation:  getThreatLevel(score),
			KeyPoints:       []string{"Automated threat analysis", "Keyword-based detection"},
			ConfidenceLevel: 0.8,
		},
		CreatedAt:   now,
		CompletedAt: &now,
	}, nil
}

func (m *MockAIClient) GetReviewStatus(ctx context.Context, reviewID string) (*input.CommunityReview, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockAIClient) GetCommunityMembers(ctx context.Context, communityID string) ([]input.PersonaInfo, error) {
	return []input.PersonaInfo{}, nil
}

func getThreatLevel(score float64) string {
	if score >= 0.9 {
		return "critical"
	} else if score >= 0.8 {
		return "high"
	} else if score >= 0.6 {
		return "medium"
	} else if score >= 0.4 {
		return "low"
	}
	return "minimal"
}

func main() {
	// Create mock AI client
	aiClient := &MockAIClient{}
	
	// Create processors with default configs
	ircConfig := &input.IRCConfig{
		Server:            "irc.example.com",
		Port:              6667,
		CommunityTopic:    "irc_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedNicks:      []string{"trusted_admin"},
		IgnoredNicks:      []string{"bot", "system"},
	}
	
	smsConfig := &input.SMSConfig{
		Provider:          "twilio",
		CommunityTopic:    "sms_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedNumbers:    []string{"+15554567890"},
		BlockedNumbers:    []string{"+15559999999"},
	}
	
	voiceConfig := &input.VoiceConfig{
		Provider:          "twilio",
		CommunityTopic:    "voice_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedNumbers:    []string{"+15554567890"},
		BlockedNumbers:    []string{"+15559999999"},
	}
	
	esmtpConfig := &input.ESMTPConfig{
		Server:            "mail.example.com",
		Port:              587,
		CommunityTopic:    "email_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedDomains:    []string{"github.com", "company.com"},
		BlockedDomains:    []string{"fake-bank.com", "malicious-site.com"},
	}
	
	// Create processors
	ircProcessor, _ := input.NewIRCProcessor(ircConfig, aiClient)
	smsProcessor, _ := input.NewSMSProcessor(smsConfig, aiClient)
	voiceProcessor, _ := input.NewVoiceProcessor(voiceConfig, aiClient)
	esmtpProcessor, _ := input.NewESMTPProcessor(esmtpConfig, aiClient)
	
	// Create router
	router := mux.NewRouter()
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")
	
	// Webhook endpoints
	router.HandleFunc("/webhook/irc", createWebhookHandler(ircProcessor)).Methods("POST")
	router.HandleFunc("/webhook/sms", createWebhookHandler(smsProcessor)).Methods("POST")
	router.HandleFunc("/webhook/voice", createWebhookHandler(voiceProcessor)).Methods("POST")
	router.HandleFunc("/webhook/esmtp", createWebhookHandler(esmtpProcessor)).Methods("POST")
	
	// Start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	
	// Start server in goroutine
	go func() {
		log.Printf("Webhook server starting on :8080")
		log.Printf("Available endpoints:")
		log.Printf("  GET  /health")
		log.Printf("  POST /webhook/irc")
		log.Printf("  POST /webhook/sms")
		log.Printf("  POST /webhook/voice")
		log.Printf("  POST /webhook/esmtp")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

func createWebhookHandler(processor input.WebhookProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var requestBody interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		// Create webhook request
		webhookReq := &input.WebhookRequest{
			ID:        fmt.Sprintf("req_%d", time.Now().UnixNano()),
			Source:    "webhook_server",
			Tag:       processor.GetTag(),
			Timestamp: time.Now(),
			Headers:   make(map[string]string),
			Body:      requestBody,
			Metadata:  make(map[string]interface{}),
		}
		
		// Copy headers
		for key, values := range r.Header {
			if len(values) > 0 {
				webhookReq.Headers[key] = values[0]
			}
		}
		
		// Process webhook
		response, err := processor.ProcessWebhook(r.Context(), webhookReq)
		if err != nil {
			log.Printf("Webhook processing error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Return response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Response encoding error: %v", err)
			http.Error(w, "Response encoding failed", http.StatusInternalServerError)
		}
	}
}
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"fr0g-ai-master-control/internal/mastercontrol/input"
)

// MockAIClient implements a mock AI client for testing
type MockAIClient struct{}

func (m *MockAIClient) CreateCommunity(ctx context.Context, topic string, personaCount int) (*input.Community, error) {
	return &input.Community{
		ID:        fmt.Sprintf("community_%d", time.Now().UnixNano()),
		Topic:     topic,
		Members:   []input.PersonaInfo{},
		CreatedAt: time.Now(),
		Status:    "active",
	}, nil
}

func (m *MockAIClient) SubmitForReview(ctx context.Context, communityID string, content string) (*input.CommunityReview, error) {
	// Mock threat analysis based on content keywords
	score := 0.3 // Default low threat
	
	// Simple keyword-based threat detection for demo
	threatKeywords := []string{
		"urgent", "click here", "verify", "suspended", "bitcoin", "giveaway",
		"microsoft", "technical support", "irs", "arrest", "wire transfer",
		"malware", "virus", "crack", "download", "free money", "investment",
		"cryptocurrency", "phishing", "scam", "fraud", "suspicious",
	}
	
	contentLower := strings.ToLower(content)
	for _, keyword := range threatKeywords {
		if strings.Contains(contentLower, keyword) {
			score += 0.15
		}
	}
	
	if score > 1.0 {
		score = 0.95
	}
	
	now := time.Now()
	return &input.CommunityReview{
		ReviewID:    fmt.Sprintf("review_%d", time.Now().UnixNano()),
		Topic:       "threat_analysis",
		Content:     content,
		Consensus: &input.Consensus{
			OverallScore:    score,
			Agreement:       0.85,
			Recommendation:  getThreatLevel(score),
			KeyPoints:       []string{"Automated threat analysis", "Keyword-based detection"},
			ConfidenceLevel: 0.8,
		},
		CreatedAt:   now,
		CompletedAt: &now,
	}, nil
}

func (m *MockAIClient) GetReviewStatus(ctx context.Context, reviewID string) (*input.CommunityReview, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockAIClient) GetCommunityMembers(ctx context.Context, communityID string) ([]input.PersonaInfo, error) {
	return []input.PersonaInfo{}, nil
}

func getThreatLevel(score float64) string {
	if score >= 0.9 {
		return "critical"
	} else if score >= 0.8 {
		return "high"
	} else if score >= 0.6 {
		return "medium"
	} else if score >= 0.4 {
		return "low"
	}
	return "minimal"
}

func main() {
	// Create mock AI client
	aiClient := &MockAIClient{}
	
	// Create processors with default configs
	ircConfig := &input.IRCConfig{
		Server:            "irc.example.com",
		Port:              6667,
		CommunityTopic:    "irc_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedNicks:      []string{"trusted_admin"},
		IgnoredNicks:      []string{"bot", "system"},
	}
	
	smsConfig := &input.SMSConfig{
		Provider:          "twilio",
		CommunityTopic:    "sms_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedNumbers:    []string{"+15554567890"},
		BlockedNumbers:    []string{"+15559999999"},
	}
	
	voiceConfig := &input.VoiceConfig{
		Provider:          "twilio",
		CommunityTopic:    "voice_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedNumbers:    []string{"+15554567890"},
		BlockedNumbers:    []string{"+15559999999"},
	}
	
	esmtpConfig := &input.ESMTPConfig{
		Server:            "mail.example.com",
		Port:              587,
		CommunityTopic:    "email_threat_analysis",
		PersonaCount:      5,
		RequiredConsensus: 0.7,
		TrustedDomains:    []string{"github.com", "company.com"},
		BlockedDomains:    []string{"fake-bank.com", "malicious-site.com"},
	}
	
	// Create processors
	ircProcessor, _ := input.NewIRCProcessor(ircConfig, aiClient)
	smsProcessor, _ := input.NewSMSProcessor(smsConfig, aiClient)
	voiceProcessor, _ := input.NewVoiceProcessor(voiceConfig, aiClient)
	esmtpProcessor, _ := input.NewESMTPProcessor(esmtpConfig, aiClient)
	
	// Create router
	router := mux.NewRouter()
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")
	
	// Webhook endpoints
	router.HandleFunc("/webhook/irc", createWebhookHandler(ircProcessor)).Methods("POST")
	router.HandleFunc("/webhook/sms", createWebhookHandler(smsProcessor)).Methods("POST")
	router.HandleFunc("/webhook/voice", createWebhookHandler(voiceProcessor)).Methods("POST")
	router.HandleFunc("/webhook/esmtp", createWebhookHandler(esmtpProcessor)).Methods("POST")
	
	// Start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	
	// Start server in goroutine
	go func() {
		log.Printf("Webhook server starting on :8080")
		log.Printf("Available endpoints:")
		log.Printf("  GET  /health")
		log.Printf("  POST /webhook/irc")
		log.Printf("  POST /webhook/sms")
		log.Printf("  POST /webhook/voice")
		log.Printf("  POST /webhook/esmtp")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

func createWebhookHandler(processor input.WebhookProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var requestBody interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		// Create webhook request
		webhookReq := &input.WebhookRequest{
			ID:        fmt.Sprintf("req_%d", time.Now().UnixNano()),
			Source:    "webhook_server",
			Tag:       processor.GetTag(),
			Timestamp: time.Now(),
			Headers:   make(map[string]string),
			Body:      requestBody,
			Metadata:  make(map[string]interface{}),
		}
		
		// Copy headers
		for key, values := range r.Header {
			if len(values) > 0 {
				webhookReq.Headers[key] = values[0]
			}
		}
		
		// Process webhook
		response, err := processor.ProcessWebhook(r.Context(), webhookReq)
		if err != nil {
			log.Printf("Webhook processing error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Return response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Response encoding error: %v", err)
			http.Error(w, "Response encoding failed", http.StatusInternalServerError)
		}
	}
}
