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
	"time"
	"gorm.io/gorm/clause"
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

type PasswordRecovery struct {
	Email string `json:"email"`
}

type UpdatePassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type UserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
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
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "An error when decoding request", http.StatusInternalServerError)
		return
	}
	log.Printf("You are entered to system by Google accout: %v\n", userInfo)

	// Save data into database
	userID, err := SaveUpdateGoogleUser(userInfo)
	if err != nil {
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
		AccessToken: accessToken,
	}

	SaveRefreshToken(w, refreshToken)

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
	emailSvc := config.App.EmailSvc

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

	token := auth.GenerateToken()
	user := db.User{
		Name:                   userData.Name,
		Email:                  userData.Email,
		PasswordHash:           string(hashedPassword),
		EmailVerificationToken: token,
	}

	// Save user in DB
	if err := DBConn.Create(&user).Error; err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Send email for registration confirmation
	err = emailSvc.SendRegistrationConfirmation(user.Email, token)
	if err != nil {
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
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

	// Check account is verified
	if !user.EmailVerified {
		http.Error(w, "Please complete registration process by email verification", http.StatusUnauthorized)
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

	SaveRefreshToken(w, refreshToken)

	// Response data
	authResponse := AuthResponse{
		AccessToken: accessToken,
	}

	// Send tokens in response
	responseJSON(w, authResponse)
}

// Refresh token
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Refresh token missing", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value
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

func HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	DBConn := config.App.DB
	EmailSvc := config.App.EmailSvc
	var user db.User

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Find user in DB
	if err := DBConn.Where("email_verification_token = ?", token).First(&user).Error; err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Make user as verified
	user.EmailVerified = true
	user.EmailVerificationToken = ""
	DBConn.Save(&user)

	EmailSvc.SendTransactionalEmail(user.Email, "Registration is complete", "Email verified successfully")
	fmt.Println(w, "Email verified successfully")
}

func HandleRequestResetPassword(w http.ResponseWriter, r *http.Request) {
	DBConn := config.App.DB
	EmailSvc := config.App.EmailSvc
	var user db.User
	var requestedUser PasswordRecovery

	if err := json.NewDecoder(r.Body).Decode(&requestedUser); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Find user in DB
	if err := DBConn.Where("email = ?", requestedUser.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Generate token
	token := auth.GenerateToken()
	expiry := time.Now().Add(1 * time.Hour)

	// Save token in database
	user.ResetToken = token
	user.ResetTokenExpiry = expiry
	DBConn.Save(user)

	EmailSvc.SendPasswordReset(user.Email, token)
	fmt.Println(w, "Password reset email sent")
}

func HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	DBConn := config.App.DB
	var user db.User
	var requestedUser UpdatePassword

	if err := json.NewDecoder(r.Body).Decode(&requestedUser); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check token
	err := DBConn.Where("reset_token = ?", requestedUser.Token).Find(&user).Error
	if err != nil || time.Now().After(user.ResetTokenExpiry) {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// Get password hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	// Update password and remove token
	user.PasswordHash = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}
	DBConn.Save(&user)

	fmt.Println(w, "New password successfully set")
}

func SaveRefreshToken(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func SaveUpdateGoogleUser(userInfo UserInfo) (uint, error) {
	DBConn := config.App.DB
	user := db.User{
		GoogleID:      &userInfo.ID,
		Email:         userInfo.Email,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
		EmailVerified: true,
		UpdatedAt:     time.Now(),
	}

	// Use create with conflict method to create or update the user
	result := DBConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"google_id", "name", "picture", "updated_at"}),
	}).Create(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	log.Printf("User with email %s was created", user.Email)
	return user.ID, nil
}