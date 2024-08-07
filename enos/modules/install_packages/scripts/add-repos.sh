#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${PACKAGE_MANAGER}" ]] && fail "PACKAGE_MANAGER env variable has not been set"
[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

# Add any repositories that have have been passed in
add_repos() {
  # If we don't have any repos on the list for this distro, no action needed.
  if [ ${#DISTRO_REPOS[@]} -lt 1 ]; then
    echo "DISTRO_REPOS is empty; No repos required for the packages for this Linux distro."
    return 0
  fi

  case $PACKAGE_MANAGER in
    apt)
      # NOTE: We do not currently add any apt repositories in our scenarios. I suspect if that time
      # comes we'll need to add support for apt-key here.
      for repo in ${DISTRO_REPOS}; do
        if [ "$repo" == "__none" ]; then
          continue
        fi
        sudo add-apt-repository "${repo}"
      done
    ;;
    dnf)
      for repo in ${DISTRO_REPOS}; do
        if [ "$repo" == "__none" ]; then
          continue
        fi
        sudo dnf install -y "${repo}"
        sudo dnf makecache -y
      done
    ;;
    yum)
      for repo in ${DISTRO_REPOS}; do
        if [ "$repo" == "__none" ]; then
          continue
        fi
        sudo yum install -y "${repo}"
        sudo yum makecache -y
      done
    ;;
    zypper)
      # Add each repo
      for repo in ${DISTRO_REPOS}; do
        if [ "$repo" == "__none" ]; then
          continue
        fi
        if sudo zypper lr "${repo}"; then
          echo "A repo named ${repo} already exists, skipping..."
          continue
        fi
        sudo zypper --gpg-auto-import-keys --non-interactive addrepo "${repo}"
      done
      sudo zypper --gpg-auto-import-keys ref
      sudo zypper --gpg-auto-import-keys refs
    ;;
    *)
      fail "Unsupported package manager: ${PACKAGE_MANAGER}"
  esac
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if add_repos; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for distro repos to be set up"
