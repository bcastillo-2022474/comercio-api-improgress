// Package handler exposes the identity module over HTTP with Huma. It is the only place that maps domain errors to HTTP statuses.
package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/identity/use_cases/command"
	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared/apperr"
	"joao-api-comercio/internal/shared/httperr"
)

type Identity struct {
	registerUC *command.RegisterUseCase
	loginUC    *command.LoginUseCase
	log        *slog.Logger
}

func NewIdentity(registerUC *command.RegisterUseCase, loginUC *command.LoginUseCase, log *slog.Logger) *Identity {
	return &Identity{registerUC: registerUC, loginUC: loginUC, log: log}
}

// Mount adds the identity operations to api.
func (h *Identity) Mount(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Method:        http.MethodPost,
		Path:          "/register",
		Summary:       "Create an account",
		Tags:          []string{"identity"},
		DefaultStatus: http.StatusCreated,
	}, h.registerUser)

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Log in and receive an access token",
		Tags:        []string{"identity"},
	}, h.loginUser)
}

type registerInput struct {
	Body struct {
		Email    string `json:"email" doc:"Account email. Stored lowercase."`
		Password string `json:"password" doc:"Account password. Length limits are enforced by the domain."`
	}
}

type userOutput struct {
	Body struct {
		Id        string    `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}
}

type loginInput struct {
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}

type loginOutput struct {
	Body struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
}

func (h *Identity) registerUser(ctx context.Context, in *registerInput) (*userOutput, error) {
	u, err := h.registerUC.Execute(ctx, command.RegisterCommand{Email: in.Body.Email, Password: in.Body.Password})
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "register"), err)
	}

	out := &userOutput{}
	out.Body.Id = u.Id.String()
	out.Body.Email = string(u.Email)
	out.Body.CreatedAt = u.CreatedAt
	return out, nil
}

func (h *Identity) loginUser(ctx context.Context, in *loginInput) (*loginOutput, error) {
	res, err := h.loginUC.Execute(ctx, command.LoginCommand{Email: in.Body.Email, Password: in.Body.Password})
	if err != nil {
		if apperr.CodeOf(err) == user.CodeInvalidCredentials {
			h.log.WarnContext(ctx, "login failed")
		}
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "login"), err)
	}

	out := &loginOutput{}
	out.Body.Token = res.Token
	out.Body.ExpiresAt = res.ExpiresAt
	return out, nil
}
