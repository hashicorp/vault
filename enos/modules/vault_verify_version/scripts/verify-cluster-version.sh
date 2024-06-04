#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


# Verify the Vault "version" includes the correct base version, build date,
# revision SHA, and edition metadata.
set -e

binpath=${VAULT_INSTALL_DIR}/vault
edition=${VAULT_EDITION}
version=${VAULT_VERSION}
sha=${VAULT_REVISION}
build_date=${VAULT_BUILD_DATE}
# VAULT_TOKEN must also be set

fail() {
  echo "$1" 1>&2
  exit 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

version_expected="Vault v$version ($sha), built $build_date"

case "$edition" in
  *ce) ;;
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
  msg="$(printf "\nThe Vault cluster did not match the expected version, expected:\n%s\nor\n%s\ngot:\n%s" "$version_expected" "$version_expected_nosha" "$version_output")"
  if type diff &> /dev/null; then
    # Diff exits non-zero if we have a diff, which we want, so we'll guard against failing early.
    if ! version_diff=$(diff  <(echo "$version_expected" ) <(echo "$version_output") -u -L expected -L got); then
      msg="$(printf "\nThe Vault cluster did not match the expected version:\n%s" "$version_diff")"
    fi
  fi

  fail "$msg"
fi
