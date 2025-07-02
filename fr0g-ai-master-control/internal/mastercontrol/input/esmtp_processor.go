package input

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/mail"
	"net/textproto"
	"strings"
	"time"
)


// ESMTPProcessor handles incoming ESMTP connections and processes emails
// Acting as the "intelligent front desk" for email threat vector analysis
type ESMTPProcessor struct {
	config            *ESMTPConfig
	aiCommunityClient AIPersonaCommunityClient
	listener          net.Listener
	tlsConfig         *tls.Config
	shutdown          chan struct{}
}

type ESMTPConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	TLSPort        int           `yaml:"tls_port"`
	Hostname       string        `yaml:"hostname"`
	MaxMessageSize int64         `yaml:"max_message_size"`
	Timeout        time.Duration `yaml:"timeout"`
	EnableTLS      bool          `yaml:"enable_tls"`
	CertFile       string        `yaml:"cert_file"`
	KeyFile        string        `yaml:"key_file"`
	
	// Community settings
	CommunityTopic    string        `yaml:"community_topic"`
	PersonaCount      int           `yaml:"persona_count"`
	ReviewTimeout     time.Duration `yaml:"review_timeout"`
	RequiredConsensus float64       `yaml:"required_consensus"`
}

// EmailThreatVector represents an intercepted email for analysis
type EmailThreatVector struct {
	ID          string            `json:"id"`
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	Attachments []AttachmentInfo  `json:"attachments"`
	Timestamp   time.Time         `json:"timestamp"`
	ThreatLevel string            `json:"threat_level"` // "unknown", "low", "medium", "high", "critical"
	Source      string            `json:"source"`       // "esmtp"
}

type AttachmentInfo struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Hash        string `json:"hash"`
}


// NewESMTPProcessor creates a new ESMTP processor instance
func NewESMTPProcessor(config *ESMTPConfig, aiClient AIPersonaCommunityClient) (*ESMTPProcessor, error) {
	processor := &ESMTPProcessor{
		config:            config,
		aiCommunityClient: aiClient,
		shutdown:          make(chan struct{}),
	}

	// Setup TLS if enabled
	if config.EnableTLS {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificates: %w", err)
		}
		processor.tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   config.Hostname,
		}
	}

	return processor, nil
}

// Start begins listening for ESMTP connections
func (p *ESMTPProcessor) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	p.listener = listener

	log.Printf("🐸 fr0g.ai ESMTP Threat Vector Interceptor listening on %s", addr)
	log.Printf("📧 Email intelligence gathering active - eliminating human-computer interaction vulnerabilities")

	go p.acceptConnections(ctx)
	return nil
}

// acceptConnections handles incoming ESMTP connections
func (p *ESMTPProcessor) acceptConnections(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.shutdown:
			return
		default:
			conn, err := p.listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}

			// Handle each connection in a goroutine
			go p.handleConnection(ctx, conn)
		}
	}
}

// handleConnection processes a single ESMTP connection
func (p *ESMTPProcessor) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(p.config.Timeout))

	session := &ESMTPSession{
		conn:      conn,
		processor: p,
		reader:    textproto.NewReader(bufio.NewReader(conn)),
		writer:    textproto.NewWriter(bufio.NewWriter(conn)),
	}

	if err := session.handle(ctx); err != nil {
		log.Printf("ESMTP session error: %v", err)
	}
}

// ESMTPSession represents a single ESMTP session
type ESMTPSession struct {
	conn      net.Conn
	processor *ESMTPProcessor
	reader    *textproto.Reader
	writer    *textproto.Writer
	from      string
	to        []string
	data      []byte
}

