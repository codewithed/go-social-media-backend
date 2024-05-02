package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccessToken(user *User) (string, error) {
	claims := &jwt.MapClaims{
		"userID":    user.ID,
		"expiresAt": time.Now().Add(time.Minute * 15),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Check if the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func generateHash(pw string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func authoriseCurrentUser(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token string
		tokenString := r.Header.Get("x-jwt-token")

		// get the token and validate it
		token, err := ValidateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		// get user from database
		username := getUserName(r)
		user, err := s.GetUserByName(username)
		if err != nil {
			permissionDenied(w)
			return
		}

		// compare userID with that in the jwt token
		claims := token.Claims.(jwt.MapClaims)

		//check if token is expired	or not
		exp := claims["expiresAt"].(time.Time)
		if time.Now().After(exp) {
			WriteJson(w, http.StatusUnauthorized, "token is expired")
			return
		}

		if user.ID != int64(claims["userID"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	WriteJson(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
}

func resourceBasedJWTauth(handlerfunc http.HandlerFunc, s Storage, resourceType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")

		token, err := ValidateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		// validate ownership of resource
		resourceID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := int64(claims["userID"].(float64))

		//check if token is expired	or not
		exp := claims["expiresAt"].(time.Time)
		if time.Now().After(exp) {
			WriteJson(w, http.StatusUnauthorized, "token is expired, please log in again")
			return
		}

		ok, _ := validateOwnership(userID, resourceID, resourceType, s)
		if !ok {
			permissionDenied(w)
			return
		}

		handlerfunc(w, r)
	}
}

func verifyUser(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := ValidateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := int64(claims["userID"].(float64))

		//check if token is expired	or not
		exp := claims["expiresAt"].(time.Time)
		if time.Now().After(exp) {
			WriteJson(w, http.StatusUnauthorized, "token is expired, please refresh")
			return
		}
		if _, err := s.GetUserByID(userID); err != nil {
			permissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

func validateOwnership(userID, resourceID int64, resourceType string, s Storage) (bool, error) {
	if resourceType == "post" {
		post, err := s.GetPost(resourceID)
		if err != nil {
			return false, err
		}
		if post == nil {
			return false, fmt.Errorf("couldn't get post")
		}
		if post.UserID != int64(userID) {
			return false, nil
		}
		return true, nil
	}

	if resourceType == "comment" {
		comment, err := s.GetComment(resourceID)
		if err != nil {
			return false, err
		}
		if comment == nil {
			return false, fmt.Errorf("couldn't get comment")
		}
		if comment.UserID != int64(userID) {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("invalid resource type: %v", resourceType)
}

func (user *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pw)) == nil
}
