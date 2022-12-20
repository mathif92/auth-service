package handlers

import (
	"net/http"

	"github.com/mathif92/auth-service/internal/services"
)

type Actions struct {
	actionsService *services.Actions
}

func NewActions(actionsService *services.Actions) *Actions {
	return &Actions{actionsService: actionsService}
}

func (a *Actions) CreateAction(w http.ResponseWriter, r *http.Request) {

}

func (a *Actions) UpdateAction(w http.ResponseWriter, r *http.Request) {

}
