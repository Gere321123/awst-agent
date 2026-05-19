// internal/cache/manager.go (új fájl)
package cache

import (
	"awst-agent/pkg/models"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type CacheManager struct {
	mu    sync.RWMutex
	cache *models.AgentCache
	path  string
}

func NewCacheManager(cachePath string) *CacheManager {
	return &CacheManager{
		path: cachePath,
		cache: &models.AgentCache{
			APIKeys: make(map[string]*models.CachedAPIKey),
		},
	}
}

func (cm *CacheManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, cm.cache)
}

func (cm *CacheManager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.path, data, 0600)
}

func (cm *CacheManager) ValidateAPIKey(publicKey, requestedURL string) (bool, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	apiKey, exists := cm.cache.APIKeys[publicKey]
	if !exists {
		return false, fmt.Errorf("API key not found")
	}

	// Validációk
	if !apiKey.IsActive || apiKey.IsRevoked {
		return false, fmt.Errorf("API key is inactive or revoked")
	}

	if time.Now().After(apiKey.ValidUntil) {
		return false, fmt.Errorf("API key expired")
	}

	// Domain check
	// URL check
	// Rate limit check

	return true, nil
}

func (cm *CacheManager) UpdateFromServer(loginResp *models.LoginResponse) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Frissíti a cache-t a szervertől kapott adatokkal
	// ...

	return cm.Save()
}
