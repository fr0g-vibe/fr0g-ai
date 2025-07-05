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

// SMTPServer implements a full SMTP server
type SMTPServer struct {
	config      *EmailConfig
	listener    net.Listener
	tlsListener net.Listener
	handler     EmailHandler
	mu          sync.RWMutex
	running     bool
	connections map[string]*SMTPConnection
	wg          sync.WaitGroup
}

// EmailHandler is a function type for handling received emails
type EmailHandler func(rawEmail []byte, from string, to []string) error

// SMTPConnection represents an active SMTP connection
type SMTPConnection struct {
	conn       net.Conn
	reader     *textproto.Reader
	writer     *textproto.Writer
	remoteAddr string
	state      SMTPState
	helo       string
	from       string
	to         []string
	data       []byte
	startTime  time.Time
}

// SMTPState represents the current state of an SMTP session
type SMTPState int

const (
	StateConnected SMTPState = iota
	StateHelo
	StateMailFrom
	StateRcptTo
	StateData
	StateQuit
)

// SMTP response codes
const (
	StatusReady             = 220
	StatusClosing           = 221
	StatusOK                = 250
	StatusStartData         = 354
	StatusSyntaxError       = 500
	StatusParameterError    = 501
	StatusCommandNotImpl    = 502
	StatusBadSequence       = 503
	StatusParameterNotImpl  = 504
	StatusMailboxUnavail    = 550
	StatusUserNotLocal      = 551
	StatusInsufficientSpace = 552
	StatusMailboxSyntax     = 553
	StatusTransactionFailed = 554
)

// NewSMTPServer creates a new SMTP server instance
func NewSMTPServer(config *EmailConfig, handler EmailHandler) *SMTPServer {
	return &SMTPServer{
		config:      config,
		handler:     handler,
		connections: make(map[string]*SMTPConnection),
	}
}

// Start starts the SMTP server on configured ports
func (s *SMTPServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("SMTP server already running")
	}

	// Start SMTP listener
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.config.SMTPPort))
	if err != nil {
		return fmt.Errorf("failed to start SMTP listener: %w", err)
	}

	// Start SMTPS listener if TLS is enabled
	if s.config.EnableTLS && s.config.TLSCertFile != "" && s.config.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(s.config.TLSCertFile, s.config.TLSKeyFile)
		if err != nil {
			s.listener.Close()
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		s.tlsListener, err = tls.Listen("tcp", fmt.Sprintf(":%d", s.config.SMTPSPort), tlsConfig)
		if err != nil {
			s.listener.Close()
			return fmt.Errorf("failed to start SMTPS listener: %w", err)
		}
	}

	s.running = true

	// Start accepting connections
	s.wg.Add(1)
	go s.acceptConnections(ctx, s.listener, false)

	if s.tlsListener != nil {
		s.wg.Add(1)
		go s.acceptConnections(ctx, s.tlsListener, true)
	}

	log.Printf("SMTP server started on port %d", s.config.SMTPPort)
	if s.tlsListener != nil {
		log.Printf("SMTPS server started on port %d", s.config.SMTPSPort)
	}

	return nil
}

// Stop stops the SMTP server
func (s *SMTPServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false

	// Close listeners
	if s.listener != nil {
		s.listener.Close()
	}
	if s.tlsListener != nil {
		s.tlsListener.Close()
	}

	// Close all active connections
	for _, conn := range s.connections {
		conn.conn.Close()
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("SMTP server stopped gracefully")
	case <-ctx.Done():
		log.Println("SMTP server stop timed out")
		return ctx.Err()
	}

	return nil
}

// acceptConnections accepts incoming connections
func (s *SMTPServer) acceptConnections(ctx context.Context, listener net.Listener, isTLS bool) {
	defer s.wg.Done()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if s.isRunning() {
				log.Printf("Failed to accept connection: %v", err)
			}
			return
		}

		// Check connection limit
		if len(s.connections) >= s.config.MaxConnections {
			conn.Close()
			continue
		}

		// Handle connection
		s.wg.Add(1)
		go s.handleConnection(ctx, conn, isTLS)
	}
}

// handleConnection handles a single SMTP connection
func (s *SMTPServer) handleConnection(ctx context.Context, conn net.Conn, isTLS bool) {
	defer s.wg.Done()
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(s.config.ConnectionTimeout))

	remoteAddr := conn.RemoteAddr().String()
	log.Printf("New SMTP connection from %s (TLS: %v)", remoteAddr, isTLS)

	// Create SMTP connection
	smtpConn := &SMTPConnection{
		conn:       conn,
		reader:     textproto.NewReader(bufio.NewReader(conn)),
		writer:     textproto.NewWriter(bufio.NewWriter(conn)),
		remoteAddr: remoteAddr,
		state:      StateConnected,
		startTime:  time.Now(),
	}

	// Register connection
	s.mu.Lock()
	s.connections[remoteAddr] = smtpConn
	s.mu.Unlock()

	// Unregister connection when done
	defer func() {
		s.mu.Lock()
		delete(s.connections, remoteAddr)
		s.mu.Unlock()
	}()

	// Send greeting
	if err := smtpConn.writeResponse(StatusReady, fmt.Sprintf("%s ESMTP Service Ready", s.config.Hostname)); err != nil {
		log.Printf("Failed to send greeting: %v", err)
		return
	}

	// Handle SMTP session
	s.handleSMTPSession(ctx, smtpConn)
}

