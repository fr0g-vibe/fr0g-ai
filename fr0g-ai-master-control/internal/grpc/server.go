package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"fr0g-ai-master-control/internal/mastercontrol/input"
)

// Fr0gIOInputServer implements the gRPC server for receiving input events from fr0g-ai-io
type Fr0gIOInputServer struct {
	server       *grpc.Server
	inputHandler input.Fr0gIOInputHandler
	config       *ServerConfig
	listener     net.Listener
}

// ServerConfig holds gRPC server configuration
type ServerConfig struct {
	Host string
	Port int
}

// NewFr0gIOInputServer creates a new gRPC server for input events
func NewFr0gIOInputServer(config *ServerConfig, inputHandler input.Fr0gIOInputHandler) *Fr0gIOInputServer {
	server := grpc.NewServer()
	
	inputServer := &Fr0gIOInputServer{
		server:       server,
		inputHandler: inputHandler,
		config:       config,
	}

	// For now, we'll simulate the gRPC server without actual protobuf registration
	log.Printf("gRPC Input Server: Initialized server for %s:%d", config.Host, config.Port)

	return inputServer
}

// Start starts the gRPC server
func (s *Fr0gIOInputServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	s.listener = listener
	log.Printf("gRPC Input Server: Starting server on %s (simulated)", address)
	
	// For now, just simulate the server running
	go func() {
		log.Printf("gRPC Input Server: Server running on %s (ready to receive fr0g-ai-io events)", address)
		// In a real implementation, this would be: s.server.Serve(listener)
		// For now, we'll just keep the listener open
		for {
			time.Sleep(30 * time.Second)
			log.Printf("gRPC Input Server: Heartbeat - ready for fr0g-ai-io connections")
		}
	}()

	return nil
}

// Stop stops the gRPC server
func (s *Fr0gIOInputServer) Stop() {
	log.Println("gRPC Input Server: Stopping server...")
	if s.listener != nil {
		s.listener.Close()
	}
	s.server.GracefulStop()
}

// SimulateInputEvent simulates receiving an input event from fr0g-ai-io
func (s *Fr0gIOInputServer) SimulateInputEvent(eventType, source, content string) error {
	log.Printf("gRPC Input Server: Simulating %s input event from %s", eventType, source)

	// Create a simulated input event
	event := &input.InputEvent{
		ID:        fmt.Sprintf("sim_%d", time.Now().UnixNano()),
		Type:      eventType,
		Source:    source,
		Content:   content,
		Timestamp: time.Now(),
		Priority:  1,
		Metadata:  map[string]interface{}{"simulated": true},
	}

	// Process the event through the input handler
	_, err := s.inputHandler.HandleInputEvent(context.Background(), event)
	if err != nil {
		log.Printf("gRPC Input Server: Failed to process simulated event: %v", err)
		return err
	}

	log.Printf("gRPC Input Server: Successfully processed simulated %s event", eventType)
	return nil
}

// StartSimulation starts simulating input events for demonstration
func (s *Fr0gIOInputServer) StartSimulation() {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		events := []struct {
			eventType string
			source    string
			content   string
		}{
			{"sms", "+1234567890", "Hello from SMS simulation"},
			{"discord", "user#1234", "Discord message simulation"},
			{"irc", "testuser@irc.example.com", "IRC channel message simulation"},
			{"voice", "+1987654321", "Voice message transcription simulation"},
		}

		eventIndex := 0
		for {
			select {
			case <-ticker.C:
				event := events[eventIndex%len(events)]
				s.SimulateInputEvent(event.eventType, event.source, event.content)
				eventIndex++
			}
		}
	}()
}
