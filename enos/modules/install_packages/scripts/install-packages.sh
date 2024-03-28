#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "${PACKAGES}" ]] && fail "PACKAGES env variable has not been set"
[[ -z "${PACKAGE_MANAGER}" ]] && fail "PACKAGE_MANAGER env variable has not been set"

install_packages() {
  if [[ "${PACKAGES}" = "__skip" ]]; then
    return 0
  fi 

  set -x
  echo "Installing Dependencies: ${PACKAGES}"

  # Use the default package manager of the current Linux distro to install packages
  if [[ "${PACKAGE_MANAGER}" = "apt" ]] ; then
    set -x
    sudo apt update
    for package in ${PACKAGES}; do
      if dpkg -s "${package}"; then
        continue
      else
        echo "Installing ${package}"
        sudo apt install -y "${package}"
      fi
    done
  elif [[ "${PACKAGE_MANAGER}" = "yum" ]] ; then
    set -x
    for package in ${PACKAGES}; do
      if sudo yum list installed | rpm -q "${package}"; then
        continue
      else
        echo "Installing ${package}"
        sudo yum -y install "${package}"
      fi
    done
  elif [[ "${PACKAGE_MANAGER}" = "zypper" ]]; then
    set -x
    cd /tmp
    sudo zypper --gpg-auto-import-keys ref
    for package in ${PACKAGES}; do
      if rpm -q "${package}"; then
        continue
      else
        echo "Installing ${package}"
        sudo zypper --non-interactive install "${package}"
        date
      fi
      sudo zypper search -i
    done
  else
    fail "No matching package manager provided."
  fi
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [[ "$(date +%s)" -lt "${end_time}" ]]; do
  if install_packages; then
    exit 0
  fi

  sleep "${RETRY_INTERVAL}"
done

fail "Timed out waiting for packages to install"
