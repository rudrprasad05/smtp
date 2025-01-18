package main

import (
	"log"

	"rudrprasad.com/server/smtp"
)

func main() {
	// Load configuration
	config := smtp.LoadConfig()

	// Create and start the SMTP server
	smtpServer := smtp.NewSMTPServer(config)
	log.Printf("Starting SMTP server on %s:%d\n", config.Host, config.Port)
	if err := smtpServer.Start(); err != nil {
		log.Fatalf("Failed to start SMTP server: %v", err)
	}
}
