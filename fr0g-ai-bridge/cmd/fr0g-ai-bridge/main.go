package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/api"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/client"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/config"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Command line flags
	var (
		configPath = flag.String("config", "", "Path to configuration file")
		httpOnly   = flag.Bool("http-only", false, "Run only HTTP REST server")
		grpcOnly   = flag.Bool("grpc-only", false, "Run only gRPC server")
		version    = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *version {
		fmt.Println("fr0g-ai-bridge v1.0.0")
		return
	}

	// Default to running both HTTP and gRPC if no specific mode is chosen
	if !*httpOnly && !*grpcOnly && len(os.Args) == 1 {
		log.Println("No mode specified, starting both HTTP and gRPC servers...")
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log security configuration warnings
	if cfg.Security.RequireAPIKey && len(cfg.Security.AllowedAPIKeys) == 0 {
		log.Println("WARNING: API key authentication enabled but no keys configured")
	}
	if cfg.Security.EnableReflection {
		log.Println("WARNING: gRPC reflection is enabled - disable in production")
	}
	if len(cfg.Security.AllowedOrigins) == 1 && cfg.Security.AllowedOrigins[0] == "*" {
		log.Println("WARNING: CORS allows all origins - restrict in production")
	}

	// Create OpenWebUI client
	openWebUIClient := client.NewOpenWebUIClient(
		cfg.OpenWebUI.BaseURL,
		cfg.OpenWebUI.APIKey,
		time.Duration(cfg.OpenWebUI.Timeout)*time.Second,
	)

	// TODO: Add service discovery client initialization here
	// discoveryClient := discovery.NewClient(&discovery.ClientConfig{
	//     RegistryURL:    "http://localhost:8500",
	//     ServiceName:    "fr0g-ai-bridge",
	//     ServiceID:      "bridge-001",
	//     ServiceAddress: "localhost",
	//     ServicePort:    cfg.Server.HTTPPort,
	//     Tags:           []string{"bridge", "ai", "chat"},
	// })

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to collect errors from servers
	errChan := make(chan error, 2)

	// Start HTTP REST server (unless grpc-only is specified)
	if !*grpcOnly {
		go func() {
			log.Printf("Starting HTTP REST server on %s:%d", cfg.Server.Host, cfg.Server.HTTPPort)

			restServer := api.NewRESTServer(openWebUIClient, cfg)

			// Add health check endpoint
			mux := http.NewServeMux()
			mux.Handle("/", restServer.GetRouter())
			mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{
					"status": "healthy",
					"service": "fr0g-ai-bridge",
					"timestamp": "%s",
					"version": "1.0.0"
				}`, time.Now().Format(time.RFC3339))
			})

			httpServer := &http.Server{
				Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort),
				Handler: mux,
			}

			// Start server in goroutine
			go func() {
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					errChan <- fmt.Errorf("HTTP server error: %w", err)
				}
			}()

			// Wait for context cancellation
			<-ctx.Done()

			// Graceful shutdown
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()

			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				log.Printf("HTTP server shutdown error: %v", err)
			} else {
				log.Println("HTTP server stopped gracefully")
			}
		}()
	}

	// Start gRPC server (unless http-only is specified)
	if !*httpOnly {
		go func() {
			log.Printf("Starting gRPC server on %s:%d", cfg.Server.Host, cfg.Server.GRPCPort)

			lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort))
			if err != nil {
				errChan <- fmt.Errorf("failed to listen on gRPC port: %w", err)
				return
			}

			grpcServer := grpc.NewServer()

			bridgeServer := api.NewGRPCServer(openWebUIClient)
			pb.RegisterFr0GAiBridgeServiceServer(grpcServer, bridgeServer)
			log.Printf("gRPC bridge service registered successfully")

			// Enable gRPC reflection for debugging and tools like grpcurl (if enabled)
			if cfg.Security.EnableReflection {
				reflection.Register(grpcServer)
				log.Printf("gRPC reflection enabled")
			}

			log.Printf("gRPC server ready with Fr0gAiBridgeService")

			// Start server in goroutine
			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					errChan <- fmt.Errorf("gRPC server error: %w", err)
				}
			}()

			// Wait for context cancellation
			<-ctx.Done()

			// Graceful shutdown
			log.Println("Shutting down gRPC server...")
			grpcServer.GracefulStop()
			log.Println("gRPC server stopped gracefully")
		}()
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
		cancel()
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		cancel()
	}

	// Give servers time to shut down gracefully
	time.Sleep(2 * time.Second)
	log.Println("fr0g-ai-bridge stopped")
}
