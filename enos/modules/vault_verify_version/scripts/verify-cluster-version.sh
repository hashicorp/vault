#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# Verify the Vault "version" includes the correct base version, build date,
# revision SHA, and edition metadata.
set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_BUILD_DATE" ]] && fail "VAULT_BUILD_DATE env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_VERSION" ]] && fail "VAULT_VERSION env variable has not been set"

# The sys/version-history endpoint only includes major.minor.patch, any other semver fields need to
# be stripped out.
if ! version=$(cut -d + -f1 <<< "$VAULT_VERSION" | cut -d - -f1); then
  fail "failed to parse the expected version: $version"
fi

if ! vh=$(curl -s -X LIST -H "X-Vault-Token: $VAULT_TOKEN" http://127.0.0.1:8200/v1/sys/version-history | jq -eMc '.data'); then
  fail "failed to Vault cluster version history: $vh"
fi

if ! out=$(jq -eMc --arg version "$version" '.keys | contains([$version])' <<< "$vh"); then
  fail "cluster version history does not include our expected version: expected: $version, versions: $(jq -eMc '.keys' <<< "$vh"): output: $out"
fi

if ! out=$(jq -eMc --arg version "$version" --arg bd "$VAULT_BUILD_DATE" '.key_info[$version].build_date == $bd' <<< "$vh"); then
  fail "cluster version history build date is not the expected date: expected: true, expected date: $VAULT_BUILD_DATE, key_info: $(jq -eMc '.key_info' <<< "$vh"), output: $out"
fi

printf "Cluster version information is valid!: %s\n" "$vh"
