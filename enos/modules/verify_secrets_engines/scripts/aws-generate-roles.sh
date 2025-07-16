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
[[ -z "$AWS_POLICY_ARN" ]] && fail "AWS_POLICY_ARN env variable has not been set"
[[ -z "$AWS_ROLE_ARN" ]] && fail "AWS_ROLE_ARN env variable has not been set"
[[ -z "$AWS_USER_NAME" ]] && fail "AWS_USER_NAME env variable has not been set"
[[ -z "$AWS_ACCESS_KEY_ID" ]] && fail "AWS_ACCESS_KEY_ID env variable has not been set"
[[ -z "$AWS_SECRET_ACCESS_KEY" ]] && fail "AWS_SECRET_ACCESS_KEY env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

while true; do
  echo -e "Waiting for IAM user to be done setting up...\n"
  # Fetch the IAM user creation date and convert it to a Unix timestamp
  create_timestamp=$(aws iam get-user --user-name "${AWS_USER_NAME}" --query 'User.CreateDate' --output text | sed 's/\([+-][0-9]\{2\}:[0-9]\{2\}\)$//' | date -f - "+%s")
  if (($(date +%s) - create_timestamp > 75)); then
    break
  fi
  sleep 2
done

echo -e "Configuring Vault AWS \n"
USERNAME_TEMPLATE="{{ if (eq .Type \"STS\") }}{{ printf \"${AWS_USER_NAME}-%s-%s\" (random 20) (unix_time) | truncate 32 }}{{ else }}{{ printf \"${AWS_USER_NAME}-%s-%s\" (unix_time) (random 20) | truncate 60 }}{{ end }}"
"$binpath" write "${MOUNT}/config/root" access_key="${AWS_ACCESS_KEY_ID}" secret_key="${AWS_SECRET_ACCESS_KEY}" region="${AWS_REGION}" username_template="${USERNAME_TEMPLATE}"

echo -e "Creating Role to create user \n"
"$binpath" write "aws/roles/${VAULT_AWS_ROLE}" \
    credential_type=iam_user \
    permissions_boundary_arn="${AWS_POLICY_ARN}" \
    policy_document=- << EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["ec2:DescribeRegions"],
      "Resource": ["*"]
    }
  ]
}
EOF

echo -e "Verifying root config \n"
"$binpath" read "${MOUNT}/config/root"
ROOT_USERNAME_TEMPLATE=$("$binpath" read "${MOUNT}/config/root" | jq -r '.data.username_template')
[[ "$ROOT_USERNAME_TEMPLATE" == *"$AWS_USER_NAME"* ]] || fail "Uername Template does not include the current role"
