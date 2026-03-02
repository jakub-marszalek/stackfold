# Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
# Licensed under the MIT License. See LICENSE file in the project root.

FROM node:22-alpine

WORKDIR /var/www/html

RUN npm install -g npm@latest

# Copy packages
COPY package*.json ./

# Install dependencies
RUN npm install

CMD ["npm", "run", "dev"]
