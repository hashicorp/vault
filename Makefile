# Base image
FROM alpine:3 AS default

ARG BIN_NAME
ARG NAME=vault
ARG PRODUCT_VERSION
ARG TARGETOS TARGETARCH

LABEL name="Vault" \
      maintainer="Vault Team <vault@hashicorp.com>" \
      vendor="HashiCorp" \
      version=${PRODUCT_VERSION} \
      description="Vault is a tool for securely accessing secrets."

# Add non-root user
RUN addgroup ${NAME} && adduser -S -G ${NAME} ${NAME}

# Install necessary tools
RUN apk add --no-cache libcap dumb-init tzdata curl && \
    mkdir -p /vault/logs /vault/file /vault/config && \
    chown -R ${NAME}:${NAME} /vault

# Copy Vault binary
COPY dist/$TARGETOS/$TARGETARCH/$BIN_NAME /bin/

# Set entrypoint
COPY .release/docker/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

# Expose ports and volumes
VOLUME /vault/logs /vault/file
EXPOSE 8200

# Default to dev server for QA. Use production configurations when appropriate.
CMD ["server", "-dev"]

# QA-Specific Recommendations
# For QA testing, inject a non-sensitive token for Vault operations:
ENV VAULT_DEV_ROOT_TOKEN_ID="qa-testing-token"

# Simulate real configurations for integration testing:
# Uncomment the following lines for QA:
# COPY test-config.hcl /vault/config/config.hcl
# CMD ["server", "-config=/vault/config/config.hcl"]
