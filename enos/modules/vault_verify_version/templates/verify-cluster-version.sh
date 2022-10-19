#!/usr/bin/env bash

# The Vault smoke test to verify the Vault version installed

set -e

binpath=${vault_install_dir}/vault
edition=${vault_edition}
version=${vault_version}
sha=${vault_revision}
builddate=${vault_build_date}
release="$version+$edition"

fail() {
	echo "$1" 1>&2
	exit 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='${vault_token}'

if [[ "$builddate" != "" ]]; then
  build_date=$builddate
else
  build_date=$("$binpath" status -format=json | jq -Mr .build_date)
fi

if [[ "$(echo $version |awk -F'.' '{print $2}')" -ge 11 ]]; then
  version_expected="Vault v$version ($sha), built $build_date"
else
  version_expected="Vault v$version ($sha)"
fi

case "$release" in
  *+oss) ;;
  *+ent) ;;
  *+ent.hsm) version_expected="$version_expected (cgo)";;
  *+ent.fips1402) version_expected="$version_expected (cgo)" ;;
  *+ent.hsm.fips1402) version_expected="$version_expected (cgo)" ;;
  *) fail "($release) file doesn't match any known license types"
esac

version_expected_nosha=$(echo "$version_expected" | awk '!($3="")' | sed 's/  / /' | sed -e 's/[[:space:]]*$//')
version_output=$("$binpath" version)

if [[ "$version_output" == "$version_expected_nosha" ]] || [[ "$version_output" == "$version_expected" ]]; then
  echo "Version verification succeeded!"
else
  fail "expected Version=$version_expected or $version_expected_nosha, got: $version_output"
fi
