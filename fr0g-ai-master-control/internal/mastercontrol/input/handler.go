package input

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Fr0gIOInputHandlerImpl implements the Fr0gIOInputHandler interface
type Fr0gIOInputHandlerImpl struct {
	workflowEngine WorkflowEngine
}

// WorkflowEngine interface for processing input events
type WorkflowEngine interface {
	ProcessInputEvent(ctx context.Context, event *InputEvent) (*InputEventResponse, error)
}

// NewFr0gIOInputHandler creates a new input handler
func NewFr0gIOInputHandler(workflowEngine WorkflowEngine) *Fr0gIOInputHandlerImpl {
	return &Fr0gIOInputHandlerImpl{
		workflowEngine: workflowEngine,
	}
}

// SetWorkflowEngine sets the workflow engine (used for circular dependency resolution)
func (h *Fr0gIOInputHandlerImpl) SetWorkflowEngine(workflowEngine WorkflowEngine) {
	h.workflowEngine = workflowEngine
}

// HandleInputEvent processes generic input events
func (h *Fr0gIOInputHandlerImpl) HandleInputEvent(ctx context.Context, event *InputEvent) (*InputEventResponse, error) {
	log.Printf("Input Handler: Processing input event %s of type %s from %s", event.ID, event.Type, event.Source)

	// Delegate to workflow engine for processing
	if h.workflowEngine != nil {
		return h.workflowEngine.ProcessInputEvent(ctx, event)
	}

	// Fallback processing if no workflow engine available
	return &InputEventResponse{
		EventID:     event.ID,
		Processed:   true,
		Actions:     []OutputAction{},
		Metadata:    map[string]interface{}{"handler": "fallback"},
		ProcessedAt: time.Now(),
	}, nil
}

// HandleSMSMessage processes SMS messages
func (h *Fr0gIOInputHandlerImpl) HandleSMSMessage(ctx context.Context, message *SMSMessage) error {
	log.Printf("Input Handler: Processing SMS message %s from %s", message.ID, message.From)

	// Convert SMS message to generic input event
	event := &InputEvent{
		ID:      message.ID,
		Type:    "sms",
		Source:  message.From,
		Content: message.Body,
		Metadata: map[string]interface{}{
			"to":           message.To,
			"message_sid":  message.MessageSID,
			"message_type": message.MessageType,
			"status":       message.Status,
			"direction":    message.Direction,
		},
		Timestamp: message.Timestamp,
		Priority:  1, // Default priority for SMS
	}

	// Process through workflow engine
	_, err := h.HandleInputEvent(ctx, event)
	return err
}

// HandleVoiceMessage processes voice messages
func (h *Fr0gIOInputHandlerImpl) HandleVoiceMessage(ctx context.Context, message *VoiceMessage) error {
	log.Printf("Input Handler: Processing voice message %s from %s (duration: %.2fs)", message.ID, message.From, message.RecordingDuration)

	// Convert voice message to generic input event
	event := &InputEvent{
		ID:      message.ID,
		Type:    "voice",
		Source:  message.From,
		Content: message.Transcription, // Use transcription as content
		Metadata: map[string]interface{}{
			"to":                message.To,
			"duration":          message.RecordingDuration,
			"format":            message.AudioFormat,
			"call_sid":          message.CallSID,
			"recording_url":     message.RecordingURL,
			"confidence":        message.Confidence,
			"language":          message.Language,
			"file_size":         message.FileSize,
			"direction":         message.Direction,
			"status":            message.Status,
		},
		Timestamp: message.Timestamp,
		Priority:  2, // Higher priority for voice messages
	}

	// Process through workflow engine
	_, err := h.HandleInputEvent(ctx, event)
	return err
}

// HandleIRCMessage processes IRC messages
func (h *Fr0gIOInputHandlerImpl) HandleIRCMessage(ctx context.Context, message *IRCMessage) error {
	log.Printf("Input Handler: Processing IRC message %s from %s in %s", message.ID, message.From, message.Channel)

	// Convert IRC message to generic input event
	event := &InputEvent{
		ID:      message.ID,
		Type:    "irc",
		Source:  fmt.Sprintf("%s@%s", message.From, message.Server),
		Content: message.Message,
		Metadata: map[string]interface{}{
			"server":     message.Server,
			"channel":    message.Channel,
			"from":       message.From,
			"to":         message.To,
			"type":       message.Type,
			"is_private": message.IsPrivate,
		},
		Timestamp: message.Timestamp,
		Priority:  1, // Default priority for IRC
	}

	// Process through workflow engine
	_, err := h.HandleInputEvent(ctx, event)
	return err
}

// HandleDiscordMessage processes Discord messages
func (h *Fr0gIOInputHandlerImpl) HandleDiscordMessage(ctx context.Context, message *DiscordMessage) error {
	log.Printf("Input Handler: Processing Discord message %s from %s in guild %s", message.ID, message.Username, message.GuildID)

	// Convert Discord message to generic input event
	event := &InputEvent{
		ID:      message.ID,
		Type:    "discord",
		Source:  fmt.Sprintf("%s#%s", message.Username, message.UserID),
		Content: message.Content,
		Metadata: map[string]interface{}{
			"guild_id":     message.GuildID,
			"channel_id":   message.ChannelID,
			"user_id":      message.UserID,
			"username":     message.Username,
			"message_type": message.MessageType,
		},
		Timestamp: message.Timestamp,
		Priority:  1, // Default priority for Discord
	}

	// Process through workflow engine
	_, err := h.HandleInputEvent(ctx, event)
	return err
}
