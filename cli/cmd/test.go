/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Performs a health check of the stack services",
	Long: `Runs a diagnostic test to ensure the environment is fully operational.

The test includes:
- Docker container status and health verification (nginx, php, node, db, redis)
- Nginx HTTP response check
- PHP Stack internal health checks
- Vite Dev Server availability
- Host port accessibility check`,

	Run: func(cmd *cobra.Command, args []string) {
		HandleSimpleScript(PathStackTest)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
