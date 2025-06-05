#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$GOARCH" ]] && fail "A GOARCH has not been defined"
[[ -z "$GITHUB_TOKEN" ]] && fail "A GITHUB_TOKEN has not been defined"

host_arch="$(dpkg --print-architecture)"
host_arch="${host_arch##*-}"
if [[ "$host_arch" != "$GOARCH" ]]; then
  # We're building for a different architecture than our target host OS so
  # we have to tell the Go compiler to use the correct C cross-compiler for
  # our target instead of relying on the host C compiler.
  #
  # https://packages.ubuntu.com/search?suite=noble&section=all&arch=any&keywords=linux-gnu-gcc&searchon=contents
  case "$GOARCH" in
    amd64)
      export CC=x86_64-linux-gnu-gcc
      ;;
    arm64)
      export CC=aarch64-linux-gnu-gcc
      ;;
    s390x)
      export CC=s390x-linux-gnu-gcc
      ;;
    *)
      fail "Building for $GOARCH has not been implemented"
      ;;
  esac
fi

# Assume that /build is where we've mounted the vault repo.
git config --global --add safe.directory /build
git config --global url."https://${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"

# Exec our command
cd build || exit 1
exec "$@"
