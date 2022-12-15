#!/usr/bin/env bash

set -e

binpath=${vault_install_dir}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

$binpath login -method=userpass username=testuser password=passuser1
$binpath kv get secret/test
