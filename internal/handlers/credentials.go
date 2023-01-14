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

type credentialsCtxKey string

type Credentials struct {
	credentialsService *services.Credentials
}

func NewCredentialsHandler(credentialsService *services.Credentials) *Credentials {
	return &Credentials{credentialsService: credentialsService}
}

func (c *Credentials) CreateCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var credentials CredentialsInput
	if err := render.Bind(r, &credentials); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if err := c.credentialsService.SaveCredentials(ctx, services.CredentialsModel{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: credentials.Password,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusCreated)
}

func (c *Credentials) ValidateCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var credentials CredentialsInput
	if err := render.Bind(r, &credentials); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	token, err := c.credentialsService.Authenticate(ctx, services.CredentialsModel{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, AuthenticationResponse{Token: token}, http.StatusOK)
}

func (c *Credentials) CredentialsContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		credentialsIDStr := chi.URLParam(req, "credentialsID")
		credentialsID, err := strconv.Atoi(credentialsIDStr)
		if err != nil {
			api.RespondError(ctx, w, errors.New("credentialsID must be a number", http.StatusBadRequest))
			return
		}
		credentials, err := c.credentialsService.GetCredentials(ctx, int64(credentialsID))
		if err != nil {
			api.RespondError(ctx, w, err)
			return
		}
		ctx = context.WithValue(req.Context(), credentialsCtxKey("credentials"), ConvertCredentialsFromDBModel(credentials))
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
