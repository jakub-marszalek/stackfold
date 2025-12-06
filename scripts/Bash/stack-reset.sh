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
STACK_DOWN_FILE="$SCRIPT_DIR/stack-down.sh"
STACK_DOWN_PATH="$SCRIPT_DIR/Bash/stack-down.sh"
STACK_BUILD_FILE="$SCRIPT_DIR/stack-build.sh"
STACK_BUILD_PATH="$SCRIPT_DIR/Bash/stack-build.sh"

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

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
echo ""

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
