package utils

import (
	"os"
)

type EnvConfig struct {
	CentralServer string
	BackendServer string
	ProxyPort     string
}

func LoadEnv() (*EnvConfig, error) {
	config := &EnvConfig{
		CentralServer: getEnv("CENTRAL_SERVER", "http://localhost:5000"),
		BackendServer: getEnv("BACKEND_SERVER", "http://localhost:3000"),
		ProxyPort:     getEnv("PROXY_PORT", ":1133"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
