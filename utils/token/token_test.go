package token

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGenerateToken(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("TOKEN_HOUR_LIFESPAN", "1")
	os.Setenv("API_SECRET", "test-secret-key")

	userID := uint(123)
	token, err := GenerateToken(userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("Expected token to be generated, got empty string")
	}

	// Test with invalid environment
	os.Setenv("TOKEN_HOUR_LIFESPAN", "invalid")
	_, err = GenerateToken(userID)
	if err == nil {
		t.Fatal("Expected error with invalid TOKEN_HOUR_LIFESPAN, got nil")
	}
}

func TestTokenValidation(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("TOKEN_HOUR_LIFESPAN", "1")
	os.Setenv("API_SECRET", "test-secret-key")

	userID := uint(123)
	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Create a mock gin context with the token
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Test with Authorization header
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	err = TokenValid(c)
	if err != nil {
		t.Fatalf("Expected token to be valid, got error: %v", err)
	}

	// Test ExtractTokenId
	extractedID, err := ExtractTokenId(c)
	if err != nil {
		t.Fatalf("Failed to extract token ID: %v", err)
	}

	if extractedID != userID {
		t.Fatalf("Expected user ID %d, got %d", userID, extractedID)
	}
}

func TestTokenValidationWithInvalidToken(t *testing.T) {
	os.Setenv("API_SECRET", "test-secret-key")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Test with invalid token
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	c.Request = req

	err := TokenValid(c)
	if err == nil {
		t.Fatal("Expected error with invalid token, got nil")
	}

	// Test with no token
	req2, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req2

	err = TokenValid(c)
	if err == nil {
		t.Fatal("Expected error with no token, got nil")
	}
}