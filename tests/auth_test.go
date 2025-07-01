package tests

import (
	"testing"

	"server/auth"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == password {
		t.Error("Password hash should not equal original password")
	}

	if len(hash) == 0 {
		t.Error("Password hash should not be empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	if !auth.CheckPasswordHash(password, hash) {
		t.Error("Password check should return true for correct password")
	}

	// Test incorrect password
	if auth.CheckPasswordHash("wrongpassword", hash) {
		t.Error("Password check should return false for incorrect password")
	}
}

func TestGenerateToken(t *testing.T) {
	userID := int64(123)
	secret := "test-secret"

	token, err := auth.GenerateToken(userID, secret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}

	// Token should have 3 parts separated by dots
	parts := len(token)
	if parts < 10 {
		t.Error("Generated token should have reasonable length")
	}
}

func TestValidateToken(t *testing.T) {
	userID := int64(123)
	secret := "test-secret"

	token, err := auth.GenerateToken(userID, secret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Test valid token
	validatedUserID, err := auth.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, validatedUserID)
	}
}

func TestValidateTokenInvalid(t *testing.T) {
	secret := "test-secret"

	// Test invalid token format
	_, err := auth.ValidateToken("invalid.token.format", secret)
	if err == nil {
		t.Error("Should return error for invalid token format")
	}

	// Test wrong secret
	userID := int64(123)
	token, _ := auth.GenerateToken(userID, secret)

	_, err = auth.ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Error("Should return error for wrong secret")
	}
}
