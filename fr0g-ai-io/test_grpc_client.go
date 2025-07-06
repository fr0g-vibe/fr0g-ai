package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/grpc"
)

func main() {
	fmt.Println("=== gRPC MCP Client Test for fr0g-ai-io ===")
	fmt.Println()

	// Create MCP client configuration
	config := &grpc.MCPClientConfig{
		Host:        "localhost",
		Port:        9092, // Master-control gRPC port
		Timeout:     30 * time.Second,
		MaxRetries:  3,
		ServiceName: "fr0g-ai-io",
	}

	// Create and connect MCP client
	fmt.Println("Test 1: Connecting to master-control...")
	client, err := grpc.NewMCPGRPCClient(config)
	if err != nil {
		log.Printf("Failed to connect to master-control: %v", err)
		fmt.Println("Note: Make sure master-control is running on localhost:9092")
		return
	}
	defer client.Close()

	fmt.Printf("✓ Connected to master-control at %s:%d\n", config.Host, config.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 2: Send SMS input event
	fmt.Println("\nTest 2: Sending SMS input event...")
	smsEvent := &grpc.InputEvent{
		ID:      fmt.Sprintf("test-sms-%d", time.Now().Unix()),
		Type:    "sms",
		Source:  "+1234567890",
		Content: "Test SMS message from gRPC client",
		Metadata: map[string]interface{}{
			"test_type": "grpc_client",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		Timestamp: time.Now(),
		Priority:  1,
	}

	smsResp, err := client.SendInputEvent(ctx, smsEvent)
	if err != nil {
		log.Printf("Failed to send SMS event: %v", err)
	} else {
		fmt.Printf("✓ SMS event processed: %s\n", smsResp.EventID)
		fmt.Printf("  Actions generated: %d\n", len(smsResp.Actions))
		for _, action := range smsResp.Actions {
			fmt.Printf("  - %s to %s: %s\n", action.Type, action.Target, action.Content)
		}
	}

	// Test 3: Send IRC input event
	fmt.Println("\nTest 3: Sending IRC input event...")
	ircEvent := &grpc.InputEvent{
		ID:      fmt.Sprintf("test-irc-%d", time.Now().Unix()),
		Type:    "irc",
		Source:  "testuser",
		Content: "!help command test",
		Metadata: map[string]interface{}{
			"channel":   "#test",
			"server":    "irc.example.com",
			"test_type": "grpc_client",
		},
		Timestamp: time.Now(),
		Priority:  2,
	}

	ircResp, err := client.SendInputEvent(ctx, ircEvent)
	if err != nil {
		log.Printf("Failed to send IRC event: %v", err)
	} else {
		fmt.Printf("✓ IRC event processed: %s\n", ircResp.EventID)
		fmt.Printf("  Actions generated: %d\n", len(ircResp.Actions))
		for _, action := range ircResp.Actions {
			fmt.Printf("  - %s to %s: %s\n", action.Type, action.Target, action.Content)
		}
	}

	// Test 4: Send Discord input event
	fmt.Println("\nTest 4: Sending Discord input event...")
	discordEvent := &grpc.InputEvent{
		ID:      fmt.Sprintf("test-discord-%d", time.Now().Unix()),
		Type:    "discord",
		Source:  "user123",
		Content: "Hello from Discord test",
		Metadata: map[string]interface{}{
			"channel_id": "123456789",
			"guild_id":   "987654321",
			"test_type":  "grpc_client",
		},
		Timestamp: time.Now(),
		Priority:  1,
	}

	discordResp, err := client.SendInputEvent(ctx, discordEvent)
	if err != nil {
		log.Printf("Failed to send Discord event: %v", err)
	} else {
		fmt.Printf("✓ Discord event processed: %s\n", discordResp.EventID)
		fmt.Printf("  Actions generated: %d\n", len(discordResp.Actions))
		for _, action := range discordResp.Actions {
			fmt.Printf("  - %s to %s: %s\n", action.Type, action.Target, action.Content)
		}
	}

	// Test 5: Send Voice input event
	fmt.Println("\nTest 5: Sending Voice input event...")
	voiceEvent := &grpc.InputEvent{
		ID:      fmt.Sprintf("test-voice-%d", time.Now().Unix()),
		Type:    "voice",
		Source:  "+1987654321",
		Content: "Voice message transcription test",
		Metadata: map[string]interface{}{
			"duration":  "15.5",
			"quality":   "high",
			"test_type": "grpc_client",
		},
		Timestamp: time.Now(),
		Priority:  3,
	}

	voiceResp, err := client.SendInputEvent(ctx, voiceEvent)
	if err != nil {
		log.Printf("Failed to send Voice event: %v", err)
	} else {
		fmt.Printf("✓ Voice event processed: %s\n", voiceResp.EventID)
		fmt.Printf("  Actions generated: %d\n", len(voiceResp.Actions))
		for _, action := range voiceResp.Actions {
			fmt.Printf("  - %s to %s: %s\n", action.Type, action.Target, action.Content)
		}
	}

	// Test 6: Check client status
	fmt.Println("\nTest 6: Checking client status...")
	status := client.GetStatus()
	fmt.Printf("✓ Client status:\n")
	for key, value := range status {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// Test 7: Test connection status
	fmt.Println("\nTest 7: Testing connection status...")
	if client.IsConnected() {
		fmt.Printf("✓ Client is connected to master-control\n")
	} else {
		fmt.Printf("⚠ Client is not connected to master-control\n")
	}

	fmt.Println("\n=== gRPC MCP Client Tests Completed ===")
}
