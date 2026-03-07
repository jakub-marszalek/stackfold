#!/usr/bin/env bash
# Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
# Licensed under the MIT License. See LICENSE file in the project root.

VERSION=${1:-dev}
OUTPUT="../bin/stackfold"

echo "Building stackfold v$VERSION..."

go build \
  -ldflags "-X gitlab.com/hesoyamcc/stackfold/cli/cmd.version=$VERSION" \
  -o "$OUTPUT" \
  .

echo "Done: $OUTPUT"