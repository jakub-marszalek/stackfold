#!/usr/bin/env bash
#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

#
# Configuration / Constants
#
# Path related constants
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env.docker"

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

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

if [ "$REMOVE_VOLUMES" = true ]; then
  printSection "1. Whole stack shutting down (removing volumes)..."
  output=$(docker-compose --env-file "$ENV_FILE" down -v)
  status=$?
else
  printSection "1. Whole stack shutting down (keeping volumes)..."
  output=$(docker-compose --env-file "$ENV_FILE" down)
  status=$?
fi

echo "$output"

if [ "$status" -eq 0 ]; then
  echo -e "\n${GREEN}Docker stack shut down successfully.${RESET}"
else
  echo -e "\n${RED}Failed to shut down Docker stack!${RESET}"
  exit 1
fi
echo ""

# Summary
printHeader "Shutdown completed"
