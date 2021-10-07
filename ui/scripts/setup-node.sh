#!/usr/bin/env bash

# This script is mostly intended for use in the builder image
# defined in <reporoot>/packages-oss.yml

set -euo pipefail

NODE_VERSION="$(cat .nvmrc)"

# shellcheck disable=SC1090
command -v nvm || { source ~/.bashrc; }
command -v nvm || { echo "ERROR: nvm not installed"; exit 1; }

nvm install "$NODE_VERSION"
nvm use "$NODE_VERSION"
