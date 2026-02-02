#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$OIDC_ISSUER_URL" ]] && fail "OIDC_ISSUER_URL env variable has not been set"
[[ -z "$OIDC_KEY_NAME" ]] && fail "OIDC_KEY_NAME env variable has not been set"
[[ -z "$OIDC_KEY_ROTATION_PERIOD" ]] && fail "OIDC_KEY_ROTATION_PERIOD env variable has not been set"
[[ -z "$OIDC_KEY_VERIFICATION_TTL" ]] && fail "OIDC_KEY_VERIFICATION_TTL env variable has not been set"
[[ -z "$OIDC_KEY_ALGORITHM" ]] && fail "OIDC_KEY_ALGORITHM env variable has not been set"
[[ -z "$OIDC_ROLE_NAME" ]] && fail "OIDC_ROLE_NAME env variable has not been set"
[[ -z "$OIDC_ROLE_TTL" ]] && fail "OIDC_ROLE_TTL env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json

# Verify that we have the correct issuer URL
if ! cfg=$("$binpath" read identity/oidc/config); then
  fail "failed to read identity/oidc/config: $cfg"
elif ! jq -Merc --arg URL "$OIDC_ISSUER_URL" '.data.issuer == $URL' <<< "$cfg"; then
  fail "oidc issuer URL is incorrect, expected: $OIDC_ISSUER_URL, got $(jq -Mrc '.data.issuer' <<< "$cfg")"
fi

# Verify that our token algorithm, rotation period and verification TTL are correct
if ! key_res=$("$binpath" read "identity/oidc/key/$OIDC_KEY_NAME"); then
  fail "failed to read identity/oidc/key/$OIDC_KEY_NAME: $key_res"
fi

if ! jq -Merc --arg ALG "$OIDC_KEY_ALGORITHM" '.data.algorithm == $ALG' <<< "$key_res"; then
  fail "oidc token algorithm is incorrect, expected: $OIDC_KEY_ALGORITHM, got $(jq -Mrc '.data.algorithm' <<< "$key_res")"
fi

if ! jq -Merc --argjson RP "$OIDC_KEY_ROTATION_PERIOD" '.data.rotation_period == $RP' <<< "$key_res"; then
  fail "oidc token rotation_period is incorrect, expected: $OIDC_KEY_ROTATION_PERIOD, got $(jq -Mrc '.data.rotation_period' <<< "$key_res")"
fi

if ! jq -Merc --argjson TTL "$OIDC_KEY_VERIFICATION_TTL" '.data.verification_ttl == $TTL' <<< "$key_res"; then
  fail "oidc token verification_ttl is incorrect, expected: $OIDC_KEY_VERIFICATION_TTL, got $(jq -Mrc '.data.verification_ttl' <<< "$key_res")"
fi

# Verify that our role key and TTL are correct.
if ! role_res=$("$binpath" read "identity/oidc/role/$OIDC_ROLE_NAME"); then
  fail "failed to read identity/oidc/role/$OIDC_ROLE_NAME: $role_res"
fi

if ! jq -Merc --arg KEY "$OIDC_KEY_NAME" '.data.key == $KEY' <<< "$role_res"; then
  fail "oidc role key is incorrect, expected: $OIDC_KEY_NAME, got $(jq -Mrc '.data.key' <<< "$role_res")"
fi

if ! jq -Merc --argjson TTL "$OIDC_ROLE_TTL" '.data.ttl == $TTL' <<< "$role_res"; then
  fail "oidc role ttl is incorrect, expected: $OIDC_ROLE_TTL, got $(jq -Mrc '.data.ttl' <<< "$role_res")"
fi
