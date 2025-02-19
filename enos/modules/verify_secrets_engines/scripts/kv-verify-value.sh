#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$MOUNT" ]] && fail "MOUNT env variable has not been set"
[[ -z "$SECRET_PATH" ]] && fail "SECRET_PATH env variable has not been set"
[[ -z "$KEY" ]] && fail "KEY env variable has not been set"
[[ -z "$VALUE" ]] && fail "VALUE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if res=$("$binpath" kv get -mount="$MOUNT" "$SECRET_PATH"); then
  # Note that this expects KVv2 response payloads. KVv1 does not include doubly nested .data
  if jq -Merc --arg VALUE "$VALUE" --arg KEY "$KEY" '.data.data[$KEY] == $VALUE' <<< "$res"; then
    printf "kv %s/%s %s=%s is valid\n" "$MOUNT" "$SECRET_PATH" "$KEY" "$VALUE"
    exit 0
  fi
  fail "kv $MOUNT/$SECRET_PATH $KEY=$VALUE invalid! Got: $(jq -Mrc --arg KEY "$KEY" '.data[$KEY]' <<< "$res")"
else
  fail "failed to read kv data for $MOUNT/$SECRET_PATH: $res"
fi
