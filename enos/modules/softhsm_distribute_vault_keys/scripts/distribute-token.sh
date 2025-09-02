#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -ex

fail() {
  echo "$1" 1>&2
  exit 1
}

# If we're not given keys we'll short circuit. This should only happen if we're skipping distribution
# because we haven't created a token or keys.
if [ -z "$TOKEN_BASE64" ]; then
  echo "TOKEN_BASE64 environment variable was unset. Assuming we don't need to distribute our token" 1>&2
  exit 0
fi

[[ -z "$SOFTHSM_GROUP" ]] && fail "SOFTHSM_GROUP env variable has not been set"
[[ -z "$TOKEN_DIR" ]] && fail "TOKEN_DIR env variable has not been set"

# Convert our base64 encoded gzipped tarball of the softhsm token back into a tarball.
base64 --decode - > token.tgz <<< "$TOKEN_BASE64"

# Expand it. We assume it was written with the correct directory metadata. Do this as a superuser
# because the token directory should be owned by root.
sudo tar -xvf token.tgz -C "$TOKEN_DIR"

# Make sure the vault user is in the softhsm group to get access to the tokens.
sudo usermod -aG "$SOFTHSM_GROUP" vault
sudo chown -R "vault:$SOFTHSM_GROUP" "$TOKEN_DIR"
