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
)

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

// Install a fresh Laravel project into the app/ directory using official Composer Docker image.
func InstallLaravel(cfg EnvConfig) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Could not determine working directory: %w", err)
	}

	projectRoot, err := FindProjectRoot(workingDir)
	if err != nil {
		return fmt.Errorf("Could not locate project root: %w", err)
	}

	appDir := filepath.Join(projectRoot, "app")

	// Ensure app/ dir exists so Docker can bind-mount it
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("could not create app/ directory: %w", err)
	}

	// Run composer create-project via the official Docker image.
	// --user uid:gid ensures generated files are owned by the host user, not root.
	uid := fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())

	dockerArgs := []string{
		"run", "--rm",
		"--user", uid,
		"-v", appDir + ":/app",
		"composer:latest",
		"composer", "create-project", "laravel/laravel", ".", "^12.0",
		"--prefer-dist",
		"--no-interaction",
	}

	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("composer create-project failed: %w", err)
	}

	// Sync DB credentials from .env.docker into the Laravel .env
	fmt.Println()
	PrintSection("Syncing database credentials to Laravel .env...")
	WriteLaravelEnv(cfg)
	LogSuccess("Database credentials synced to app/.env.")

	// Patch vite.config.js for Docker HMR support
	fmt.Println()
	PrintSection("Patching vite.config.js for Docker...")
	if err := patchViteConfig(appDir); err != nil {
		LogWarning(fmt.Sprintf("Could not patch vite.config.js: %v", err))
		LogWarning("You may need to configure vite.config.js manually.")
	} else {
		LogSuccess("vite.config.js patched for Docker HMR.")
	}

	return nil
}

// Patch Laravel vite.config.js file to work with docker stack
func patchViteConfig(appDir string) error {
	viteConfig := `import { defineConfig } from 'vite';
import laravel from 'laravel-vite-plugin';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
    plugins: [
        laravel({
            input: ['resources/css/app.css', 'resources/js/app.js'],
            refresh: true,
        }),
        tailwindcss(),
    ],
    server: {
        host: '0.0.0.0',
        port: 5173,
        hmr: {
            host: 'localhost',
        },
        watch: {
            ignored: ['**/storage/framework/views/**'],
        },
    },
});
`
	configPath := filepath.Join(appDir, "vite.config.js")
	return os.WriteFile(configPath, []byte(viteConfig), 0644)
}
