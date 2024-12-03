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
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

# ------ Generate and sign certificate ------
CRT_NAME="${MOUNT}.crt"
CSR_NAME="${MOUNT}.csr"
PRIV_KEY_NAME="${MOUNT}_priv.key"
SIGNED_CRT_NAME="${MOUNT}_signed.pem"
ROLE_NAME="${COMMON_NAME}-role"
mkdir "${TMP_TEST_RESULTS}"

# Generating root CA.crt
"$binpath" write ${MOUNT}/root/generate/internal common_name="${COMMON_NAME}.com" ttl="${TTL}" -format=json | jq -r '.data.certificate' > "${TMP_TEST_RESULTS}/${CRT_NAME}"
# Creating a role
"$binpath" write ${MOUNT}/roles/${ROLE_NAME} allowed_domains="${COMMON_NAME}.com" allow_subdomains=true max_ttl="${TTL+5}"
# Issue Certificate
openssl req -new -newkey rsa:2048 -nodes -subj "/CN=www.${COMMON_NAME}.com" -keyout "${TMP_TEST_RESULTS}/${PRIV_KEY_NAME}" -out "${TMP_TEST_RESULTS}/${CSR_NAME}"
# Sign Certificate
"$binpath" write ${MOUNT}/sign/${ROLE_NAME} csr="@${TMP_TEST_RESULTS}/${CSR_NAME}" format=pem ttl="${TTL+5}" | jq -r '.data.certificate' > "${TMP_TEST_RESULTS}/${SIGNED_CRT_NAME}"

# ------ Generate and sign intermediate ------
INTERMEDIATE_COMMON_NAME="intermediate_${COMMON_NAME}"
INTERMEDIATE_CSR_NAME="${MOUNT}_${INTERMEDIATE_COMMON_NAME}.csr"
INTERMEDIATE_SIGNED_CRT_NAME="${MOUNT}_${INTERMEDIATE_COMMON_NAME}_signed.crt"
# Setting AIA fields for Certificate
"$binpath" write ${MOUNT}/config/urls issuing_certificates="${VAULT_ADDR}/v1/pki/ca" crl_distribution_points="${VAULT_ADDR}/v1/pki/crl"
# Generate Intermediate Certificate
"$binpath" write ${MOUNT}/intermediate/generate/internal common_name="${INTERMEDIATE_COMMON_NAME}.com" ttl="${TTL}" | jq -r '.data.csr' > "${TMP_TEST_RESULTS}/${INTERMEDIATE_CSR_NAME}"
# Sign Intermediate Certificate
"$binpath" write ${MOUNT}/root/sign-intermediate csr="@${TMP_TEST_RESULTS}/${INTERMEDIATE_CSR_NAME}" format=pem_bundle ttl="${TTL}" | jq -r '.data.certificate' > "${TMP_TEST_RESULTS}/${INTERMEDIATE_SIGNED_CRT_NAME}"
# Import Signed Intermediate Certificate into Vault
"$binpath" write ${MOUNT}/intermediate/set-signed certificate="@${TMP_TEST_RESULTS}/${INTERMEDIATE_SIGNED_CRT_NAME}"










