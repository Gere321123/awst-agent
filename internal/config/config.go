package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"awst-agent/pkg/models"
)

const (
	ConfigFile = "/etc/awst/config.json"
	CacheFile  = "/etc/awst/cache.json"
)

type Manager struct {
	ConfigPath string
	CachePath  string
}

func NewManager() *Manager {
	return &Manager{
		ConfigPath: ConfigFile,
		CachePath:  CacheFile,
	}
}

func (m *Manager) Save(response *models.LoginResponse) error {
	configDir := filepath.Dir(m.ConfigPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	config := models.AgentConfig{
		Email:   response.Email,
		Token:   response.Token,
		AgentID: response.AgentID,
		Cache:   response.CacheData,
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(m.ConfigPath, jsonData, 0600); err != nil {
		return err
	}

	// Külön cache fájl is (opcionális)
	if response.CacheData != nil {
		cacheData, err := json.MarshalIndent(response.CacheData, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(m.CachePath, cacheData, 0600); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Load() (*models.AgentConfig, error) {
	jsonData, err := os.ReadFile(m.ConfigPath)
	if err != nil {
		return nil, err
	}

	var config models.AgentConfig
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (m *Manager) LoadCache() (*models.AgentCache, error) {
	jsonData, err := os.ReadFile(m.CachePath)
	if err != nil {
		return nil, err
	}

	var cache models.AgentCache
	if err := json.Unmarshal(jsonData, &cache); err != nil {
		return nil, err
	}

	// Inicializáljuk a map-et ha üres
	if cache.APIKeys == nil {
		cache.APIKeys = make(map[string]*models.CachedAPIKey)
	}

	return &cache, nil
}

func (m *Manager) Exists() bool {
	_, err := os.Stat(m.ConfigPath)
	return !os.IsNotExist(err)
}

func (m *Manager) SaveCache(cache *models.AgentCache) error {
	configDir := filepath.Dir(m.CachePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.CachePath, jsonData, 0600)
}
