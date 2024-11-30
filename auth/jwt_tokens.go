package auth

/* JWT generation functions */

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token secrets for generating access tokens
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
var jwtRefreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET_KEY"))

type UserClaims struct {
	UserID string
	Email  string
	jwt.RegisteredClaims
}

// Generate short-lived JWT token
func GenerateAccessToken(userID, email string) (string, error) {
	mins, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_TIME"))	
	log.Printf("Minutes: %v", mins)
	expirationTime := time.Now().Local().Add(time.Minute * time.Duration(mins))
	log.Printf("Expiration time: %v", expirationTime)
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Generate refresh token (long lived JWT token)
func GenerateRefreshToken(userID, email string) (string, error) {
	hours, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_TIME"))	
	expirationTime := time.Now().Local().Add(time.Hour * time.Duration(hours))
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtRefreshSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Refresh token
func RefreshToken(refreshToken string) (string, error) {
	// Decode refresh token
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtRefreshSecret, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	// Generate new access token
	newAccessToken, err := GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func GenerateToken() string {
	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte(os.Getenv("JWT_SECRET_KEY"))
	t = jwt.New(jwt.SigningMethodHS256)
	s, _ = t.SignedString(key)
	return s
}
