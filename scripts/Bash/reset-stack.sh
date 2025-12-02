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
SHUTDOWN_STACK_FILE="$SCRIPT_DIR/shutdown-stack.sh"
SHUTDOWN_STACK_PATH="$SCRIPT_DIR/Bash/shutdown-stack.sh"
BUILD_STACK_FILE="$SCRIPT_DIR/build-stack.sh"
BUILD_STACK_PATH="$SCRIPT_DIR/Bash/build-stack.sh"

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

#
# Reset sequence
#
printHeader "Resetting docker stack..."

# Step 1: Shutdown stack and remove volumes
printSection "1. Whole stack shutting down (removing volumes)..."
if [ -f "$SHUTDOWN_STACK_PATH" ]; then
  bash "$SHUTDOWN_STACK_PATH" --remove
elif [ -f "$SHUTDOWN_STACK_FILE" ]; then
  bash "$SHUTDOWN_STACK_FILE" --remove
else
  echo -e "${RED}Shutdown script not found! Expected $SHUTDOWN_STACK_PATH or $SHUTDOWN_STACK_FILE${RESET}"
  exit 1
fi
echo ""

# Step 2: Clean build with startup
printSection "2. Rebuilding and starting stack..."
if [ -f "$BUILD_STACK_PATH" ]; then
  bash "$BUILD_STACK_PATH" --build-and-start
elif [ -f "$BUILD_STACK_FILE" ]; then
  bash "$BUILD_STACK_FILE" --build-and-start
else
  echo -e "${RED}Build script not found! Expected $BUILD_STACK_PATH or $BUILD_STACK_FILE${RESET}"
  exit 1
fi
echo ""

# Summary
printHeader "Reset completed"
