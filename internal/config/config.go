package config

import (
	"log"
	"os"
	"strconv"
)

// AppConfig holds the configuration settings for the application.
type AppConfig struct {
	ServerPort     int    `json:"serverPort"`
	DockerEndpoint string `json:"dockerEndpoint"`
}

// LoadConfig loads configuration settings from environment variables.
func LoadConfig() *AppConfig {
	config := &AppConfig{
		ServerPort:     getEnvAsInt("SERVER_PORT", 8080),
		DockerEndpoint: getEnvAsString("DOCKER_ENDPOINT", "unix:///var/run/docker.sock"),
	}
	return config
}

// Helper function to read an environment variable as a string or return a default value
func getEnvAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as an integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
		log.Printf("Warning: Could not parse %s as integer, defaulting to %d", key, defaultValue)
	}
	return defaultValue
}
