# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

## DOCKERHUB DOCKERFILE ##
FROM alpine:3 AS default

ARG BIN_NAME
# NAME and PRODUCT_VERSION are the name of the software in releases.hashicorp.com
# and the version to download. Example: NAME=vault PRODUCT_VERSION=1.2.3.
ARG NAME=vault
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH
# LICENSE_SOURCE is the path to IBM license documents, which may be architecture-specific.
ARG LICENSE_SOURCE
# LICENSE_DEST is the path where license files are installed in the container
ARG LICENSE_DEST

# Additional metadata labels used by container registries, platforms
# and certification scanners.
LABEL name="Vault" \
      maintainer="Vault Team <vault@hashicorp.com>" \
      vendor="HashiCorp" \
      version=${PRODUCT_VERSION} \
      release=${PRODUCT_REVISION} \
      revision=${PRODUCT_REVISION} \
      summary="Vault is a tool for securely accessing secrets." \
      description="Vault is a tool for securely accessing secrets. A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, and more. Vault provides a unified interface to any secret, while providing tight access control and recording a detailed audit log."

# Copy the license file as per Legal requirement
COPY ${LICENSE_SOURCE} ${LICENSE_DEST}

# Set ARGs as ENV so that they can be used in ENTRYPOINT/CMD
ENV NAME=$NAME
ENV VERSION=$VERSION

# Create a non-root user to run the software.
RUN addgroup ${NAME} && adduser -S -G ${NAME} ${NAME}

RUN apk add --no-cache libcap su-exec dumb-init tzdata

COPY dist/$TARGETOS/$TARGETARCH/$BIN_NAME /bin/

# /vault/logs is made available to use as a location to store audit logs, if
# desired; /vault/file is made available to use as a location with the file
# storage backend, if desired; the server will be started with /vault/config as
# the configuration directory so you can add additional config files in that
# location.
RUN mkdir -p /vault/logs && \
    mkdir -p /vault/file && \
    mkdir -p /vault/config && \
    chown -R ${NAME}:${NAME} /vault

# Expose the logs directory as a volume since there's potentially long-running
# state in there
VOLUME /vault/logs

# Expose the file directory as a volume since there's potentially long-running
# state in there
VOLUME /vault/file

# 8200/tcp is the primary interface that applications use to interact with
# Vault.
EXPOSE 8200

# The entry point script uses dumb-init as the top-level process to reap any
# zombie processes created by Vault sub-processes.
#
# For production derivatives of this container, you should add the IPC_LOCK
# capability so that Vault can mlock memory.
COPY .release/docker/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]


# # By default you'll get a single-node development server that stores everything
# # in RAM and bootstraps itself. Don't use this configuration for production.
CMD ["server", "-dev"]


## UBI DOCKERFILE ##
FROM registry.access.redhat.com/ubi8/ubi-minimal AS ubi

ARG BIN_NAME
# NAME and PRODUCT_VERSION are the name of the software in releases.hashicorp.com
# and the version to download. Example: NAME=vault PRODUCT_VERSION=1.2.3.
ARG NAME=vault
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH
# LICENSE_SOURCE is the path to IBM license documents, which may be architecture-specific.
ARG LICENSE_SOURCE
# LICENSE_DEST is the path where license files are installed in the container
ARG LICENSE_DEST

# Additional metadata labels used by container registries, platforms
# and certification scanners.
LABEL name="Vault" \
      maintainer="Vault Team <vault@hashicorp.com>" \
      vendor="HashiCorp" \
      version=${PRODUCT_VERSION} \
      release=${PRODUCT_REVISION} \
      revision=${PRODUCT_REVISION} \
      summary="Vault is a tool for securely accessing secrets." \
      description="Vault is a tool for securely accessing secrets. A secret is anything that you want to tightly control access to, such as API keys, passwords, certificates, and more. Vault provides a unified interface to any secret, while providing tight access control and recording a detailed audit log."

# Set ARGs as ENV so that they can be used in ENTRYPOINT/CMD
ENV NAME=$NAME
ENV VERSION=$VERSION

# Copy the license file as per Legal requirement
COPY ${LICENSE_SOURCE} ${LICENSE_DEST}/

# We must have a copy of the license in this directory to comply with the HasLicense Redhat requirement
# Note the trailing slash on the first argument -- plain files meet the requirement but directories do not.
COPY ${LICENSE_SOURCE}/ /licenses/

# Set up certificates, our base tools, and Vault. Unlike the other version of
# this (https://github.com/hashicorp/docker-vault/blob/master/ubi/Dockerfile),
# we copy in the Vault binary from CRT.
RUN set -eux; \
    microdnf install -y ca-certificates gnupg openssl libcap tzdata procps shadow-utils util-linux tar

