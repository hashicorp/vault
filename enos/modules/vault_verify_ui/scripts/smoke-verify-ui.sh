#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

url_effective=$(curl -w "%{url_effective}\n" -I -L -s -S "${VAULT_ADDR}" -o /dev/null)
expected="${VAULT_ADDR}/ui/"
if [ "${url_effective}" != "${expected}" ]; then
  fail "Expecting Vault to redirect to UI.\nExpected: ${expected}\nGot: ${url_effective}"
fi

if curl -s "${VAULT_ADDR}/ui/" | grep -q 'Vault UI is not available'; then
  fail "Vault UI is not available"
fi
