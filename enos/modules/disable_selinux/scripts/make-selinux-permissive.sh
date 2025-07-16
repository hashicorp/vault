#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

if ! type getenforce &> /dev/null; then
  exit 0
fi

if sudo getenforce | grep Enforcing; then
  sudo setenforce 0
fi
