package user_test

import (
	"errors"
	"strings"
	"testing"

	"joao-api-comercio/internal/modules/identity/user"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    user.Email
		wantErr bool
	}{
		{name: "normalizes case and spaces", raw: "  Ana@Example.COM ", want: "ana@example.com"},
		{name: "missing at sign", raw: "ana.example.com", wantErr: true},
		{name: "display name is not an address", raw: "Ana <ana@example.com>", wantErr: true},
		{name: "empty", raw: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewEmail(tt.raw)
			if tt.wantErr {
				var target *user.InvalidEmailError
				if !errors.As(err, &target) {
					t.Fatalf("err = %v, want *InvalidEmailError", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "minimum length", password: strings.Repeat("a", user.MinPasswordLength)},
		{name: "maximum bytes", password: strings.Repeat("a", user.MaxPasswordBytes)},
		{name: "too short", password: strings.Repeat("a", user.MinPasswordLength-1), wantErr: true},
		{name: "too long", password: strings.Repeat("a", user.MaxPasswordBytes+1), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ValidatePassword(tt.password)
			var target *user.WeakPasswordError
			if got := errors.As(err, &target); got != tt.wantErr {
				t.Fatalf("weak password error = %v, want %v (err: %v)", got, tt.wantErr, err)
			}
		})
	}
}
