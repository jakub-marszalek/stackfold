/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Write .env.docker from the given config
func WriteDockerEnv(cfg EnvConfig) {
	dockerEnv := fmt.Sprintf(`PROJECT_NAME=%s

CONTAINER_NGINX=${PROJECT_NAME}-nginx
CONTAINER_PHP=${PROJECT_NAME}-php
CONTAINER_POSTGRES=${PROJECT_NAME}-postgres
CONTAINER_REDIS=${PROJECT_NAME}-redis
CONTAINER_NODE=${PROJECT_NAME}-node

HTTP_PORT=%s
DB_PORT=%s

POSTGRES_DATABASE=%s
POSTGRES_USER=%s
POSTGRES_PASSWORD=%s
`, cfg.ProjectName, cfg.HttpPort, cfg.DbPort, cfg.DbName, cfg.DbUser, cfg.DbPass)

	if err := os.WriteFile(".env.docker", []byte(dockerEnv), 0644); err != nil {
		LogError(fmt.Sprintf("Failed to write .env.docker file: %v", err))
		os.Exit(1)
	}
}

// Sync DB credentials into app/.env
// If the file doesn't exist yet, a minimal stub is created
// If it does exist, only the relevant keys are updated (existing content is preserved)
func WriteLaravelEnv(cfg EnvConfig) {
	const laravelEnvPath = "app/.env"

	updates := map[string]string{
		"DB_CONNECTION": "pgsql",
		"DB_HOST":       "db",
		"DB_PORT":       cfg.DbPort,
		"DB_DATABASE":   cfg.DbName,
		"DB_USERNAME":   cfg.DbUser,
		"DB_PASSWORD":   cfg.DbPass,
		"REDIS_CLIENT":  "phpredis",
		"REDIS_HOST":    "redis",
	}

	if _, err := os.Stat(laravelEnvPath); os.IsNotExist(err) {
		if err := os.MkdirAll("app", 0755); err != nil {
			LogError(fmt.Sprintf("Failed to create app/ directory: %v", err))
			os.Exit(1)
		}

		var sb strings.Builder
		for k, v := range updates {
			sb.WriteString(k + "=" + v + "\n")
		}

		if err := os.WriteFile(laravelEnvPath, []byte(sb.String()), 0644); err != nil {
			LogError(fmt.Sprintf("Failed to write %s: %v", laravelEnvPath, err))
			os.Exit(1)
		}
		return
	}

	if err := updateEnvFile(laravelEnvPath, updates); err != nil {
		LogError(fmt.Sprintf("Failed to update %s: %v", laravelEnvPath, err))
		os.Exit(1)
	}
}

// Write both .env.docker and app/.env and print a success summary
func WriteEnvFiles(cfg EnvConfig) {
	PrintSection("Writing configuration files...")
	WriteDockerEnv(cfg)
	WriteLaravelEnv(cfg)
	LogSuccess("Project name and database credentials have been synced across files.")
	fmt.Println()
}

// Read database credentials and container name from .env.docker.
func ResolveDBConfig() (container, dbName, dbUser string) {
	envPath := filepath.Join(MustFindProjectRoot(), ".env.docker")
	env, err := parseEnvFile(envPath)
	if err != nil {
		LogError(fmt.Sprintf("Could not read .env.docker: %v", err))
		os.Exit(1)
	}

	container, err = ResolveContainerName("CONTAINER_POSTGRES")
	if err != nil {
		LogError(fmt.Sprintf("Could not resolve container name: %v", err))
		os.Exit(1)
	}

	dbName = env["POSTGRES_DATABASE"]
	dbUser = env["POSTGRES_USER"]

	if dbName == "" || dbUser == "" {
		LogError("Missing POSTGRES_DATABASE or POSTGRES_USER in .env.docker")
		os.Exit(1)
	}

	return container, dbName, dbUser
}
