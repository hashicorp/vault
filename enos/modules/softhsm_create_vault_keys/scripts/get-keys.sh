#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$TOKEN_DIR" ]] && fail "TOKEN_DIR env variable has not been set"

# Tar up our token. We have to do this as a superuser because softhsm is owned by root.
sudo tar -czf token.tgz -C "$TOKEN_DIR" .
me="$(whoami)"
sudo chown "$me:$me" token.tgz

# Write the value STDOUT as base64 so we can handle binary data as a string
base64 -i token.tgz
