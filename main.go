package main

import (
	"backend/auth"
	"backend/cmd"
	"backend/config"
	"backend/db"
	"log"
)


func init() {
	app := &config.App
	app.DB = db.Connect()
	db.Migrate(app.DB)
	app.OAuthConfig = auth.InitOAuthConfig()
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
