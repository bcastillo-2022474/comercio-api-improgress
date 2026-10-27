// Package contracts holds the behavior every shared.UnitOfWork must have. The in-memory fake and
// the SQL adapter both run it, so the two cannot drift apart.
package contracts

import (
	"context"
	"errors"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared"
	"joao-api-comercio/internal/shared/outbox"
)

// Harness is one unit of work wired to a fresh store, plus how to observe what it saved.
type Harness struct {
	UoW    shared.UnitOfWork[user.Repository]
	Users  user.Repository
	Outbox outbox.Repository
}

var created = time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)

// RunUnitOfWork runs the contract. newHarness gets a fresh, empty store for each subtest.
func RunUnitOfWork(t *testing.T, newHarness func(t *testing.T) Harness) {
	t.Run("a committed write and its event are both saved", func(t *testing.T) {
		h := newHarness(t)
		u := user.NewUser("ana@example.com", "hash-1", created)

		err := h.UoW.Execute(context.Background(), func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
			if err := users.Save(ctx, u); err != nil {
				return err
			}
			pub.Publish(user.NewUserRegisteredEvent(u))
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := findUser(t, h.Users, "ana@example.com"); !ok {
			t.Error("the write was not saved")
		}
		if got := outboxCount(t, h.Outbox); got != 1 {
			t.Errorf("outbox records = %d, want 1", got)
		}
	})

	t.Run("a failing callback saves neither the write nor the event", func(t *testing.T) {
		h := newHarness(t)
		u := user.NewUser("ana@example.com", "hash-1", created)
		boom := errors.New("boom")

		err := h.UoW.Execute(context.Background(), func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
			if err := users.Save(ctx, u); err != nil {
				return err
			}
			pub.Publish(user.NewUserRegisteredEvent(u))
			return boom
		})
		if !errors.Is(err, boom) {
			t.Fatalf("err = %v, want the callback's error", err)
		}
		if _, ok := findUser(t, h.Users, "ana@example.com"); ok {
			t.Error("the write was saved despite the failure")
		}
		if got := outboxCount(t, h.Outbox); got != 0 {
			t.Errorf("outbox records = %d, want 0", got)
		}
	})

	t.Run("a rejected write saves no event and leaves the first user intact", func(t *testing.T) {
		h := newHarness(t)
		first := user.NewUser("ana@example.com", "hash-1", created)
		mustCommit(t, h, first)

		second := user.NewUser("ana@example.com", "hash-2", created)
		err := h.UoW.Execute(context.Background(), func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
			if err := users.Save(ctx, second); err != nil {
				return err
			}
			pub.Publish(user.NewUserRegisteredEvent(second))
			return nil
		})
		var taken *user.EmailTakenError
		if !errors.As(err, &taken) {
			t.Fatalf("err = %v, want *EmailTakenError", err)
		}
		if got := outboxCount(t, h.Outbox); got != 1 {
			t.Errorf("outbox records = %d, want 1 (only the first commit)", got)
		}
		kept, ok := findUser(t, h.Users, "ana@example.com")
		if !ok || kept.Id != first.Id || kept.PasswordHash != "hash-1" {
			t.Errorf("first user changed after a rejected write: %+v", kept)
		}
	})

	t.Run("a write without an event saves no outbox record", func(t *testing.T) {
		h := newHarness(t)
		u := user.NewUser("ana@example.com", "hash-1", created)

		err := h.UoW.Execute(context.Background(), func(ctx context.Context, users user.Repository, _ shared.Publisher) error {
			return users.Save(ctx, u)
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := findUser(t, h.Users, "ana@example.com"); !ok {
			t.Error("the write was not saved")
		}
		if got := outboxCount(t, h.Outbox); got != 0 {
			t.Errorf("outbox records = %d, want 0", got)
		}
	})
}

func mustCommit(t *testing.T, h Harness, u user.User) {
	t.Helper()
	err := h.UoW.Execute(context.Background(), func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
		if err := users.Save(ctx, u); err != nil {
			return err
		}
		pub.Publish(user.NewUserRegisteredEvent(u))
		return nil
	})
	if err != nil {
		t.Fatalf("committing: %v", err)
	}
}

func findUser(t *testing.T, users user.Repository, email user.Email) (user.User, bool) {
	t.Helper()
	got, err := users.FindByEmail(context.Background(), user.Email(email))
	if err != nil {
		t.Fatalf("finding user: %v", err)
	}
	return got.Value, got.IsPresent
}

func outboxCount(t *testing.T, records outbox.Repository) int {
	t.Helper()
	n, err := records.Count(context.Background())
	if err != nil {
		t.Fatalf("counting outbox records: %v", err)
	}
	return n
}
