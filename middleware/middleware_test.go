package middleware

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
)

func TestGenerateAndValidateToken(t *testing.T) {
	os.Setenv("JWT_KEY", "test-secret")

	token, err := GenerateToken("admin-123")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	parsed, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		t.Fatal("invalid token/claims")
	}
	if claims["id"] != "admin-123" {
		t.Fatalf("id mismatch: got %v", claims["id"])
	}
	if _, ok := claims["exp"]; !ok {
		t.Fatal("missing exp claim")
	}
}

func TestValidateTokenRejectsGarbage(t *testing.T) {
	os.Setenv("JWT_KEY", "test-secret")
	if _, err := ValidateToken("not-a-valid-token"); err == nil {
		t.Fatal("expected error for invalid token")
	}
}
