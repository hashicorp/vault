#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

retry() {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep 30
    else
      return "$exit"
    fi
  done

  return 0
}

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_IBM_LICENSE_EDITION" ]] && fail "VAULT_IBM_LICENSE_EDITION env variable has not been set"

binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

license_get() {
  issuer=$($binpath license get -format=json | jq -r '.data.autoloaded.issuer')
  edition=$($binpath license get -format=json | jq -r '.data.autoloaded.edition')

  if [ "$issuer" == "pao.ibm.com" ] && [ "$edition" == "$VAULT_IBM_LICENSE_EDITION" ]; then
    echo "License updated; using an IBM license with $VAULT_IBM_LICENSE_EDITION entitlement"
    return 0
  else
    fail "Expected an IBM PAO license with $VAULT_IBM_LICENSE_EDITION entitlement, got issuer: $issuer and edition: $edition"
  fi
}

retry 10 license_get
