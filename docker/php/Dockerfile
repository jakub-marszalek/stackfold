FROM php:8.4-fpm-alpine

# Install packages
RUN apk add --no-cache \
    postgresql-dev \
    libpq \
    icu-dev \
    zip \
    unzip \
    git \
    curl \
    $PHPIZE_DEPS \
    && docker-php-ext-install pdo pdo_pgsql intl \
    && pecl install redis \
    && docker-php-ext-enable redis

# Copy composer
COPY --from=composer:2 /usr/bin/composer /usr/local/bin/composer
WORKDIR /var/www/html

COPY app/composer.json app/composer.lock ./

# Run composer install
RUN composer install --no-interaction --no-scripts --no-autoloader --prefer-dist
COPY app/ .
RUN composer dump-autoload --optimize

# Copy entrypoint
COPY docker/php/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Run entrypoint
ENTRYPOINT ["/entrypoint.sh"]
CMD ["php-fpm"]
