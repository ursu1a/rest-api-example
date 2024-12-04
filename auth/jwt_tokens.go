package auth

/* JWT generation functions */

import (
	"backend/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token secrets for generating access tokens
var jwtSecret = []byte(utils.GetEnv("JWT_SECRET_KEY", "secret"))
var jwtRefreshSecret = []byte(utils.GetEnv("JWT_REFRESH_SECRET_KEY", "refresh"))

type UserClaims struct {
	UserID string
	Email  string
	jwt.RegisteredClaims
}

// Generate short-lived JWT token
func GenerateAccessToken(userID, email string) (string, error) {
	mins, _ := strconv.Atoi(utils.GetEnv("JWT_ACCESS_TOKEN_TIME", "15"))
	expirationTime := time.Now().Local().Add(time.Minute * time.Duration(mins))
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
	hours, _ := strconv.Atoi(utils.GetEnv("JWT_REFRESH_TOKEN_TIME", "24"))
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

// Generate token for authorization processes
func GenerateToken() string {
	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = jwtSecret
	t = jwt.New(jwt.SigningMethodHS256)
	s, _ = t.SignedString(key)
	return s
}