// handleSMTPSession handles the SMTP protocol session
func (s *SMTPServer) handleSMTPSession(ctx context.Context, conn *SMTPConnection) {
	for {
		// Check if server is still running
		if !s.isRunning() {
			break
		}

		// Read command
		line, err := conn.reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read command: %v", err)
			}
			break
		}

		// Parse command
		parts := strings.Fields(line)
		if len(parts) == 0 {
			conn.writeResponse(StatusSyntaxError, "Syntax error")
			continue
		}

		command := strings.ToUpper(parts[0])
		args := parts[1:]

		// Handle command
		if err := s.handleSMTPCommand(conn, command, args); err != nil {
			log.Printf("Error handling command %s: %v", command, err)
			break
		}

		// Check if session should end
		if conn.state == StateQuit {
			break
		}
	}
}

// handleSMTPCommand handles individual SMTP commands
func (s *SMTPServer) handleSMTPCommand(conn *SMTPConnection, command string, args []string) error {
	switch command {
	case "HELO", "EHLO":
		return s.handleHelo(conn, command, args)
	case "MAIL":
		return s.handleMail(conn, args)
	case "RCPT":
		return s.handleRcpt(conn, args)
	case "DATA":
		return s.handleData(conn)
	case "RSET":
		return s.handleRset(conn)
	case "QUIT":
		return s.handleQuit(conn)
	case "NOOP":
		return conn.writeResponse(StatusOK, "OK")
	default:
		return conn.writeResponse(StatusCommandNotImpl, "Command not implemented")
	}
}

// handleHelo handles HELO/EHLO commands
func (s *SMTPServer) handleHelo(conn *SMTPConnection, command string, args []string) error {
	if len(args) == 0 {
		return conn.writeResponse(StatusParameterError, "Syntax error")
	}

	conn.helo = args[0]
	conn.state = StateHelo

	if command == "EHLO" {
		// Send ESMTP capabilities
		responses := []string{
			fmt.Sprintf("250-%s Hello %s", s.config.Hostname, conn.helo),
			"250-SIZE " + fmt.Sprintf("%d", s.config.MaxMessageSize),
			"250 HELP",
		}
		return conn.writeMultiResponse(responses)
	} else {
		return conn.writeResponse(StatusOK, fmt.Sprintf("%s Hello %s", s.config.Hostname, conn.helo))
	}
}

// handleMail handles MAIL FROM commands
func (s *SMTPServer) handleMail(conn *SMTPConnection, args []string) error {
	if conn.state != StateHelo {
		return conn.writeResponse(StatusBadSequence, "Bad sequence of commands")
	}

	if len(args) == 0 || !strings.HasPrefix(strings.ToUpper(args[0]), "FROM:") {
		return conn.writeResponse(StatusParameterError, "Syntax error")
	}

	// Extract email address
	fromAddr := strings.TrimPrefix(strings.ToUpper(args[0]), "FROM:")
	fromAddr = strings.Trim(fromAddr, "<>")
	
	conn.from = fromAddr
	conn.to = nil // Reset recipients
	conn.state = StateMailFrom

	return conn.writeResponse(StatusOK, "OK")
}

// handleRcpt handles RCPT TO commands
func (s *SMTPServer) handleRcpt(conn *SMTPConnection, args []string) error {
	if conn.state != StateMailFrom && conn.state != StateRcptTo {
		return conn.writeResponse(StatusBadSequence, "Bad sequence of commands")
	}

	if len(args) == 0 || !strings.HasPrefix(strings.ToUpper(args[0]), "TO:") {
		return conn.writeResponse(StatusParameterError, "Syntax error")
	}

	// Extract email address
	toAddr := strings.TrimPrefix(strings.ToUpper(args[0]), "TO:")
	toAddr = strings.Trim(toAddr, "<>")
	
	conn.to = append(conn.to, toAddr)
	conn.state = StateRcptTo

	return conn.writeResponse(StatusOK, "OK")
}

// handleData handles DATA command
func (s *SMTPServer) handleData(conn *SMTPConnection) error {
	if conn.state != StateRcptTo {
		return conn.writeResponse(StatusBadSequence, "Bad sequence of commands")
	}

	if err := conn.writeResponse(StatusStartData, "Start mail input; end with <CRLF>.<CRLF>"); err != nil {
		return err
	}

	// Read email data until "."
	var data []byte
	for {
		line, err := conn.reader.ReadLine()
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
		if int64(len(data)) > s.config.MaxMessageSize {
			return conn.writeResponse(StatusInsufficientSpace, "Message too large")
		}
	}

	conn.data = data

	// Process the email
	if s.handler != nil {
		if err := s.handler(data, conn.from, conn.to); err != nil {
			log.Printf("Email handler error: %v", err)
			return conn.writeResponse(StatusTransactionFailed, "Transaction failed")
		}
	}

	// Reset for next message
	conn.from = ""
	conn.to = nil
	conn.data = nil
	conn.state = StateHelo

	return conn.writeResponse(StatusOK, "OK")
}

// handleRset handles RSET command
func (s *SMTPServer) handleRset(conn *SMTPConnection) error {
	conn.from = ""
	conn.to = nil
	conn.data = nil
	conn.state = StateHelo
	return conn.writeResponse(StatusOK, "OK")
}

// handleQuit handles QUIT command
func (s *SMTPServer) handleQuit(conn *SMTPConnection) error {
	conn.state = StateQuit
	return conn.writeResponse(StatusClosing, "Bye")
}

// writeResponse writes a single SMTP response
func (conn *SMTPConnection) writeResponse(code int, message string) error {
	return conn.writer.PrintfLine("%d %s", code, message)
}

// writeMultiResponse writes multiple SMTP response lines
func (conn *SMTPConnection) writeMultiResponse(responses []string) error {
	for _, response := range responses {
		if err := conn.writer.PrintfLine("%s", response); err != nil {
			return err
		}
	}
	return nil
}

// isRunning checks if the server is running (thread-safe)
func (s *SMTPServer) isRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
