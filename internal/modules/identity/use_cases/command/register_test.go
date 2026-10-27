package command

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	userfakes "joao-api-comercio/internal/modules/identity/user/fakes"
	sharedfakes "joao-api-comercio/internal/shared/fakes"
	outboxfakes "joao-api-comercio/internal/shared/outbox/fakes"
)

var fixedNow = time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

// newRegister wires the use case over repo and returns the outbox the unit of work writes to, so tests can inspect events.
func newRegister(repo *userfakes.Repository) (*RegisterUseCase, *outboxfakes.Repository) {
	outbox := outboxfakes.NewRepository()
	uow := sharedfakes.NewUnitOfWork[user.Repository](repo, outbox, func() time.Time { return fixedNow })
	return NewRegisterUseCase(uow, userfakes.Hasher{}, func() time.Time { return fixedNow }, discard), outbox
}

func TestRegister_StoresHashedUserAndPublishesEvent(t *testing.T) {
	repo := userfakes.NewRepository()
	register, outbox := newRegister(repo)

	u, err := register.Execute(context.Background(), RegisterCommand{
		Email:    "Ana@Example.com",
		Password: "supersecret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "ana@example.com" {
		t.Errorf("email = %q, want normalized form", u.Email)
	}
	if u.PasswordHash == "supersecret" {
		t.Error("plain password was stored")
	}
	if !u.CreatedAt.Equal(fixedNow) {
		t.Errorf("created at = %v, want %v", u.CreatedAt, fixedNow)
	}

	stored, err := repo.FindByEmail(context.Background(), "ana@example.com")
	if err != nil || !stored.IsPresent {
		t.Fatalf("user not persisted: present=%v err=%v", stored.IsPresent, err)
	}
	if stored.Value.Id != u.Id {
		t.Errorf("stored ID %v, want %v", stored.Value.Id, u.Id)
	}

	records := outbox.Records()
	if len(records) != 1 || records[0].EventName != "user.registered" {
		t.Fatalf("outbox = %v, want one user.registered record", records)
	}
	if want := `{"user_id":"` + u.Id.String() + `"}`; string(records[0].Payload) != want {
		t.Errorf("payload = %s, want %s", records[0].Payload, want)
	}
}

func TestRegister_Rejections(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		seed     bool
		target   any
	}{
		{name: "invalid email", email: "not-an-email", password: "supersecret", target: new(*user.InvalidEmailError)},
		{name: "weak password", email: "ana@example.com", password: strings.Repeat("a", user.MinPasswordLength-1), target: new(*user.WeakPasswordError)},
		{name: "email already taken", email: "ana@example.com", password: "supersecret", seed: true, target: new(*user.EmailTakenError)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := userfakes.NewRepository()
			if tt.seed {
				seed, _ := newRegister(repo)
				if _, err := seed.Execute(context.Background(), RegisterCommand{Email: tt.email, Password: "firstpassword"}); err != nil {
					t.Fatalf("seeding user: %v", err)
				}
			}

			register, outbox := newRegister(repo)
			before := len(outbox.Records())
			_, err := register.Execute(context.Background(), RegisterCommand{Email: tt.email, Password: tt.password})
			if !errors.As(err, tt.target) {
				t.Fatalf("err = %v, want %T", err, tt.target)
			}
			if len(outbox.Records()) != before {
				t.Error("a rejected registration wrote an outbox record")
			}
		})
	}
}
