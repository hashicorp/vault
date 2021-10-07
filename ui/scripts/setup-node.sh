#!/usr/bin/env bash

# This script is mostly intended for use in the builder image
# defined in <reporoot>/packages-oss.yml

set -euo pipefail

NODE_VERSION="$(cat .nvmrc)"
NODE_VERSION="${NODE_VERSION#v}"

echo "==> Setting up node v$NODE_VERSION (from .nvmrc)"

# shellcheck disable=SC1090
command -v nvm || {
	echo "==> Sourcing .bashrc..."
	source ~/.bashrc
}
command -v nvm || { echo "ERROR: nvm not installed"; exit 1; }

echo "==> nvm install $NODE_VERSION"
nvm install "$NODE_VERSION"

echo "==> nvm use $NODE_VERSION"
nvm use "$NODE_VERSION"
