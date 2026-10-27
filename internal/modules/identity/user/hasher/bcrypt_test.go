package hasher_test

import (
	"testing"

	"joao-api-comercio/internal/modules/identity/user/hasher"
)

func TestBcryptRoundTrip(t *testing.T) {
	h := hasher.NewBcrypt()

	hash, err := h.Hash("supersecret")
	if err != nil {
		t.Fatalf("hashing: %v", err)
	}
	if hash == "supersecret" {
		t.Fatal("hash equals the plain password")
	}
	if !h.Verify(hash, "supersecret") {
		t.Error("correct password did not verify")
	}
	if h.Verify(hash, "wrongpassword") {
		t.Error("wrong password verified")
	}
}
