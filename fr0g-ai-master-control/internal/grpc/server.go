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
}

// ServerConfig holds gRPC server configuration
type ServerConfig struct {
	Host string
	Port int
}

// gRPC message types for input events
type InputEventRequest struct {
	Event *InputEvent `json:"event"`
}

type InputEventResponse struct {
	Success     bool           `json:"success"`
	Message     string         `json:"message"`
	EventId     string         `json:"event_id"`
	Actions     []*OutputAction `json:"actions"`
	ProcessedAt int64          `json:"processed_at"`
}

type InputEvent struct {
	Id        string            `json:"id"`
	Type      string            `json:"type"`
	Source    string            `json:"source"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata"`
	Timestamp int64             `json:"timestamp"`
	Priority  int32             `json:"priority"`
}

type OutputAction struct {
	Type     string            `json:"type"`
	Target   string            `json:"target"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

type SMSMessageRequest struct {
	Message *SMSMessage `json:"message"`
}

type SMSMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SMSMessage struct {
	Id          string `json:"id"`
	From        string `json:"from"`
	To          string `json:"to"`
	Content     string `json:"content"`
	Timestamp   int64  `json:"timestamp"`
	Provider    string `json:"provider"`
	MessageType string `json:"message_type"`
}

type VoiceMessageRequest struct {
	Message *VoiceMessage `json:"message"`
}

type VoiceMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VoiceMessage struct {
	Id            string  `json:"id"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	AudioData     []byte  `json:"audio_data"`
	Transcription string  `json:"transcription"`
	Duration      float64 `json:"duration"`
	Timestamp     int64   `json:"timestamp"`
	Format        string  `json:"format"`
}

type IRCMessageRequest struct {
	Message *IRCMessage `json:"message"`
}

type IRCMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type IRCMessage struct {
	Id        string `json:"id"`
	Server    string `json:"server"`
	Channel   string `json:"channel"`
	Nick      string `json:"nick"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	IsPrivate bool   `json:"is_private"`
}

type DiscordMessageRequest struct {
	Message *DiscordMessage `json:"message"`
}

type DiscordMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DiscordMessage struct {
	Id          string `json:"id"`
	GuildId     string `json:"guild_id"`
	ChannelId   string `json:"channel_id"`
	UserId      string `json:"user_id"`
	Username    string `json:"username"`
	Content     string `json:"content"`
	Timestamp   int64  `json:"timestamp"`
	MessageType string `json:"message_type"`
}

// Fr0gIOInputServiceServer defines the gRPC service interface
type Fr0gIOInputServiceServer interface {
	HandleInputEvent(context.Context, *InputEventRequest) (*InputEventResponse, error)
	HandleSMSMessage(context.Context, *SMSMessageRequest) (*SMSMessageResponse, error)
	HandleVoiceMessage(context.Context, *VoiceMessageRequest) (*VoiceMessageResponse, error)
	HandleIRCMessage(context.Context, *IRCMessageRequest) (*IRCMessageResponse, error)
	HandleDiscordMessage(context.Context, *DiscordMessageRequest) (*DiscordMessageResponse, error)
}

// NewFr0gIOInputServer creates a new gRPC server for input events
func NewFr0gIOInputServer(config *ServerConfig, inputHandler input.Fr0gIOInputHandler) *Fr0gIOInputServer {
	server := grpc.NewServer()
	
	inputServer := &Fr0gIOInputServer{
		server:       server,
		inputHandler: inputHandler,
		config:       config,
	}

	// Register the service (in a real implementation, this would use generated code)
	// For now, we'll implement a mock registration
	log.Printf("gRPC Input Server: Registered Fr0gIOInputService")

	return inputServer
}

// Start starts the gRPC server
func (s *Fr0gIOInputServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	log.Printf("gRPC Input Server: Starting server on %s", address)
	
	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Printf("gRPC Input Server: Server error: %v", err)
		}
	}()

	return nil
}

// Stop stops the gRPC server
func (s *Fr0gIOInputServer) Stop() {
	log.Println("gRPC Input Server: Stopping server...")
	s.server.GracefulStop()
}

// HandleInputEvent handles generic input events from fr0g-ai-io
func (s *Fr0gIOInputServer) HandleInputEvent(ctx context.Context, req *InputEventRequest) (*InputEventResponse, error) {
	log.Printf("gRPC Input Server: Received input event %s of type %s", req.Event.Id, req.Event.Type)

	// Convert gRPC message to internal type
	event := &input.InputEvent{
		ID:        req.Event.Id,
		Type:      req.Event.Type,
		Source:    req.Event.Source,
		Content:   req.Event.Content,
		Timestamp: time.Unix(req.Event.Timestamp, 0),
		Priority:  int(req.Event.Priority),
		Metadata:  make(map[string]interface{}),
	}

	// Convert metadata
	for k, v := range req.Event.Metadata {
		event.Metadata[k] = v
	}

	// Process the event
	response, err := s.inputHandler.HandleInputEvent(ctx, event)
	if err != nil {
		return &InputEventResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to process input event: %v", err),
			EventId: req.Event.Id,
		}, nil
	}

	// Convert response actions
	grpcActions := make([]*OutputAction, len(response.Actions))
	for i, action := range response.Actions {
		grpcActions[i] = &OutputAction{
			Type:     action.Type,
			Target:   action.Target,
			Content:  action.Content,
			Metadata: make(map[string]string),
		}
		
		// Convert metadata
		for k, v := range action.Metadata {
			if str, ok := v.(string); ok {
				grpcActions[i].Metadata[k] = str
			}
		}
	}

	return &InputEventResponse{
		Success:     true,
		Message:     "Input event processed successfully",
		EventId:     response.EventID,
		Actions:     grpcActions,
		ProcessedAt: response.ProcessedAt.Unix(),
	}, nil
}

// HandleSMSMessage handles SMS messages from fr0g-ai-io
func (s *Fr0gIOInputServer) HandleSMSMessage(ctx context.Context, req *SMSMessageRequest) (*SMSMessageResponse, error) {
	log.Printf("gRPC Input Server: Received SMS message %s from %s", req.Message.Id, req.Message.From)

	// Convert gRPC message to internal type
	smsMessage := &input.SMSMessage{
		ID:          req.Message.Id,
		From:        req.Message.From,
		To:          req.Message.To,
		Content:     req.Message.Content,
		Timestamp:   time.Unix(req.Message.Timestamp, 0),
		Provider:    req.Message.Provider,
		MessageType: req.Message.MessageType,
	}

	// Process the SMS message
	err := s.inputHandler.HandleSMSMessage(ctx, smsMessage)
	if err != nil {
		return &SMSMessageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to process SMS message: %v", err),
		}, nil
	}

	return &SMSMessageResponse{
		Success: true,
		Message: "SMS message processed successfully",
	}, nil
}

// HandleVoiceMessage handles voice messages from fr0g-ai-io
func (s *Fr0gIOInputServer) HandleVoiceMessage(ctx context.Context, req *VoiceMessageRequest) (*VoiceMessageResponse, error) {
	log.Printf("gRPC Input Server: Received voice message %s from %s", req.Message.Id, req.Message.From)

	// Convert gRPC message to internal type
	voiceMessage := &input.VoiceMessage{
		ID:            req.Message.Id,
		From:          req.Message.From,
		To:            req.Message.To,
		AudioData:     req.Message.AudioData,
		Transcription: req.Message.Transcription,
		Duration:      req.Message.Duration,
		Timestamp:     time.Unix(req.Message.Timestamp, 0),
		Format:        req.Message.Format,
	}

	// Process the voice message
	err := s.inputHandler.HandleVoiceMessage(ctx, voiceMessage)
	if err != nil {
		return &VoiceMessageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to process voice message: %v", err),
		}, nil
	}

	return &VoiceMessageResponse{
		Success: true,
		Message: "Voice message processed successfully",
	}, nil
}

// HandleIRCMessage handles IRC messages from fr0g-ai-io
func (s *Fr0gIOInputServer) HandleIRCMessage(ctx context.Context, req *IRCMessageRequest) (*IRCMessageResponse, error) {
	log.Printf("gRPC Input Server: Received IRC message %s from %s in %s", req.Message.Id, req.Message.Nick, req.Message.Channel)

	// Convert gRPC message to internal type
	ircMessage := &input.IRCMessage{
		ID:        req.Message.Id,
		Server:    req.Message.Server,
		Channel:   req.Message.Channel,
		Nick:      req.Message.Nick,
		Content:   req.Message.Content,
		Timestamp: time.Unix(req.Message.Timestamp, 0),
		IsPrivate: req.Message.IsPrivate,
	}

	// Process the IRC message
	err := s.inputHandler.HandleIRCMessage(ctx, ircMessage)
	if err != nil {
		return &IRCMessageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to process IRC message: %v", err),
		}, nil
	}

	return &IRCMessageResponse{
		Success: true,
		Message: "IRC message processed successfully",
	}, nil
}

// HandleDiscordMessage handles Discord messages from fr0g-ai-io
func (s *Fr0gIOInputServer) HandleDiscordMessage(ctx context.Context, req *DiscordMessageRequest) (*DiscordMessageResponse, error) {
	log.Printf("gRPC Input Server: Received Discord message %s from %s in guild %s", req.Message.Id, req.Message.Username, req.Message.GuildId)

	// Convert gRPC message to internal type
	discordMessage := &input.DiscordMessage{
		ID:          req.Message.Id,
		GuildID:     req.Message.GuildId,
		ChannelID:   req.Message.ChannelId,
		UserID:      req.Message.UserId,
		Username:    req.Message.Username,
		Content:     req.Message.Content,
		Timestamp:   time.Unix(req.Message.Timestamp, 0),
		MessageType: req.Message.MessageType,
	}

	// Process the Discord message
	err := s.inputHandler.HandleDiscordMessage(ctx, discordMessage)
	if err != nil {
		return &DiscordMessageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to process Discord message: %v", err),
		}, nil
	}

	return &DiscordMessageResponse{
		Success: true,
		Message: "Discord message processed successfully",
	}, nil
}
