#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

binpath=${VAULT_INSTALL_DIR}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# Activate the primary
$binpath write -f sys/replication/dr/primary/enable
