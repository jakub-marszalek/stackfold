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
STACK_TEST_FILE="$SCRIPT_DIR/stack-test.sh"
STACK_TEST_PATH="$SCRIPT_DIR/Bash/stack-test.sh"

# Waiting time before running healthcheck (in seconds)
SLEEP_SECONDS=10

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

#
# Stack start sequence
#
printHeader "Starting docker stack..."

# Step 1: Start whole docker stack
printSection "$START_TAG Starting whole stack..."
if docker-compose --env-file "$ENV_FILE" up -d; then
  echo -e "${GREEN}Docker stack started successfully.${RESET}\n"
else
  echo -e "${RED}Failed to start Docker stack!${RESET}\n"
  exit 1
fi

# Step 2: Wait for services to start
printSection "$START_TAG Waiting for services to start (${SLEEP_SECONDS}s)..."
sleep "$SLEEP_SECONDS"
echo ""

# Step 3: Run stack healthcheck
printSection "$START_TAG Running stack's healthcheck..."
if [ -f "$STACK_TEST_PATH" ]; then
  bash "$STACK_TEST_PATH"
elif [ -f "$STACK_TEST_FILE" ]; then
  bash "$STACK_TEST_FILE"
else
  echo -e "${RED}Healthcheck script not found! Expected $STACK_TEST_PATH or $STACK_TEST_FILE${RESET}"
fi

# Summary
printHeader "Start completed"
