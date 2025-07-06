package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/processors"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/queue"
	"google.golang.org/grpc"
)

// Server represents the fr0g-ai-io server
type Server struct {
	config         *sharedconfig.Config
	httpServer     *http.Server
	grpcServer     *grpc.Server
	processorMgr   *processors.Manager
	outputMgr      *outputs.Manager
	queueMgr       *queue.Manager
	mu             sync.RWMutex
	isRunning      bool
}

// NewServer creates a new server instance
func NewServer(cfg *sharedconfig.Config) (*Server, error) {
	// Create processor manager
	processorMgr, err := processors.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create processor manager: %w", err)
	}

	// Create output manager
	outputMgr, err := outputs.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create output manager: %w", err)
	}

	// Create queue manager
	queueMgr, err := queue.NewManager(&cfg.Queue)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue manager: %w", err)
	}

	// Create HTTP server
	httpMux := http.NewServeMux()
	httpAddr := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	httpServer := &http.Server{
		Addr:         httpAddr,
		Handler:      httpMux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	server := &Server{
		config:       cfg,
		httpServer:   httpServer,
		grpcServer:   grpcServer,
		processorMgr: processorMgr,
		outputMgr:    outputMgr,
		queueMgr:     queueMgr,
	}

	// Setup HTTP routes
	server.setupHTTPRoutes(httpMux)

	// Setup gRPC services
	server.setupGRPCServices()

	return server, nil
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	log.Printf("Starting fr0g-ai-io server...")

	// Start queue manager
	if err := s.queueMgr.Start(ctx); err != nil {
		return fmt.Errorf("failed to start queue manager: %w", err)
	}

	// Start processor manager
	if err := s.processorMgr.Start(ctx); err != nil {
		return fmt.Errorf("failed to start processor manager: %w", err)
	}

	// Start output manager
	if err := s.outputMgr.Start(ctx); err != nil {
		return fmt.Errorf("failed to start output manager: %w", err)
	}

	// Start HTTP server
	go func() {
		httpAddr := fmt.Sprintf("%s:%s", s.config.HTTP.Host, s.config.HTTP.Port)
		log.Printf("HTTP server listening on %s", httpAddr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Start gRPC server
	go func() {
		grpcAddr := fmt.Sprintf("%s:%s", s.config.GRPC.Host, s.config.GRPC.Port)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Printf("Failed to listen on gRPC address %s: %v", grpcAddr, err)
			return
		}

		log.Printf("gRPC server listening on %s", grpcAddr)
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	log.Printf("fr0g-ai-io server started successfully")
	return nil
}

// Stop stops the server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("server is not running")
	}

	log.Printf("Stopping fr0g-ai-io server...")

	// Stop processor manager
	if err := s.processorMgr.Stop(); err != nil {
		log.Printf("Error stopping processor manager: %v", err)
	}

	// Stop output manager
	if err := s.outputMgr.Stop(); err != nil {
		log.Printf("Error stopping output manager: %v", err)
	}

	// Stop queue manager
	if err := s.queueMgr.Stop(); err != nil {
		log.Printf("Error stopping queue manager: %v", err)
	}

	// Stop HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Error stopping HTTP server: %v", err)
	}

	// Stop gRPC server
	s.grpcServer.GracefulStop()

	s.isRunning = false
	log.Printf("fr0g-ai-io server stopped")
	return nil
}

// setupHTTPRoutes sets up HTTP routes
func (s *Server) setupHTTPRoutes(mux *http.ServeMux) {
	// Health check endpoint
	mux.HandleFunc("/health", s.healthHandler)

	// Processor status endpoints
	mux.HandleFunc("/processors", s.processorsHandler)
	mux.HandleFunc("/processors/", s.processorHandler)

	// Queue status endpoints
	mux.HandleFunc("/queue/status", s.queueStatusHandler)
	mux.HandleFunc("/queue/stats", s.queueStatsHandler)

	// Output endpoints
	mux.HandleFunc("/outputs", s.outputsHandler)
	mux.HandleFunc("/send", s.sendHandler)
	mux.HandleFunc("/broadcast", s.broadcastHandler)

	// Webhook endpoints for external services
	mux.HandleFunc("/webhook/sms", s.smsWebhookHandler)
	mux.HandleFunc("/webhook/voice", s.voiceWebhookHandler)
	mux.HandleFunc("/webhook/discord", s.discordWebhookHandler)
	mux.HandleFunc("/webhook/generic", s.genericWebhookHandler)
}

// setupGRPCServices sets up gRPC services
func (s *Server) setupGRPCServices() {
	// TODO: Register gRPC services
	// This will include services for:
	// - Input message processing
	// - Output message sending
	// - Processor management
	// - Queue management
}

// HTTP Handlers

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	isRunning := s.isRunning
	s.mu.RUnlock()

	status := map[string]interface{}{
		"status":     "healthy",
		"service":    "fr0g-ai-io",
		"version":    "1.0.0",
		"timestamp":  time.Now().UTC(),
		"is_running": isRunning,
		"processors": s.processorMgr.GetStatus(),
		"outputs":    s.outputMgr.GetStatus(),
		"queue":      s.queueMgr.GetStatus(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Simple JSON response without external dependencies
	fmt.Fprintf(w, `{
		"status": "%s",
		"service": "%s", 
		"version": "%s",
		"timestamp": "%s",
		"is_running": %t
	}`, 
		status["status"], 
		status["service"], 
		status["version"], 
		status["timestamp"], 
		status["is_running"])
}

func (s *Server) processorsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := s.processorMgr.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"processors": %v}`, status)
}

func (s *Server) processorHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle individual processor operations
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) queueStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := s.queueMgr.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"queue": %v}`, status)
}

func (s *Server) queueStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := s.queueMgr.GetStats()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"stats": %v}`, stats)
}

// Webhook Handlers

func (s *Server) smsWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle SMS webhook
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) voiceWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle Voice webhook
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) discordWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle Discord webhook
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) genericWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Handle generic webhook
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) outputsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := s.outputMgr.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"outputs": %v}`, status)
}

func (s *Server) sendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Parse request body and send message
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) broadcastHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Parse request body and broadcast message
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
