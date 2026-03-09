/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stops all containers in the stack",
	Long: `Stops and removes all containers associated with the stack. 
By default, it keeps the volumes intact, so your database data persists.

If the --` + RemoveFlag + ` flag is used, it will also destroy Docker volumes. 
Use this with caution — run 'stackfold backup' first to avoid data loss.`,
	Run: func(cmd *cobra.Command, args []string) {
		HandleFlaggedScript(PathStackDown, cmd, RemoveFlag)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().BoolP(
		RemoveFlag,
		"r",
		false,
		"Also remove volumes. WARNING: this will remove all data from your database! Run 'stackfold backup' first.",
	)
}
