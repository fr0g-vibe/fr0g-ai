package email

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"
)

// SMTPServer implements an SMTP server for email processing
type SMTPServer struct {
	config        *EmailConfig
	emailHandler  func([]byte, string, []string) error
	listener      net.Listener
	tlsListener   net.Listener
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	connections   map[net.Conn]bool
	connMutex     sync.RWMutex
	maxConns      int
	currentConns  int
}

// SMTPSession represents an active SMTP session
type SMTPSession struct {
	conn         net.Conn
	reader       *textproto.Reader
	writer       *textproto.Writer
	server       *SMTPServer
	helo         string
	from         string
	recipients   []string
	data         []byte
	authenticated bool
	tlsEnabled   bool
}

// NewSMTPServer creates a new SMTP server instance
func NewSMTPServer(config *EmailConfig, emailHandler func([]byte, string, []string) error) *SMTPServer {
	return &SMTPServer{
		config:       config,
		emailHandler: emailHandler,
		connections:  make(map[net.Conn]bool),
		maxConns:     config.MaxConnections,
	}
}

// Start starts the SMTP server on configured ports
func (s *SMTPServer) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	// Start SMTP server (port 25)
	if err := s.startSMTPListener(); err != nil {
		return fmt.Errorf("failed to start SMTP listener: %w", err)
	}

	// Start SMTPS server (port 465) if TLS is enabled
	if s.config.EnableTLS {
		if err := s.startSMTPSListener(); err != nil {
			s.listener.Close()
			return fmt.Errorf("failed to start SMTPS listener: %w", err)
		}
	}

	log.Printf("SMTP server started on port %d", s.config.SMTPPort)
	if s.config.EnableTLS {
		log.Printf("SMTPS server started on port %d", s.config.SMTPSPort)
	}

	return nil
}

// Stop stops the SMTP server
func (s *SMTPServer) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}

	// Close listeners
	if s.listener != nil {
		s.listener.Close()
	}
	if s.tlsListener != nil {
		s.tlsListener.Close()
	}

	// Close all active connections
	s.connMutex.Lock()
	for conn := range s.connections {
		conn.Close()
	}
	s.connMutex.Unlock()

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout waiting for connections to close")
	}
}

// startSMTPListener starts the SMTP listener on the configured port
func (s *SMTPServer) startSMTPListener() error {
	addr := fmt.Sprintf(":%d", s.config.SMTPPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = listener
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		s.acceptConnections(listener, false)
	}()

	return nil
}

// startSMTPSListener starts the SMTPS listener with TLS
func (s *SMTPServer) startSMTPSListener() error {
	cert, err := tls.LoadX509KeyPair(s.config.TLSCertFile, s.config.TLSKeyFile)
	if err != nil {
		// Use self-signed cert for development if files don't exist
		cert, err = s.generateSelfSignedCert()
		if err != nil {
			return fmt.Errorf("failed to load or generate TLS certificate: %w", err)
		}
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   s.config.Hostname,
	}

	addr := fmt.Sprintf(":%d", s.config.SMTPSPort)
	listener, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	s.tlsListener = listener
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		s.acceptConnections(listener, true)
	}()

	return nil
}

// acceptConnections accepts and handles incoming connections
func (s *SMTPServer) acceptConnections(listener net.Listener, isTLS bool) {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			if s.ctx.Err() != nil {
				return // Server is shutting down
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Check connection limits
		s.connMutex.Lock()
		if s.currentConns >= s.maxConns {
			s.connMutex.Unlock()
			conn.Close()
			log.Printf("Connection rejected: maximum connections reached")
			continue
		}
		s.currentConns++
		s.connections[conn] = true
		s.connMutex.Unlock()

		// Handle connection in goroutine
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer func() {
				s.connMutex.Lock()
				delete(s.connections, conn)
				s.currentConns--
				s.connMutex.Unlock()
				conn.Close()
			}()

			s.handleConnection(conn, isTLS)
		}()
	}
}

// handleConnection handles a single SMTP connection
func (s *SMTPServer) handleConnection(conn net.Conn, isTLS bool) {
	// Set connection timeout
	conn.SetDeadline(time.Now().Add(s.config.ConnectionTimeout))

	session := &SMTPSession{
		conn:       conn,
		reader:     textproto.NewReader(bufio.NewReader(conn)),
		writer:     textproto.NewWriter(bufio.NewWriter(conn)),
		server:     s,
		tlsEnabled: isTLS,
	}

	// Send greeting
	session.writer.PrintfLine("220 %s ESMTP Service Ready", s.config.Hostname)

	// Handle SMTP commands
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		line, err := session.reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v", err)
			}
			return
		}

		if err := session.handleCommand(line); err != nil {
			log.Printf("Error handling command: %v", err)
			return
		}
	}
}

// handleCommand processes SMTP commands
func (s *SMTPSession) handleCommand(line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		s.writer.PrintfLine("500 Command not recognized")
		return nil
	}

	command := strings.ToUpper(parts[0])
	args := parts[1:]

	switch command {
	case "HELO":
		return s.handleHELO(args)
	case "EHLO":
		return s.handleEHLO(args)
	case "MAIL":
		return s.handleMAIL(line)
	case "RCPT":
		return s.handleRCPT(line)
	case "DATA":
		return s.handleDATA()
	case "RSET":
		return s.handleRSET()
	case "QUIT":
		return s.handleQUIT()
	case "NOOP":
		return s.handleNOOP()
	case "STARTTLS":
		return s.handleSTARTTLS()
	default:
		s.writer.PrintfLine("502 Command not implemented")
	}

	return nil
}

