#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

#MOUNT=pki_secret
#ISSUER=issuer
#COMMON_NAME=common
#TTL=8760h
#VAULT_ADDR=http://127.0.0.1:8200
#VAULT_INSTALL_DIR=/opt/homebrew/bin
#VAULT_TOKEN=root
#TMP_TEST_RESULTS="pki_tmp_results"
#vault secrets enable --path=${MOUNT} pki > /dev/null 2>&1  || echo "PKI already enabled!"

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$COMMON_NAME" ]] && fail "COMMON_NAME env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$TMP_TEST_RESULTS" ]] && fail "TMP_TEST_RESULTS env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath" || fail "The certificate appears to be improperly configured or contains errors"

export VAULT_FORMAT=json

# Validate cert details:
SIGNED_CRT_NAME="${MOUNT}_signed.pem"
openssl x509 -in "${TMP_TEST_RESULTS}/${SIGNED_CRT_NAME}" -text -noout || fail "The certificate appears to be improperly configured or contains errors"

# Validate intermediate cert details:
INTERMEDIATE_COMMON_NAME="intermediate_${COMMON_NAME}"
INTERMEDIATE_SIGNED_CRT_NAME="${MOUNT}_${INTERMEDIATE_COMMON_NAME}_signed.crt"
openssl x509 -in "${TMP_TEST_RESULTS}/${INTERMEDIATE_SIGNED_CRT_NAME}" -text -noout || fail "The intermediate certificate appears to be improperly configured or contains errors"









