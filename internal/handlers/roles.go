package handlers

import (
	"net/http"

	"github.com/mathif92/auth-service/internal/api"
	"github.com/mathif92/auth-service/internal/services"
)

type Roles struct {
	rolesService *services.Roles
}

func NewRoles(rolesService *services.Roles) *Roles {
	return &Roles{rolesService: rolesService}
}

func (r *Roles) CreateRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var role CreateRoleInput
	if err := role.Bind(req); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	roleID, err := r.rolesService.SaveRole(ctx, services.RoleModel{Name: role.Name})
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, ResourceCreatedResponse{ID: roleID}, http.StatusCreated)
}

func (r *Roles) UpdateRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var role UpdateRoleInput
	if err := role.Bind(req); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.UpdateRole(ctx, services.RoleModel{Name: role.Name}, role.Enabled); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusCreated)
}

func (r *Roles) DeleteRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	// TODO: Not implemented yet
	// TODO: Decide whether to delete the row phisically or logically
	api.Respond(ctx, w, nil, http.StatusNotImplemented)
}

func (r *Roles) AddActionToRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var roleAction RoleActionInput
	if err := roleAction.Bind(req); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.AddActionToRole(ctx, services.RolesActionsModel{
		RoleID:   roleAction.RoleID,
		ActionID: roleAction.ActionID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) RemoveActionFromRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var roleAction RoleActionInput
	if err := roleAction.Bind(req); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.RemoveActionFromRole(ctx, services.RolesActionsModel{
		RoleID:   roleAction.RoleID,
		ActionID: roleAction.ActionID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) AssignRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var roleCreds RoleCredentialsInput
	if err := roleCreds.Bind(req); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.AddRoleToCredentials(ctx, services.RolesCredentialsModel{
		RoleID:        roleCreds.RoleID,
		CredentialsID: roleCreds.CredentialsID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) UnassignRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var roleCreds RoleCredentialsInput
	if err := roleCreds.Bind(req); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.UnassignRole(ctx, services.RolesCredentialsModel{
		RoleID:        roleCreds.RoleID,
		CredentialsID: roleCreds.CredentialsID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}
