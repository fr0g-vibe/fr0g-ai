package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fr0g-ai-master-control/internal/mastercontrol"
)

func main() {
	fmt.Println("🌐 Webhook Input System Demo")
	fmt.Println("============================")
	
	// Load MCP configuration
	config := mastercontrol.DefaultMCPConfig()
	
	// Create and start MCP
	mcp := mastercontrol.NewMasterControlProgram(config)
	
	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start MCP: %v", err)
	}
	
	fmt.Println("✅ MCP with webhook input system started successfully")
	fmt.Printf("🌐 Webhook server listening on http://%s:%d\n", 
		config.Input.Webhook.Host, config.Input.Webhook.Port)
	fmt.Println()
	
	// Wait a moment for server to start
	time.Sleep(time.Second * 2)
	
	// Demonstrate webhook functionality
	demonstrateWebhooks(config.Input.Webhook.Port)
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	fmt.Println("🎯 Webhook system is running...")
	fmt.Println("   Available endpoints:")
	fmt.Printf("   - POST http://localhost:%d/webhook/discord\n", config.Input.Webhook.Port)
	fmt.Printf("   - GET  http://localhost:%d/health\n", config.Input.Webhook.Port)
	fmt.Printf("   - GET  http://localhost:%d/status\n", config.Input.Webhook.Port)
	fmt.Println("   Press Ctrl+C to shutdown")
	fmt.Println()
	
	// Wait for shutdown signal
	<-c
	
	fmt.Println("\n🛑 Shutdown signal received...")
	
	if err := mcp.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	
	fmt.Println("👋 Webhook demo complete")
}

func demonstrateWebhooks(port int) {
	fmt.Println("🔍 Demonstrating Webhook Functionality:")
	fmt.Println("---------------------------------------")
	
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	
	// Test health endpoint
	fmt.Println("📊 Testing health endpoint...")
	if err := testHealthEndpoint(baseURL); err != nil {
		log.Printf("Health check failed: %v", err)
	}
	
	// Test status endpoint
	fmt.Println("📋 Testing status endpoint...")
	if err := testStatusEndpoint(baseURL); err != nil {
		log.Printf("Status check failed: %v", err)
	}
	
	// Test Discord webhook
	fmt.Println("💬 Testing Discord webhook...")
	if err := testDiscordWebhook(baseURL); err != nil {
		log.Printf("Discord webhook test failed: %v", err)
	}
	
	// Test unknown webhook tag
	fmt.Println("❓ Testing unknown webhook tag...")
	if err := testUnknownWebhook(baseURL); err != nil {
		log.Printf("Unknown webhook test completed: %v", err)
	}
	
	fmt.Println()
}

func testHealthEndpoint(baseURL string) error {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Health endpoint responding correctly")
	} else {
		fmt.Printf("❌ Health endpoint returned status: %d\n", resp.StatusCode)
	}
	
	return nil
}

func testStatusEndpoint(baseURL string) error {
	resp, err := http.Get(baseURL + "/status")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		var status map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&status); err == nil {
			fmt.Println("✅ Status endpoint responding correctly")
			if processors, ok := status["processors"].(map[string]interface{}); ok {
				fmt.Printf("   Registered processors: %d\n", len(processors))
				for tag, desc := range processors {
					fmt.Printf("   - %s: %v\n", tag, desc)
				}
			}
		}
	} else {
		fmt.Printf("❌ Status endpoint returned status: %d\n", resp.StatusCode)
	}
	
	return nil
}

func testDiscordWebhook(baseURL string) error {
	// Create a mock Discord message
	discordMessage := map[string]interface{}{
		"content": "Hello from Discord! This is a test message that should be reviewed by the AI community.",
		"author": map[string]interface{}{
			"username": "testuser",
			"id":       "123456789",
		},
		"channel_id": "987654321",
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	
	jsonData, err := json.Marshal(discordMessage)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(baseURL+"/webhook/discord", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Discord webhook processed successfully")
		if success, ok := response["success"].(bool); ok && success {
			if data, ok := response["data"].(map[string]interface{}); ok {
				if action, ok := data["action"].(string); ok {
					fmt.Printf("   Action determined: %s\n", action)
				}
				if personaCount, ok := data["persona_count"].(float64); ok {
					fmt.Printf("   AI personas consulted: %.0f\n", personaCount)
				}
			}
		}
	} else {
		fmt.Printf("❌ Discord webhook failed with status: %d\n", resp.StatusCode)
		fmt.Printf("   Response: %v\n", response)
	}
	
	return nil
}

func testUnknownWebhook(baseURL string) error {
	testData := map[string]interface{}{
		"message": "This is a test for an unknown webhook tag",
	}
	
	jsonData, err := json.Marshal(testData)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(baseURL+"/webhook/unknown", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("✅ Unknown webhook tag correctly rejected")
	} else {
		fmt.Printf("❌ Expected 404 for unknown tag, got: %d\n", resp.StatusCode)
	}
	
	return nil
}
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fr0g-ai-master-control/internal/mastercontrol"
)

