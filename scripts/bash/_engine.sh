#!/usr/bin/env bash

#
# Paths
#
SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPTS_DIR/../.." && pwd)"
# shellcheck disable=SC2034
ENV_FILE="$PROJECT_ROOT/.env.docker"

# Script paths
# shellcheck disable=SC2034
STACK_BUILD_PATH="$SCRIPTS_DIR/stack-build.sh"
# shellcheck disable=SC2034
STACK_DOWN_PATH="$SCRIPTS_DIR/stack-down.sh"
# shellcheck disable=SC2034
STACK_TEST_PATH="$SCRIPTS_DIR/stack-test.sh"
# shellcheck disable=SC2034
STACK_UP_PATH="$SCRIPTS_DIR/stack-up.sh"

#
# Dynamic compose detection
#
if docker compose version >/dev/null 2>&1; then
    DOCKER_COMPOSE="docker compose"
elif docker-compose --version >/dev/null 2>&1; then
    # shellcheck disable=SC2034
    DOCKER_COMPOSE="docker-compose"
else
    source "$SCRIPTS_DIR/_ui.sh"
    echo -e "${RED}[SYSTEM ERROR] Neither 'docker compose' nor 'docker-compose' found.${RESET}"
    echo -e "${YELLOW}Stackfold requires Docker Compose to function.${RESET}"
    exit 1
fi
