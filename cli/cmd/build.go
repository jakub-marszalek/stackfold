/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"github.com/spf13/cobra"
)

const (
	BuildFlag     = "--build-and-start"
	BuildFlagName = "start"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images for the stack",
	Long: `Builds or rebuilds Docker images defined in the stack configuration using --no-cache.

If the --start flag is provided, the stack will automatically start after 
a successful build by invoking the stack-up sequence.`,

	Run: func(cmd *cobra.Command, args []string) {
		HandleFlaggedScript(PathStackBuild, cmd, BuildFlagName, BuildFlag)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolP(
		BuildFlagName,
		"s",
		false,
		"Build images and start containers immediately after",
	)
}
