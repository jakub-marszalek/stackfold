package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Update the environment configuration of an existing stack",
	Long: `Updates the project name and database credentials in .env.docker and app/.env.

This command is intended for updating an EXISTING Stackfold environment.
It will:
  1. Stop any running containers to avoid naming/port conflicts.
  2. Ask for new project name and database credentials.
  3. Overwrite .env.docker and app/.env with the new values.

IMPORTANT: It does NOT touch the app/ directory or your Laravel codebase.

WARNING: Changing database credentials does NOT migrate existing data.
If you change the database name or password, you must also reset your
Docker volumes via 'stackfold reset --force' to apply the new credentials.
Your database will be lost after a reset.`,
	Run: func(cmd *cobra.Command, args []string) {
		runConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig() {
	// Guard: require .env.docker to already exist
	if _, err := os.Stat(".env.docker"); err != nil {
		DrawLine(UIConfig.Red)
		LogError("No configuration found!")
		LogWarning(".env.docker does not exist.")
		fmt.Println(UIConfig.Reset + "To create a new project, use: stackfold init")
		DrawLine(UIConfig.Red)
		os.Exit(1)
	}

	DrawLine(UIConfig.Yellow)
	LogWarning("You are about to update your stack configuration.")
	fmt.Println(UIConfig.Reset + "This will NOT affect your Laravel codebase in app/")
	fmt.Println("except database credentials in the .env file.")
	DrawLine(UIConfig.Yellow)
	fmt.Println()

	// Stop containers before changing env to avoid conflicts
	PrintSection("Maintenance")
	LogWarning("Stopping active containers to avoid naming and port conflicts...")
	fmt.Println()
	HandleSimpleScript(PathStackDown)
	fmt.Println()

	PrintHeader("Updating stack configuration...")
	cfg := RunConfigWizard()
	WriteEnvFiles(cfg)

    // Remind about volume implications
    DrawLine(UIConfig.Yellow)
    LogWarning("REMINDER: Changing database name or password won't update existing volumes.")
    fmt.Println("If you changed those values, run 'stackfold reset --force' afterward.")
    fmt.Println("Read 'stackfold reset --help' before doing so — all data will be lost.")
    DrawLine(UIConfig.Yellow)
    fmt.Println()

	PrintHeader("Configuration updated")
}
