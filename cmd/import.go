package cmd

import (
	"backend/config"
	"backend/db"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

/**
* Command for creating a database from a JSON dump
* @param filename string
*
* Example: backend import fixtures/<file_name>.json
 */
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import data from a JSON file into the database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		DBConn := config.App.DB

		// Read the file
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}

		var users []db.User
		if err := json.Unmarshal(data, &users); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Insert data to database
		if err := DBConn.Create(&users).Error; err != nil {
			log.Fatalf("Failed to insert data into the database: %v", err)
		}

		fmt.Println("Data imported successfully from", filePath)
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}
