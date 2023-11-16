#!/usr/bin/env sh
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

status=$(${VAULT_BIN_PATH} status -format=json)
version=$(${VAULT_BIN_PATH} version)

echo "{\"status\": ${status}, \"version\": \"${version}\"}"
