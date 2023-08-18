#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

function retry {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

function fail {
	echo "$1" 1>&2
	exit 1
}

binpath=${VAULT_INSTALL_DIR}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# To keep the authentication method and module verification consistent between all
# Enos scenarios we authenticate using testuser created by vault_verify_write_data module
retry 5 $binpath login -method=userpass username=testuser password=passuser1
retry 5 $binpath kv get secret/test
