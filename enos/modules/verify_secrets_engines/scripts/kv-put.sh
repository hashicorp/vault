#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$KEY" ]] && fail "KEY env variable has not been set"
[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$SECRET_PATH" ]] && fail "SECRET_PATH env variable has not been set"
[[ -z "$VALUE" ]] && fail "VALUE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

"$binpath" kv put -mount="$MOUNT" "$SECRET_PATH" "$KEY=$VALUE"
