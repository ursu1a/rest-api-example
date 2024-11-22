package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "My CLI for interacting with the database",
	Long:  `CLI tool for importing and exporting data from the database.`,
}

func Execute() error {
	return RootCmd.Execute()
}
