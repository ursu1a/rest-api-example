package main

import (
	"backend/cmd"
	"backend/config"
	"backend/db"
	"log"
)

func init() {
	config.App = config.InitAppContext()
	db.Migrate(config.App.DB)
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
