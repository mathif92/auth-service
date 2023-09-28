package handlers

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/mathif92/auth-service/internal/api"
)

type Health struct {
	db *sqlx.DB
}

func NewHealth(db *sqlx.DB) *Health {
	return &Health{
		db: db,
	}
}

// Health checks the app status by checking the DB connection, returns an error if the DB ping is not successful and status OK otherwise
func (h *Health) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.db.PingContext(ctx); err != nil {
		api.Respond(ctx, w, HealthResponse{Status: fmt.Sprintf("NOT_HEALTHY: %s", err.Error())}, http.StatusInternalServerError)
		return
	}
	api.Respond(ctx, w, HealthResponse{Status: "OK"}, http.StatusOK)
}
