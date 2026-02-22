#!/usr/bin/env bash

#
# Configuration / Constants
#
# Path related constants
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env.docker"

# Stack related constants
PHP_CONTAINER="php"
PHP_HEALTH_CHECK_FILE="/var/www/html/src/Health/Stack/stack_health_check.php"
NGINX_URL="http://localhost:8080"
VITE_CONTAINER="node"
VITE_LOG_SEARCH="ready in"

# Load UI file
source "$SCRIPT_DIR"/_ui.sh

#
# Read dynamic ports from .env.docker
# Redis and Vite stay conventional — only HTTP and DB are user-configurable.
#
HTTP_PORT=$(grep -E '^HTTP_PORT=' "$ENV_FILE" | cut -d '=' -f2)
DB_PORT=$(grep -E '^DB_PORT=' "$ENV_FILE" | cut -d '=' -f2)
HTTP_PORT=${HTTP_PORT:-8080}
DB_PORT=${DB_PORT:-5432}

NGINX_URL="http://localhost:${HTTP_PORT}"

declare -A PORTS=(
  ["Nginx"]=$HTTP_PORT
  ["PostgreSQL"]=$DB_PORT
  ["Redis"]=6379
  ["Vite"]=5173
)

#
# Stack tests sequence
#
printHeader "Testing docker stack..."

# 1. Docker Compose Status
printSection "$TEST_TAG Testing Docker Compose Status..."
docker compose --env-file "$ENV_FILE" ps
echo ""

# 2. Nginx server HTTP check
printSection "$TEST_TAG Testing Nginx server (HTTP)..."
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
printSection "$TEST_TAG Testing PHP Stack..."
php_result=$(docker compose --env-file "$ENV_FILE" exec -T "$PHP_CONTAINER" php "$PHP_HEALTH_CHECK_FILE" 2>/dev/null)

if [ -n "$php_result" ]; then
  json=$(echo "$php_result" | awk '/^\{/{found=1} found{print}')

  if [ -n "$json" ] && echo "$json" | jq . >/dev/null 2>&1; then
    echo "$json" | jq .
  else
    echo -e "${RED}PHP health check returned no valid JSON.${RESET}"
    echo "Raw output:"
    echo "$php_result"
  fi
else
  echo -e "${RED}PHP health check failed or returned no data.${RESET}"
fi
echo ""

# 4. Vite Dev Server
printSection "$TEST_TAG Testing Vite Dev Server..."
vite_status=$(docker compose --env-file "$ENV_FILE" logs --tail=5 "$VITE_CONTAINER" | grep "$VITE_LOG_SEARCH")
if [ -n "$vite_status" ]; then
  echo -e "${GREEN}Vite: Running${RESET}"
  echo "$vite_status"
else
  echo -e "${YELLOW}WARNING! Vite: Check logs with: docker-compose --env-file $ENV_FILE logs node${RESET}"
fi
echo ""

# 5. Ports
printSection "$TEST_TAG Testing Open Ports..."
for name in "${!PORTS[@]}"; do
  port=${PORTS[$name]}
  if nc -z 127.0.0.1 "$port" &>/dev/null; then
    echo -e "${GREEN}$name: Listening on port $port - OK${RESET}"
  else
    echo -e "${RED}$name: NOT listening on port $port - ERROR${RESET}"
  fi
done

# Summary
printHeader "Test completed"
