#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$AUTH_PATH" ]] && fail "AUTH_PATH env variable has not been set"
[[ -z "$USERNAME" ]] && fail "USERNAME env variable has not been set"
[[ -z "$PASSWORD" ]] && fail "PASSWORD env variable has not been set"
[[ -z "$VERIFY_PKI_CERTS" ]] && fail "VERIFY_CERT_DETAILS env variable has not been set"
[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$COMMON_NAME" ]] && fail "COMMON_NAME env variable has not been set"
[[ -z "$ISSUER_NAME" ]] && fail "ISSUER_NAME env variable has not been set"
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$TEST_DIR" ]] && fail "TEST_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath" || fail "The certificate appears to be improperly configured or contains errors"
export VAULT_FORMAT=json

# Log in so this vault instance have access to the primary pki roles, issuers, and etc
if [ "${VERIFY_PKI_CERTS}" = false ]; then
  echo "Logging in Vault with username and password: ${USERNAME}"
  VAULT_TOKEN=$("$binpath" write "auth/$AUTH_PATH/login/$USERNAME" password="$PASSWORD" | jq -r '.auth.client_token')
fi

# Verifying List Roles
ROLE=$("$binpath" list "${MOUNT}/roles" | jq -r '.[]')
[[ -z "$ROLE" ]] && fail "No roles created!"

# Verifying List Issuer
ISSUER=$("$binpath" list "${MOUNT}/issuers" | jq -r '.[]')
[[ -z "$ISSUER" ]] && fail "No issuers created!"

# Verifying Root CA Certificate
ROOT_CA_CERT=$("$binpath" read pki/cert/ca | jq -r '.data.certificate')
[[ -z "$ROOT_CA_CERT" ]] && fail "No root ca certificate generated"

# Verifying Certificates
if [ "${VERIFY_PKI_CERTS}" = true ]; then
  if [ ! -d "${TEST_DIR}" ]; then
      echo "Directory does not exist. Creating it now."
      mkdir -p "${TEST_DIR}" # Need to create this directory for Enterprise test
  fi
  TMP_FILE="tmp-vault-cert.pem"

  # Verify List Certificate
  VAULT_CERTS=$("$binpath" list "${MOUNT}/certs" | jq -r '.[]')
  [[ -z "$VAULT_CERTS" ]] && fail "VAULT_CERTS should include vault certificates"
  for CERT in $VAULT_CERTS; do
    echo "Getting certificate from Vault PKI: ${CERT}"
    "$binpath" read "${MOUNT}/cert/${CERT}" | jq -r '.data.certificate' > "${TEST_DIR}/${TMP_FILE}"
    echo "Verifying certificate contents..."
    openssl x509 -in "${TEST_DIR}/${TMP_FILE}" -text -noout || fail "The certificate appears to be improperly configured or contains errors"
    CURR_CERT_SERIAL=$(echo "${CERT}" | tr -d ':' | tr '[:lower:]' '[:upper:]')
    TMP_CERT_SUBJECT=$(openssl x509 -in "${TEST_DIR}/${TMP_FILE}" -noout -subject | awk -F'= ' '{print $2}')
    TMP_CERT_ISSUER=$(openssl x509 -in "${TEST_DIR}/${TMP_FILE}" -noout -issuer | awk -F'= ' '{print $2}')
    TMP_CERT_SERIAL=$(openssl x509 -in "${TEST_DIR}/${TMP_FILE}" -noout -serial | awk -F'=' '{print $2}')
    [[ "${TMP_CERT_SUBJECT}" == *"${COMMON_NAME}.com"* ]] || fail "Subject is incorrect. Actual Subject: ${TMP_CERT_SUBJECT}"
    [[ "${TMP_CERT_ISSUER}" == *"${COMMON_NAME}.com"* ]] || fail "Issuer is incorrect. Actual Issuer: ${TMP_CERT_ISSUER}"
    [[ "${TMP_CERT_SERIAL}" == *"${CURR_CERT_SERIAL}"* ]] || fail "Certificate Serial is incorrect. Actual certificate Serial: ${CURR_CERT_SERIAL},${TMP_CERT_SERIAL}"
    echo "Successfully verified certificate contents."

    # Setting up variables for types of certificates
    IS_CA=$(openssl x509 -in "${TEST_DIR}/${TMP_FILE}" -text -noout | grep -q "CA:TRUE" && echo "TRUE" || echo "FALSE")
    if [[ "${IS_CA}" == "TRUE" ]]; then
      if [[ "${COMMON_NAME}.com" == "${TMP_CERT_SUBJECT}" ]]; then
        CA_CERT=${CERT}
      elif [[ "intermediate-${COMMON_NAME}.com" == "${TMP_CERT_SUBJECT}" ]]; then
        INTERMEDIATE_CA_CERT=${CERT}
      fi
    elif [[ "${IS_CA}" == "FALSE" ]]; then
      INTERMEDIATE_ISSUED_CERT=${CERT}
    fi

  done

  echo "Verifying that Vault PKI has successfully generated valid certificates for the CA, Intermediate CA, and issued certificates..."
  if [[ -n "${CA_CERT}" ]] && [[ -n "${INTERMEDIATE_CA_CERT}" ]] && [[ -n "${INTERMEDIATE_ISSUED_CERT}" ]]; then
    CA_NAME="ca.pem"
    INTERMEDIATE_CA_NAME="intermediate-ca.pem"
    ISSUED_NAME="issued.pem"
    "$binpath" read "${MOUNT}/cert/${CA_CERT}" | jq -r '.data.certificate' > "${TEST_DIR}/${CA_NAME}"
    "$binpath" read "${MOUNT}/cert/${INTERMEDIATE_CA_CERT}" | jq -r '.data.certificate' > "${TEST_DIR}/${INTERMEDIATE_CA_NAME}"
    "$binpath" read "${MOUNT}/cert/${INTERMEDIATE_ISSUED_CERT}" | jq -r '.data.certificate' > "${TEST_DIR}/${ISSUED_NAME}"
    openssl verify --CAfile "${TEST_DIR}/${CA_NAME}" -untrusted "${TEST_DIR}/${INTERMEDIATE_CA_NAME}" "${TEST_DIR}/${ISSUED_NAME}" || fail "One or more Certificate is not valid."
  else
    echo "CA Cert: ${CA_CERT}, Intermedidate Cert: ${INTERMEDIATE_CA_CERT}, Issued Cert: ${INTERMEDIATE_ISSUED_CERT}"
  fi

  echo "Revoking certificate: ${INTERMEDIATE_ISSUED_CERT}"
  "$binpath" write "${MOUNT}/revoke" serial_number="${INTERMEDIATE_ISSUED_CERT}" || fail "Could not revoke certificate ${INTERMEDIATE_ISSUED_CERT}"
  echo "Verifying Revoked Certificate"
  REVOKED_CERT_FROM_LIST=$("$binpath" list "${MOUNT}/certs/revoked" | jq -r '.[0]')
  [[ "${INTERMEDIATE_ISSUED_CERT}" == "${REVOKED_CERT_FROM_LIST}" ]] || fail "Expected: ${INTERMEDIATE_ISSUED_CERT}, actual: ${REVOKED_CERT_FROM_LIST}"
  echo "Successfully verified revoked certificate"
else
  echo "Skipping verify certificates!"
fi
