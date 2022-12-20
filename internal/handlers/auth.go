package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/mathif92/auth-service/internal/api"
	"github.com/mathif92/auth-service/internal/services"
)

type Authentication struct {
	authService *services.Authentication
}

func NewAuthenticationHandler(authService *services.Authentication) *Authentication {
	return &Authentication{authService: authService}
}

func (a *Authentication) CreateCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var credentials CredentialsInput
	if err := render.Bind(r, &credentials); err != nil {
		api.RespondError(ctx, w, err)
		return
	}
	if err := a.authService.SaveCredentials(ctx, services.CredentialsModel{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: credentials.Password,
	}); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	api.Respond(ctx, w, nil, http.StatusCreated)
}

func (a *Authentication) ValidateCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var credentials CredentialsInput
	if err := render.Bind(r, &credentials); err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	token, err := a.authService.Authenticate(ctx, services.CredentialsModel{
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
