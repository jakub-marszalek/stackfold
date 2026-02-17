package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Check if certain bool flag is set
func GetBoolFlagOrExit(cmd *cobra.Command, flagName string) bool {
	value, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		fmt.Printf("Error when rading dlag '%s': %v\n", flagName, err)
		os.Exit(1)
	}
	return value
}

// Find project root directory
func FindProjectRoot(startPath string) (string, error) {
	currentPath := startPath

	for {
		// Check for docker-compose.yaml
		if _, err := os.Stat(filepath.Join(currentPath, "docker-compose.yaml")); err == nil {
			return currentPath, nil
		}

		// Go one level up
		parentPath := filepath.Dir(currentPath)

		// If there is no more parent directories
		if parentPath == currentPath {
			return "", fmt.Errorf("Project root not found")
		}

		currentPath = parentPath
	}
}

// Resolve absolute path from script root path
func ResolvePathFromRoot(scriptPathFromRoot string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error while getting working directory: %w", err)
	}

	projectRoot, err := FindProjectRoot(workingDir)
	if err != nil {
		return "", fmt.Errorf("Project root dir not found: %w", err)
	}

	scriptPath := filepath.Join(projectRoot, scriptPathFromRoot)
	return scriptPath, nil
}

// Run script using bash
func RunBashScript(scriptPath string, args ...string) error {
	bashCmd := exec.Command("bash", scriptPath)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr

	if len(args) > 0 {
		bashCmd.Args = append(bashCmd.Args, args...)
	}

	return bashCmd.Run()
}

// Handle script that has no flags
func HandleSimpleScript(scriptPathFromRoot string, args ...string) {
	fullPath, err := ResolvePathFromRoot(scriptPathFromRoot)
	if err != nil {
		LogError(fmt.Sprintf("Path error: %v\n", err))
		os.Exit(1)
	}

	if err := RunBashScript(fullPath, args...); err != nil {
		LogError("Script execution failed!")
		os.Exit(1)
	}
}

// Handle script that uses optional flag
func HandleFlaggedScript(scriptPath string, cmd *cobra.Command, flagName string, bashArg string) {
	var bashArgs []string

	if GetBoolFlagOrExit(cmd, flagName) {
		bashArgs = append(bashArgs, bashArg)
	}

	HandleSimpleScript(scriptPath, bashArgs...)
}