// handle processes the ESMTP protocol for this session
func (s *ESMTPSession) handle(ctx context.Context) error {
	// Send greeting
	if err := s.writer.PrintfLine("220 %s fr0g.ai ESMTP Threat Vector Interceptor Ready", s.processor.config.Hostname); err != nil {
		return err
	}

	for {
		line, err := s.reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		cmd := strings.ToUpper(strings.Fields(line)[0])
		
		switch cmd {
		case "HELO", "EHLO":
			if err := s.handleHelo(line); err != nil {
				return err
			}
		case "MAIL":
			if err := s.handleMail(line); err != nil {
				return err
			}
		case "RCPT":
			if err := s.handleRcpt(line); err != nil {
				return err
			}
		case "DATA":
			if err := s.handleData(ctx); err != nil {
				return err
			}
		case "QUIT":
			s.writer.PrintfLine("221 %s closing connection", s.processor.config.Hostname)
			return nil
		case "RSET":
			s.reset()
			s.writer.PrintfLine("250 OK")
		case "NOOP":
			s.writer.PrintfLine("250 OK")
		default:
			s.writer.PrintfLine("502 Command not implemented")
		}
	}
}

func (s *ESMTPSession) handleHelo(line string) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return s.writer.PrintfLine("501 Syntax error")
	}

	// Enhanced SMTP response with fr0g.ai capabilities
	if strings.ToUpper(parts[0]) == "EHLO" {
		s.writer.PrintfLine("250-%s Hello %s", s.processor.config.Hostname, parts[1])
		s.writer.PrintfLine("250-SIZE %d", s.processor.config.MaxMessageSize)
		s.writer.PrintfLine("250-8BITMIME")
		if s.processor.config.EnableTLS {
			s.writer.PrintfLine("250-STARTTLS")
		}
		s.writer.PrintfLine("250 HELP")
	} else {
		s.writer.PrintfLine("250 %s Hello %s", s.processor.config.Hostname, parts[1])
	}
	return nil
}

func (s *ESMTPSession) handleMail(line string) error {
	// Parse MAIL FROM command
	if !strings.HasPrefix(strings.ToUpper(line), "MAIL FROM:") {
		return s.writer.PrintfLine("501 Syntax error")
	}

	from := strings.TrimSpace(line[10:])
	from = strings.Trim(from, "<>")
	s.from = from

	return s.writer.PrintfLine("250 OK")
}

func (s *ESMTPSession) handleRcpt(line string) error {
	// Parse RCPT TO command
	if !strings.HasPrefix(strings.ToUpper(line), "RCPT TO:") {
		return s.writer.PrintfLine("501 Syntax error")
	}

	to := strings.TrimSpace(line[8:])
	to = strings.Trim(to, "<>")
	s.to = append(s.to, to)

	return s.writer.PrintfLine("250 OK")
}

func (s *ESMTPSession) handleData(ctx context.Context) error {
	if s.from == "" || len(s.to) == 0 {
		return s.writer.PrintfLine("503 Bad sequence of commands")
	}

	if err := s.writer.PrintfLine("354 Start mail input; end with <CRLF>.<CRLF>"); err != nil {
		return err
	}

	// Read email data until "."
	var data []byte
	for {
		line, err := s.reader.ReadLine()
		if err != nil {
			return err
		}
		if line == "." {
			break
		}
		// Handle dot-stuffing
		if strings.HasPrefix(line, ".") {
			line = line[1:]
		}
		data = append(data, []byte(line+"\r\n")...)
	}

	s.data = data

	// Process the email through fr0g.ai threat analysis
	if err := s.processEmailThreatVector(ctx); err != nil {
		log.Printf("Error processing email threat vector: %v", err)
		return s.writer.PrintfLine("451 Temporary failure - threat analysis error")
	}

	return s.writer.PrintfLine("250 OK: Message accepted for threat analysis")
}

