#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -ex -o pipefail

# Wait for cloud-init to finish so it doesn't race with any of our package installations.
# Note: Amazon Linux 2 throws Python 2.7 errors when running `cloud-init status` as
# non-root user (known bug).
sudo cloud-init status --wait

if [ "$PACKAGE_MANAGER" == "zypper" ]; then
  echo "Installing dependencies libcap-progs and openssl"
  sudo zypper install --no-confirm libcap-progs
  sudo zypper install --no-confirm openssl
else
  exit 0
fi
