// Package user is the identity context's user aggregate: what a user is and
// what makes one valid. It knows nothing about storage, hashing or tokens.
package user

import (
	"net/mail"
	"strings"
	"time"

	"joao-api-comercio/internal/shared/types"
)

const (
	MinPasswordLength = 3
	// bcrypt silently truncates input past 72 bytes, so longer passwords would
	// be accepted but only partly checked. Reject them instead.
	MaxPasswordBytes = 72
)

// Email is a normalized (trimmed, lowercase) and syntactically valid address.
type Email string

// NewEmail normalizes raw and checks that it is a plain address.
func NewEmail(raw string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	addr, err := mail.ParseAddress(normalized)
	if err != nil || addr.Address != normalized {
		return "", &InvalidEmailError{Email: raw}
	}
	return Email(normalized), nil
}

func ValidatePassword(plain string) error {
	if len(plain) < MinPasswordLength || len(plain) > MaxPasswordBytes {
		return &WeakPasswordError{MinLength: MinPasswordLength, MaxBytes: MaxPasswordBytes}
	}
	return nil
}

type User struct {
	Id           types.UserId
	Email        Email
	PasswordHash string
	CreatedAt    time.Time
}

// NewUser builds a new user with a fresh ID. The caller hashes the password first.
func NewUser(email Email, passwordHash string, now time.Time) User {
	return User{
		Id:           types.NewUserId(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
	}
}
