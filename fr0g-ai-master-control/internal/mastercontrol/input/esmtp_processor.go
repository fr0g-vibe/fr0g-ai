package input

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ESMTPProcessor handles incoming ESMTP connections and processes emails
// Acting as the "intelligent front desk" for email threat vector analysis
type ESMTPProcessor struct {
	config     *ESMTPConfig
	mcpClient  MasterControlClient
	listener   net.Listener
	tlsConfig  *tls.Config
	shutdown   chan struct{}
}

type ESMTPConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	TLSPort        int           `yaml:"tls_port"`
	Hostname       string        `yaml:"hostname"`
	MaxMessageSize int64         `yaml:"max_message_size"`
	Timeout        time.Duration `yaml:"timeout"`
	MCPAddress     string        `yaml:"mcp_address"`
	EnableTLS      bool          `yaml:"enable_tls"`
	CertFile       string        `yaml:"cert_file"`
	KeyFile        string        `yaml:"key_file"`
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

// MasterControlClient interface for communicating with MCP
type MasterControlClient interface {
	ProcessThreatVector(ctx context.Context, vector *EmailThreatVector) error
	HealthCheck(ctx context.Context) error
}

// grpcMasterControlClient implements MasterControlClient using gRPC
type grpcMasterControlClient struct {
	conn   *grpc.ClientConn
	client interface{} // Will be the actual gRPC client when proto is defined
}

// NewESMTPProcessor creates a new ESMTP processor instance
func NewESMTPProcessor(config *ESMTPConfig) (*ESMTPProcessor, error) {
	processor := &ESMTPProcessor{
		config:   config,
		shutdown: make(chan struct{}),
	}

	// Initialize MCP client
	mcpClient, err := newMCPClient(config.MCPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	processor.mcpClient = mcpClient

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
		ThreatLevel: "unknown", // Will be analyzed by MCP
		Source:      "esmtp",
	}

	// Send to Master Control for processing
	return s.processor.mcpClient.ProcessThreatVector(ctx, vector)
}

func (s *ESMTPSession) reset() {
	s.from = ""
	s.to = nil
	s.data = nil
}

// newMCPClient creates a new Master Control Protocol client
func newMCPClient(address string) (MasterControlClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MCP: %w", err)
	}

	return &grpcMasterControlClient{
		conn: conn,
		// client: will be initialized with actual gRPC client
	}, nil
}

func (c *grpcMasterControlClient) ProcessThreatVector(ctx context.Context, vector *EmailThreatVector) error {
	// TODO: Implement actual gRPC call to MCP
	// For now, log the threat vector
	vectorJSON, _ := json.MarshalIndent(vector, "", "  ")
	log.Printf("🚨 THREAT VECTOR INTERCEPTED:\n%s", string(vectorJSON))
	
	// Simulate processing
	log.Printf("📡 Forwarding to Master Control Protocol for cognitive analysis...")
	return nil
}

func (c *grpcMasterControlClient) HealthCheck(ctx context.Context) error {
	// TODO: Implement actual health check
	return nil
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
