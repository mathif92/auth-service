package handlers

import (
	"net/http"

	"github.com/mathif92/auth-service/internal/errors"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type CredentialsInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

type CreateRoleInput struct {
	Name string `json:"name"`
}

func (r *CreateRoleInput) Bind(req *http.Request) error {
	if r.Name == "" {
		return errors.New("name must be provided", http.StatusBadRequest)
	}

	return nil
}

type UpdateRoleInput struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Enabled *bool  `json:"enabled,omitempty"`
}

func (r *UpdateRoleInput) Bind(req *http.Request) error {
	if r.ID <= 0 {
		return errors.New("id must be provided", http.StatusBadRequest)
	}

	if r.Name == "" {
		return errors.New("name must be provided", http.StatusBadRequest)
	}

	return nil
}

type CreateActionInput struct {
	Action string `json:"action"`
	Entity string `json:"entity"`
}

func (a *CreateActionInput) Bind(r *http.Request) error {
	if a.Action == "" {
		return errors.New("action must be provided", http.StatusBadRequest)
	}

	if a.Entity == "" {
		return errors.New("entity must be provided", http.StatusBadRequest)
	}

	return nil
}

type UpdateActionInput struct {
	ID      int64  `json:"id"`
	Action  string `json:"action"`
	Entity  string `json:"entity"`
	Enabled *bool  `json:"enabled,omitempty"`
}

func (a *UpdateActionInput) Bind(r *http.Request) error {
	if a.ID <= 0 {
		return errors.New("id must be provided", http.StatusBadRequest)
	}
	if a.Action == "" {
		return errors.New("action must be provided", http.StatusBadRequest)
	}

	if a.Entity == "" {
		return errors.New("entity must be provided", http.StatusBadRequest)
	}

	return nil
}

type RoleActionInput struct {
	RoleID   int64 `json:"role_id"`
	ActionID int64 `json:"action_id"`
}

func (r *RoleActionInput) Bind(req *http.Request) error {
	if r.RoleID <= 0 {
		return errors.New("role_id must be provided", http.StatusBadRequest)
	}

	if r.ActionID <= 0 {
		return errors.New("action_id must be provided", http.StatusBadRequest)
	}

	return nil
}

type RoleCredentialsInput struct {
	RoleID        int64 `json:"role_id"`
	CredentialsID int64 `json:"credentials_id"`
}

func (r *RoleCredentialsInput) Bind(req *http.Request) error {
	if r.RoleID <= 0 {
		return errors.New("role_id must be provided", http.StatusBadRequest)
	}

	if r.CredentialsID <= 0 {
		return errors.New("credentials_id must be provided", http.StatusBadRequest)
	}

	return nil
}

type AuthenticationResponse struct {
	Token string `json:"token"`
}

type ResourceCreatedResponse struct {
	ID int64 `json:"id"`
}
