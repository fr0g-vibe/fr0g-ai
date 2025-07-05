package smtp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/types"
)

// Processor handles SMTP email output
type Processor struct {
	config *sharedconfig.SMTPConfig
}

// NewProcessor creates a new SMTP output processor
func NewProcessor(cfg *sharedconfig.Config) *Processor {
	return &Processor{
		config: cfg.SMTP,
	}
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "smtp"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.config != nil && p.config.Enabled
}

// Process sends an email using SMTP
func (p *Processor) Process(command *types.OutputCommand) (*types.OutputResult, error) {
	startTime := time.Now()

	if !p.IsEnabled() {
		return &types.OutputResult{
			CommandID:      command.ID,
			Success:        false,
			ErrorMessage:   "SMTP processor is disabled",
			Metadata:       map[string]string{"error": "processor_disabled"},
			CompletedAt:    time.Now(),
			AttemptCount:   1,
			ProcessingTime: time.Since(startTime),
		}, nil
	}

	// Extract email details from command
	to := command.Target
	subject := p.extractSubject(command)
	body := command.Content

	// Send email
	err := p.sendEmail(to, subject, body)
	if err != nil {
		return &types.OutputResult{
			CommandID:      command.ID,
			Success:        false,
			ErrorMessage:   fmt.Sprintf("Failed to send email: %v", err),
			Metadata:       map[string]string{"error": "send_failed", "target": to},
			CompletedAt:    time.Now(),
			AttemptCount:   1,
			ProcessingTime: time.Since(startTime),
		}, nil
	}

	return &types.OutputResult{
		CommandID:      command.ID,
		Success:        true,
		Metadata:       map[string]string{"target": to, "subject": subject},
		CompletedAt:    time.Now(),
		AttemptCount:   1,
		ProcessingTime: time.Since(startTime),
	}, nil
}

// sendEmail sends an email using SMTP
func (p *Processor) sendEmail(to, subject, body string) error {
	// Prepare message
	msg := p.formatMessage(to, subject, body)

	// Setup authentication
	auth := smtp.PlainAuth("", p.config.Username, p.config.Password, p.config.Host)

	// Determine server address
	addr := fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)

	// Send email based on TLS configuration
	if p.config.TLS {
		return p.sendWithTLS(addr, auth, to, msg)
	} else {
		return smtp.SendMail(addr, auth, p.config.From, []string{to}, []byte(msg))
	}
}

// sendWithTLS sends email with TLS encryption
func (p *Processor) sendWithTLS(addr string, auth smtp.Auth, to, msg string) error {
	// Create TLS connection
	tlsConfig := &tls.Config{
		ServerName:         p.config.Host,
		InsecureSkipVerify: p.config.InsecureSkipVerify,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to create TLS connection: %v", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, p.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Authenticate
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %v", err)
		}
	}

	// Set sender
	if err := client.Mail(p.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// Set recipient
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return writer.Close()
}

// formatMessage formats the email message with headers
func (p *Processor) formatMessage(to, subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = p.config.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=utf-8"
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	// Build message
	var msg strings.Builder
	for key, value := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return msg.String()
}

// extractSubject extracts subject from command metadata or uses default
func (p *Processor) extractSubject(command *types.OutputCommand) string {
	if subject, exists := command.Metadata["subject"]; exists {
		if subjectStr, ok := subject.(string); ok {
			return subjectStr
		}
	}

	// Default subject based on command type
	switch command.Type {
	case "alert":
		return "Security Alert from fr0g.ai"
	case "report":
		return "Threat Analysis Report"
	case "notification":
		return "Notification from fr0g.ai"
	default:
		return "Message from fr0g.ai"
	}
}

// GetStatus returns the processor status
func (p *Processor) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"type":    p.GetType(),
		"enabled": p.IsEnabled(),
		"host":    p.config.Host,
		"port":    p.config.Port,
		"tls":     p.config.TLS,
	}
}
