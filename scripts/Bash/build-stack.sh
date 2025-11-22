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
START_STACK_FILE="$SCRIPT_DIR/start-stack.sh"
START_STACK_PATH="$SCRIPT_DIR/Bash/start-stack.sh"

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

#
# Parse Arguments
#
BUILD_AND_START=false

for arg in "$@"; do
  case $arg in
    --build-and-start|--bs|--B)
      BUILD_AND_START=true
      shift
      ;;
    *)
      ;;
  esac
done

#
# Build sequence
#
printHeader "Building docker stack..."

if [ "$BUILD_AND_START" = true ]; then
  echo -e "${CYAN}Option: Build and start${RESET}\n"
fi

# Step 1: Build docker stack
printSection "1. Building docker stack..."

if docker-compose --env-file "$ENV_FILE" build --no-cache; then
  echo -e "${GREEN}Docker build completed successfully.${RESET}\n"
else
  echo -e "${RED}Docker build failed!${RESET}\n"
  exit 1
fi

# Step 2: Optionally start stack
if [ "$BUILD_AND_START" = true ]; then
  printSection "2. Starting docker stack..."

  if [ -f "$START_STACK_PATH" ]; then
    bash "$START_STACK_PATH"
  elif [ -f "$START_STACK_FILE" ]; then
    bash "$START_STACK_FILE"
  else
    echo -e "${RED}Start script not found! Expected $START_STACK_PATH or $START_STACK_FILE${RESET}"
  fi
  echo ""
fi

# Summary
printHeader "Build completed"
