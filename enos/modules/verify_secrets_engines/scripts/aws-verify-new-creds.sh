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
[[ -z "$VAULT_AWS_ROLE" ]] && fail "VAULT_AWS_ROLE env variable has not been set"
[[ -z "$AWS_REGION" ]] && fail "AWS_REGION env variable has not been set"
[[ -z "$AWS_USER_NAME" ]] && fail "AWS_USER_NAME env variable has not been set"
[[ -z "$AWS_ACCESS_KEY_ID" ]] && fail "AWS_ACCESS_KEY_ID env variable has not been set"
[[ -z "$AWS_SECRET_ACCESS_KEY" ]] && fail "AWS_SECRET_ACCESS_KEY env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

echo -e "Configuring Vault AWS \n"
USERNAME_TEMPLATE="{{ if (eq .Type \"STS\") }}{{ printf \"${AWS_USER_NAME}-%s-%s\" (random 20) (unix_time) | truncate 32 }}{{ else }}{{ printf \"${AWS_USER_NAME}-%s-%s\" (unix_time) (random 20) | truncate 60 }}{{ end }}"
"$binpath" write "${MOUNT}/config/root" access_key="${AWS_ACCESS_KEY_ID}" secret_key="${AWS_SECRET_ACCESS_KEY}" region="${AWS_REGION}" username_template="${USERNAME_TEMPLATE}"

echo -e "Verifying root config \n"
"$binpath" read "${MOUNT}/config/root"
ROOT_USERNAME_TEMPLATE=$("$binpath" read "${MOUNT}/config/root" | jq -r '.data.username_template')
[[ "$ROOT_USERNAME_TEMPLATE" == *"$AWS_USER_NAME"* ]] || fail "Uername Template does not include the current role"

echo -e "Verifying roles list \n"
"$binpath" list "${MOUNT}/roles"
ROLE=$("$binpath" list "${MOUNT}/roles" | jq -r '.[]')
[[ -z "$ROLE" ]] && fail "No AWS roles created!"

echo -e "Generate New Credentials \n"
TEMP_IAM_USER=$("$binpath" read "${MOUNT}/creds/${VAULT_AWS_ROLE}") || fail "Failed to generate new credentials for iam user: ${VAULT_AWS_ROLE}"
TEMP_ACCESS_KEY=$(echo ${TEMP_IAM_USER} | jq -r '.data.access_key') || fail "Failed to get access key from: ${VAULT_AWS_ROLE}"
if [[ -z "$TEMP_ACCESS_KEY" && "$TEMP_ACCESS_KEY" != "$AWS_USER_NAME" ]]; then
  failed "The new access key is empty or is matching the old one: ${TEMP_ACCESS_KEY}"
fi
