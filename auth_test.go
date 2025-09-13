package main

import (
	"math/rand"
	"os"
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

func init() {
	// set JWT_SECRET for tests
	os.Setenv("JWT_SECRET", "test-secret")
}

func TestCreateAccessToken(t *testing.T) {
	tokenString, err := CreateAccessToken(&user)
	if err != nil {
		t.Fatalf("CreateAccessToken returned an error: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Fatal("Generated token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	// userID check
	userID := int64(claims["userID"].(float64))
	if userID != user.ID {
		t.Errorf("UserID claim is incorrect. Expected: %d, Got: %d", user.ID, userID)
	}

	// exp check
	expUnix := int64(claims["exp"].(float64))
	expiresAt := time.Unix(expUnix, 0)

	if time.Now().Add(15*time.Minute).Before(expiresAt) == false {
		t.Error("ExpiresAt claim is not within 15 minutes")
	}
}

func TestValidateJWT(t *testing.T) {
	tokenString, err := CreateAccessToken(&user)
	if err != nil {
		t.Fatalf("CreateAccessToken returned an error: %v", err)
	}

	token, err := ValidateJWT(tokenString)
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	if !token.Valid {
		t.Fatal("Validated token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["userID"].(float64))
	if userID != user.ID {
		t.Errorf("UserID claim is incorrect. Expected: %d, Got: %d", user.ID, userID)
	}
}

func TestGenerateHash(t *testing.T) {
	password := "test_password"

	hashedPassword, err := generateHash(password)
	if err != nil {
		t.Fatalf("generateHash returned an error: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password is empty")
	}

	if hashedPassword == password {
		t.Error("Hashed password is the same as the original password")
	}
}
