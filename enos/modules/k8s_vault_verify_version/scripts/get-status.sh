#!/usr/bin/env sh
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

status=$(${VAULT_BIN_PATH} status -format=json)
version=$(${VAULT_BIN_PATH} version)

echo "{\"status\": ${status}, \"version\": \"${version}\"}"
