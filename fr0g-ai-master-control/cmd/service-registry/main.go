package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fr0g-ai-master-control/internal/registry"
)

func main() {
	var (
		addr = flag.String("addr", ":8500", "Server address")
		help = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create and start registry server
	server := registry.NewServer()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received, stopping service registry...")
		cancel()
	}()

	log.Printf("fr0g.ai fr0g.ai Service Registry starting on %s", *addr)
	if err := server.Start(ctx, *addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("COMPLETED fr0g.ai Service Registry stopped")
}
