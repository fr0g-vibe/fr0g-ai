package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// IRCProcessor processes IRC chat messages for threat analysis
type IRCProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *IRCConfig
}

// IRCConfig holds IRC processor configuration
type IRCConfig struct {
	Server            string        `yaml:"server"`
	Port              int           `yaml:"port"`
	UseSSL            bool          `yaml:"use_ssl"`
	Nickname          string        `yaml:"nickname"`
	Username          string        `yaml:"username"`
	RealName          string        `yaml:"real_name"`
	Password          string        `yaml:"password,omitempty"`
	Channels          []string      `yaml:"channels"`
	CommunityTopic    string        `yaml:"community_topic"`
	PersonaCount      int           `yaml:"persona_count"`
	ReviewTimeout     time.Duration `yaml:"review_timeout"`
	RequiredConsensus float64       `yaml:"required_consensus"`
	MonitorPrivateMsg bool          `yaml:"monitor_private_msg"`
	IgnoredNicks      []string      `yaml:"ignored_nicks"`
	TrustedNicks      []string      `yaml:"trusted_nicks"`
	AutoJoinChannels  bool          `yaml:"auto_join_channels"`
}

// IRCMessage represents an IRC message
type IRCMessage struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`      // "PRIVMSG", "NOTICE", "JOIN", "PART", "QUIT", "KICK", "MODE"
	From      string            `json:"from"`      // Nickname
	To        string            `json:"to"`        // Channel or nickname
	Message   string            `json:"message"`
	Channel   string            `json:"channel"`
	Server    string            `json:"server"`
	Timestamp time.Time         `json:"timestamp"`
	IsPrivate bool              `json:"is_private"`
	UserInfo  *IRCUserInfo      `json:"user_info,omitempty"`
	Metadata  map[string]string `json:"metadata"`
}

// IRCUserInfo represents IRC user information
type IRCUserInfo struct {
	Nickname string   `json:"nickname"`
	Username string   `json:"username"`
	Hostname string   `json:"hostname"`
	RealName string   `json:"real_name"`
	Channels []string `json:"channels"`
	IsOp     bool     `json:"is_op"`
	IsVoice  bool     `json:"is_voice"`
	IdleTime int      `json:"idle_time"`
}

// NewIRCProcessor creates a new IRC processor
func NewIRCProcessor(config *IRCConfig, aiClient AIPersonaCommunityClient) (*IRCProcessor, error) {
	return &IRCProcessor{
		aiClient: aiClient,
		config:   config,
	}, nil
}

// GetTag returns the processor tag
func (i *IRCProcessor) GetTag() string {
	return "irc"
}

// GetDescription returns the processor description
func (i *IRCProcessor) GetDescription() string {
	return fmt.Sprintf("IRC Chat Threat Vector Interceptor on %s:%d - Chat intelligence gathering for AI community review on topic: %s", 
		i.config.Server, i.config.Port, i.config.CommunityTopic)
}

// ProcessWebhook processes an IRC message webhook
func (i *IRCProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("IRC Processor: Processing IRC message threat vector webhook")
	
	// Parse IRC message from request body
	ircMsg, err := i.parseIRCMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IRC message: %w", err)
	}
	
	// Check if user is ignored
	if i.isUserIgnored(ircMsg.From) {
		log.Printf("IRC Processor: Ignored user %s, skipping analysis", ircMsg.From)
		return &WebhookResponse{
			Success:   true,
			Message:   "IRC message from ignored user skipped",
			RequestID: request.ID,
			Data: map[string]interface{}{
				"action":   "ignored",
				"reason":   "ignored_user",
				"nickname": ircMsg.From,
			},
			Timestamp: time.Now(),
		}, nil
	}
	
	// Check if user is trusted (lower threat threshold)
	isTrusted := i.isUserTrusted(ircMsg.From)
	
	// Analyze IRC message for threats using AI community
	threatLevel, consensus, err := i.analyzeIRCThreats(ctx, ircMsg, isTrusted)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze IRC threats: %w", err)
	}
	
	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "IRC message threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"threat_level":  threatLevel,
			"consensus":     consensus,
			"review_id":     fmt.Sprintf("irc_review_%d", time.Now().UnixNano()),
			"nickname":      ircMsg.From,
			"channel":       ircMsg.Channel,
			"message_type":  ircMsg.Type,
			"is_private":    ircMsg.IsPrivate,
			"is_trusted":    isTrusted,
		},
		Timestamp: time.Now(),
	}
	
	log.Printf("IRC Processor: Message analyzed - From: %s, Channel: %s, Threat Level: %s, Consensus: %.2f", 
		ircMsg.From, ircMsg.Channel, threatLevel, consensus)
	
	return response, nil
}

// analyzeIRCThreats analyzes IRC message for threats using AI community
func (i *IRCProcessor) analyzeIRCThreats(ctx context.Context, ircMsg *IRCMessage, isTrusted bool) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
IRC Chat Message Threat Analysis Request:
From: %s (nickname)
To: %s
Channel: %s
Message Type: %s
Is Private Message: %v
Is Trusted User: %v
Server: %s
Timestamp: %s

Message Content:
%s

`, ircMsg.From, ircMsg.To, ircMsg.Channel, ircMsg.Type, ircMsg.IsPrivate, isTrusted, 
	ircMsg.Server, ircMsg.Timestamp.Format(time.RFC3339), ircMsg.Message)
	
	// Add user info if available
	if ircMsg.UserInfo != nil {
		analysisContent += fmt.Sprintf(`
User Information:
- Username: %s
- Hostname: %s
- Real Name: %s
- Is Operator: %v
- Is Voice: %v
- Idle Time: %d seconds
- Channels: %s

`, ircMsg.UserInfo.Username, ircMsg.UserInfo.Hostname, ircMsg.UserInfo.RealName,
			ircMsg.UserInfo.IsOp, ircMsg.UserInfo.IsVoice, ircMsg.UserInfo.IdleTime,
			strings.Join(ircMsg.UserInfo.Channels, ", "))
	}
	
	analysisContent += `
Please analyze this IRC message for potential threats including:
- Malicious links and URLs
- Social engineering attempts
- Phishing attempts
- Malware distribution
- Spam and unwanted content
- Harassment or abusive behavior
- Bot activity and automated spam
- Channel flooding or disruption
- Private message threats
- Suspicious file sharing
- Cryptocurrency scams
- Identity theft attempts
- Doxxing or privacy violations
`
	
	// Create AI community for threat analysis
	community, err := i.aiClient.CreateCommunity(ctx, i.config.CommunityTopic, i.config.PersonaCount)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to create AI community: %w", err)
	}
	
	// Submit for AI community review
	review, err := i.aiClient.SubmitForReview(ctx, community.ID, analysisContent)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to submit for review: %w", err)
	}
	
	// Determine threat level based on consensus
	threatLevel := "unknown"
	consensus := 0.0
	
	if review.Consensus != nil {
		consensus = review.Consensus.OverallScore
		
		// Adjust thresholds for trusted users
		thresholds := map[string]float64{
			"critical": 0.9,
			"high":     0.8,
			"medium":   0.6,
			"low":      0.4,
		}
		
		if isTrusted {
			// Higher thresholds for trusted users
			thresholds["critical"] = 0.95
			thresholds["high"] = 0.85
			thresholds["medium"] = 0.7
			thresholds["low"] = 0.5
		}
		
		if consensus >= thresholds["critical"] {
			threatLevel = "critical"
		} else if consensus >= thresholds["high"] {
			threatLevel = "high"
		} else if consensus >= thresholds["medium"] {
			threatLevel = "medium"
		} else if consensus >= thresholds["low"] {
			threatLevel = "low"
		} else {
			threatLevel = "minimal"
		}
	}
	
	return threatLevel, consensus, nil
}

