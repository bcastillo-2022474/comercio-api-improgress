// Package contracts holds the behavior every user.Repository must have. The in-memory fake
// and the SQL adapter both run it, so the two cannot drift apart.
package contracts

import (
	"context"
	"errors"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
)

var created = time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)

// RunRepository runs the contract. newRepo gets a fresh, empty repository for each subtest.
func RunRepository(t *testing.T, newRepo func(t *testing.T) user.Repository) {
	t.Run("a saved user is found by email", func(t *testing.T) {
		repo := newRepo(t)
		u := user.NewUser("ana@example.com", "hash-1", created)
		mustSave(t, repo, u)

		got := mustFind(t, repo, u.Email)
		if got.Id != u.Id || got.Email != u.Email || got.PasswordHash != u.PasswordHash || !got.CreatedAt.Equal(u.CreatedAt) {
			t.Fatalf("found %+v, want %+v", got, u)
		}
	})

	t.Run("an unknown email is absent, not an error", func(t *testing.T) {
		repo := newRepo(t)

		got, err := repo.FindByEmail(context.Background(), "nobody@example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.IsPresent {
			t.Fatalf("found %+v for an unknown email", got.Value)
		}
	})

	t.Run("a second user cannot take an existing email", func(t *testing.T) {
		repo := newRepo(t)
		first := user.NewUser("ana@example.com", "hash-1", created)
		mustSave(t, repo, first)

		err := repo.Save(context.Background(), user.NewUser("ana@example.com", "hash-2", created))
		var taken *user.EmailTakenError
		if !errors.As(err, &taken) {
			t.Fatalf("err = %v, want *EmailTakenError", err)
		}

		got := mustFind(t, repo, "ana@example.com")
		if got.Id != first.Id || got.PasswordHash != "hash-1" {
			t.Fatalf("first user changed after a rejected save: %+v", got)
		}
	})

	t.Run("saving a known user updates it", func(t *testing.T) {
		repo := newRepo(t)
		u := user.NewUser("ana@example.com", "hash-1", created)
		mustSave(t, repo, u)

		u.PasswordHash = "hash-2"
		mustSave(t, repo, u)

		got := mustFind(t, repo, u.Email)
		if got.Id != u.Id || got.PasswordHash != "hash-2" {
			t.Fatalf("after update found %+v, want the new hash on the same id", got)
		}
	})
}

func mustSave(t *testing.T, repo user.Repository, u user.User) {
	t.Helper()
	if err := repo.Save(context.Background(), u); err != nil {
		t.Fatalf("saving user: %v", err)
	}
}

func mustFind(t *testing.T, repo user.Repository, email user.Email) user.User {
	t.Helper()
	got, err := repo.FindByEmail(context.Background(), email)
	if err != nil {
		t.Fatalf("finding user: %v", err)
	}
	if !got.IsPresent {
		t.Fatalf("%s not found", email)
	}
	return got.Value
}