# Create a non-root user to run the software.
RUN groupadd --gid 1000 vault && \
    adduser --uid 100 --system -g vault vault && \
    usermod -a -G root vault

# Copy in the new Vault from CRT pipeline, rather than fetching it from our
# public releases.
COPY dist/$TARGETOS/$TARGETARCH/$BIN_NAME /bin/

# /vault/logs is made available to use as a location to store audit logs, if
# desired; /vault/file is made available to use as a location with the file
# storage backend, if desired; the server will be started with /vault/config as
# the configuration directory so you can add additional config files in that
# location.
ENV HOME=/home/vault
RUN mkdir -p /vault/logs && \
    mkdir -p /vault/file && \
    mkdir -p /vault/config && \
    mkdir -p $HOME && \
    chown -R vault /vault && chown -R vault $HOME && \
    chgrp -R 0 $HOME && chmod -R g+rwX $HOME && \
    chgrp -R 0 /vault && chmod -R g+rwX /vault

# Expose the logs directory as a volume since there's potentially long-running
# state in there
VOLUME /vault/logs

# Expose the file directory as a volume since there's potentially long-running
# state in there
VOLUME /vault/file

# 8200/tcp is the primary interface that applications use to interact with
# Vault.
EXPOSE 8200

# The entry point script uses dumb-init as the top-level process to reap any
# zombie processes created by Vault sub-processes.
#
# For production derivatives of this container, you should add the IPC_LOCK
# capability so that Vault can mlock memory.
COPY .release/docker/ubi-docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

# Use the Vault user as the default user for starting this container.
USER vault

# # By default you'll get a single-node development server that stores everything
# # in RAM and bootstraps itself. Don't use this configuration for production.
CMD ["server", "-dev"]

FROM ubi AS ubi-fips

FROM ubi AS ubi-hsm

FROM ubi AS ubi-hsm-fips

## Builder:
#
#  A build container used to build the Vault binary. We use focal because the
#  version of glibc is old enough for all of our supported distros for editions
#  that require CGO. This container is used in CI to build all binaries that
#  require CGO.
#
#  To run it locally, first build the builder container:
#    docker build -t builder --build-arg GO_VERSION=$(cat .go-version) .
#
#  Then build Vault using the builder container:
#    docker run -it -v $(pwd):/build -v GITHUB_TOKEN=$GITHUB_TOKEN --env GO_TAGS='ui enterprise cgo hsm venthsm' --env GOARCH=s390x --env GOOS=linux --env VERSION=1.20.0-beta1 --env VERSION_METADATA=ent.hsm --env CGO_ENABLED=1 builder make ci-build
#
#  You can also share your local Go modules with the container to avoid downloading
#  them every time:
#    docker run -it -v $(pwd):/build -v $(go env GOMODCACHE):/go-mod-cache --env GITHUB_TOKEN=$GITHUB_TOKEN --env GO_TAGS='ui enterprise cgo hsm venthsm' --env GOARCH=s390x --env GOOS=linux --env VERSION=1.20.0-beta1 --env VERSION_METADATA=ent.hsm --env GOMODCACHE=/go-mod-cache --env CGO_ENABLED=1 builder make ci-build
#
#  If you have a linux machine you can also share the tools
#    GOBIN="$(go env GOPATH)/bin" make tools
#    docker run -it -v $(pwd):/build -v $(go env GOMODCACHE):/go-mod-cache -v "$(go env GOPATH)/bin":/opt/tools/bin --env GITHUB_TOKEN=$GITHUB_TOKEN --env GO_TAGS='ui enterprise cgo hsm venthsm' --env GOARCH=s390x --env GOOS=linux --env VERSION=1.20.0-beta1 --env VERSION_METADATA=ent.hsm --env GOMODCACHE=/go-mod-cache --env CGO_ENABLED=1 builder make ci-build
FROM ubuntu:focal AS builder

# Pass in the GO_VERSION as a build-arg
ARG GO_VERSION

# Set our environment
ENV PATH="/root/go/bin:/opt/go/bin:/opt/tools/bin:$PATH"
ENV GOPRIVATE='github.com/hashicorp/*'

# Install the necessary system tooling to cross compile vault for our various
# CGO targets. Do this separately from branch specific Go and build toolchains
# so our various builder image layers can share cache.
COPY .build/system.sh .
RUN chmod +x system.sh && ./system.sh && rm -rf system.sh

# Install the correct Go toolchain
COPY .build/go.sh .
RUN chmod +x go.sh && ./go.sh && rm -rf go.sh

# Install the vault tools installer. It might be required during build if the
# pre-build tools are not mounted into the container.
COPY tools/tools.sh .
RUN chmod +x tools.sh

# Run the build
COPY .build/entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
