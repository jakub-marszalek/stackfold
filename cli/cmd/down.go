/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

const (
	RemoveFlag     = "--remove"
	RemoveFlagName = "remove"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stops all containers in the stack",
	Long: `Stops and removes all containers associated with the stack. 
By default, it keeps the volumes intact, so your database data persists.

If the --remove flag is used, it will also destroy Docker volumes. 
Use this with caution!`,

	Run: func(cmd *cobra.Command, args []string) {
		HandleFlaggedScript(PathStackDown, cmd, RemoveFlagName, RemoveFlag)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().BoolP(
		"remove",
		"r",
		false,
		"Also remove volumes. WARNING: this will remove all data from your database!",
	)
}
