package main

import (
	"backend/auth"
	"backend/cmd"
	"backend/config"
	"backend/db"
	"backend/emails"
	"log"
)


func init() {
	app := &config.App
	app.DB = db.Connect()
	db.Migrate(app.DB)
	app.OAuthConfig = auth.InitOAuthConfig()
	app.EmailSvc = emails.InitEmailService()
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
