package command

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared"
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token     string
	ExpiresAt time.Time
}

type LoginUseCase struct {
	uow    shared.UnitOfWork[user.Repository]
	hasher user.PasswordHasher
	tokens TokenIssuer
	log    *slog.Logger
}

func NewLoginUseCase(uow shared.UnitOfWork[user.Repository], hasher user.PasswordHasher, tokens TokenIssuer, log *slog.Logger) *LoginUseCase {
	return &LoginUseCase{uow: uow, hasher: hasher, tokens: tokens, log: log.With("use_case", "login")}
}

// Execute checks the credentials and issues a token. Every credential failure
// returns *user.InvalidCredentialsError. The log never holds the email.
func (uc *LoginUseCase) Execute(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	uc.log.InfoContext(ctx, "starting login")

	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		return nil, &user.InvalidCredentialsError{}
	}

	var result *LoginResult
	err = uc.uow.Execute(ctx, func(ctx context.Context, users user.Repository, _ shared.Publisher) error {
		found, err := users.FindByEmail(ctx, email)
		if err != nil {
			return fmt.Errorf("finding user: %w", err)
		}
		if !found.IsPresent || !uc.hasher.Verify(found.Value.PasswordHash, cmd.Password) {
			return &user.InvalidCredentialsError{}
		}

		token, expiresAt, err := uc.tokens.Issue(found.Value.Id, found.Value.Email)
		if err != nil {
			return fmt.Errorf("issuing token: %w", err)
		}
		result = &LoginResult{Token: token, ExpiresAt: expiresAt}
		uc.log.InfoContext(ctx, "login succeeded", "user_id", found.Value.Id.String())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
