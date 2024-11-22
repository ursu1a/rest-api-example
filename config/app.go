package config

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AppContext struct {
	Server *http.Server
	DB *gorm.DB
	Router *mux.Router
	OAuthConfig *oauth2.Config
}

var App AppContext