// handleHELO handles the HELO command
func (s *SMTPSession) handleHELO(args []string) error {
	if len(args) != 1 {
		s.writer.PrintfLine("501 Syntax error in parameters")
		return nil
	}

	s.helo = args[0]
	s.writer.PrintfLine("250 %s Hello %s", s.server.config.Hostname, s.helo)
	return nil
}

// handleEHLO handles the EHLO command (Extended SMTP)
func (s *SMTPSession) handleEHLO(args []string) error {
	if len(args) != 1 {
		s.writer.PrintfLine("501 Syntax error in parameters")
		return nil
	}

	s.helo = args[0]
	s.writer.PrintfLine("250-%s Hello %s", s.server.config.Hostname, s.helo)
	s.writer.PrintfLine("250-SIZE %d", s.server.config.MaxMessageSize)
	if s.server.config.EnableTLS && !s.tlsEnabled {
		s.writer.PrintfLine("250-STARTTLS")
	}
	s.writer.PrintfLine("250 HELP")
	return nil
}

// handleMAIL handles the MAIL FROM command
func (s *SMTPSession) handleMAIL(line string) error {
	if s.helo == "" {
		s.writer.PrintfLine("503 Bad sequence of commands")
		return nil
	}

	// Parse MAIL FROM:<address>
	if !strings.HasPrefix(strings.ToUpper(line), "MAIL FROM:") {
		s.writer.PrintfLine("501 Syntax error in parameters")
		return nil
	}

	from := strings.TrimSpace(line[10:])
	if strings.HasPrefix(from, "<") && strings.HasSuffix(from, ">") {
		from = from[1 : len(from)-1]
	}

	s.from = from
	s.recipients = nil
	s.writer.PrintfLine("250 OK")
	return nil
}

// handleRCPT handles the RCPT TO command
func (s *SMTPSession) handleRCPT(line string) error {
	if s.from == "" {
		s.writer.PrintfLine("503 Bad sequence of commands")
		return nil
	}

	// Parse RCPT TO:<address>
	if !strings.HasPrefix(strings.ToUpper(line), "RCPT TO:") {
		s.writer.PrintfLine("501 Syntax error in parameters")
		return nil
	}

	to := strings.TrimSpace(line[8:])
	if strings.HasPrefix(to, "<") && strings.HasSuffix(to, ">") {
		to = to[1 : len(to)-1]
	}

	s.recipients = append(s.recipients, to)
	s.writer.PrintfLine("250 OK")
	return nil
}

// handleDATA handles the DATA command
func (s *SMTPSession) handleDATA() error {
	if len(s.recipients) == 0 {
		s.writer.PrintfLine("503 Bad sequence of commands")
		return nil
	}

	s.writer.PrintfLine("354 Start mail input; end with <CRLF>.<CRLF>")

	// Read email data until "." on a line by itself
	var data []byte
	for {
		line, err := s.reader.ReadLine()
		if err != nil {
			return err
		}

		if line == "." {
			break
		}

		// Handle dot-stuffing (remove leading dot if present)
		if strings.HasPrefix(line, ".") {
			line = line[1:]
		}

		data = append(data, []byte(line+"\r\n")...)

		// Check message size limit
		if int64(len(data)) > s.server.config.MaxMessageSize {
			s.writer.PrintfLine("552 Message size exceeds maximum allowed")
			return nil
		}
	}

	s.data = data

	// Process the email
	if err := s.server.emailHandler(s.data, s.from, s.recipients); err != nil {
		log.Printf("Error processing email: %v", err)
		s.writer.PrintfLine("451 Requested action aborted: local error in processing")
		return nil
	}

	s.writer.PrintfLine("250 OK: Message accepted for delivery")
	
	// Reset session state
	s.from = ""
	s.recipients = nil
	s.data = nil

	return nil
}

// handleRSET handles the RSET command
func (s *SMTPSession) handleRSET() error {
	s.from = ""
	s.recipients = nil
	s.data = nil
	s.writer.PrintfLine("250 OK")
	return nil
}

// handleQUIT handles the QUIT command
func (s *SMTPSession) handleQUIT() error {
	s.writer.PrintfLine("221 %s Service closing transmission channel", s.server.config.Hostname)
	return fmt.Errorf("client quit")
}

// handleNOOP handles the NOOP command
func (s *SMTPSession) handleNOOP() error {
	s.writer.PrintfLine("250 OK")
	return nil
}

// handleSTARTTLS handles the STARTTLS command
func (s *SMTPSession) handleSTARTTLS() error {
	if s.tlsEnabled {
		s.writer.PrintfLine("503 Already using TLS")
		return nil
	}

	if !s.server.config.EnableTLS {
		s.writer.PrintfLine("502 Command not implemented")
		return nil
	}

	s.writer.PrintfLine("220 Ready to start TLS")

	// Upgrade connection to TLS
	cert, err := s.server.generateSelfSignedCert()
	if err != nil {
		s.writer.PrintfLine("454 TLS not available due to temporary reason")
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   s.server.config.Hostname,
	}

	tlsConn := tls.Server(s.conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return err
	}

	// Update session with TLS connection
	s.conn = tlsConn
	s.reader = textproto.NewReader(bufio.NewReader(tlsConn))
	s.writer = textproto.NewWriter(bufio.NewWriter(tlsConn))
	s.tlsEnabled = true

	// Reset session state after STARTTLS
	s.helo = ""
	s.from = ""
	s.recipients = nil

	return nil
}

// generateSelfSignedCert generates a self-signed certificate for development
func (s *SMTPServer) generateSelfSignedCert() (tls.Certificate, error) {
	// This is a simplified implementation for development
	// In production, you should use proper certificates
	return tls.Certificate{}, fmt.Errorf("self-signed certificate generation not implemented")
}
