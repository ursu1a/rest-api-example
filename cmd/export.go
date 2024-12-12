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
*
* Command to create a database dump
* @param filename string
*
* Example: backend export /tmp/<file_name>.json
 */
var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export data from the database to a JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		DBConn := config.App.DB

		var users []db.User
		result := DBConn.Find(&users)
		if result.Error != nil {
			log.Fatalf("Failed to fetch data: %v", result.Error)
		}

		// Serialize data to JSON
		jsonData, err := json.MarshalIndent(users, "", "    ")
		if err != nil {
			log.Fatalf("Failed to marshal data to JSON: %v", err)
		}

		// Write to file
		err = os.WriteFile(filePath, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write JSON data to file: %v", err)
		}

		fmt.Println("Data exported successfully to", filePath)
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
