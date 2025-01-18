package smtp

import (
	"os"
	"strconv"

	"github.com/rudrprasad05/go-logs/logs"
)

// Config holds server configuration details
type Config struct {
	Host string
	Port int
	LOG  *logs.Logger
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() Config {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := 587 // Default port for testing
	if portEnv := os.Getenv("SMTP_PORT"); portEnv != "" {
		if parsedPort, err := strconv.Atoi(portEnv); err == nil {
			port = parsedPort
		}
	}
	logger, _ := logs.NewLogger()

	return Config{
		Host: host,
		Port: port,
		LOG:  logger,
	}
}
