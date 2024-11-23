package auth

import (
	"backend/config"
	"backend/db"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"time"
)

type UserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func InitOAuthConfig() *oauth2.Config {
	var baseUrl = fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"))
	var redirectUrl = fmt.Sprintf("%s%s", baseUrl, os.Getenv("GOOGLE_OAUTH_CALLBACK_PATH"))

	config := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		RedirectURL:  redirectUrl,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return &config
}

func SaveUpdateUser(userInfo UserInfo) (uint, error) {
	DBConn := config.App.DB
	user := db.User{
		GoogleID:  userInfo.ID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		Picture:   userInfo.Picture,
		UpdatedAt: time.Now(),
	}

	// Use create with conflict method to create or update the user
	result := DBConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"google_id", "name", "picture", "updated_at"}),
	}).Create(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	// Generate/update access token
	strUserID := fmt.Sprintf("%d", user.ID)
	refreshToken, err := GenerateRefreshToken(strUserID, userInfo.Email)
	if err != nil {
		return user.ID, err
	}

	user.RefreshToken = refreshToken
	DBConn.Save(&user)

	log.Printf("User with email %s was created", user.Email)
	return user.ID, nil
}
