package service

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "password123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hashedPassword == password {
		t.Fatal("hashed password must not equal plain password")
	}

	if !CheckPassword(hashedPassword, password) {
		t.Fatal("CheckPassword() should return true")
	}

	if CheckPassword(hashedPassword, "wrong-password") {
		t.Fatal("CheckPassword() should return false")
	}
}
