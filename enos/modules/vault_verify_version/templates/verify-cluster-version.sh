#!/usr/bin/env bash

# The Vault smoke test to verify the Vault version installed

set -e

binpath=${vault_install_dir}/vault

fail() {
	echo "$1" 1>&2
	exit 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='${vault_token}'

found_version=$($binpath version)
if [[ "$found_version" != '${expected_version}' ]]; then
  fail "expected version ${expected_version}, got $found_version"
fi
