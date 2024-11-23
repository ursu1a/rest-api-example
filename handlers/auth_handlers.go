package handlers

/* Authentication handlers */

import (
	"backend/auth"
	"backend/config"
	"backend/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.App.OAuthConfig
	// Redirect to Google authentication page
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthConfig := config.App.OAuthConfig
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Receiving a token is failed", http.StatusBadRequest)
		return
	}

	// Use token for User's information receiving
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "User data receiving is failed", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	// Proceed for User's information
	var userInfo auth.UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "An error when decoding request", http.StatusInternalServerError)
		return
	}
	log.Printf("You are entered to system by Google accout: %v\n", userInfo)

	// Save data into database
	userID, err := auth.SaveUpdateUser(userInfo); if err != nil {
		http.Error(w, "An error when saving user", http.StatusInternalServerError)
		return
	}

	// Generate Access & Refresh tokens
	strUserID := fmt.Sprintf("%d", userID)	
	accessToken, err := auth.GenerateAccessToken(strUserID, userInfo.Email)
	if err != nil {
		http.Error(w, "Access token generation failed", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(strUserID, userInfo.Email)
	if err != nil {
		http.Error(w, "Refresh token generation failed", http.StatusInternalServerError)
		return
	}

	// Response data
	authResponse := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Insert JSON-response into client script
	responseScript := `
		<script>
			window.opener.postMessage(%s, '*');
			window.close();
		</script>
	`

	responseJSON, _ := json.Marshal(authResponse)
	responseText(w, fmt.Sprintf(responseScript, responseJSON))
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	DBConn := config.App.DB
	var userData NewUser
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Register user's data: %v", userData)

	// Get password hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	user := db.User{
		Name:         userData.Name,
		Email:        userData.Email,
		PasswordHash: string(hashedPassword),
	}

	// Save user in DB
	if err := DBConn.Create(&user).Error; err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("User with email: %v registered successfully", user.Email)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	DBConn := config.App.DB
	var user db.User
	var credentials Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Login credentials: %v", credentials)

	// Find user in DB
	if err := DBConn.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate access & refresh tokens
	strUserID := fmt.Sprintf("%d", user.ID)
	accessToken, err := auth.GenerateAccessToken(strUserID, user.Email)
	if err != nil {
		http.Error(w, "Could not generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(strUserID, user.Email)
	if err != nil {
		http.Error(w, "Could not generate refresh token", http.StatusInternalServerError)
		return
	}

	// Response data
	authResponse := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Send tokens in response
	responseJSON(w, authResponse)
}

// Refresh token
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	refreshToken := body.RefreshToken
	if refreshToken == "" {
		http.Error(w, "Refresh token is missing", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := auth.RefreshToken(refreshToken)

	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid refresh token. Reason: %v", err), http.StatusUnauthorized)
		return
	}

	responseJSON(w, map[string]interface{}{"access_token": newAccessToken})
}
