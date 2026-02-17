/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const asciiLogo = `    _____ __             __   ____      __    __ 
   / ___// /_____ ______/ /__/ __/___  / /___/ / 
   \__ \/ __/ __ '/ ___/ //_/ /_/ __ \/ / __  /  
  ___/ / /_/ /_/ / /__/ ,< / __/ /_/ / / /_/ /   
 /____/\__/\__,_/\___/_/|_/_/  \____/_/\__,_/`

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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
