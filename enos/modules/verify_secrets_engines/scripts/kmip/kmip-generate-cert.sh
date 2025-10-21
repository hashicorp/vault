#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${MOUNT}" ]] && fail "MOUNT env variable has not been set"
[[ -z "${VAULT_ADDR}" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "${VAULT_INSTALL_DIR}" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "${VAULT_TOKEN}" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "${SCOPE_NAME}" ]] && fail "SCOPE_NAME env variable has not been set"
[[ -z "${ROLE_NAME}" ]] && fail "ROLE_NAME env variable has not been set"
[[ -z "${CERT_FORMAT}" ]] && fail "CERT_FORMAT env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "${binpath}" || fail "unable to locate vault binary at ${binpath}"

export VAULT_FORMAT=json

"${binpath}" write -format=json \
    "${MOUNT}"/scope/"${SCOPE_NAME}"/role/"${ROLE_NAME}"/credential/generate \
    format="${CERT_FORMAT}" > credential.json

jq -r .data.certificate < credential.json > cert.pem
jq -r .data.private_key < credential.json > key.pem

cat cert.pem key.pem > client.pem
cat client.pem
