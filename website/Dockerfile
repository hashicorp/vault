FROM docker.mirror.hashicorp.services/node:14.17.0-alpine
RUN apk add --update --no-cache git make g++ automake autoconf libtool nasm libpng-dev

COPY ./package.json /website/package.json
COPY ./package-lock.json /website/package-lock.json
WORKDIR /website
RUN npm install -g npm@latest
RUN npm install
