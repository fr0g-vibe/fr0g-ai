package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/processors"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/types"
)

func main() {
	fmt.Println("Testing fr0g-ai-io processors...")

	// Load configuration with defaults
	cfg := sharedconfig.GetDefaults()

	// Create processor manager
	processorMgr, err := processors.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create processor manager: %v", err)
	}

	fmt.Printf("Processor manager initialized with %d processors\n", len(processorMgr.GetProcessors()))

	// Test each processor type
	testCases := []struct {
		eventType string
		content   string
		metadata  map[string]interface{}
	}{
		{
			eventType: "sms",
			content:   "URGENT! You've won $1000! Click here to claim your prize now! Reply STOP to opt out. bit.ly/scam123",
			metadata: map[string]interface{}{
				"from": "+1234567890",
				"to":   "+0987654321",
			},
		},
		{
			eventType: "voice",
			content:   "This is the IRS calling about your tax refund. You owe money and will be arrested if you don't pay immediately. Press 1 to speak with an agent or face legal action.",
			metadata: map[string]interface{}{
				"from":     "+18001234567",
				"to":       "+1234567890",
				"duration": 30 * time.Second,
			},
		},
		{
			eventType: "irc",
			content:   "Check out this amazing deal! Download malware.exe from http://bit.ly/malware123 - FREE MONEY GUARANTEED!",
			metadata: map[string]interface{}{
				"server":       "irc.example.com",
				"channel":      "#general",
				"nick":         "spammer123",
				"message_type": "PRIVMSG",
			},
		},
		{
			eventType: "discord",
			content:   "Free Discord Nitro! Click here for your free gift: discord.gg/scam123",
			metadata: map[string]interface{}{
				"guild_id":     "123456789",
				"channel_id":   "987654321",
				"user_id":      "555666777",
				"username":     "scammer",
				"message_type": "DEFAULT",
			},
		},
		{
			eventType: "esmtp",
			content:   "Your account has been suspended due to suspicious activity. Click here to verify your credentials immediately: http://bit.ly/phishing123 or your account will be permanently deleted.",
			metadata: map[string]interface{}{
				"from":    "noreply@phishing-site.com",
				"to":      []string{"victim@example.com"},
				"subject": "URGENT: Account Verification Required - Immediate Action Needed",
				"headers": map[string]string{
					"X-Mailer": "Suspicious Mailer 1.0",
				},
			},
		},
	}

	fmt.Println("\nTesting processors with sample threat data:")
	fmt.Println(strings.Repeat("=", 60))

	for _, testCase := range testCases {
		fmt.Printf("\nTesting %s processor:\n", testCase.eventType)
		fmt.Printf("Content: %s\n", testCase.content)

		// Create test event
		event := &types.InputEvent{
			ID:        fmt.Sprintf("test-%s-%d", testCase.eventType, time.Now().Unix()),
			Type:      testCase.eventType,
			Source:    "test-source",
			Content:   testCase.content,
			Metadata:  testCase.metadata,
			Timestamp: time.Now(),
			Priority:  1,
		}

		// Process the event
		response, err := processorMgr.ProcessEvent(event)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		// Display results
		fmt.Printf("Processed: %t\n", response.Processed)
		if response.Processed {
			fmt.Printf("Actions generated: %d\n", len(response.Actions))
			for i, action := range response.Actions {
				fmt.Printf("  Action %d: %s -> %s (Priority: %d)\n", 
					i+1, action.Type, action.Target, action.Priority)
			}
			
			if threatLevel, ok := response.Metadata["threat_level"].(string); ok {
				fmt.Printf("Threat Level: %s\n", threatLevel)
			}
			if confidence, ok := response.Metadata["confidence"].(float64); ok {
				fmt.Printf("Confidence: %.2f\n", confidence)
			}
		} else {
			if errorMsg, ok := response.Metadata["error"].(string); ok {
				fmt.Printf("Processing Error: %s\n", errorMsg)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Processor testing completed!")

	// Display processor status
	status := processorMgr.GetStatus()
	fmt.Printf("\nProcessor Status Summary:\n")
	fmt.Printf("Total processors: %v\n", status["total_processors"])
	
	if processors, ok := status["processors"].(map[string]interface{}); ok {
		for name, info := range processors {
			if procInfo, ok := info.(map[string]interface{}); ok {
				enabled := procInfo["enabled"].(bool)
				procType := procInfo["type"].(string)
				fmt.Printf("- %s (%s): %s\n", name, procType, 
					map[bool]string{true: "ENABLED", false: "DISABLED"}[enabled])
			}
		}
	}
}
