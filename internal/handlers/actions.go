package handlers

import (
	"net/http"

	"github.com/mathif92/auth-service/internal/api"
	"github.com/mathif92/auth-service/internal/services"
)

type Actions struct {
	actionsService *services.Actions
}

func NewActions(actionsService *services.Actions) *Actions {
	return &Actions{actionsService: actionsService}
}

func (a *Actions) CreateAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var action CreateActionInput
	if err := action.Bind(r); err != nil {
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

func (a *Actions) UpdateAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var action UpdateActionInput
	if err := action.Bind(r); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	actionModel := services.ActionModel{
		Entity: action.Entity,
		Action: action.Action,
	}
	if err := a.actionsService.UpdateAction(ctx, actionModel, action.Enabled); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusOK)
}
