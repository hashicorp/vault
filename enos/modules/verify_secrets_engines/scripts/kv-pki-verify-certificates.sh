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
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$COMMON_NAME" ]] && fail "COMMON_NAME env variable has not been set"
[[ -z "$ISSUER_NAME" ]] && fail "ISSUER_NAME env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$TMP_TEST_RESULTS" ]] && fail "TMP_TEST_RESULTS env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath" || fail "The certificate appears to be improperly configured or contains errors"
export VAULT_FORMAT=json

# Verifying List Roles
ROLE=$("$binpath" list -format=json "${MOUNT}/roles" | jq -r '.[]')
[[ -z "$ROLE" ]] && fail "No roles created!"

# Verifying List Issuer
ISSUER=$("$binpath" list -format=json "${MOUNT}/issuers" | jq -r '.[]')
[[ -z "$ISSUER" ]] && fail "No issuers created!"

# Verifying Root CA Certificate
ROOT_CA_CERT=$("$binpath" read -format=json pki/cert/ca | jq -r '.data.certificate')
[[ -z "$ROOT_CA_CERT" ]] && fail "No root ca certificate generated"

# Verify List Certificate
VAULT_CERTS=$("$binpath" list -format=json "${MOUNT}/certs" | jq -r '.[]')
[[ -z "$VAULT_CERTS" ]] && fail "VAULT_CERTS should include vault certificates"

# Verifying Certificates
for CERT in $VAULT_CERTS; do
  echo "Getting Certificate from Vault PKI: ${CERT}"
  "$binpath" read "${MOUNT}/cert/${CERT}" | jq -r '.data.certificate' > "${TMP_TEST_RESULTS}/tmp_vault_cert.pem"
  echo "Verifying Certificate..."
  openssl x509 -in "${TMP_TEST_RESULTS}/tmp_vault_cert.pem" -text -noout || fail "The certificate appears to be improperly configured or contains errors"
  echo "Successfully Verified Certificate"

  IS_CA=$(openssl x509 -in "${TMP_TEST_RESULTS}/tmp_vault_cert.pem" -text -noout | grep -q "CA:TRUE" && echo "TRUE" || echo "FALSE")
  if [[ "${IS_CA}" == "FALSE" ]]; then
    echo "Revoking Certificate: ${CERT}"
    "$binpath" write "${MOUNT}/revoke" serial_number="${CERT}" || fail "Could not revoke certificate ${CERT}"
  else
    echo "Skipping revoking step for this certificate to being a root CA Cert: ${CERT}"
  fi
done

# Verify List Revoked Certificate
"$binpath" list -format=json "${MOUNT}/certs/revoked" | jq -r '.[]' || fail "There are no revoked certificate listed"