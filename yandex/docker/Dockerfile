ARG BASE_VAULT_VERSION=1.11.5

# Clone sources from https://github.com/yandex-cloud/vault.git and build vault binary with support for Yandex KMS
FROM --platform=linux/amd64 golang:1.19
ARG BASE_VAULT_VERSION
RUN set -eux \
  && cd / \
  && apt update -y \
  && apt install -y build-essential \
  && wget https://nodejs.org/dist/v16.18.1/node-v16.18.1-linux-x64.tar.xz \
  && tar -oxf node-v16.18.1-linux-x64.tar.xz \
  && PATH=$PATH:/node-v16.18.1-linux-x64/bin \
  && npm install -g yarn \
  && git clone https://github.com/yandex-cloud/vault.git \
  && cd vault \
  && git checkout v${BASE_VAULT_VERSION}+yckms \
  && make bootstrap static-dist dev-ui

# Replace vault binary in official vault image with the binary built on the previous stage
FROM vault:${BASE_VAULT_VERSION}
COPY --from=0 /vault/bin/vault /bin/vault
