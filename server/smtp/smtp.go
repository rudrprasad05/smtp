package smtp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// SMTPServer represents the SMTP server
type SMTPServer struct {
	Host string
	Port int
}

// NewSMTPServer creates a new SMTPServer instance
func NewSMTPServer(config Config) *SMTPServer {
	return &SMTPServer{
		Host: config.Host,
		Port: config.Port,
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
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

// handleConnection handles the SMTP session with a single client
func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s\n", conn.RemoteAddr())

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send initial 220 response
	writer.WriteString("220 rudrprasad.com ESMTP Service Ready\r\n")
	writer.Flush()

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
