// Package fakes holds in-memory implementations of the user ports for unit tests.
package fakes

import (
	"context"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared/types"
)

// Repository enforces the same unique email as the real store.
type Repository struct {
	users map[types.UserId]user.User
}

func NewRepository() *Repository {
	return &Repository{users: map[types.UserId]user.User{}}
}

func (r *Repository) Save(_ context.Context, u user.User) error {
	for id, existing := range r.users {
		if existing.Email == u.Email && id != u.Id {
			return &user.EmailTakenError{Email: u.Email}
		}
	}
	r.users[u.Id] = u
	return nil
}

func (r *Repository) FindByEmail(_ context.Context, email user.Email) (types.Maybe[user.User], error) {
	for _, u := range r.users {
		if u.Email == email {
			return types.Some(u), nil
		}
	}
	return types.None[user.User](), nil
}
