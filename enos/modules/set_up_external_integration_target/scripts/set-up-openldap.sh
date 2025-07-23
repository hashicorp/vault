#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$LDAP_DOMAIN" ]] && fail "LDAP_DOMAIN env variable has not been set"
[[ -z "$LDAP_ORG" ]] && fail "LDAP_ORG env variable has not been set"
[[ -z "$LDAP_ADMIN_PW" ]] && fail "LDAP_ADMIN_PW env variable has not been set"
[[ -z "$LDAP_CONTAINER_VERSION" ]] && fail "LDAP_CONTAINER_VERSION env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAPS_PORT" ]] && fail "LDAPS_PORT env variable has not been set"

# Pulling image
CONTAINER_CMD="sudo podman"
LDAP_DOCKER_NAME="docker.io/osixia/openldap:${LDAP_CONTAINER_VERSION}"
echo "Pulling image: ${LDAP_DOCKER_NAME}"
$CONTAINER_CMD pull "${LDAP_DOCKER_NAME}"

# Run OpenLDAP container
echo "Starting OpenLDAP container..."
$CONTAINER_CMD run -d \
  --name openldap \
  -p "${LDAP_PORT}:${LDAP_PORT}" \
  -p "${LDAPS_PORT}:${LDAPS_PORT}" \
  -e LDAP_ORGANISATION="${LDAP_ORG}" \
  -e LDAP_DOMAIN="${LDAP_DOMAIN}" \
  -e LDAP_ADMIN_PASSWORD="${LDAP_ADMIN_PW}" \
  "${LDAP_DOCKER_NAME}"

echo "OpenLDAP server is now running in Docker!"
