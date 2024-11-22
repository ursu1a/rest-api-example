package middleware

import (
	"backend/auth"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = os.Getenv("JWT_SECRET_KEY")

// Checks if the user has a valid JWT token
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from headers
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &auth.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Receive user data from token and save it to request header
		if claims, ok := token.Claims.(*auth.UserClaims); ok && token.Valid {
			r.Header.Set("UserID", claims.UserID)
			r.Header.Set("Email", claims.Email)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid token claims", http.StatusForbidden)
		}
	})
}
