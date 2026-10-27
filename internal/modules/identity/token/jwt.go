// Package token issues and verifies the JWT access tokens of the identity module.
package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared/types"
)

// MinSecretBytes keeps HS256 keys at the length of the hash, so they can't be brute-forced cheaply.
const MinSecretBytes = 32

type JWTIssuer struct {
	key []byte
	ttl time.Duration
	now func() time.Time
}

func NewJWTIssuer(secret string, ttl time.Duration, now func() time.Time) (*JWTIssuer, error) {
	if len(secret) < MinSecretBytes {
		return nil, fmt.Errorf("JWT secret must be at least %d bytes", MinSecretBytes)
	}
	return &JWTIssuer{key: []byte(secret), ttl: ttl, now: now}, nil
}

type claims struct {
	Email user.Email `json:"email"`
	jwt.RegisteredClaims
}

// Issue signs an access token for the user. The subject is the user id.
func (j *JWTIssuer) Issue(userId types.UserId, email user.Email) (string, time.Time, error) {
	issuedAt := j.now()
	expiresAt := issuedAt.Add(j.ttl)
	c := claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(j.key)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("signing token: %w", err)
	}
	return signed, expiresAt, nil
}

// Verify checks the signature and expiry, and returns the user the token was issued for.
func (j *JWTIssuer) Verify(tokenString string) (types.UserId, user.Email, error) {
	var c claims
	_, err := jwt.ParseWithClaims(tokenString, &c,
		func(*jwt.Token) (any, error) { return j.key, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(j.now),
	)
	if err != nil {
		return types.UserId{}, "", fmt.Errorf("verifying token: %w", err)
	}

	id, err := types.ParseUserId(c.Subject)
	if err != nil {
		return types.UserId{}, "", fmt.Errorf("parsing token subject: %w", err)
	}
	return id, c.Email, nil
}
