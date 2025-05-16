#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$AUTH_PATH" ]] && fail "AUTH_PATH env variable has not been set"
[[ -z "$GROUPATTR" ]] && fail "GROUPATTR env variable has not been set"
[[ -z "$GROUPDN" ]] && fail "GROUPDN env variable has not been set"
[[ -z "$INSECURE_TLS" ]] && fail "INSECURE_TLS env variable has not been set"
[[ -z "$UPNDOMAIN" ]] && fail "UPNDOMAIN env variable has not been set"
[[ -z "$URL" ]] && fail "URL env variable has not been set"
[[ -z "$USERATTR" ]] && fail "USERATTR env variable has not been set"
[[ -z "$USERDN" ]] && fail "USERDN env variable has not been set"

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
"$binpath" write "auth/$AUTH_PATH/config" \
    url="$URL" \
    userdn="$USERDN" \
    userattr="$USERATTR" \
    groupdn="$GROUPDN" \
    groupattr="$GROUPATTR" \
    upndomain="$UPNDOMAIN" \
    insecure_tls="$INSECURE_TLS"
