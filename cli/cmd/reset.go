/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	ForceFlagName = "force"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Performs a full stack reset: shutdown, build, and start",
	Long: `Performs a complete environment reset:
1. Shuts down the stack and removes all Docker volumes (WIPES ALL DATA).
2. Rebuilds Docker images from scratch.
3. Starts all services.

WARNING: This command is DESTRUCTIVE. You must use the --force flag to proceed.`,

	Run: func(cmd *cobra.Command, args []string) {
		force := GetBoolFlagOrExit(cmd, ForceFlagName)

		if !force {
			DrawLine(UIConfig.Red)
			LogError("DESTRUCTIVE OPERATION DENIED")
			LogWarning("This will PERMANENTLY DELETE your database volumes.")
			fmt.Println(UIConfig.Reset + "To proceed, use: sf reset --force")
			DrawLine(UIConfig.Red)

			os.Exit(1)
		}

		HandleSimpleScript(PathStackReset)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().BoolP(
		ForceFlagName,
		"f",
		false,
		"Force the reset and delete all data volumes",
	)
}
