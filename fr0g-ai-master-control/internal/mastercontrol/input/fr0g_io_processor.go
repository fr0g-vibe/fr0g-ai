package input

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Fr0gIOWebhookProcessor processes webhooks from fr0g-ai-io service
type Fr0gIOWebhookProcessor struct {
	ioManager IOManagerInterface
}

// IOManagerInterface defines the interface for I/O manager operations
type IOManagerInterface interface {
	SendOutputCommand(command *OutputCommand) error
	SendThreatAnalysis(result *ThreatAnalysisResult) error
}

// NewFr0gIOWebhookProcessor creates a new fr0g-ai-io webhook processor
func NewFr0gIOWebhookProcessor(ioManager IOManagerInterface) *Fr0gIOWebhookProcessor {
	return &Fr0gIOWebhookProcessor{
		ioManager: ioManager,
	}
}

// ProcessWebhook processes incoming webhooks from fr0g-ai-io
func (p *Fr0gIOWebhookProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("Fr0g-AI-IO Processor: Processing webhook from fr0g-ai-io, ID: %s", request.ID)

	// Parse the webhook body to determine the event type
	bodyMap, ok := request.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid webhook body format")
	}

	eventType, ok := bodyMap["event_type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing event_type in webhook body")
	}

	var err error
	switch eventType {
	case "input_event":
		err = p.processInputEvent(ctx, bodyMap)
	case "status_update":
		err = p.processStatusUpdate(ctx, bodyMap)
	case "error_notification":
		err = p.processErrorNotification(ctx, bodyMap)
	default:
		log.Printf("Fr0g-AI-IO Processor: Unknown event type: %s", eventType)
		err = fmt.Errorf("unknown event type: %s", eventType)
	}

	if err != nil {
		return &WebhookResponse{
			Success:   false,
			Message:   fmt.Sprintf("Error processing webhook: %v", err),
			RequestID: request.ID,
			Timestamp: time.Now(),
		}, nil
	}

	return &WebhookResponse{
		Success:   true,
		Message:   "Webhook processed successfully",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"event_type": eventType,
			"processed":  true,
		},
		Timestamp: time.Now(),
	}, nil
}

// GetTag returns the processor tag
func (p *Fr0gIOWebhookProcessor) GetTag() string {
	return "fr0g-ai-io"
}

// GetDescription returns the processor description
func (p *Fr0gIOWebhookProcessor) GetDescription() string {
	return "Processes webhooks from fr0g-ai-io service for bidirectional I/O communication"
}

// processInputEvent processes input events from fr0g-ai-io
func (p *Fr0gIOWebhookProcessor) processInputEvent(ctx context.Context, data map[string]interface{}) error {
	// Extract input event data
	eventData, ok := data["event"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing event data")
	}

	// Convert to InputEvent struct
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	var inputEvent InputEvent
	if err := json.Unmarshal(eventJSON, &inputEvent); err != nil {
		return fmt.Errorf("failed to unmarshal input event: %w", err)
	}

	log.Printf("Fr0g-AI-IO Processor: Received input event %s of type %s", inputEvent.ID, inputEvent.Type)

	// Process the event based on its type
	response := p.generateEventResponse(&inputEvent)

	// Send output commands if any
	for _, action := range response.Actions {
		command := &OutputCommand{
			ID:       fmt.Sprintf("webhook_%s_%d", inputEvent.ID, time.Now().UnixNano()),
			Type:     action.Type,
			Target:   action.Target,
			Content:  action.Content,
			Metadata: action.Metadata,
			Priority: inputEvent.Priority,
		}

		if err := p.ioManager.SendOutputCommand(command); err != nil {
			log.Printf("Fr0g-AI-IO Processor: Error sending output command: %v", err)
		}
	}

	// Send threat analysis if available
	if response.Analysis != nil {
		if err := p.ioManager.SendThreatAnalysis(response.Analysis); err != nil {
			log.Printf("Fr0g-AI-IO Processor: Error sending threat analysis: %v", err)
		}
	}

	return nil
}

// processStatusUpdate processes status updates from fr0g-ai-io
func (p *Fr0gIOWebhookProcessor) processStatusUpdate(ctx context.Context, data map[string]interface{}) error {
	status := getStringFromMap(data, "status")
	service := getStringFromMap(data, "service")
	
	log.Printf("Fr0g-AI-IO Processor: Status update from %s: %s", service, status)
	
	// Handle status updates (could trigger alerts, logging, etc.)
	return nil
}

// processErrorNotification processes error notifications from fr0g-ai-io
func (p *Fr0gIOWebhookProcessor) processErrorNotification(ctx context.Context, data map[string]interface{}) error {
	errorMsg := getStringFromMap(data, "error")
	service := getStringFromMap(data, "service")
	severity := getStringFromMap(data, "severity")
	
	log.Printf("Fr0g-AI-IO Processor: Error notification from %s [%s]: %s", service, severity, errorMsg)
	
	// Handle error notifications (could trigger alerts, recovery actions, etc.)
	return nil
}

