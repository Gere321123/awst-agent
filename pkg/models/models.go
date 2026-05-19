package models

type AgentConfig struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	AgentID string `json:"agent_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	AgentID string `json:"agent_id"`
	Token   string `json:"token"`
	Email   string `json:"email"`
	Message string `json:"message"`
}
