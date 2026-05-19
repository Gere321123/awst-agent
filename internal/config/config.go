package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Gere321123/awst-agent/pkg/models"
)

const (
	ConfigFile = "/etc/awst/config.json"
)

type Manager struct {
	ConfigPath string
}

func NewManager() *Manager {
	return &Manager{
		ConfigPath: ConfigFile,
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
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(m.ConfigPath, jsonData, 0600); err != nil {
		return err
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

func (m *Manager) Exists() bool {
	_, err := os.Stat(m.ConfigPath)
	return !os.IsNotExist(err)
}