func (s *ESMTPSession) processEmailThreatVector(ctx context.Context) error {
	// Parse the email
	msg, err := mail.ReadMessage(strings.NewReader(string(s.data)))
	if err != nil {
		return fmt.Errorf("failed to parse email: %w", err)
	}

	// Extract headers
	headers := make(map[string]string)
	for key, values := range msg.Header {
		headers[key] = strings.Join(values, ", ")
	}

	// Read body
	bodyBytes, err := io.ReadAll(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to read email body: %w", err)
	}

	// Create threat vector
	vector := &EmailThreatVector{
		ID:          generateThreatVectorID(),
		From:        s.from,
		To:          s.to,
		Subject:     msg.Header.Get("Subject"),
		Body:        string(bodyBytes),
		Headers:     headers,
		Attachments: []AttachmentInfo{}, // TODO: Parse attachments
		Timestamp:   time.Now(),
		ThreatLevel: "unknown", // Will be analyzed by community
		Source:      "esmtp",
	}

	// Submit to AI community for review
	_, err = s.processor.submitForCommunityReview(ctx, vector)
	return err
}

func (s *ESMTPSession) reset() {
	s.from = ""
	s.to = nil
	s.data = nil
}

// WebhookProcessor interface implementation
func (p *ESMTPProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	// Convert webhook request to email threat vector
	emailData, ok := request.Body.(map[string]interface{})
	if !ok {
		return &WebhookResponse{
			Success:   false,
			Message:   "Invalid email data format",
			RequestID: request.ID,
			Timestamp: time.Now(),
		}, nil
	}

	// Extract email fields
	from, _ := emailData["from"].(string)
	to, _ := emailData["to"].([]string)
	subject, _ := emailData["subject"].(string)
	body, _ := emailData["body"].(string)

	// Create threat vector for community review
	vector := &EmailThreatVector{
		ID:          generateThreatVectorID(),
		From:        from,
		To:          to,
		Subject:     subject,
		Body:        body,
		Headers:     make(map[string]string),
		Attachments: []AttachmentInfo{},
		Timestamp:   time.Now(),
		ThreatLevel: "unknown",
		Source:      "esmtp",
	}

	// Submit to AI community for review
	review, err := p.submitForCommunityReview(ctx, vector)
	if err != nil {
		return &WebhookResponse{
			Success:   false,
			Message:   fmt.Sprintf("Community review failed: %v", err),
			RequestID: request.ID,
			Timestamp: time.Now(),
		}, nil
	}

	return &WebhookResponse{
		Success:   true,
		Message:   "Email threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"review_id":    review.ReviewID,
			"threat_level": vector.ThreatLevel,
			"consensus":    review.Consensus.OverallScore,
		},
		Timestamp: time.Now(),
	}, nil
}

func (p *ESMTPProcessor) GetTag() string {
	return "esmtp"
}

func (p *ESMTPProcessor) GetDescription() string {
	return fmt.Sprintf("ESMTP Threat Vector Interceptor on %s:%d - Email intelligence gathering for AI community review on topic: %s", 
		p.config.Host, p.config.Port, p.config.CommunityTopic)
}

// submitForCommunityReview submits email threat vector to AI community
func (p *ESMTPProcessor) submitForCommunityReview(ctx context.Context, vector *EmailThreatVector) (*CommunityReview, error) {
	// Create or get community for email threat analysis
	community, err := p.aiCommunityClient.CreateCommunity(ctx, p.config.CommunityTopic, p.config.PersonaCount)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	// Format content for review
	content := fmt.Sprintf(`Email Threat Vector Analysis Request:

From: %s
To: %s
Subject: %s

Body:
%s

Headers: %v
Attachments: %d

Please analyze this email for potential threats, social engineering attempts, phishing indicators, and overall risk assessment.`,
		vector.From,
		strings.Join(vector.To, ", "),
		vector.Subject,
		vector.Body,
		vector.Headers,
		len(vector.Attachments))

	// Submit for community review
	return p.aiCommunityClient.SubmitForReview(ctx, community.ID, content)
}

// generateThreatVectorID creates a unique ID for threat vectors
func generateThreatVectorID() string {
	return fmt.Sprintf("etv_%d", time.Now().UnixNano())
}

// Stop gracefully shuts down the ESMTP processor
func (p *ESMTPProcessor) Stop() error {
	close(p.shutdown)
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}
