package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	user = User{
		ID:         int64(rand.Intn(10000)),
		UserName:   "eddicus",
		Name:       "Ed Icus",
		Email:      "edicus@gmail.com",
		Bio:        "A ship in a harbour is safe",
		Created_at: time.Now().UTC(),
	}
)

func TestCreateAccessToken(t *testing.T) {
	// Call CreateAccessToken function
	tokenString, err := CreateAccessToken(&user)
	if err != nil {
		t.Errorf("CreateAccessToken returned an error: %v", err)
	}

	// Parse the token
	secret := []byte(getTestSecret())
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		t.Error("Generated token is not valid")
	}

	// Check if the userID claim is correct
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["userID"].(float64))
	if userID != user.ID {
		t.Errorf("UserID claim is incorrect. Expected: %d, Got: %d", user.ID, userID)
	}

	// Check if the exp claim is within ~15 minutes
	expUnix := int64(claims["exp"].(float64))
	expiresAt := time.Unix(expUnix, 0)
	expected := time.Now().Add(15 * time.Minute)

	// allow ±10s drift
	diff := expected.Sub(expiresAt)
	if diff > 10*time.Second || diff < -10*time.Second {
		t.Errorf("ExpiresAt claim mismatch. Expected around: %v, Got: %v", expected, expiresAt)
	}
}

func TestValidateJWT(t *testing.T) {
	// Create a token for the mock user
	tokenString, _ := CreateAccessToken(&user)

	// Call ValidateJWT function with the generated token
	token, err := ValidateJWT(tokenString)
	if err != nil {
		t.Errorf("ValidateJWT returned an error: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		t.Error("Validated token is not valid")
	}

	// Check if the token's userID claim is correct
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["userID"].(float64))
	if userID != user.ID {
		t.Errorf("UserID claim is incorrect. Expected: %d, Got: %d", user.ID, userID)
	}
}

func TestGenerateHash(t *testing.T) {
	// Define a password for testing
	password := "test_password"

	// Call generateHash function
	hashedPassword, err := generateHash(password)
	if err != nil {
		t.Errorf("generateHash returned an error: %v", err)
	}

	// Check if the hashed password is not empty
	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}

	// Check if the hashed password is different from the original password
	if hashedPassword == password {
		t.Error("Hashed password is the same as the original password")
	}
}

// helper to provide a consistent JWT secret for tests
func getTestSecret() string {
	secret := "test_secret"
	return secret
}
