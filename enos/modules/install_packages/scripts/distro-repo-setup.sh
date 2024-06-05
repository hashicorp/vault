#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$DISTRO" ]] && fail "DISTRO env variable has not been set"
[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

setup_repos() {
  # If we don't have any repos on the list for this distro, no action needed.
  if [ ${#DISTRO_REPOS[@]} -lt 1 ]; then
    echo "DISTRO_REPOS is empty; No repos required for the packages for this Linux distro."
    return 0
  fi

  # Wait for cloud-init to finish so it doesn't race with any of our package installations.
  # Note: Amazon Linux 2 throws Python 2.7 errors when running `cloud-init status` as
  # non-root user (known bug).
  sudo cloud-init status --wait

  case $DISTRO in
    "sles")
      for repo in ${DISTRO_REPOS}; do
        sudo zypper addrepo "${repo}"
      done
      ;;
    "rhel")
      for repo in ${DISTRO_REPOS}; do
        sudo rm -r /var/cache/dnf
        sudo dnf install -y "${repo}"
        sudo dnf update -y --refresh
      done
      ;;
    *)
      return
      ;;
  esac
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if setup_repos; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for distro repos to install"
