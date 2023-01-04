#!/usr/bin/env bash

set -e

binpath=${VAULT_INSTALL_DIR}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# To keep the authentication method and module verification consistent between all
# Enos scenarios we authenticate using testuser created by vault_verify_write_data module
$binpath login -method=userpass username=testuser password=passuser1
$binpath kv get secret/test
