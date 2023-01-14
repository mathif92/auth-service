package handlers

import (
	"net/http"

	"github.com/mathif92/auth-service/internal/api"
)

type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Health(w http.ResponseWriter, r *http.Request) {
	api.Respond(r.Context(), w, HealthResponse{Status: "OK"}, http.StatusOK)
}
