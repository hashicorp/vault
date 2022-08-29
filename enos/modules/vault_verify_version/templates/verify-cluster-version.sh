#!/usr/bin/env bash

# The Vault smoke test to verify the Vault version installed

set -e

binpath=${vault_install_dir}/vault

fail() {
	echo "$1" 1>&2
	exit 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

binary_version_full=$($binpath version)
# Get the Vault build tag
binary_version=$(cut -d ' ' -f2 <<< $binary_version_full)
# Strip the leading v
semantic=$${binary_version:1}
# Get the build timestamp
build_date=$(cut -d ' ' -f5 <<< $binary_version_full)

export VAULT_ADDR='http://127.0.0.1:8200'

# Ensure that the cluster version and build time match the binary installed
vault_status=$("$binpath" status -format json)
result=$(jq -Mr \
  --arg version "$semantic" \
  --arg build_date "$build_date" \
  'select(.version == $version) | .build_date == $build_date' \
  <<< $vault_status
)

if [[ "$result" != "true" ]]; then
  fail "expected version $binary_version with build_date $build_date, got status $vault_status"
fi
