# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

FROM docker.mirror.hashicorp.services/node:18.18.2-alpine
RUN apk add --update --no-cache git make g++ automake autoconf libtool nasm libpng-dev

COPY ./package.json /website/package.json
COPY ./package-lock.json /website/package-lock.json
WORKDIR /website
RUN npm install -g npm@latest
RUN npm install
