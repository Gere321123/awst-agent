package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Gere321123/awst-agent/internal/client"
	"github.com/Gere321123/awst-agent/internal/config"
	"github.com/Gere321123/awst-agent/internal/proxy"
	"github.com/Gere321123/awst-agent/pkg/utils"
)

func main() {
	// Load environment variables
	envConfig, err := utils.LoadEnv()
	if err != nil {
		fmt.Printf("Error loading environment: %v\n", err)
		os.Exit(1)
	}

	configManager := config.NewManager()
	centralClient := client.NewCentralClient(envConfig.CentralServer)

	// Check login status
	if err := checkLogin(configManager, centralClient, envConfig); err != nil {
		fmt.Printf("Login error: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := configManager.Load()
	if err != nil {
		fmt.Printf("Configuration loading error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Agent running: %s\n", cfg.Email)
	fmt.Printf("Connecting to central server: %s\n", envConfig.CentralServer)
	fmt.Printf("Proxying to backend: %s\n", envConfig.BackendServer)

	// Setup proxy
	agentProxy, err := proxy.NewAgentProxy(envConfig.BackendServer, cfg.Token)
	if err != nil {
		fmt.Printf("Proxy initialization error: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", agentProxy.Handler())

	fmt.Printf("Proxy running on: http://localhost%s\n", envConfig.ProxyPort)
	fmt.Println("Press CTRL+C to stop")

	if err := http.ListenAndServe(envConfig.ProxyPort, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func checkLogin(configManager *config.Manager, centralClient *client.CentralClient, envConfig *utils.EnvConfig) error {
	if !configManager.Exists() {
		fmt.Println("\n========================================")
		fmt.Println("   AWST Agent Setup and Login")
		fmt.Println("========================================\n")

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Print("Password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if email == "" || password == "" {
			return fmt.Errorf("email and password are required")
		}

		fmt.Printf("\nLogging in to central server: %s...\n", envConfig.CentralServer)

		response, err := centralClient.Login(email, password)
		if err != nil {
			return fmt.Errorf("network error: %v\nCheck if central server is reachable: %s", err, envConfig.CentralServer)
		}

		if !response.Success {
			return fmt.Errorf("login failed: %s", response.Message)
		}

		if err := configManager.Save(response); err != nil {
			return fmt.Errorf("configuration save error: %v", err)
		}

		fmt.Println("\n✓ Login successful!")
		fmt.Printf("  User: %s\n", response.Email)
		fmt.Printf("  Agent ID: %s\n", response.AgentID)
		fmt.Printf("  Configuration saved: %s\n\n", configManager.ConfigPath)
	}

	return nil
}
