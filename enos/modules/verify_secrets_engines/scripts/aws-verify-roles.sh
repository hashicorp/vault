#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

## # -------PKI TESTING
# MOUNT=aws
# AWS_ROLE=test-role
# VAULT_ADDR=http://127.0.0.1:8200
# VAULT_INSTALL_DIR=/opt/homebrew/bin
# VAULT_TOKEN=root
# vault secrets enable --path=${MOUNT} aws > /dev/null 2>&1  || echo "AWS Engine already enabled!"
echo "------------${AWS_REGION}-----------${AWS_ACCESS_KEY_ID}"

[[ -z "$AWS_REGION" ]] && fail "AWS_REGION env variable has not been set"
[[ -z "$AWS_ACCESS_KEY_ID" ]] && fail "AWS_ACCESS_KEY_ID env variable has not been set"
[[ -z "$AWS_SECRET_ACCESS_KEY" ]] && fail "AWS_SECRET_ACCESS_KEY env variable has not been set"
[[ -z "$AWS_ROLE" ]] && fail "AWS_ROLE env variable has not been set"
[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

echo "Verifying roles list"
ROLE=$("$binpath" list "${MOUNT}/roles" | jq -r '.[]')
[[ -z "$ROLE" ]] && fail "No AWS roles created!"

echo "Verifying Root Access Key"
"$binpath" read "${MOUNT}/config/root" | jq -r '.data.access_key'
ROOT_ACCESS_KEY=$("$binpath" read "${MOUNT}/config/root" | jq -r '.data.access_key')
echo "----------------${ROOT_ACCESS_KEY}---------${AWS_ACCESS_KEY_ID}"
[[ "$ROOT_ACCESS_KEY" == "$AWS_ACCESS_KEY_ID" ]] && fail "AWS Access Key does not match: $ROOT_ACCESS_KEY, $AWS_ACCESS_KEY_ID"

# Read role
"$binpath" read "${MOUNT}/sts/creds/${AWS_ROLE}"
