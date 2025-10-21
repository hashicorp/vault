#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  return 1
}

[[ -z "$VAULT_PROXY_ADDRESS" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_PROXY_PIDFILE" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# Will cause the Vault CLI to communicate with the Vault Proxy, since it
# is listening at port 8100.
export VAULT_ADDR="http://${VAULT_PROXY_ADDRESS}"

# Explicitly unsetting VAULT_TOKEN to make sure that the Vault Proxy's token
# is used.
unset VAULT_TOKEN

# Use the Vault CLI to communicate with the Vault Proxy (via the VAULT_ADDR env
# var) to lookup the details of the Proxy's token and make sure that the
# .data.path field contains 'auth/approle/login', thus confirming that the Proxy
# automatically authenticated itself.
if ! $binpath token lookup -format=json | jq -Mer --arg expected "auth/approle/login" '.data.path == $expected'; then
  fail "expected proxy to automatically authenticate using 'auth/approle/login', got: '$($binpath token lookup -format=json | jq -r '.data.path')'"
fi

# Now that we're done, kill the proxy
pkill -F "${VAULT_PROXY_PIDFILE}" || true
