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
# AWS_REGION=us-east-1
# AWS_ROLE=test-role
# VAULT_ADDR=http://127.0.0.1:8200
# VAULT_INSTALL_DIR=/opt/homebrew/bin
# VAULT_TOKEN=root
# vault secrets enable --path=${MOUNT} aws > /dev/null 2>&1  || echo "AWS Engine already enabled!"
echo -e "------------|${AWS_REGION}|-----------|${AWS_ACCESS_KEY_ID}|--------\n"
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

echo "Configuring Vault AWS"
"$binpath" write "${MOUNT}/config/root" access_key="${AWS_ACCESS_KEY_ID}" secret_key="${AWS_SECRET_ACCESS_KEY}" region=${AWS_REGION} || fail "Cannot set vault AWS credentials"

echo "Creating AWS Role"
#"$binpath" write "${MOUNT}/roles/${AWS_ROLE}" credential_type=iam_user policy_arns="arn:aws:iam::aws:policy/AdministratorAccess" ttl="1h" max_ttl="24h" || fail "Cannot create AWS role"
"$binpath" write "aws/roles/${AWS_ROLE}" \
    credential_type=iam_user \
    ttl="1h" max_ttl="24h" \
    policy_document=-<<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ec2:*",
      "Resource": "*"
    }
  ]
}
EOF

echo "Verifying roles list"
ROLE=$("$binpath" list "${MOUNT}/roles" | jq -r '.[]')
[[ -z "$ROLE" ]] && fail "No AWS roles created!"

echo "Verifying Root Access Key"
"$binpath" read "${MOUNT}/config/root" | jq -r '.data.access_key'
ROOT_ACCESS_KEY=$("$binpath" read "${MOUNT}/config/root" | jq -r '.data.access_key')
[[ "$ROOT_ACCESS_KEY" != "$AWS_ACCESS_KEY_ID" ]] && fail "AWS Access Key does not match: $ROOT_ACCESS_KEY, $AWS_ACCESS_KEY_ID"

# Read role
"$binpath" read "${MOUNT}/creds/${AWS_ROLE}"