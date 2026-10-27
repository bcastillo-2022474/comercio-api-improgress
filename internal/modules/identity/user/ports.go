package user

import (
	"context"

	"joao-api-comercio/internal/shared/types"
)

type Repository interface {
	// Save inserts or updates u. It returns *EmailTakenError if another user owns the email.
	Save(ctx context.Context, u User) error
	// FindByEmail reports absence through Maybe rather than an error.
	FindByEmail(ctx context.Context, email Email) (types.Maybe[User], error)
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	// Verify reports whether plain produces hash.
	Verify(hash, plain string) bool
}
