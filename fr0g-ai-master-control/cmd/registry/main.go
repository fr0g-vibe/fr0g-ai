package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/registry"
)

func main() {
	var (
		port = flag.Int("port", 8500, "Registry port")
		host = flag.String("host", "0.0.0.0", "Registry host")
	)
	flag.Parse()

	// Create registry configuration
	config := &registry.RegistryConfig{
		Port:           *port,
		Host:           *host,
		HealthInterval: 30 * time.Second,
		ServiceTTL:     2 * time.Minute,
		EnableHTTPAPI:  true,
	}

	// Create and start registry
	serviceRegistry := registry.NewServiceRegistry(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := serviceRegistry.Start(ctx); err != nil {
		log.Fatalf("Failed to start registry: %v", err)
	}

	log.Printf("🔍 fr0g.ai Service Registry started on %s:%d", *host, *port)

	// Wait for shutdown signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("🛑 Shutting down service registry...")
	cancel()
}
