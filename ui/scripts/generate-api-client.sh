#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Generate the API client for Vault UI using the OpenAPI spec
# Vault Enterprise binary required to generate complete API client
# Run 'make entdev' in vault-enterprise directory prior to running this script or an error will occur

# exit if a command fails
set -e

# check if vault version contains +ent
VERSION=$(vault version)
if [[ ! "$VERSION" == *"+ent"* ]]; then
  echo "Vault Enterprise binary not found and is required to generate the complete API client"
  echo $VERSION
  echo "Run 'make entdev' in vault-enterprise directory to build the binary and try again"
  exit 1
fi

echo "Generating OpenAPI spec..."
echo

# directory where this script is located
SOURCE_DIR=$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")
# parent vault directory
VAULT_DIR="$( cd -P "$SOURCE_DIR/../../" && pwd )"
# scripts directory
SCRIPTS_DIR="$VAULT_DIR/scripts"
# api client directory
API_CLIENT_DIR="$VAULT_DIR/ui/api-client"

# cleanup to remove the generated openapi.json file on exit
cleanup() {
  rm -f "$SCRIPTS_DIR/openapi.json"
  rm -f "$API_CLIENT_DIR/openapi.json"
}
trap 'cleanup' INT TERM EXIT

# change into scripts directory and execute gen_openapi.sh
cd "$SCRIPTS_DIR"
./gen_openapi.sh
echo
echo "OpenAPI spec generated successfully!"

# create the api client directory
echo "Creating API client directory..."
mkdir "$API_CLIENT_DIR" || echo

# move the generated openapi spec to the api-client directory
mv "$SCRIPTS_DIR/openapi.json" "$API_CLIENT_DIR/openapi.json"

# generate the typescript-fetch API client
echo "Generating API client..."
echo
cd "$API_CLIENT_DIR"
# config args specific to typescript-fetch generator
GEN_CONFIG="prefixParameterInterfaces=true,useSingleRequestParameter=false,useSquareBracketsInArrayNames=true"
npx @openapitools/openapi-generator-cli generate -i openapi.json -g typescript-fetch -o . --skip-validate-spec --additional-properties=${GEN_CONFIG}
echo
echo "API client generated successfully!"