#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

MOUNT=pki_secret
ISSUER=issuer
COMMON_NAME=common
TTL=8760h
VAULT_ADDR=http://127.0.0.1:8200
VAULT_INSTALL_DIR=/opt/homebrew/bin
VAULT_TOKEN=root

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$COMMON_NAME" ]] && fail "COMMON_NAME env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

# Generating root CA.crt
"$binpath" write ${MOUNT}/root/generate/internal common_name="${COMMON_NAME}.com" ttl="${TTL}" -format=json | jq -r '.data.certificate' > ${MOUNT}.crt

# Creating a role
"$binpath" write ${MOUNT}/roles/${COMMON_NAME}-dot-com allowed_domains="${COMMON_NAME}.com" allow_subdomains=true max_ttl="72h"

# Issue Certificate
openssl req -new -newkey rsa:2048 -nodes -subj "/CN=www.${COMMON_NAME}.com" -keyout ${MOUNT}_private_key.key -out ${MOUNT}.csr

# Sign Certificate
"$binpath" write ${MOUNT}/sign/${COMMON_NAME}-dot-com csr="@${MOUNT}.csr" format=pem ttl=24h > ${MOUNT}_signed.crt








