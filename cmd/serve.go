package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend/api"
	"backend/config"
	"backend/utils"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runServer(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	},
}

func runServer() error {
	port := utils.GetEnv("SERVER_PORT", "8080")

	// Init routes
	router := api.InitRouter()
	config.App.Router = router

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	config.App.Server = server

	// Graceful shutdown: channel to process signal of shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Server is starting on port %s...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe error: %v", err)
		}
	}()

	// Waiting for shutdown
	<-stop
	log.Println("Shutting down server...")

	// Context with timeout for stop active queries
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server stopped.")
	return nil
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
