// Package handler exposes the identity module over HTTP with Huma. It is the only place that maps domain errors to HTTP statuses.
package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/identity/use_cases/command"
	"joao-api-comercio/internal/modules/identity/user"
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
		return nil, h.toHumaError(ctx, err)
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
		return nil, h.toHumaError(ctx, err)
	}

	out := &loginOutput{}
	out.Body.Token = res.Token
	out.Body.ExpiresAt = res.ExpiresAt
	return out, nil
}

// toHumaError turns a domain error into an HTTP error. Anything unknown is logged in full and sent
// to the client as a plain 500, so no internal detail leaks.
func (h *Identity) toHumaError(ctx context.Context, err error) error {
	var invalidEmail *user.InvalidEmailError
	var weakPassword *user.WeakPasswordError
	var taken *user.EmailTakenError
	var badCredentials *user.InvalidCredentialsError

	switch {
	case errors.As(err, &invalidEmail), errors.As(err, &weakPassword):
		return huma.Error422UnprocessableEntity(err.Error())
	case errors.As(err, &taken):
		return huma.Error409Conflict(err.Error())
	case errors.As(err, &badCredentials):
		return huma.Error401Unauthorized("invalid credentials")
	default:
		h.log.ErrorContext(ctx, "identity request failed", "err", err)
		return huma.Error500InternalServerError("internal error")
	}
}
