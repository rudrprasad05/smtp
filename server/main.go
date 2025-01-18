package main

import (
	"fmt"
	"log"

	"github.com/rudrprasad05/go-logs/logs"
	"rudrprasad.com/server/smtp"
)

func main() {
	// Load configuration
	config := smtp.LoadConfig()
	logger, err := logs.NewLogger()

	if err != nil {
		logger.Info(fmt.Sprintf("err %v", err))
		return
	}

	// Create and start the SMTP server
	smtpServer := smtp.NewSMTPServer(config)
	logger.Info(fmt.Sprintf("Starting SMTP server on %s:%d\n", config.Host, config.Port))
	if err := smtpServer.Start(); err != nil {
		log.Fatalf("Failed to start SMTP server: %v", err)
		logger.Info(fmt.Sprintf("Failed to start SMTP server: %v", err))

	}
}
