/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Creates a backup of the PostgreSQL database",
	Long: `Dumps the current PostgreSQL database to a .sql file in the backups/ directory.

The backup file is named with a timestamp: backup_YYYY-MM-DD_HH-MM-SS.sql`,
	Run: func(cmd *cobra.Command, args []string) {
		runBackup()
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}

func runBackup() {
	PrintHeader("Creating database backup...")

	// Create backups/ directory if it doesn't exist
	backupsDir := filepath.Join(MustFindProjectRoot(), "backups")
	if err := os.MkdirAll(backupsDir, 0755); err != nil {
		LogError(fmt.Sprintf("Could not create backups/ directory: %v", err))
		os.Exit(1)
	}

	// Generate timestamped filename
	filename := fmt.Sprintf("backup_%s.sql", time.Now().Format("2006-01-02_15-04-05"))
	outputPath := filepath.Join(backupsDir, filename)

	// Run pg_dump inside the postgres container
	PrintSection("Running pg_dump...")

	container, dbName, dbUser := ResolveDBConfig()
	dumpCmd := exec.Command(
		"docker", "exec", container,
		"pg_dump", "-U", dbUser, dbName,
	)

	outFile, err := os.Create(outputPath)
	if err != nil {
		LogError(fmt.Sprintf("Could not create backup file: %v", err))
		os.Exit(1)
	}
	defer outFile.Close()

	dumpCmd.Stdout = outFile
	dumpCmd.Stderr = os.Stderr

	if err := dumpCmd.Run(); err != nil {
		LogError(fmt.Sprintf("pg_dump failed: %v", err))
		os.Remove(outputPath)
		os.Exit(1)
	}

	LogSuccess(fmt.Sprintf("Backup saved to: backups/%s", filename))
	PrintHeader("Backup completed")
}
