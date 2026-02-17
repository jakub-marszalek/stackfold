#!/bin/sh

set -e

cd /var/www/html

if [ ! -d vendor ]; then
  echo "[php] vendor missing on mounted volume – running composer install"
  composer install --no-interaction --prefer-dist
else
  echo "[php] vendor exists – refreshing autoloader"
  composer dump-autoload
fi

exec "$@"
