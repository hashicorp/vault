#!/usr/bin/env sh

set -e

status=$(${VAULT_BIN_PATH} status -format=json)
version=$(${VAULT_BIN_PATH} version)

echo "{\"status\": ${status}, \"version\": \"${version}\"}"
