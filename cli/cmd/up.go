/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Starts the stack and runs health checks",
	Long: `Starts all Docker containers defined in the stack (Nginx, PHP, Node, Postgres, Redis).
After starting the containers, it automatically triggers a health check sequence 
to ensure all services are ready to handle requests.

Usage:
  stackfold up`,

	Run: func(cmd *cobra.Command, args []string) {
		HandleSimpleScript(PathStackUp)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
