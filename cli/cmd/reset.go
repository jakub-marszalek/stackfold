/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Performs a full stack reset: shutdown, build, and start",
	Long: `Performs a complete environment reset:
1. Shuts down the stack and removes all Docker volumes (WIPES ALL DATA).
2. Rebuilds Docker images from scratch.
3. Starts all services.

WARNING: This command is DESTRUCTIVE. You must use the --` + ForceFlag + ` flag to proceed.`,
	Run: func(cmd *cobra.Command, args []string) {
		force := GetBoolFlagOrExit(cmd, ForceFlag)

		if !force {
			DrawLine(UIConfig.Red)
			LogError("DESTRUCTIVE OPERATION DENIED")
			LogWarning("This will PERMANENTLY DELETE your database volumes.")
			fmt.Println(UIConfig.Reset + "To proceed, use: sf reset --" + ForceFlag)
			DrawLine(UIConfig.Red)

			DrawLine(UIConfig.Yellow)
			LogWarning("RECOMMENDATION: Back up your database before proceeding.")
			fmt.Println(UIConfig.Reset + "Run 'stackfold backup' to create a backup first.")
			DrawLine(UIConfig.Yellow)
			fmt.Println()

			os.Exit(1)
		}

		HandleSimpleScript(PathStackReset)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().BoolP(
		ForceFlag,
		"f",
		false,
		"Force the reset and delete all data volumes",
	)
}
