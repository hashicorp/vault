#!/bin/bash
set -eu -o pipefail

# Set up the environment for building Vault.
root_dir="$(git rev-parse --show-toplevel)"

pushd "$root_dir" > /dev/null

IFS="-" read -r BASE_VERSION _other <<< "$(make version)"
export VAULT_VERSION=$BASE_VERSION
echo $VAULT_VERSION
