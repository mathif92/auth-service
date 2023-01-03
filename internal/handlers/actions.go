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

type actionCtxKey string

type Actions struct {
	actionsService *services.Actions
}

func NewActions(actionsService *services.Actions) *Actions {
	return &Actions{actionsService: actionsService}
}

func (a *Actions) CreateAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var action CreateActionInput
	if err := render.Bind(r, &action); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	actionID, err := a.actionsService.SaveAction(ctx, services.ActionModel{
		Entity: action.Entity,
		Action: action.Action,
	})
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, ResourceCreatedResponse{ID: actionID}, http.StatusCreated)
}

func (a *Actions) GetAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	action, ok := ctx.Value(actionCtxKey("action")).(ActionResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("action not found", http.StatusBadRequest))
		return
	}

	api.Respond(ctx, w, action, http.StatusOK)
}

func (a *Actions) UpdateAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	action, ok := ctx.Value(actionCtxKey("action")).(ActionResponse)
	if !ok {
		api.RespondError(ctx, w, errors.New("action not found", http.StatusBadRequest))
		return
	}

	var actionInput UpdateActionInput
	if err := render.Bind(r, &actionInput); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	actionModel := services.ActionModel{
		ID:     action.ID,
		Entity: actionInput.Entity,
		Action: actionInput.Action,
	}
	if err := a.actionsService.UpdateAction(ctx, actionModel, actionInput.Enabled); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}

func (a *Actions) ActionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		actionIDStr := chi.URLParam(req, "actionID")
		actionID, err := strconv.Atoi(actionIDStr)
		if err != nil {
			api.RespondError(ctx, w, errors.New("actionID must be a number", http.StatusBadRequest))
			return
		}
		action, err := a.actionsService.GetAction(ctx, int64(actionID))
		if err != nil {
			api.RespondError(ctx, w, err)
			return
		}
		ctx = context.WithValue(req.Context(), actionCtxKey("action"), ConvertActionFromDBModel(action))
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
