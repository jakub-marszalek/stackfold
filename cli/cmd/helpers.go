package cmd

import (
	"bufio"
	"fmt"
	"net"
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
		fmt.Printf("Error when rading dlag '%s': %v\n", flagName, err)
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

/*
Interactively collect project name and DB credentials,
then write .env.docker and app/.env.
Shared by both 'init' (first install) and 'config' (update existing).
*/
func RunConfigWizard() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(UIConfig.SeparatorColor + "Tip: Leave blank to accept the [default] values." + UIConfig.Reset)
	fmt.Println()

	// --- Environment Configuration ---
	PrintSection("Environment Configuration")
	projectName := askQuestion(reader, "Project Name", "myapp")
	fmt.Println()

	// --- Database Credentials ---
	PrintSection("Database Credentials")
	dbName := askQuestion(reader, "Database Name", "mydatabase")
	dbUser := askQuestion(reader, "Database User", "admin")
	dbPass := askQuestion(reader, "Database Password", "admin123")
	fmt.Println()

	// --- Ports ---
	PrintSection("Port configuration")
	fmt.Println(UIConfig.SeparatorColor + "Change these if the default ports are already in use on your machine." + UIConfig.Reset)
	httpPort := askPort(reader, "HTTP port (Nginx)", "8080")
	dbPort := askPort(reader, "PostgreSQL port", "5432")
	fmt.Println()

	// --- Write env files ---
	PrintSection("Writing configuration files...")
	CreateEnvFile(projectName, dbName, dbUser, dbPass, httpPort, dbPort)
	fmt.Println()
}

// Write .env.docker and app/.env with provided credentials
func CreateEnvFile(projectName, dbName, dbUser, dbPass, httpPort, dbPort string) {
	dockerEnv := fmt.Sprintf(`# Docker Stack Configuration
PROJECT_NAME=%s

# Containers
CONTAINER_NGINX=${PROJECT_NAME}-nginx
CONTAINER_PHP=${PROJECT_NAME}-php
CONTAINER_POSTGRES=${PROJECT_NAME}-postgres
CONTAINER_REDIS=${PROJECT_NAME}-redis
CONTAINER_NODE=${PROJECT_NAME}-node

# Ports
HTTP_PORT=%s
DB_PORT=%s

# Postgres Credentials
POSTGRES_DATABASE=%s
POSTGRES_USER=%s
POSTGRES_PASSWORD=%s
`, projectName, httpPort, dbPort, dbName, dbUser, dbPass)

	appEnv := fmt.Sprintf(`# Laravel Configuration
DATABASE_NAME=%s
DATABASE_USER=%s
DATABASE_PASSWORD=%s
`, dbName, dbUser, dbPass)

	if err := os.WriteFile(".env.docker", []byte(dockerEnv), 0644); err != nil {
		LogError(fmt.Sprintf("Failed to write .env.docker file: %v", err))
		os.Exit(1)
	}

	if err := os.MkdirAll("app", 0755); err != nil {
		LogError(fmt.Sprintf("Failed to create app directory: %v", err))
		os.Exit(1)
	}

	if err := os.WriteFile("app/.env", []byte(appEnv), 0644); err != nil {
		LogError(fmt.Sprintf("Failed to write app/.env: %v", err))
		os.Exit(1)
	}

	LogSuccess("Project name and database credentials have been synced across files.")
}

// Check if given TCP port is free on the host machine
func IsPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// Check if a Laravel project already exists in app/ directory
func IsLaravelInstalled() bool {
	indicators := []string{"app/artisan", "app/composer.json"}
	for _, path := range indicators {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

// Prompt the user with a label and default value, returning the input
func askQuestion(reader *bufio.Reader, label string, defaultValue string) string {
	fmt.Printf("%s%s%s [default: %s]: ", UIConfig.Yellow, label, UIConfig.Reset, defaultValue)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

// Ask the user for a port number and keep prompting until they provide one that is currently free on the host.
func askPort(reader *bufio.Reader, label string, defaultPort string) string {
	for {
		port := askQuestion(reader, label, defaultPort)
		if IsPortAvailable(port) {
			return port
		}
		LogWarning(fmt.Sprintf("Port %s is already in use. Please choose a different port.", port))
	}
}
