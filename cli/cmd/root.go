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

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "stackfold",
	Short: "Stackfold - Laravel Docker Automation",
	Long: fmt.Sprintf("%s\n\nProject:     Stackfold - Laravel Docker Automation\nEnvironment: Development\n\n%s",
		asciiLogo,
		"Stackfold is an automation tool designed to streamline the setup,\n"+
			"build, and management of new Laravel projects using Docker.\n"+
			"It handles container orchestration, health checks, and environment\n"+
			"setup, enabling developers to start new projects quickly and consistently."),

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
}
