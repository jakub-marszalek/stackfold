#!/usr/bin/env bash
#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

#
# Default configuration
#

# Colors
# shellcheck disable=SC2034
CYAN="\033[0;36m"
# shellcheck disable=SC2034
YELLOW="\033[1;33m"
# shellcheck disable=SC2034
GREEN="\033[0;32m"
# shellcheck disable=SC2034
RED="\033[0;31m"
# shellcheck disable=SC2034
WHITE="\033[1;37m"
# shellcheck disable=SC2034
DARK_GREY="\033[90m"

# UI constants
declare -A UIConfig=(
    [LineChar]="-"
    [EdgeChar]="+"
    [LineWidth]=80
    [ShortLineWidth]=55
    [TitleColor]=${CYAN}
    [SeparatorLineColor]=${DARK_GREY}
)

RESET="\033[0m"

#
# UI functions
#

# Print Stackr logo
showLogo() {
    echo -e "${UIConfig[TitleColor]}   ____  _             _       ____   ${RESET}"
    echo -e "${UIConfig[TitleColor]}  / ___|| |_ ____  ___| | ___ |  _ \  ${RESET}"
    echo -e "${UIConfig[TitleColor]}  \___ \| __/ _  |/ __| |/ __|| |_)_| ${RESET}"
    echo -e "${UIConfig[TitleColor]}   ___) | |( (_| | (__|   (__ |    \  ${RESET}"
    echo -e "${UIConfig[TitleColor]}  |____/ \__\__,_|\___|_|\___|\_|\__\ ${RESET}"
    echo
}

# Draw simple line
drawLine() {
    local width="${1:-${UIConfig[LineWidth]}}"
    local edge="${2:-${UIConfig[EdgeChar]}}"
    local char="${3:-${UIConfig[LineChar]}}"
    local color="${4:-${UIConfig[TitleColor]}}"

    local innerWidth=$((width - 2))
    (( innerWidth < 1 )) && innerWidth=1

    local line
    line="${edge}$(printf '%*s' "$innerWidth" "" | tr ' ' "$char")${edge}"

    echo -e "${color}${line}${RESET}"
}

# Draw short line
drawShortLine() {
    local width="${1:-${UIConfig[ShortLineWidth]}}"
    local char="${2:-${UIConfig[LineChar]}}"
    local color="${3:-${UIConfig[SeparatorLineColor]}}"

    local line
    line="$(printf '%*s' "$width" ""| tr ' ' "$char")"

    echo -e "${color}${line}${RESET}"
}

# Write centered header text
printHeader() {
    local text="$1"
    local width="${2:-${UIConfig[LineWidth]}}"
    local color="${3:-${UIConfig[TitleColor]}}"

    drawLine "$width"

    local padding=$(( (width - ${#text}) / 2 ))
    (( padding < 0 )) && padding=0

    printf -v line "%*s%s" "$padding" "" "$text"
    echo -e "${color}${line}${RESET}"

    drawLine "$width"
}

# Print section
printSection() {
    echo -e "${YELLOW}$1${RESET}"
}
