#!/usr/bin/env bash
#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

# Load UI file
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR"/_ui.sh
source "$SCRIPT_DIR"/_list-config.sh

# Scripts info list
scripts=(
  "stack-test.sh|Performs a full health check of all services in the stack.|None"
  "stack-up.sh|Starts the stack based on already built images and runs the health check.|None"
  "stack-build.sh|Builds Docker images for the entire stack.|--build-and-start (build and start stack)"
  "stack-down.sh|Stops all containers in the stack.|--remove (also removes volumes)"
  "stack-reset.sh|Performs a full stack reset: shutdown + build + start.|None"
  "list-scripts.sh|Lists all available scripts.|None"
)

#
# List sequence
#
drawLine
  showLogo
drawLine
# Variables loaded from _list-config.sh
# shellcheck disable=SC2154
echo -e "${WHITE}Project     : $project ${RESET}"
# shellcheck disable=SC2154
echo -e "${WHITE}Author      : $author ${RESET}"
# shellcheck disable=SC2154
echo -e "${WHITE}Environment : $environment ${RESET}"
# shellcheck disable=SC2154
echo -e "${WHITE}Version     : $version ${RESET}"
echo -e ""
# shellcheck disable=SC2154
echo -e "${WHITE}$description ${RESET}"
drawLine

# Display table
for entry in "${scripts[@]}"; do
  IFS="|" read -r script desc flags <<< "$entry"
  echo -e "${YELLOW}Script      :${RESET} ${script}"
  echo -e "${WHITE}Description :${RESET} ${desc}"
  echo -e "${GREEN}Flags       :${RESET} ${flags}"
  drawShortLine
done

# Display footer
echo -e "\n${CYAN}Usage:${RESET} ./scripts/Bash/<script-name>.sh [flags]"
echo -e "${CYAN}Example:${RESET} ./scripts/Bash/stack-build.sh --bs"
drawLine
