#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1
ARTIFACTORY_NPM_URL="https://artifactory.hashicorp.engineering/artifactory/npm/"

# Show error if doormat CLI or jq doesn't exist

if ! command -v doormat &> /dev/null; then
  echo "Error: doormat CLI is not installed."
  exit 1
fi

if ! command -v jq &> /dev/null; then
  echo "Error: jq is not installed."
  exit 1
fi

doormat login
# Get artifactory token from doormat
AF_TOKEN=$(doormat artifactory create-token | jq -r .access_token)

# Get package info from artifactory API
PACKAGE_INFO=$(curl -s -L -H "Authorization: Bearer $AF_TOKEN" "https://artifactory.hashicorp.engineering/artifactory/api/npm/npm/@hashicorp-internal/vault-reporting")

# Extract the latest version and tarball URL
LATEST_VERSION=$(echo "$PACKAGE_INFO" | jq -r '."dist-tags".latest')
TARBALL_URL=$(echo "$PACKAGE_INFO" | jq -r --arg version "$LATEST_VERSION" '.versions[$version].dist.tarball')

echo "LATEST_VERSION: $LATEST_VERSION"
echo "TARBALL_URL: $TARBALL_URL"

# Download the tarball
echo "Downloading vault-reporting v$LATEST_VERSION..."
curl -L -H "Authorization: Bearer $AF_TOKEN" -o "vault-reporting/$LATEST_VERSION.tgz" "$TARBALL_URL"

if [ $? -eq 0 ]; then
    echo "✅ Successfully downloaded vault-reporting-$LATEST_VERSION.tgz"
    echo "📦 File size: $(ls -lh vault-reporting-$LATEST_VERSION.tgz | awk '{print $5}')"
else
    echo "❌ Failed to download tarball"
    exit 1
fi

# Delete all but the latest version
find vault-reporting -name "*.tgz" ! -name "$LATEST_VERSION.tgz" -delete

# Update package.json to point to file:vault-reporting/$LATEST_VERSION.tgz
FILENAME="$LATEST_VERSION.tgz"
jq --arg version "$LATEST_VERSION" --arg filename "$FILENAME" '.dependencies["@hashicorp-internal/vault-reporting"] = "file:vault-reporting/" + $filename' package.json > tmp.json && mv tmp.json package.json
# Install dependencies
yarn install