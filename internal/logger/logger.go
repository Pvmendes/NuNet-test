package logger

import (
	"log"
	"os"
)

// Setup initializes the application's logger.
func Setup() {
	log.SetOutput(os.Stdout) // Sets the output destination for the logger
	log.SetPrefix("[deployer-manager] ") // Sets the prefix to appear on each log line
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile) // Configures the logging properties
}
