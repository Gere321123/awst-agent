package models

import "time"

type AgentConfig struct {
	Email   string      `json:"email"`
	Token   string      `json:"token"`
	AgentID string      `json:"agent_id"`
	Cache   *AgentCache `json:"cache,omitempty"`
}

type AgentCache struct {
	SellerDomain string                   `json:"seller_domain"`
	Domains      []string                 `json:"domains"`
	APIKeys      map[string]*CachedAPIKey `json:"api_keys"`
	LastUpdate   time.Time                `json:"last_update"`
}

type CachedAPIKey struct {
	PublicKey       string      `json:"public_key"`
	PrivateKey      string      `json:"private_key"`
	Domain          string      `json:"domain"`
	AllowedURLs     []string    `json:"allowed_urls"`
	ValidUntil      time.Time   `json:"valid_until"`
	ValidFrom       time.Time   `json:"valid_from"`
	RateLimitCount  int         `json:"rate_limit_count"`
	RateLimitPeriod string      `json:"rate_limit_period"`
	UsageStats      *UsageStats `json:"usage_stats"`
	IsActive        bool        `json:"is_active"`
	IsRevoked       bool        `json:"is_revoked"`
}

type UsageStats struct {
	TotalRequests int            `json:"total_requests"`
	DailyUsage    map[string]int `json:"daily_usage"`
	LastReset     time.Time      `json:"last_reset"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success   bool        `json:"success"`
	AgentID   string      `json:"agent_id"`
	Token     string      `json:"token"`
	Email     string      `json:"email"`
	Message   string      `json:"message"`
	CacheData *AgentCache `json:"cache_data"`
}
