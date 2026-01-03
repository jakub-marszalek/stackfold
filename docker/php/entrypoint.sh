#!/bin/sh
#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

set -e

cd /var/www/html

if [ ! -d vendor ]; then
  echo "[php] vendor missing – running composer install"
  composer install \
    --no-interaction \
    --prefer-dist \
    --no-progress
else
  echo "[php] vendor exists – skipping composer install"
fi

exec "$@"
