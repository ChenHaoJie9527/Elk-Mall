package pkg

import (
	"testing"
	"time"
)

func TestGenerateAndParseJWTToken(t *testing.T) {
	secret := []byte("test-secret")
	expiresIn := time.Hour * 24
	userID := uint(1)
	token, err := GenerateJWTToken(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("GenerateJWTToken failed: %v", err)
	}
	t.Logf("token: %s", token)
	parsedUserID, err := ParseJWTToken(token, secret)
	if err != nil {
		t.Fatalf("ParseJWTToken failed: %v", err)
	}
	t.Logf("parsedUserID: %d", parsedUserID)
	if parsedUserID != userID {
		t.Fatalf("parsedUserID != userID: %d != %d", parsedUserID, userID)
	}
}
