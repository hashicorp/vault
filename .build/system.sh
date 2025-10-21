#!/bin/bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

export DEBIAN_FRONTEND=noninteractive

install() {
  apt-get install -y "$@"
}

# Install our cross building tools
# https://packages.ubuntu.com/search?suite=noble&section=all&arch=any&keywords=crossbuild-essential&searchon=names

apt-get update
apt-get install -y --no-install-recommends build-essential \
  gcc-s390x-linux-gnu \
  crossbuild-essential-s390x \
  ca-certificates \
  curl \
  git

host_arch="$(dpkg --print-architecture)"
host_arch="${host_arch##*-}"
case "$host_arch" in
  amd64)
    install crossbuild-essential-arm64 gcc-aarch64-linux-gnu
    ;;
  arm64)
    install gcc-x86-64-linux-gnu
    ;;
  *)
    echo "Building on $host_arch has not been implemented" 1>&2
    exit 1
    ;;
esac

# Clean up after ourselves for a minimal image
apt-get clean
rm -rf /var/lib/apt/lists/*
