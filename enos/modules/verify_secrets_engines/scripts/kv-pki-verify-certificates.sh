#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
#[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$COMMON_NAME" ]] && fail "COMMON_NAME env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$TMP_TEST_RESULTS" ]] && fail "TMP_TEST_RESULTS env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath" || fail "The certificate appears to be improperly configured or contains errors"
export VAULT_FORMAT=json

# Getting Certificates
VAULT_CERTS=$("$binpath" list -format=json "${MOUNT}/certs" | jq -r '.[]')
[[ -z "$VAULT_CERTS" ]] && fail "VAULT_CERTS should include vault certificates"

# Verifying Certificates
for CERT in $VAULT_CERTS; do
  echo "Getting Certificate from Vault PKI: ${CERT}"
  "$binpath" read "${MOUNT}/cert/${CERT}" | jq -r '.data.certificate' > "${TMP_TEST_RESULTS}/tmp_vault_cert.pem"
  echo "Verifying Certificate..."
  openssl x509 -in "${TMP_TEST_RESULTS}/tmp_vault_cert.pem" -text -noout || fail "The certificate appears to be improperly configured or contains errors"
done
