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
[[ -z "$TEST_DIR" ]] && fail "TEST_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"
export VAULT_FORMAT=json

# ------ Generate and sign certificate ------
CA_NAME="${MOUNT}-ca.pem"
ISSUED_CERT_NAME="${MOUNT}-issued.pem"
ROLE_NAME="${COMMON_NAME}-role"
SUBJECT="test.${COMMON_NAME}"
TMP_TTL="1h"
rm -rf "${TEST_DIR}"
mkdir "${TEST_DIR}"

## Setting AIA fields for Certificate
"$binpath" write "${MOUNT}/config/urls" issuing_certificates="${VAULT_ADDR}/v1/pki/ca" crl_distribution_points="${VAULT_ADDR}/v1/pki/crl"

# Generating CA Certificate
"$binpath" write "${MOUNT}/root/generate/internal" common_name="${COMMON_NAME}.com" issuer_name="${ISSUER_NAME}" ttl="${TTL}" | jq -r '.data.issuing_ca' > "${TEST_DIR}/${CA_NAME}"
# Creating a role
"$binpath" write "${MOUNT}/roles/${ROLE_NAME}" allowed_domains="${COMMON_NAME}.com" allow_subdomains=true max_ttl="${TMP_TTL}"
# Issuing Signed Certificate
"$binpath" write "${MOUNT}/issue/${ROLE_NAME}" common_name="${SUBJECT}.com" ttl="${TMP_TTL}" | jq -r '.data.certificate' > "${TEST_DIR}/${ISSUED_CERT_NAME}"

# ------ Generate and sign intermediate ------
INTERMEDIATE_COMMON_NAME="intermediate-${COMMON_NAME}"
INTERMEDIATE_ISSUER_NAME="intermediate-${ISSUER_NAME}"
INTERMEDIATE_ROLE_NAME="intermediate-${COMMON_NAME}-role"
INTERMEDIATE_CA_NAME="${MOUNT}-${INTERMEDIATE_COMMON_NAME}.pem"
INTERMEDIATE_SIGNED_NAME="${MOUNT}-${INTERMEDIATE_COMMON_NAME}-ca.pem"
INTERMEDIATE_ISSUED_NAME="${MOUNT}-${INTERMEDIATE_COMMON_NAME}-issued.pem"

# Generate Intermediate CSR
"$binpath" write "${MOUNT}/intermediate/generate/internal" common_name="${INTERMEDIATE_COMMON_NAME}.com" issuer_name="${INTERMEDIATE_ISSUER_NAME}" ttl="${TTL}" | jq -r '.data.csr' > "${TEST_DIR}/${INTERMEDIATE_CA_NAME}"
# Creating a intermediate role
"$binpath" write "${MOUNT}/roles/${INTERMEDIATE_ROLE_NAME}" allowed_domains="${INTERMEDIATE_COMMON_NAME}.com" allow_subdomains=true max_ttl="${TMP_TTL}"
# Sign Intermediate Certificate
"$binpath" write "${MOUNT}/root/sign-intermediate" csr="@${TEST_DIR}/${INTERMEDIATE_CA_NAME}" format=pem_bundle ttl="${TMP_TTL}" | jq -r '.data.certificate' > "${TEST_DIR}/${INTERMEDIATE_SIGNED_NAME}"
# Import Signed Intermediate Certificate into Vault
"$binpath" write "${MOUNT}/intermediate/set-signed" certificate="@${TEST_DIR}/${INTERMEDIATE_SIGNED_NAME}"
# Issuing Signed Certificate with the intermediate role
"$binpath" write "${MOUNT}/issue/${INTERMEDIATE_ROLE_NAME}" common_name="www.${INTERMEDIATE_COMMON_NAME}.com" ttl="${TMP_TTL}" | jq -r '.data.certificate' > "${TEST_DIR}/${INTERMEDIATE_ISSUED_NAME}"
