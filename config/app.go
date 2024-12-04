package config

import (
	"backend/auth"
	"backend/db"
	"backend/emails"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"net/http"
)

type AppContext struct {
	Server      *http.Server
	DB          *gorm.DB
	Router      *mux.Router
	OAuthConfig *oauth2.Config
	EmailSvc    *emails.EmailService
}

func InitAppContext() AppContext {
	return AppContext{
		DB:          db.Connect(),
		OAuthConfig: auth.InitOAuthConfig(),
		EmailSvc:    emails.InitEmailService(),
	}
}

var App AppContext
