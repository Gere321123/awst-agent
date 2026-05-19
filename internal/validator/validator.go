package validator

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"awst-agent/pkg/models"
)

type Validator struct {
	cache *models.AgentCache
}

func NewValidator(cache *models.AgentCache) *Validator {
	return &Validator{
		cache: cache,
	}
}

func (v *Validator) UpdateCache(cache *models.AgentCache) {
	v.cache = cache
}

// ValidateAPIKey ellenőrzi, hogy az API kulcs érvényes-e a kért URL-hez
func (v *Validator) ValidateAPIKey(publicKey, requestedURL string) (*models.CachedAPIKey, error) {
	if v.cache == nil || v.cache.APIKeys == nil {
		return nil, fmt.Errorf("cache not initialized")
	}

	// 1. API kulcs keresése
	apiKey, exists := v.cache.APIKeys[publicKey]
	if !exists {
		return nil, fmt.Errorf("API key not found")
	}

	// 2. Aktív-e?
	if !apiKey.IsActive {
		return nil, fmt.Errorf("API key is inactive")
	}

	// 3. Visszavonták?
	if apiKey.IsRevoked {
		return nil, fmt.Errorf("API key has been revoked")
	}

	// 4. Lejárt?
	if time.Now().After(apiKey.ValidUntil) {
		return nil, fmt.Errorf("API key expired on %s", apiKey.ValidUntil.Format("2006-01-02"))
	}

	// 5. Domain ellenőrzés
	if !strings.Contains(requestedURL, apiKey.Domain) {
		return nil, fmt.Errorf("domain mismatch: key is for %s, requested %s", apiKey.Domain, requestedURL)
	}

	// 6. Allowed URLs ellenőrzés
	if len(apiKey.AllowedURLs) > 0 {
		allowed := false
		for _, allowedURL := range apiKey.AllowedURLs {
			if strings.Contains(requestedURL, allowedURL) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("URL not allowed: %s (allowed: %v)", requestedURL, apiKey.AllowedURLs)
		}
	}

	// 7. Rate limit ellenőrzés
	if err := v.checkRateLimit(apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// checkRateLimit ellenőrzi a rate limitet
func (v *Validator) checkRateLimit(apiKey *models.CachedAPIKey) error {
	if apiKey.UsageStats == nil {
		apiKey.UsageStats = &models.UsageStats{
			DailyUsage: make(map[string]int),
			LastReset:  time.Now(),
		}
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	// Reseteljük a napi számlálót ha új nap van
	if apiKey.UsageStats.LastReset.Format("2006-01-02") != today {
		apiKey.UsageStats.DailyUsage = make(map[string]int)
		apiKey.UsageStats.LastReset = now
	}

	// Ellenőrizzük a limitet
	currentUsage := apiKey.UsageStats.DailyUsage[today]
	if currentUsage >= apiKey.RateLimitCount {
		return fmt.Errorf("rate limit exceeded: %d/%d requests today", currentUsage, apiKey.RateLimitCount)
	}

	return nil
}

// UpdateUsage frissíti a használati statisztikákat
func (v *Validator) UpdateUsage(publicKey, url string) error {
	if v.cache == nil || v.cache.APIKeys == nil {
		return fmt.Errorf("cache not initialized")
	}

	apiKey, exists := v.cache.APIKeys[publicKey]
	if !exists {
		return fmt.Errorf("API key not found")
	}

	if apiKey.UsageStats == nil {
		apiKey.UsageStats = &models.UsageStats{
			DailyUsage: make(map[string]int),
			LastReset:  time.Now(),
		}
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	// Reset if new day
	if apiKey.UsageStats.LastReset.Format("2006-01-02") != today {
		apiKey.UsageStats.DailyUsage = make(map[string]int)
		apiKey.UsageStats.LastReset = now
	}

	// Increment usage
	apiKey.UsageStats.DailyUsage[today]++
	apiKey.UsageStats.TotalRequests++

	return nil
}

// IsInternalRequest ellenőrzi, hogy a kérés belső (frontend) -e
func (v *Validator) IsInternalRequest(r *http.Request) bool {
	// Belső IP címek
	internalIPs := []string{
		"127.0.0.1",
		"localhost",
		"::1",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	clientIP := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(clientIP, ":"); idx != -1 {
		clientIP = clientIP[:idx]
	}

	// Check if IP is internal (simple version)
	for _, ip := range internalIPs {
		if clientIP == ip || strings.HasPrefix(clientIP, ip) {
			return true
		}
	}

	// Check X-Forwarded-For header for internal requests
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		for _, ip := range internalIPs {
			if strings.Contains(forwarded, ip) {
				return true
			}
		}
	}

	return false
}
