package command

import (
	"context"
	"fmt"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared"
)

type RegisterCommand struct {
	Email    string
	Password string
}

type RegisterUseCase struct {
	uow    shared.UnitOfWork[user.Repository]
	hasher user.PasswordHasher
	now    func() time.Time
}

func NewRegisterUseCase(uow shared.UnitOfWork[user.Repository], hasher user.PasswordHasher, now func() time.Time) *RegisterUseCase {
	return &RegisterUseCase{uow: uow, hasher: hasher, now: now}
}

// Execute creates a user and publishes user.registered once the write commits.
// It returns the stored user, which has no plain password.
func (uc *RegisterUseCase) Execute(ctx context.Context, cmd RegisterCommand) (user.User, error) {
	email, err := user.NewEmail(cmd.Email)
	if err != nil {
		return user.User{}, err
	}
	if err := user.ValidatePassword(cmd.Password); err != nil {
		return user.User{}, err
	}
	hash, err := uc.hasher.Hash(cmd.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("hashing password: %w", err)
	}

	u := user.NewUser(email, hash, uc.now())
	err = uc.uow.Execute(ctx, func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
		if err := users.Save(ctx, u); err != nil {
			return err
		}
		pub.Publish(user.NewUserRegisteredEvent(u))
		return nil
	})
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}
