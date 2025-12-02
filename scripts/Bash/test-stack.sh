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

# Stack related constants
PHP_CONTAINER="php"
PHP_HEALTH_CHECK_FILE="/var/www/html/docker/Health/Stack/stack_health_check.php"
NGINX_URL="http://localhost:8080"
VITE_CONTAINER="node"
VITE_LOG_SEARCH="ready in"

# Ports definition (name:port)
declare -A PORTS=(
  ["Nginx"]=8080
  ["Vite"]=5173
  ["PostgreSQL"]=5432
  ["Redis"]=6379
)

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

#
# Stack tests sequence
#
printHeader "Testing docker stack..."

# 1. Docker Compose Status
printSection "1. Testing Docker Compose Status..."
  docker-compose --env-file "$ENV_FILE" ps
echo ""

# 2. Nginx server HTTP check
printSection "2. Testing Nginx server (HTTP)..."
  if command -v curl &>/dev/null; then
    status_code=$(curl -s -o /dev/null -w "%{http_code}" "$NGINX_URL")
    if [[ "$status_code" == "200" || "$status_code" == "302" ]]; then
      echo -e "${GREEN}Nginx: OK (Status $status_code)${RESET}"
    else
      echo -e "${RED}Nginx: ERROR (Status $status_code)${RESET}"
    fi
  else
    echo -e "${RED}curl not installed.${RESET}"
  fi
echo ""

# 3. PHP Health Check
printSection "3. Testing PHP Stack..."
  php_result=$(docker-compose --env-file "$ENV_FILE" exec -T "$PHP_CONTAINER" php "$PHP_HEALTH_CHECK_FILE" 2>/dev/null)
  if [ -n "$php_result" ]; then
    # pretty-print JSON if jq is available
    if command -v jq &>/dev/null; then
      echo "$php_result" | jq .
    else
      echo "$php_result"
    fi
  else
    echo -e "${RED}PHP health check failed or returned no data.${RESET}"
  fi
echo ""

# 4. Vite Dev Server
printSection "4. Testing Vite Dev Server..."
  vite_status=$(docker-compose --env-file "$ENV_FILE" logs --tail=5 "$VITE_CONTAINER" | grep "$VITE_LOG_SEARCH")
  if [ -n "$vite_status" ]; then
    echo -e "${GREEN}Vite: Running${RESET}"
    echo "$vite_status"
  else
    echo -e "${YELLOW}WARNING! Vite: Check logs with: docker-compose --env-file $ENV_FILE logs node${RESET}"
  fi
echo ""

# 5. Ports
printSection "5. Testing Open Ports..."
  for name in "${!PORTS[@]}"; do
    port=${PORTS[$name]}
    if nc -z 127.0.0.1 "$port" &>/dev/null; then
      echo -e "${GREEN}OK $name: Listening on port $port${RESET}"
    else
      echo -e "${RED}ERROR $name: NOT listening on port $port${RESET}"
    fi
  done
echo ""

# Summary
printHeader "Test completed"
