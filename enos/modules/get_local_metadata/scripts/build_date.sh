#!/bin/env bash
set -eu -o pipefail

# Set up the environment for building Vault.
root_dir="$(git rev-parse --show-toplevel)"

pushd "$root_dir" > /dev/null

IFS="-" read -r VAULT_BUILD_DATE _other <<< "$(make build-date)"
echo $VAULT_BUILD_DATE
