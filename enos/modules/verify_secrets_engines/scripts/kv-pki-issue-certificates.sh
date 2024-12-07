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
[[ -z "$TTL" ]] && fail "TTL env variable has not been set"
[[ -z "$TMP_TEST_RESULTS" ]] && fail "TMP_TEST_RESULTS env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"
export VAULT_FORMAT=json

# ------ Generate and sign certificate ------
CA_NAME="${MOUNT}.pem"
CA_CSR_NAME="${MOUNT}-csr.pem"
PRIV_KEY_NAME="${MOUNT}_priv.key"
SIGNED_CA_NAME="${MOUNT}_signed.pem"
ROLE_NAME="${COMMON_NAME}-role"
rm -rf "${TMP_TEST_RESULTS}"
mkdir "${TMP_TEST_RESULTS}"

## Setting AIA fields for Certificate
"$binpath" write "${MOUNT}/config/urls" issuing_certificates="${VAULT_ADDR}/v1/pki/ca" crl_distribution_points="${VAULT_ADDR}/v1/pki/crl"

# Generating root CA
"$binpath" write "${MOUNT}/root/generate/internal" common_name="${COMMON_NAME}.com" ttl="${TTL}" -format=json | jq -r '.data.certificate' >"${TMP_TEST_RESULTS}/${CA_NAME}"
# Creating a role
"$binpath" write "${MOUNT}/roles/${ROLE_NAME}" allowed_domains="${COMMON_NAME}.com" allow_subdomains=true max_ttl="${TTL+5}"
# Issue Certificate
openssl req -new -newkey rsa:2048 -nodes -subj "/CN=www.${COMMON_NAME}.com" -keyout "${TMP_TEST_RESULTS}/${PRIV_KEY_NAME}" -out "${TMP_TEST_RESULTS}/${CA_CSR_NAME}"
# Sign Certificate
"$binpath" write "${MOUNT}/sign/${ROLE_NAME}" csr="@${TMP_TEST_RESULTS}/${CA_CSR_NAME}" format=pem ttl="${TTL+5}" | jq -r '.data.certificate' >"${TMP_TEST_RESULTS}/${SIGNED_CA_NAME}"

# ------ Generate and sign intermediate ------
INTERMEDIATE_COMMON_NAME="intermediate-${COMMON_NAME}"
INTERMEDIATE_CA_NAME="${MOUNT}_${INTERMEDIATE_COMMON_NAME}.pem"
INTERMEDIATE_SIGNED_CA_NAME="${MOUNT}_${INTERMEDIATE_COMMON_NAME}_signed.pem"

# Generate Intermediate CA
"$binpath" write "${MOUNT}/intermediate/generate/internal" common_name="${INTERMEDIATE_COMMON_NAME}.com" ttl="${TTL}" | jq -r '.data.csr' >"${TMP_TEST_RESULTS}/${INTERMEDIATE_CA_NAME}"
# Sign Intermediate Certificate
"$binpath" write "${MOUNT}/root/sign-intermediate" csr="@${TMP_TEST_RESULTS}/${INTERMEDIATE_CA_NAME}" format=pem_bundle ttl="${TTL}" | jq -r '.data.certificate' >"${TMP_TEST_RESULTS}/${INTERMEDIATE_SIGNED_CA_NAME}"
# Import Signed Intermediate Certificate into Vault
"$binpath" write "${MOUNT}/intermediate/set-signed" certificate="@${TMP_TEST_RESULTS}/${INTERMEDIATE_SIGNED_CA_NAME}"
