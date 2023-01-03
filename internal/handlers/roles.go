package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mathif92/auth-service/internal/api"
	"github.com/mathif92/auth-service/internal/errors"
	"github.com/mathif92/auth-service/internal/services"
)

type roleCtxKey string

type Roles struct {
	rolesService *services.Roles
}

func NewRoles(rolesService *services.Roles) *Roles {
	return &Roles{rolesService: rolesService}
}

func (r *Roles) CreateRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var role CreateRoleInput
	if err := render.Bind(req, &role); err != nil {
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

func (r *Roles) GetRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	role, ok := ctx.Value(roleCtxKey("role")).(RoleResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("role not found", http.StatusNotFound))
		return
	}

	api.Respond(ctx, w, role, http.StatusCreated)
}

func (r *Roles) UpdateRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	role, ok := ctx.Value(roleCtxKey("role")).(RoleResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("role not found", http.StatusNotFound))
		return
	}

	var roleInput UpdateRoleInput
	if err := render.Bind(req, &roleInput); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.UpdateRole(ctx, services.RoleModel{ID: role.ID, Name: roleInput.Name}, roleInput.Enabled); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) DeleteRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	// TODO: Not implemented yet
	// TODO: Decide whether to delete the row phisically or logically
	api.Respond(ctx, w, nil, http.StatusNotImplemented)
}

func (r *Roles) AddActionToRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	role, ok := ctx.Value(roleCtxKey("role")).(RoleResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("role not found", http.StatusBadRequest))
		return
	}

	var roleAction RoleActionInput
	if err := render.Bind(req, &roleAction); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.AddActionToRole(ctx, services.RolesActionsModel{
		RoleID:   role.ID,
		ActionID: roleAction.ActionID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) RemoveActionFromRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	role, ok := ctx.Value(roleCtxKey("role")).(RoleResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("role not found", http.StatusBadRequest))
		return
	}

	var roleAction RoleActionInput
	if err := render.Bind(req, &roleAction); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.RemoveActionFromRole(ctx, services.RolesActionsModel{
		RoleID:   role.ID,
		ActionID: roleAction.ActionID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) AssignRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	role, ok := ctx.Value(roleCtxKey("role")).(RoleResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("role not found", http.StatusBadRequest))
		return
	}

	var roleCreds RoleCredentialsInput
	if err := render.Bind(req, &roleCreds); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.AddRoleToCredentials(ctx, services.RolesCredentialsModel{
		RoleID:        role.ID,
		CredentialsID: roleCreds.CredentialsID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) UnassignRole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	role, ok := ctx.Value(roleCtxKey("role")).(RoleResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("role not found", http.StatusBadRequest))
		return
	}

	var roleCreds RoleCredentialsInput
	if err := render.Bind(req, &roleCreds); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	if err := r.rolesService.UnassignRole(ctx, services.RolesCredentialsModel{
		RoleID:        role.ID,
		CredentialsID: roleCreds.CredentialsID,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (r *Roles) RoleContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		roleIDStr := chi.URLParam(req, "roleID")
		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			api.RespondError(ctx, w, errors.New("roleID must be a number", http.StatusBadRequest))
			return
		}
		role, err := r.rolesService.GetRole(ctx, int64(roleID))
		if err != nil {
			api.RespondError(ctx, w, err)
			return
		}
		ctx = context.WithValue(req.Context(), roleCtxKey("role"), ConvertRoleFromDBModel(role))
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
