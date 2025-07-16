#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$CONFIG_PATH" ]] && fail "CONFIG_PATH env variable has not been set"
[[ -z "$TOKEN_DIR" ]] && fail "TOKEN_DIR env variable has not been set"
[[ -z "$SKIP" ]] && fail "SKIP env variable has not been set"

if [ "$SKIP" == "true" ]; then
  exit 0
fi

cat << EOF | sudo tee "$CONFIG_PATH"
directories.tokendir = $TOKEN_DIR
objectstore.backend = file
log.level = DEBUG
slots.removable = false
slots.mechanisms = ALL
library.reset_on_fork = false
EOF

sudo mkdir -p "$TOKEN_DIR"
sudo chmod 0770 "$TOKEN_DIR"
