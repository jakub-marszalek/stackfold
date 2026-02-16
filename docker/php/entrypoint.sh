#!/bin/sh

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
