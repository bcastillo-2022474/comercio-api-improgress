// Package hasher implements user.PasswordHasher with bcrypt.
package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func NewBcrypt() Bcrypt {
	return Bcrypt{}
}

func (Bcrypt) Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(hash), nil
}

func (Bcrypt) Verify(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
