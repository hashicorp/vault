#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

retry 5 "$binpath" status > /dev/null 2>&1

# Create user policy
retry 5 "$binpath" policy write reguser -<<EOF
path "*" {
  capabilities = ["read", "list"]
}
EOF

# Enable the userpass auth method
retry 5 "$binpath" auth enable userpass > /dev/null 2>&1

# Create new user and attach reguser policy
retry 5 "$binpath" write auth/userpass/users/testuser password="passuser1" policies="reguser"

retry 5 "$binpath" secrets enable -path="secret" kv
