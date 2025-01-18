package smtp

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/rudrprasad05/go-logs/logs"
)

// SMTPServer represents the SMTP server
type SMTPServer struct {
	Host      string
	Port      int
	LOG       *logs.Logger
	TLSConfig *tls.Config
}

// NewSMTPServer creates a new SMTPServer instance
func NewSMTPServer(config Config) *SMTPServer {
	// Set up TLS configuration
	certFileLinux := "/etc/letsencrypt/live/rudrprasad.com/cert.pem"
	keyFileLinux := "/etc/letsencrypt/live/rudrprasad.com/key.pem"

	// cert, err := tls.LoadX509KeyPair("../server/certificates/smtp-cert.pem", "../server/certificates/smtp-key.pem")
	cert, err := tls.LoadX509KeyPair(certFileLinux, keyFileLinux)

	if err != nil {
		log.Fatalf("Failed to load cert and key: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true, // Skip certificate verification for testing purposes
	}

	return &SMTPServer{
		Host:      config.Host,
		Port:      config.Port,
		LOG:       config.LOG,
		TLSConfig: tlsConfig,
	}
}

// Start begins the server and listens for connections
func (s *SMTPServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("SMTP server listening on %s\n", address)
	s.LOG.Info(fmt.Sprintf("SMTP server listening on %s\n", address))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			s.LOG.Info(fmt.Sprintf("Error accepting connection: %v\n", err))

			continue
		}

		go s.handleConnection(conn)
	}
}
func (s *SMTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s\n", conn.RemoteAddr())

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send initial 220 response
	writer.WriteString("220 rudrprasad.com ESMTP Service Ready\r\n")
	writer.Flush()

	var tlsEstablished bool
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection closed: %v\n", err)
			break
		}

		command := strings.TrimSpace(line)
		log.Printf("Received: %s\n", command)

		switch {
		case strings.HasPrefix(strings.ToUpper(command), "HELO"):
			writer.WriteString("250 Hello " + command[5:] + "\r\n")
		case strings.HasPrefix(strings.ToUpper(command), "MAIL FROM:"):
			writer.WriteString("250 Sender OK\r\n")
		case strings.HasPrefix(strings.ToUpper(command), "RCPT TO:"):
			writer.WriteString("250 Recipient OK\r\n")
		case strings.HasPrefix(strings.ToUpper(command), "DATA"):
			writer.WriteString("354 Start mail input; end with <CRLF>.<CRLF>\r\n")
		case command == ".":
			writer.WriteString("250 Message accepted for delivery\r\n")
		case strings.ToUpper(command) == "STARTTLS" && !tlsEstablished:
			// Respond with readiness to upgrade to TLS
			writer.WriteString("220 Ready to start TLS\r\n")
			writer.Flush()

			// Upgrade connection to TLS
			tlsConn := tls.Server(conn, s.TLSConfig)
			if err := tlsConn.Handshake(); err != nil {
				log.Printf("TLS handshake failed: %v\n", err)
				writer.WriteString("454 TLS handshake failed\r\n")
				writer.Flush()
				return
			}

			tlsEstablished = true
			writer = bufio.NewWriter(tlsConn)
			reader = bufio.NewReader(tlsConn)
			writer.WriteString("250 OK, TLS established\r\n")
			writer.Flush()

		case strings.ToUpper(command) == "QUIT":
			writer.WriteString("221 Bye\r\n")
			writer.Flush()
			return
		default:
			writer.WriteString("500 Command not recognized\r\n")
		}

		writer.Flush()
	}
}
