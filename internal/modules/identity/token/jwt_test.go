package token_test

import (
	"strings"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/token"
	"joao-api-comercio/internal/shared/types"
)

var start = time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)

const secret = "0123456789abcdef0123456789abcdef"

func TestJWTIssuer_RoundTripAndExpiry(t *testing.T) {
	now := start
	clock := func() time.Time { return now }
	issuer, err := token.NewJWTIssuer(secret, 15*time.Minute, clock)
	if err != nil {
		t.Fatalf("creating issuer: %v", err)
	}
	userId := types.NewUserId()

	signed, expiresAt, err := issuer.Issue(userId, "ana@example.com")
	if err != nil {
		t.Fatalf("issuing: %v", err)
	}
	if !expiresAt.Equal(start.Add(15 * time.Minute)) {
		t.Errorf("expires at = %v, want %v", expiresAt, start.Add(15*time.Minute))
	}

	gotId, gotEmail, err := issuer.Verify(signed)
	if err != nil {
		t.Fatalf("verifying fresh token: %v", err)
	}
	if gotId != userId || gotEmail != "ana@example.com" {
		t.Errorf("verified (%v, %q), want (%v, ana@example.com)", gotId, gotEmail, userId)
	}

	now = start.Add(16 * time.Minute)
	if _, _, err := issuer.Verify(signed); err == nil {
		t.Error("expired token verified")
	}
}

func TestJWTIssuer_RejectsTokenSignedWithOtherSecret(t *testing.T) {
	clock := func() time.Time { return start }
	issuer, _ := token.NewJWTIssuer(secret, time.Minute, clock)
	other, _ := token.NewJWTIssuer(strings.Repeat("x", 32), time.Minute, clock)

	signed, _, err := other.Issue(types.NewUserId(), "ana@example.com")
	if err != nil {
		t.Fatalf("issuing: %v", err)
	}
	if _, _, err := issuer.Verify(signed); err == nil {
		t.Error("token signed with another secret verified")
	}
}

func TestNewJWTIssuer_RejectsShortSecret(t *testing.T) {
	if _, err := token.NewJWTIssuer("short", time.Minute, time.Now); err == nil {
		t.Error("short secret accepted")
	}
}
