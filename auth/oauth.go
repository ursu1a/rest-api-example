package auth

import (
	"backend/utils"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitOAuthConfig() *oauth2.Config {
	if err := utils.CheckEnvs([]string{"SERVER_ADDRESS", "SERVER_PORT", "GOOGLE_OAUTH_CLIENT_ID", "GOOGLE_OAUTH_CLIENT_SECRET", "GOOGLE_OAUTH_CALLBACK_PATH"}); err != nil {
		log.Fatalf("Error checking environment variables: %v", err)
		return nil
	}

	var (
		baseUrl      = fmt.Sprintf("%s:%s", utils.GetEnv("SERVER_ADDRESS", "http://localhost"), utils.GetEnv("SERVER_PORT", "8080"))
		redirectUrl  = fmt.Sprintf("%s%s", baseUrl, utils.GetEnv("GOOGLE_OAUTH_CALLBACK_PATH", "/api/auth/google/callback"))
		clientID     = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
		clientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	)

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}
