package service

import (
	"testing"
	"time"
)

func TestTokenServiceGenerateAndParse(t *testing.T) {
	tokenService := NewTokenService(
		"test-secret-that-is-long-enough",
		time.Hour,
	)

	userID := "67c123456789012345678901"

	tokenString, err := tokenService.Generate(userID)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims, err := tokenService.Parse(tokenString)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if claims.Subject != userID {
		t.Errorf(
			"Subject = %q, want %q",
			claims.Subject,
			userID,
		)
	}
}

func TestTokenServiceRejectsInvalidToken(t *testing.T) {
	tokenService := NewTokenService(
		"test-secret-that-is-long-enough",
		time.Hour,
	)

	if _, err := tokenService.Parse("invalid-token"); err == nil {
		t.Fatal("Parse() expected error")
	}
}
