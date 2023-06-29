package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func CreateJWT(user *User) (string, error) {
	claims := &jwt.MapClaims{
		"userID":    user.ID,
		"expiresAt": 15000,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Secret is a []byte containing your secret, e.g []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

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
		user, err := s.GetUser(username)
		if err != nil {
			permissionDenied(w)
			return
		}

		// compare userID with that in the jwt token
		claims := token.Claims.(jwt.MapClaims)
		if user.ID != int64(claims["id"].(int)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	WriteJson(w, http.StatusForbidden, &ApiError{Error: "permission denied"})
	return
}
