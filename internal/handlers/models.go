package handlers

import (
	"net/http"

	"github.com/mathif92/auth-service/internal/errors"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type CredentialsInput struct {
	Username string
	Email    string
	Password string
}

func (c *CredentialsInput) Bind(r *http.Request) error {
	if c.Username == "" && c.Email == "" {
		return errors.New("username or email must be provided", http.StatusBadRequest)
	}

	if c.Password == "" {
		return errors.New("password must be provided", http.StatusBadRequest)
	}

	return nil
}

type AuthenticationResponse struct {
	Token string `json:"token"`
}
