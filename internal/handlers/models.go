package handlers

import (
	"net/http"
	"time"

	"github.com/mathif92/auth-service/internal/errors"
	"github.com/mathif92/auth-service/internal/services"
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

type RoleResponse struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name"`
	Enabled   bool             `json:"enabled"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Actions   []ActionResponse `json:"actions"`
}

type ActionResponse struct {
	ID        int64     `json:"id"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ConvertRoleFromDBModel(role services.RoleWithActions) RoleResponse {
	var actions []ActionResponse
	for _, action := range role.Actions {
		actions = append(actions, ActionResponse{
			ID:        action.ID,
			Action:    action.Action,
			Entity:    action.Entity,
			Enabled:   action.Enabled,
			CreatedAt: action.CreatedAt,
			UpdatedAt: action.UpdatedAt,
		})
	}

	return RoleResponse{
		ID:        role.Role.ID,
		Name:      role.Role.Name,
		Enabled:   role.Role.Enabled,
		CreatedAt: role.Role.CreatedAt,
		UpdatedAt: role.Role.UpdatedAt,
		Actions:   actions,
	}
}
