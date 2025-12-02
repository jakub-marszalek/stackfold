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
TEST_STACK_FILE="$SCRIPT_DIR/test-stack.sh"
TEST_STACK_PATH="$SCRIPT_DIR/Bash/test-stack.sh"

# Waiting time before running healthcheck (in seconds)
SLEEP_SECONDS=10

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

#
# Stack start sequence
#
printHeader "Starting docker stack..."

# Step 1: Start whole docker stack
printSection "1. Starting whole stack..."

if docker-compose --env-file "$ENV_FILE" up -d; then
  echo -e "${GREEN}Docker stack started successfully.${RESET}\n"
else
  echo -e "${RED}Failed to start Docker stack!${RESET}\n"
  exit 1
fi

# Step 2: Wait for services to start
printSection "2. Waiting for services to start (${SLEEP_SECONDS}s)..."
sleep "$SLEEP_SECONDS"
echo ""

# Step 3: Run stack healthcheck
printSection "3. Running stack's healthcheck..."
if [ -f "$TEST_STACK_PATH" ]; then
  bash "$TEST_STACK_PATH"
elif [ -f "$TEST_STACK_FILE" ]; then
  bash "$TEST_STACK_FILE"
else
  echo -e "${RED}Healthcheck script not found! Expected $TEST_STACK_PATH or $TEST_STACK_FILE${RESET}"
fi
echo ""

# Summary
printHeader "Start completed"