// generateEventResponse generates a response for an input event
func (p *Fr0gIOWebhookProcessor) generateEventResponse(event *InputEvent) *InputEventResponse {
	response := &InputEventResponse{
		EventID:     event.ID,
		Processed:   true,
		ProcessedAt: time.Now(),
		Actions:     []OutputAction{},
		Metadata:    make(map[string]interface{}),
	}

	// Generate appropriate response based on event type and content
	switch event.Type {
	case "sms":
		response.Actions = append(response.Actions, OutputAction{
			Type:    "sms",
			Target:  event.Source,
			Content: p.generateSMSResponse(event.Content),
			Metadata: map[string]interface{}{
				"response_type": "automated",
				"original_id":   event.ID,
			},
		})

	case "voice":
		response.Actions = append(response.Actions, OutputAction{
			Type:    "voice",
			Target:  event.Source,
			Content: p.generateVoiceResponse(event.Content),
			Metadata: map[string]interface{}{
				"response_type": "automated",
				"original_id":   event.ID,
			},
		})

	case "irc":
		channel := getStringFromMap(event.Metadata, "channel")
		isPrivate := getBoolFromMap(event.Metadata, "is_private")
		
		target := event.Source
		if !isPrivate && channel != "" {
			target = channel
		}

		response.Actions = append(response.Actions, OutputAction{
			Type:    "irc",
			Target:  target,
			Content: p.generateIRCResponse(event.Content),
			Metadata: map[string]interface{}{
				"response_type": "automated",
				"original_id":   event.ID,
				"is_private":    isPrivate,
			},
		})

	case "discord":
		channelID := getStringFromMap(event.Metadata, "channel_id")
		
		response.Actions = append(response.Actions, OutputAction{
			Type:    "discord",
			Target:  channelID,
			Content: p.generateDiscordResponse(event.Content),
			Metadata: map[string]interface{}{
				"response_type": "automated",
				"original_id":   event.ID,
			},
		})
	}

	// Perform basic threat analysis
	if p.shouldAnalyzeThreat(event) {
		response.Analysis = p.performThreatAnalysis(event)
	}

	return response
}

// Response generators for different message types

func (p *Fr0gIOWebhookProcessor) generateSMSResponse(content string) string {
	// Simple response generation - in a real implementation, this would use AI
	return fmt.Sprintf("Received your message: %s. Processing...", content)
}

func (p *Fr0gIOWebhookProcessor) generateVoiceResponse(content string) string {
	return fmt.Sprintf("Voice message received and transcribed: %s", content)
}

func (p *Fr0gIOWebhookProcessor) generateIRCResponse(content string) string {
	return fmt.Sprintf("IRC message processed: %s", content)
}

func (p *Fr0gIOWebhookProcessor) generateDiscordResponse(content string) string {
	return fmt.Sprintf("Discord message received: %s", content)
}

// Threat analysis

func (p *Fr0gIOWebhookProcessor) shouldAnalyzeThreat(event *InputEvent) bool {
	// Simple heuristics for when to perform threat analysis
	content := event.Content
	
	// Check for suspicious keywords
	suspiciousKeywords := []string{"hack", "attack", "malware", "virus", "exploit", "breach"}
	for _, keyword := range suspiciousKeywords {
		if contains(content, keyword) {
			return true
		}
	}
	
	// Check for unusual patterns
	if len(content) > 1000 || containsURLs(content) {
		return true
	}
	
	return false
}

func (p *Fr0gIOWebhookProcessor) performThreatAnalysis(event *InputEvent) *ThreatAnalysisResult {
	// Basic threat analysis - in a real implementation, this would be more sophisticated
	threatLevel := "low"
	threatScore := 0.1
	threatTypes := []string{}
	
	content := event.Content
	
	// Analyze content for threats
	if contains(content, "hack") || contains(content, "attack") {
		threatLevel = "medium"
		threatScore = 0.5
		threatTypes = append(threatTypes, "potential_attack")
	}
	
	if contains(content, "malware") || contains(content, "virus") {
		threatLevel = "high"
		threatScore = 0.8
		threatTypes = append(threatTypes, "malware_reference")
	}
	
	if containsURLs(content) {
		threatTypes = append(threatTypes, "suspicious_url")
		threatScore += 0.2
	}
	
	return &ThreatAnalysisResult{
		EventID:     event.ID,
		ThreatLevel: threatLevel,
		ThreatScore: threatScore,
		ThreatTypes: threatTypes,
		Indicators: []ThreatIndicator{
			{
				Type:        "content_analysis",
				Value:       "suspicious_keywords",
				Confidence:  0.7,
				Description: "Content contains potentially suspicious keywords",
			},
		},
		Mitigation: []string{
			"Monitor for additional suspicious activity",
			"Log event for further analysis",
		},
		Confidence:  0.7,
		Analysis:    fmt.Sprintf("Automated analysis of %s event from %s", event.Type, event.Source),
		AnalyzedAt:  time.Now(),
		Metadata: map[string]interface{}{
			"analyzer":    "fr0g-io-webhook-processor",
			"event_type":  event.Type,
			"source":      event.Source,
		},
	}
}

// Utility functions

func contains(text, substring string) bool {
	return len(text) >= len(substring) && 
		   (text == substring || 
		    (len(text) > len(substring) && 
		     (text[:len(substring)] == substring || 
		      text[len(text)-len(substring):] == substring ||
		      containsSubstring(text, substring))))
}

func containsSubstring(text, substring string) bool {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

func containsURLs(text string) bool {
	return contains(text, "http://") || contains(text, "https://") || contains(text, "www.")
}

func getBoolFromMap(m map[string]interface{}, key string) bool {
	if value, ok := m[key]; ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}
