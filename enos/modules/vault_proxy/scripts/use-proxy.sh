#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

binpath=${VAULT_INSTALL_DIR}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

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
$binpath token lookup -format=json | jq -r '.data.path' | grep -q 'auth/approle/login'

# Now that we're done, kill the proxy
pkill -F "${VAULT_PROXY_PIDFILE}" || true
