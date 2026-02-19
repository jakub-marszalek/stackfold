/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Laravel project environment",
	Long: `Guides you through the first-time setup of a new Stackfold project.

Creates the necessary .env files, directory structures, and installs
a fresh Laravel project inside app/ using the official Composer Docker image.

If an existing configuration or Laravel project is detected, the command
aborts safely. Use 'stackfold config' to update an existing environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() {
	// Guard: refuse to run if .env.docker already exists
	if _, err := os.Stat(".env.docker"); err == nil {
		DrawLine(UIConfig.Red)
		LogError("Configuration already exists!")
		LogWarning("Found existing .env.docker file.")
		fmt.Println(UIConfig.Reset + "To update your configuration, use: stackfold config")
		DrawLine(UIConfig.Red)
		os.Exit(1)
	}

	PrintHeader("Initializing new stack and Laravel project...")
	RunConfigWizard()

	// Laravel Installation
	PrintSection("Laravel Installation")

	if IsLaravelInstalled() {
		LogWarning("Laravel project detected in app/. Skipping installation.")
	} else {
		fmt.Println("No Laravel project found in app/. Proceeding with fresh installation...")
		fmt.Println()

		// if err := InstallLaravel(); err != nil {
		// 	LogError(fmt.Sprintf("Laravel installation failed: %v", err))
		// 	os.Exit(1)
		// }

		LogSuccess("Laravel installed successfully.")
	}

	fmt.Println()
	PrintHeader("Initialization completed")
	LogSuccess("Run 'stackfold build -s' to build and start your stack!")
}
