/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores the PostgreSQL database from a backup",
	Long: `Lists available backups from the backups/ directory and lets you
choose which one to restore interactively.

WARNING: This will overwrite the current database contents.
You must use the --` + ForceFlag + ` flag to proceed.`,
	Run: func(cmd *cobra.Command, args []string) {
		force := GetBoolFlagOrExit(cmd, ForceFlag)
		if !force {
			DrawLine(UIConfig.Red)
			LogError("DESTRUCTIVE OPERATION DENIED")
			LogWarning("This will PERMANENTLY OVERWRITE your current database.")
			fmt.Println(UIConfig.Reset + "To proceed, use: stackfold restore --" + ForceFlag)
			DrawLine(UIConfig.Red)

			os.Exit(1)
		}

		runRestore()
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolP(
		ForceFlag,
		"f",
		false,
		"Force the restore and overwrite current database contents",
	)
}

func runRestore() {
	PrintHeader("Restoring database from backup...")

	backupsDir := filepath.Join(MustFindProjectRoot(), "backups")
	entries, err := os.ReadDir(backupsDir)
	if err != nil {
		LogError("No backups/ directory found. Run 'stackfold backup' first.")
		os.Exit(1)
	}

	var backups []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			backups = append(backups, e.Name())
		}
	}

	if len(backups) == 0 {
		LogError("No backup files found in backups/")
		os.Exit(1)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(backups)))

	PrintSection("Available backups (newest first):")
	fmt.Println()
	for i, name := range backups {
		fmt.Printf("  %s[%d]%s %s\n", UIConfig.Cyan, i+1, UIConfig.Reset, name)
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	var choice int
	for {
		fmt.Printf("%sSelect backup number [1-%d] or 0 to cancel: %s", UIConfig.Yellow, len(backups), UIConfig.Reset)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if _, err := fmt.Sscanf(input, "%d", &choice); err != nil {
			LogWarning("Invalid input. Please enter a number.")
			continue
		}

		if choice == 0 {
			LogWarning("Restore cancelled.")
			os.Exit(0)
		}

		if choice < 1 || choice > len(backups) {
			LogWarning(fmt.Sprintf("Please enter a number between 1 and %d, or 0 to cancel.", len(backups)))
			continue
		}

		break
	}

	selectedFile := backups[choice-1]
	selectedPath := filepath.Join(backupsDir, selectedFile)

	// Copy backup file into container
	PrintSection(fmt.Sprintf("Restoring %s...", selectedFile))

	container, dbName, dbUser := ResolveDBConfig()

	cpCmd := exec.Command("docker", "cp", selectedPath, fmt.Sprintf("%s:/tmp/restore.sql", container))
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	if err := cpCmd.Run(); err != nil {
		LogError(fmt.Sprintf("Failed to copy backup into container: %v", err))
		os.Exit(1)
	}

	// Run psql restore
	psqlCmd := exec.Command(
		"docker", "exec", container,
		"psql", "-U", dbUser, "-d", dbName, "-f", "/tmp/restore.sql",
	)
	psqlCmd.Stdout = os.Stdout
	psqlCmd.Stderr = os.Stderr
	if err := psqlCmd.Run(); err != nil {
		LogError(fmt.Sprintf("psql restore failed: %v", err))
		os.Exit(1)
	}

	LogSuccess(fmt.Sprintf("Database restored from: %s", selectedFile))
	PrintHeader("Restore completed")
}
