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
TMP_FILE="tmp-vault-cert.pem"
REVOKED_CERTS=()
for CERT in $VAULT_CERTS; do
  echo "Getting certificate from Vault PKI: ${CERT}"
  "$binpath" read "${MOUNT}/cert/${CERT}" | jq -r '.data.certificate' > "${TMP_TEST_RESULTS}/${TMP_FILE}"
  echo "Verifying certificate..."
  openssl x509 -in "${TMP_TEST_RESULTS}/${TMP_FILE}" -text -noout || fail "The certificate appears to be improperly configured or contains errors"
  CURR_CERT_SERIAL=$(echo "${CERT}" | tr -d ':' | tr '[:lower:]' '[:upper:]')
  TMP_CERT_SUBJECT=$(openssl x509 -in "${TMP_TEST_RESULTS}/${TMP_FILE}" -noout -subject)
  TMP_CERT_ISSUER=$(openssl x509 -in "${TMP_TEST_RESULTS}/${TMP_FILE}" -noout -issuer)
  TMP_CERT_SERIAL=$(openssl x509 -in "${TMP_TEST_RESULTS}/${TMP_FILE}" -noout -serial)
  [[ "${TMP_CERT_SUBJECT}" == *"${COMMON_NAME}.com"* ]] || fail "Subject is incorrect. Actual Subject: ${TMP_CERT_SUBJECT}"
  [[ "${TMP_CERT_ISSUER}" == *"${COMMON_NAME}.com"* ]] || fail "Issuer is incorrect. Actual Issuer: ${TMP_CERT_ISSUER}"
  [[ "${TMP_CERT_SERIAL}" == *"${CURR_CERT_SERIAL}"* ]] || fail "Certificate Serial is incorrect. Actual certificate Serial: ${CURR_CERT_SERIAL},${TMP_CERT_SERIAL}"
  echo "Certificate successfully verified"

  IS_CA=$(openssl x509 -in "${TMP_TEST_RESULTS}/${TMP_FILE}" -text -noout | grep -q "CA:TRUE" && echo "TRUE" || echo "FALSE")
  if [[ "${IS_CA}" == "FALSE" ]]; then
    echo "Revoking certificate: ${CERT}"
    "$binpath" write "${MOUNT}/revoke" serial_number="${CERT}" || fail "Could not revoke certificate ${CERT}"
    REVOKED_CERTS+=("$CERT")
  else
    echo "Skipping revoking step for this certificate to being a root CA Cert: ${CERT}"
  fi
done

echo "Verifying Revoked Certificates"
REVOKED_CERT_FROM_LIST=$("$binpath" list -format=json "${MOUNT}/certs/revoked" | jq -r '.[]')
[[ -z "$REVOKED_CERT_FROM_LIST" ]] && fail "No revoked certificates are listed."
for CERT in "${REVOKED_CERTS[@]}"; do
  [[ "${REVOKED_CERT_FROM_LIST}" == *"${CERT}"* ]] || fail "Unable to locate certificate in the Vault Revoked Certificate List: ${CERT}"
done
echo "Revoked certificate successfully verified"