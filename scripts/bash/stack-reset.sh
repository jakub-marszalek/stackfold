#!/usr/bin/env bash

# Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
# Licensed under the MIT License. See LICENSE file in the project root.

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Imports
source "$CURRENT_DIR"/_engine.sh
source "$CURRENT_DIR"/_ui.sh

#
# Reset sequence
#
printHeader "Resetting docker stack..."

# Step 1: Shutdown stack and remove volumes
printSection "$RESET_TAG Whole stack shutting down (removing volumes)..."
if [ -f "$STACK_DOWN_PATH" ]; then
  bash "$STACK_DOWN_PATH" --remove
elif [ -f "$STACK_DOWN_FILE" ]; then
  bash "$STACK_DOWN_FILE" --remove
else
  echo -e "${RED}Shutdown script not found! Expected $STACK_DOWN_PATH or $STACK_DOWN_FILE${RESET}"
  exit 1
fi

# Step 2: Clean build with startup
printSection "$RESET_TAG Rebuilding and starting stack..."
if [ -f "$STACK_BUILD_PATH" ]; then
  bash "$STACK_BUILD_PATH" --build-and-start
elif [ -f "$STACK_BUILD_FILE" ]; then
  bash "$STACK_BUILD_FILE" --build-and-start
else
  echo -e "${RED}Build script not found! Expected $STACK_BUILD_PATH or $STACK_BUILD_FILE${RESET}"
  exit 1
fi

# Summary
printHeader "Reset completed"
