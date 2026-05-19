package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"awst-agent/pkg/models"
)

type CentralClient struct {
	ServerURL  string
	HTTPClient *http.Client
}

func NewCentralClient(serverURL string) *CentralClient {
	return &CentralClient{
		ServerURL:  serverURL,
		HTTPClient: &http.Client{},
	}
}

func (c *CentralClient) Login(email, password string) (*models.LoginResponse, error) {
	loginData := models.LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(
		c.ServerURL+"/api/agent/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var loginResp models.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("response parsing error: %v", err)
	}

	return &loginResp, nil
}
