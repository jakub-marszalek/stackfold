#
# Copyright (c) 2025. Jakub Marszałek / hesoyamcc
# All Rights Reserved.
#

FROM node:22-alpine

WORKDIR /var/www/html

RUN npm install -g npm@latest

# Copy packages
COPY package*.json ./

# Install dependencies
RUN npm install

CMD ["npm", "run", "dev"]
