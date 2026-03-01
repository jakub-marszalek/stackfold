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

# Fix storage and cache permissions so Laravel can write logs, cache, sessions
echo "[php] Setting storage permissions..."
chmod -R 775 storage bootstrap/cache
chown -R www-data:www-data storage bootstrap/cache

# Run migrations automatically on every container start.
# --force is required to run in non-interactive (production-like) environments.
echo "[php] Running migrations..."
php artisan migrate --force

exec "$@"
