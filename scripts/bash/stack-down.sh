#!/usr/bin/env bash
# Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
# Licensed under the MIT License. See LICENSE file in the project root.

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Imports
source "$CURRENT_DIR"/_engine.sh
source "$CURRENT_DIR"/_ui.sh

#
# Parse arguments
#
REMOVE_VOLUMES=false

for arg in "$@"; do
  case $arg in
    --remove|--r)
      REMOVE_VOLUMES=true
      shift
      ;;
    *)
      ;;
  esac
done

#
# Shutdown sequence
#
printHeader "Stack shutting down..."

# Step 1: Check if --remove flag is set
if [ "$REMOVE_VOLUMES" = true ]; then
  echo -e "${CYAN}Option: Remove volumes${RESET}\n"
  printSection "$SHUTDOWN_TAG Whole stack shutting down (removing volumes)..."
  output=$($DOCKER_COMPOSE --env-file "$ENV_FILE" down -v)
  status=$?
else
  printSection "$SHUTDOWN_TAG Whole stack shutting down (keeping volumes)..."
  output=$($DOCKER_COMPOSE --env-file "$ENV_FILE" down)
  status=$?
fi

# Step 2: Echo output
echo "$output"

if [ "$status" -eq 0 ]; then
  echo -e "\n${GREEN}Docker stack shut down successfully.${RESET}"
else
  echo -e "\n${RED}Failed to shut down Docker stack!${RESET}"
  exit 1
fi

# Summary
printHeader "Shutdown completed"
