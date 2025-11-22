#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

FROM php:8.3-fpm-alpine

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

WORKDIR /var/www/html

CMD ["php-fpm"]
