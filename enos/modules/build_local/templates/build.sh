#!/bin/bash
set -eux -o pipefail

# Requirements
npm install --global yarn || true

# Set up the environment for building Vault.
root_dir="$(git rev-parse --show-toplevel)"

pushd "$root_dir" > /dev/null

export GO_TAGS=${build_tags}
export CGO_ENABLED=0

IFS="-" read -r BASE_VERSION _other <<< "$(make version)"
export VAULT_VERSION=$BASE_VERSION

build_date="$(make build-date)"
export VAULT_BUILD_DATE=$build_date

revision="$(git rev-parse HEAD)"
export VAULT_REVISION=$revision
popd > /dev/null

# Go to the UI directory of the Vault repo and build the UI
pushd "$root_dir/ui" > /dev/null
yarn install --ignore-optional
npm rebuild node-sass
yarn --verbose run build
popd > /dev/null

# Build for linux/amd64 and create a bundle since we're deploying it to linux/amd64
pushd "$root_dir" > /dev/null
export GOARCH=${goarch}
export GOOS=${goos}
make build

zip -r -j ${bundle_path} dist/
popd > /dev/null
