#!/usr/bin/env bash

#
# Configuration / Constants
#
# Path related constants
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env.docker"
STACK_TEST_FILE="$SCRIPT_DIR/stack-test.sh"
STACK_TEST_PATH="$SCRIPT_DIR/Bash/stack-test.sh"

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

# Function that waits for service to respond
wait_for_service() {
    local service_name="$1"
    local max_retries=30
    local counter=0
    local container_name=$(docker compose --env-file "$ENV_FILE" ps -q "$service_name")

    if [ -z "$container_name" ]; then
         echo -e "${RED}Service $service_name not found!${RESET}"
         return 1
    fi

    echo -n "Waiting for $service_name..."
    until [ "$(docker inspect -f '{{.State.Health.Status}}' "$container_name")" == "healthy" ]; do
        if [ "$counter" -gt "$max_retries" ]; then
            echo -e " ${RED}TIMEOUT!${RESET}"
            return 1
        fi
        echo -n "."
        sleep 1
        counter=$((counter+1))
    done

    echo -e " ${GREEN}OK!${RESET}"
}

#
# Stack start sequence
#
printHeader "Starting docker stack..."

# Step 1: Start whole docker stack
printSection "$START_TAG Starting whole stack..."
if docker compose --env-file "$ENV_FILE" up -d; then
  echo -e "${GREEN}Docker stack started successfully.${RESET}\n"
else
  echo -e "${RED}Failed to start Docker stack!${RESET}\n"
  exit 1
fi

# Step 2: Wait for services to start
printSection "$START_TAG Waiting for services to initialize..."
wait_for_service "db"
wait_for_service "redis"
wait_for_service "php"
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
