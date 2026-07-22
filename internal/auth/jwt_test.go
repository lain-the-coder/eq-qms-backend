package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lain-the-coder/ea-qms-backend/internal/auth"
)

func TestJWTRoundTrip(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-do-not-use-in-prod"

	token, err := auth.MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	gotID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate a freshly created token: %v", err)
	}

	if gotID != userID {
		t.Errorf("user ID did not survive the round trip: got %v, want %v", gotID, userID)
	}
}

func TestJWTRejectsWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, "the-real-secret", time.Hour)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	gotID, err := auth.ValidateJWT(token, "an-attackers-guess")
	if err == nil {
		t.Error("security failure: a token signed with a different secret was accepted")
	}
	if gotID != uuid.Nil {
		t.Errorf("expected uuid.Nil on failure, got %v", gotID)
	}
}

func TestJWTRejectsExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-do-not-use-in-prod"

	// negative duration -> the token is born already expired, no sleeping needed
	token, err := auth.MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	if _, err := auth.ValidateJWT(token, secret); err == nil {
		t.Error("security failure: an expired token was accepted")
	}
}

func TestJWTRejectsMalformedToken(t *testing.T) {
	secret := "test-secret-do-not-use-in-prod"

	for _, garbage := range []string{
		"",
		"not.a.token",
		"eyJhbGciOiJIUzI1NiJ9.garbage.signature",
	} {
		if _, err := auth.ValidateJWT(garbage, secret); err == nil {
			t.Errorf("malformed token %q was accepted", garbage)
		}
	}
}
