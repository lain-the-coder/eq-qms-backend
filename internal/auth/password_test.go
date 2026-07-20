package auth_test

import (
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/lain-the-coder/ea-qms-backend/internal/auth"
)

func TestPasswordHashingAndVerification(t *testing.T) {
	password := "super-secure-password-123"

	// use default argon2id parameters for testing
	params := argon2id.DefaultParams

	// hash the same password twice -> should produce two different strings (Unique Salts)
	hash1, err := auth.HashPassword(password, params)
	if err != nil {
		t.Fatalf("failed to generate first hash: %v", err)
	}

	hash2, err := auth.HashPassword(password, params)
	if err != nil {
		t.Fatalf("failed to generate second hash: %v", err)
	}

	if hash1 == hash2 {
		t.Errorf("security vulnerability: back-to-back hashes produced identical strings. Salt is not random!")
	}

	// checkPasswordHash with the correct password -> true
	match, err := auth.CheckPasswordHash(password, hash1)
	if err != nil {
		t.Errorf("unexpected error matching correct password: %v", err)
	}
	if !match {
		t.Errorf("expected password to match its own hash, but it failed")
	}

	// checkPasswordHash with a wrong password -> false, err == nil
	wrongPassword := "wrong-password-456"
	match, err = auth.CheckPasswordHash(wrongPassword, hash1)
	if err != nil {
		t.Errorf("unexpected error matching wrong password: %v", err)
	}
	if match {
		t.Errorf("security failure: wrong password matched a valid hash")
	}
}
