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
	"strings"

	"github.com/spf13/cobra"
)

// Check if certain bool flag is set
func GetBoolFlagOrExit(cmd *cobra.Command, flagName string) bool {
	value, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		fmt.Printf("Error when rading flag '%s': %v\n", flagName, err)
		os.Exit(1)
	}
	return value
}

// Find project root directory
func FindProjectRoot(startPath string) (string, error) {
	currentPath := startPath

	for {
		if _, err := os.Stat(filepath.Join(currentPath, "docker-compose.yaml")); err == nil {
			return currentPath, nil
		}

		parentPath := filepath.Dir(currentPath)

		if parentPath == currentPath {
			return "", fmt.Errorf("Project root not found")
		}

		currentPath = parentPath
	}
}

// Return project root or exit with error
func MustFindProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		LogError(fmt.Sprintf("Could not get working directory: %v", err))
		os.Exit(1)
	}

	projectRoot, err := FindProjectRoot(wd)
	if err != nil {
		LogError(fmt.Sprintf("Could not find project root: %v", err))
		os.Exit(1)
	}

	return projectRoot
}

// Resolve absolute path from script root path
func ResolvePathFromRoot(scriptPathFromRoot string) string {
	return filepath.Join(MustFindProjectRoot(), scriptPathFromRoot)
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
	if err := RunBashScript(ResolvePathFromRoot(scriptPathFromRoot), args...); err != nil {
		LogError("Script execution failed!")
		os.Exit(1)
	}
}

// Handle script that uses optional flag
func HandleFlaggedScript(scriptPath string, cmd *cobra.Command, flagName string) {
	var bashArgs []string

	if GetBoolFlagOrExit(cmd, flagName) {
		bashArgs = append(bashArgs, "--"+flagName)
	}

	HandleSimpleScript(scriptPath, bashArgs...)
}

// Resolve container name from .env.docker by key (e.g. "CONTAINER_POSTGRES")
// Expands ${PROJECT_NAME} references automatically
func ResolveContainerName(key string) (string, error) {
	envPath := filepath.Join(MustFindProjectRoot(), ".env.docker")
	env, err := parseEnvFile(envPath)
	if err != nil {
		return "", fmt.Errorf("could not read .env.docker: %w", err)
	}

	value, ok := env[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in .env.docker", key)
	}

	// Expand ${PROJECT_NAME} if present
	if projectName, exists := env["PROJECT_NAME"]; exists {
		value = strings.ReplaceAll(value, "${PROJECT_NAME}", projectName)
	}

	return value, nil
}

// Rewrite a .env file, replacing the values for the given keys.
func updateEnvFile(path string, updates map[string]string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	var lines []string
	applied := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check if this is a commented-out key, e.g. "# DB_HOST=127.0.0.1"
		if strings.HasPrefix(trimmed, "#") {
			stripped := strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			parts := strings.SplitN(stripped, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				if newVal, ok := updates[key]; ok {
					lines = append(lines, key+"="+newVal)
					applied[key] = true
					continue
				}
			}
			lines = append(lines, line)
			continue
		}

		// Empty line - keep as-is
		if trimmed == "" {
			lines = append(lines, line)
			continue
		}

		// Active key
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if newVal, ok := updates[key]; ok {
				lines = append(lines, key+"="+newVal)
				applied[key] = true
				continue
			}
		}

		lines = append(lines, line)
	}
	f.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	// Append keys that were not found anywhere in the file
	for key, val := range updates {
		if !applied[key] {
			lines = append(lines, key+"="+val)
		}
	}

	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0644)
}

// Read a key=value env file and return a map of the entries ignoring commented and empty lines.
func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result, scanner.Err()
}