// parseIRCMessage parses an IRC message from the request body
func (i *IRCProcessor) parseIRCMessage(body interface{}) (*IRCMessage, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}
	
	ircMsg := &IRCMessage{
		ID:        getStringFromMap(bodyMap, "id"),
		Type:      getStringFromMap(bodyMap, "type"),
		From:      getStringFromMap(bodyMap, "from"),
		To:        getStringFromMap(bodyMap, "to"),
		Message:   getStringFromMap(bodyMap, "message"),
		Channel:   getStringFromMap(bodyMap, "channel"),
		Server:    getStringFromMap(bodyMap, "server"),
		IsPrivate: getBoolFromMap(bodyMap, "is_private"),
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
	
	// Parse user info if available
	if userInfoData, ok := bodyMap["user_info"].(map[string]interface{}); ok {
		ircMsg.UserInfo = &IRCUserInfo{
			Nickname: getStringFromMap(userInfoData, "nickname"),
			Username: getStringFromMap(userInfoData, "username"),
			Hostname: getStringFromMap(userInfoData, "hostname"),
			RealName: getStringFromMap(userInfoData, "real_name"),
			IsOp:     getBoolFromMap(userInfoData, "is_op"),
			IsVoice:  getBoolFromMap(userInfoData, "is_voice"),
		}
		
		if idleTime, ok := userInfoData["idle_time"].(float64); ok {
			ircMsg.UserInfo.IdleTime = int(idleTime)
		}
		
		// Parse channels
		if channelsData, ok := userInfoData["channels"].([]interface{}); ok {
			for _, channel := range channelsData {
				if channelStr, ok := channel.(string); ok {
					ircMsg.UserInfo.Channels = append(ircMsg.UserInfo.Channels, channelStr)
				}
			}
		}
	}
	
	// Parse metadata
	if metadataData, ok := bodyMap["metadata"].(map[string]interface{}); ok {
		for key, value := range metadataData {
			if valueStr, ok := value.(string); ok {
				ircMsg.Metadata[key] = valueStr
			}
		}
	}
	
	return ircMsg, nil
}

// isUserIgnored checks if a user is in the ignored list
func (i *IRCProcessor) isUserIgnored(nickname string) bool {
	for _, ignored := range i.config.IgnoredNicks {
		if strings.EqualFold(nickname, ignored) {
			return true
		}
	}
	return false
}

// isUserTrusted checks if a user is in the trusted list
func (i *IRCProcessor) isUserTrusted(nickname string) bool {
	for _, trusted := range i.config.TrustedNicks {
		if strings.EqualFold(nickname, trusted) {
			return true
		}
	}
	return false
}
