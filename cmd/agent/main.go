package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"awst-agent/internal/client"
	"awst-agent/internal/config"
	"awst-agent/internal/proxy"
	"awst-agent/internal/validator"
	"awst-agent/pkg/models"
	"awst-agent/pkg/utils"
)

func main() {
	// Load environment variables
	envConfig, err := utils.LoadEnv()
	if err != nil {
		log.Printf("Error loading environment: %v", err)
		os.Exit(1)
	}

	configManager := config.NewManager()
	centralClient := client.NewCentralClient(envConfig.CentralServer)

	// Check login status
	if err := checkLogin(configManager, centralClient, envConfig); err != nil {
		log.Printf("Login error: %v", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := configManager.Load()
	if err != nil {
		log.Printf("Configuration loading error: %v", err)
		os.Exit(1)
	}

	// Load cache
	cache, err := configManager.LoadCache()
	if err != nil {
		log.Printf("Warning: Could not load cache: %v", err)
		cache = &models.AgentCache{
			APIKeys: make(map[string]*models.CachedAPIKey),
		}
	}

	// Initialize validator
	validatorInstance := validator.NewValidator(cache)

	// Print protected domains
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("   AWST Agent - Protecting your data")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n✅ Agent running as: %s\n", cfg.Email)
	fmt.Printf("📍 Central server: %s\n", envConfig.CentralServer)
	fmt.Printf("🔒 Proxy running on: http://localhost%s\n", envConfig.ProxyPort)

	if len(cache.Domains) > 0 {
		fmt.Println("\n🛡️  Protected domains:")
		for _, domain := range cache.Domains {
			fmt.Printf("   - %s\n", domain)
		}
	}
	fmt.Println(strings.Repeat("-", 60))

	// Setup proxy with validator
	agentProxy, err := proxy.NewAgentProxy(envConfig.BackendServer, cfg.Token)
	if err != nil {
		log.Printf("Proxy initialization error: %v", err)
		os.Exit(1)
	}

	// Main handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's an internal request (from frontend)
		if validatorInstance.IsInternalRequest(r) {
			log.Printf("🔄 Internal request from %s - forwarding to backend", r.RemoteAddr)
			agentProxy.ServeHTTP(w, r)
			return
		}

		// External request - need API key
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		}

		if apiKey == "" {
			log.Printf("❌ External request from %s - missing API key", r.RemoteAddr)
			http.Redirect(w, r, "https://www.youtube.com/watch?v=BGpzGu9Yp6Y", http.StatusTemporaryRedirect)
			return
		}

		// Validate API key
		validKey, err := validatorInstance.ValidateAPIKey(apiKey, r.URL.Path)
		if err != nil {
			log.Printf("❌ Invalid API key from %s: %v", r.RemoteAddr, err)
			http.Redirect(w, r, "https://www.youtube.com/watch?v=BGpzGu9Yp6Y", http.StatusUnauthorized)
			return
		}

		// Update usage stats
		if err := validatorInstance.UpdateUsage(apiKey, r.URL.Path); err != nil {
			log.Printf("⚠️ Could not update usage: %v", err)
		}

		log.Printf("✅ Valid request from %s - API key: %s... (domain: %s, remaining: %d/%d)",
			r.RemoteAddr,
			apiKey[:8],
			validKey.Domain,
			validKey.RateLimitCount-getDailyUsage(validKey),
			validKey.RateLimitCount)

		// Forward to backend
		agentProxy.ServeHTTP(w, r)
	})

	log.Printf("🚀 Starting proxy on %s", envConfig.ProxyPort)
	log.Println("Press CTRL+C to stop")

	if err := http.ListenAndServe(envConfig.ProxyPort, nil); err != nil {
		log.Fatal(err)
	}
}

func getDailyUsage(apiKey *models.CachedAPIKey) int {
	if apiKey.UsageStats == nil {
		return 0
	}
	today := time.Now().Format("2006-01-02")
	return apiKey.UsageStats.DailyUsage[today]
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

		if response.CacheData != nil && len(response.CacheData.APIKeys) > 0 {
			fmt.Printf("  📦 Loaded %d API keys\n", len(response.CacheData.APIKeys))
			fmt.Printf("  🌐 Protecting %d domains\n", len(response.CacheData.Domains))
		}

		fmt.Printf("  Configuration saved: %s\n\n", configManager.ConfigPath)
	}

	return nil
}
