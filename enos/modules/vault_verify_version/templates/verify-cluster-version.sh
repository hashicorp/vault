#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


# Verify the Vault "version" includes the correct base version, build date,
# revision SHA, and edition metadata.
set -e

binpath=${vault_install_dir}/vault
edition=${vault_edition}
version=${vault_version}
sha=${vault_revision}
build_date=${vault_build_date}

fail() {
	echo "$1" 1>&2
	exit 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='${vault_token}'

version_expected="Vault v$version ($sha), built $build_date"

case "$edition" in
  *oss) ;;
  *ent) ;;
  *ent.hsm) version_expected="$version_expected (cgo)";;
  *ent.fips1402) version_expected="$version_expected (cgo)" ;;
  *ent.hsm.fips1402) version_expected="$version_expected (cgo)" ;;
  *) fail "Unknown Vault edition: ($edition)" ;;
esac

version_expected_nosha=$(echo "$version_expected" | awk '!($3="")' | sed 's/  / /' | sed -e 's/[[:space:]]*$//')
version_output=$("$binpath" version)

if [[ "$version_output" == "$version_expected_nosha" ]] || [[ "$version_output" == "$version_expected" ]]; then
  echo "Version verification succeeded!"
else
  fail "expected Version=$version_expected or $version_expected_nosha, got: $version_output"
fi
