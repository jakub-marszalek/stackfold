/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Holds all values collected by the config wizard
type EnvConfig struct {
	ProjectName string
	DbName      string
	DbUser      string
	DbPass      string
	HttpPort    string
	DbPort      string
}

/*
Interactively collect project name and DB credentials,
then write .env.docker and app/.env.
Shared by both 'init' (first install) and 'config' (update existing).
*/
func RunConfigWizard() EnvConfig {
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

	return EnvConfig{
		ProjectName: projectName,
		DbName:      dbName,
		DbUser:      dbUser,
		DbPass:      dbPass,
		HttpPort:    httpPort,
		DbPort:      dbPort,
	}
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

// Ask the user for a port number and keep prompting until they provide one that is currently free on the host
func askPort(reader *bufio.Reader, label string, defaultPort string) string {
	for {
		port := askQuestion(reader, label, defaultPort)
		if isPortAvailable(port) {
			return port
		}
		LogWarning(fmt.Sprintf("Port %s is already in use. Please choose a different port.", port))
	}
}

// Check if given TCP port is free on the host machine
func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
