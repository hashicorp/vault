#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -ex -o pipefail

if [ "$PACKAGES" == "" ]
then
  echo "No dependencies to install."
  exit 0
fi

# Wait for cloud-init to finish so it doesn't race with any of our package installations.
# Note: Amazon Linux 2 throws Python 2.7 errors when running `cloud-init status` as
# non-root user (known bug).
sudo cloud-init status --wait

echo "Installing Dependencies: $PACKAGES"

# Use the default package manager of the current Linux distro to install packages
if [ "$PACKAGE_MANAGER" = "apt" ]; then
  cd /tmp
  sudo apt update
  # Disable this shellcheck rule about double-quoting array expansions; if we use
  # double quotes on ${PACKAGES[@]}, it does not take the packages as separate
  # arguments.
  # shellcheck disable=SC2068,SC2086
  sudo apt install -y ${PACKAGES[@]}
elif [ "$PACKAGE_MANAGER" = "yum" ]; then
  cd /tmp
   # shellcheck disable=SC2068,SC2086
  sudo yum -y install ${PACKAGES[@]}
elif [ "$PACKAGE_MANAGER" = "zypper" ]; then
  # Note: For some SUSE distro versions, and/or some packages, the
  # packages may not be offered in an official repo. If the first install step
  # fails, we instead attempt to register with PackageHub,SUSE's third party
  # package marketplace, and then find and install the package from there.

  # shellcheck disable=SC2068,SC2086
  sudo zypper install --no-confirm ${PACKAGES[@]} || ( sudo SUSEConnect -p PackageHub/$SLES_VERSION/$ARCH && sudo zypper install --no-confirm ${PACKAGES[@]})
  # For SUSE distros on arm64 architecture, we need to manually install these two
  # packages in order to install Vault RPM packages later.
  if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    sudo zypper install --no-confirm libcap-progs
    sudo zypper install --no-confirm openssl
  fi
else
  echo "No matching package manager provided."
  exit 1
fi
