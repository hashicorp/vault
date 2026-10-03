#!/bin/sh
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

# Sets up AppRole auth, a policy, a KV secret, and an AppRole role on the
# Vault server pod, then prints "role_id|secret_id" to stdout so the caller
# can split on "|" to extract each value.

set -eu

VAULT_BIN="${VAULT_BIN_PATH:-/bin/vault}"

# Redirect setup command output to stderr so only the credential output
# reaches stdout, keeping the "role_id|secret_id" contract.
$VAULT_BIN auth enable approle >&2 || true

# KV v1 is used (vault secrets enable -path=secret kv). The v1 API serves
# reads directly at secret/<key>, so only that path prefix is needed.
$VAULT_BIN policy write sidecar-policy - >&2 << 'POLICY'
path "secret/*" { capabilities = ["read", "list"] }
POLICY

$VAULT_BIN write auth/approle/role/sidecar-role \
  token_policies="sidecar-policy" \
  token_ttl=1h \
  token_max_ttl=4h >&2

$VAULT_BIN secrets enable -path=secret kv >&2 || true
$VAULT_BIN kv put secret/client-sidecar test_key=hello-from-client-sidecar >&2

# Print role_id and secret_id separated by a pipe so the caller can split on
# it reliably. A pipe cannot appear in a Vault UUID.
printf '%s|%s' \
  "$($VAULT_BIN read  -field=role_id   auth/approle/role/sidecar-role/role-id)" \
  "$($VAULT_BIN write -f -field=secret_id auth/approle/role/sidecar-role/secret-id)"
