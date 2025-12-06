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
STACK_UP_FILE="$SCRIPT_DIR/stack-up.sh"
STACK_UP_PATH="$SCRIPT_DIR/Bash/stack-up.sh"

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
printSection "$BUILD_TAG Building docker stack..."
if docker-compose --env-file "$ENV_FILE" build --no-cache; then
  echo -e "${GREEN}Docker build completed successfully.${RESET}\n"
else
  echo -e "${RED}Docker build failed!${RESET}\n"
  exit 1
fi

# Step 2: Optionally start stack
if [ "$BUILD_AND_START" = true ]; then
  printSection "$BUILD_TAG Starting docker stack..."

  if [ -f "$STACK_UP_PATH" ]; then
    bash "$STACK_UP_PATH"
  elif [ -f "$STACK_UP_FILE" ]; then
    bash "$STACK_UP_FILE"
  else
    echo -e "${RED}Start script not found! Expected $STACK_UP_PATH or $STACK_UP_FILE${RESET}"
  fi
fi

# Summary
printHeader "Build completed"