func main() {
	fmt.Println("🌐 Webhook Input System Demo")
	fmt.Println("============================")
	
	// Load MCP configuration
	config := mastercontrol.DefaultMCPConfig()
	
	// Create and start MCP
	mcp := mastercontrol.NewMasterControlProgram(config)
	
	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start MCP: %v", err)
	}
	
	fmt.Println("✅ MCP with webhook input system started successfully")
	fmt.Printf("🌐 Webhook server listening on http://%s:%d\n", 
		config.Input.Webhook.Host, config.Input.Webhook.Port)
	fmt.Println()
	
	// Wait a moment for server to start
	time.Sleep(time.Second * 2)
	
	// Demonstrate webhook functionality
	demonstrateWebhooks(config.Input.Webhook.Port)
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	fmt.Println("🎯 Webhook system is running...")
	fmt.Println("   Available endpoints:")
	fmt.Printf("   - POST http://localhost:%d/webhook/discord\n", config.Input.Webhook.Port)
	fmt.Printf("   - GET  http://localhost:%d/health\n", config.Input.Webhook.Port)
	fmt.Printf("   - GET  http://localhost:%d/status\n", config.Input.Webhook.Port)
	fmt.Println("   Press Ctrl+C to shutdown")
	fmt.Println()
	
	// Wait for shutdown signal
	<-c
	
	fmt.Println("\n🛑 Shutdown signal received...")
	
	if err := mcp.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	
	fmt.Println("👋 Webhook demo complete")
}

func demonstrateWebhooks(port int) {
	fmt.Println("🔍 Demonstrating Webhook Functionality:")
	fmt.Println("---------------------------------------")
	
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	
	// Test health endpoint
	fmt.Println("📊 Testing health endpoint...")
	if err := testHealthEndpoint(baseURL); err != nil {
		log.Printf("Health check failed: %v", err)
	}
	
	// Test status endpoint
	fmt.Println("📋 Testing status endpoint...")
	if err := testStatusEndpoint(baseURL); err != nil {
		log.Printf("Status check failed: %v", err)
	}
	
	// Test Discord webhook
	fmt.Println("💬 Testing Discord webhook...")
	if err := testDiscordWebhook(baseURL); err != nil {
		log.Printf("Discord webhook test failed: %v", err)
	}
	
	// Test unknown webhook tag
	fmt.Println("❓ Testing unknown webhook tag...")
	if err := testUnknownWebhook(baseURL); err != nil {
		log.Printf("Unknown webhook test completed: %v", err)
	}
	
	fmt.Println()
}

func testHealthEndpoint(baseURL string) error {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Health endpoint responding correctly")
	} else {
		fmt.Printf("❌ Health endpoint returned status: %d\n", resp.StatusCode)
	}
	
	return nil
}

func testStatusEndpoint(baseURL string) error {
	resp, err := http.Get(baseURL + "/status")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		var status map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&status); err == nil {
			fmt.Println("✅ Status endpoint responding correctly")
			if processors, ok := status["processors"].(map[string]interface{}); ok {
				fmt.Printf("   Registered processors: %d\n", len(processors))
				for tag, desc := range processors {
					fmt.Printf("   - %s: %v\n", tag, desc)
				}
			}
		}
	} else {
		fmt.Printf("❌ Status endpoint returned status: %d\n", resp.StatusCode)
	}
	
	return nil
}

func testDiscordWebhook(baseURL string) error {
	// Create a mock Discord message
	discordMessage := map[string]interface{}{
		"content": "Hello from Discord! This is a test message that should be reviewed by the AI community.",
		"author": map[string]interface{}{
			"username": "testuser",
			"id":       "123456789",
		},
		"channel_id": "987654321",
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	
	jsonData, err := json.Marshal(discordMessage)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(baseURL+"/webhook/discord", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Discord webhook processed successfully")
		if success, ok := response["success"].(bool); ok && success {
			if data, ok := response["data"].(map[string]interface{}); ok {
				if action, ok := data["action"].(string); ok {
					fmt.Printf("   Action determined: %s\n", action)
				}
				if personaCount, ok := data["persona_count"].(float64); ok {
					fmt.Printf("   AI personas consulted: %.0f\n", personaCount)
				}
			}
		}
	} else {
		fmt.Printf("❌ Discord webhook failed with status: %d\n", resp.StatusCode)
		fmt.Printf("   Response: %v\n", response)
	}
	
	return nil
}

func testUnknownWebhook(baseURL string) error {
	testData := map[string]interface{}{
		"message": "This is a test for an unknown webhook tag",
	}
	
	jsonData, err := json.Marshal(testData)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(baseURL+"/webhook/unknown", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("✅ Unknown webhook tag correctly rejected")
	} else {
		fmt.Printf("❌ Expected 404 for unknown tag, got: %d\n", resp.StatusCode)
	}
	
	return nil
}
