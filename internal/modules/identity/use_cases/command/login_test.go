package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	userfakes "joao-api-comercio/internal/modules/identity/user/fakes"
	sharedfakes "joao-api-comercio/internal/shared/fakes"
	outboxfakes "joao-api-comercio/internal/shared/outbox/fakes"
	"joao-api-comercio/internal/shared/types"
)

type fakeTokens struct {
	issued int
}

func (f *fakeTokens) Issue(_ types.UserId, email user.Email) (string, time.Time, error) {
	f.issued++
	return "token-for-" + string(email), fixedNow.Add(time.Hour), nil
}

// seededLogin registers ana@example.com with the password "supersecret" and returns a login use case.
func seededLogin(t *testing.T) (*LoginUseCase, *fakeTokens) {
	t.Helper()
	repo := userfakes.NewRepository()
	register, _ := newRegister(repo)
	if _, err := register.Execute(context.Background(), RegisterCommand{Email: "ana@example.com", Password: "supersecret"}); err != nil {
		t.Fatalf("seeding user: %v", err)
	}
	tokens := &fakeTokens{}
	uow := sharedfakes.NewUnitOfWork[user.Repository](repo, outboxfakes.NewRepository(), func() time.Time { return fixedNow })
	return NewLoginUseCase(uow, userfakes.Hasher{}, tokens), tokens
}

func TestLogin_IssuesTokenForValidCredentials(t *testing.T) {
	login, tokens := seededLogin(t)

	res, err := login.Execute(context.Background(), LoginCommand{Email: " ANA@example.com", Password: "supersecret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Token != "token-for-ana@example.com" {
		t.Errorf("token = %q", res.Token)
	}
	if tokens.issued != 1 {
		t.Errorf("tokens issued = %d, want 1", tokens.issued)
	}
}

func TestLogin_FailuresAreIndistinguishableAndIssueNoToken(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
	}{
		{name: "unknown email", email: "nobody@example.com", password: "supersecret"},
		{name: "wrong password", email: "ana@example.com", password: "wrongpassword"},
		{name: "malformed email", email: "not-an-email", password: "supersecret"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			login, tokens := seededLogin(t)

			_, err := login.Execute(context.Background(), LoginCommand{Email: tt.email, Password: tt.password})
			var target *user.InvalidCredentialsError
			if !errors.As(err, &target) {
				t.Fatalf("err = %v, want *InvalidCredentialsError", err)
			}
			if tokens.issued != 0 {
				t.Errorf("token issued on failed login")
			}
		})
	}
}
