package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/api"
	"backend/config"

	"github.com/spf13/cobra"
)

// Start server command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func runServer() {
	var port = os.Getenv("SERVER_PORT")
	var server *http.Server = &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	// Init API routes
	config.App.Router = api.InitRouter()

	// Run Server
	log.Printf("Server is starting on port %s...", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

	config.App.Server = server
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
