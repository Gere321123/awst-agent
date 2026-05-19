package utils
package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	CentralServer string
	BackendServer string
	ProxyPort     string
}

func LoadEnv() (*EnvConfig, error) {
	// Try to load .env file (optional)
	_ = godotenv.Load()

	config := &EnvConfig{
		CentralServer: getEnv("CENTRAL_SERVER", "http://localhost:5000"),
		BackendServer: getEnv("BACKEND_SERVER", "http://localhost:3000"),
		ProxyPort:     getEnv("PROXY_PORT", ":8080"),
	}

	log.Printf("Configuration loaded - Central: %s, Backend: %s, Port: %s", 
		config.CentralServer, config.BackendServer, config.ProxyPort)

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}