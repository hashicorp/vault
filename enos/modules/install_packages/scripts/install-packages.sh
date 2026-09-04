#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
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

# Install packages based on the provided packages and package manager. We assume that the repositories
# have already been synchronized by the repo setup that is a prerequisite for this script.
install_packages() {
  if [[ "${PACKAGES}" = "__skip" ]]; then
    return 0
  fi

  set -x
  echo "Installing Dependencies: ${PACKAGES}"

  # Use the default package manager of the current Linux distro to install packages
  case $PACKAGE_MANAGER in
    apt)
      for package in ${PACKAGES}; do
        if dpkg -s "${package}"; then
          echo "Skipping installation of ${package} because it is already installed"
          continue
        else
          echo "Installing ${package}"
          local output
          # Bound apt's own lock-wait to RETRY_INTERVAL so that a competing process holding
          # /var/lib/dpkg/lock-frontend (e.g. unattended-upgrades, apt-daily) doesn't cause a
          # single apt invocation to consume our whole TIMEOUT_SECONDS budget. This lets the
          # outer retry loop actually retry multiple times instead of timing out on one attempt.
          # Acquire::Retries tells apt to retry a failed package download (e.g. transient 503s
          # from a regional mirror) a few times within the same invocation before giving up.
          if ! output=$(sudo apt install -y -o DPkg::Lock::Timeout="${RETRY_INTERVAL}" -o Acquire::Retries="3" "${package}" 2>&1); then
            echo "Failed to install ${package}: ${output}" 1>&2
            # The package lists may be stale or the mirror may have recovered since
            # synchronize-repos.sh ran. Best-effort refresh before the outer loop retries so the
            # next attempt has a chance to see healed/updated metadata.
            sudo apt-get update -o DPkg::Lock::Timeout="${RETRY_INTERVAL}" -o Acquire::Retries="3" 1>&2 || true
            return 1
          fi
        fi
      done
      ;;
    dnf)
      for package in ${PACKAGES}; do
        if rpm -q "${package}"; then
          echo "Skipping installation of ${package} because it is already installed"
          continue
        else
          echo "Installing ${package}"
          local output
          if ! output=$(sudo dnf -y install "${package}" 2>&1); then
            echo "Failed to install ${package}: ${output}" 1>&2
            return 1
          fi
        fi
      done
      ;;
    yum)
      for package in ${PACKAGES}; do
        if rpm -q "${package}"; then
          echo "Skipping installation of ${package} because it is already installed"
          continue
        else
          echo "Installing ${package}"
          local output
          if ! output=$(sudo yum -y install "${package}" 2>&1); then
            echo "Failed to install ${package}: ${output}" 1>&2
            return 1
          fi
        fi
      done
      ;;
    zypper)
      for package in ${PACKAGES}; do
        if rpm -q "${package}"; then
          echo "Skipping installation of ${package} because it is already installed"
          continue
        else
          echo "Installing ${package}"
          local output
          if ! output=$(sudo zypper --non-interactive install -y -l --force-resolution "${package}" 2>&1); then
            echo "Failed to install ${package}: ${output}" 1>&2
            return 1
          fi
        fi
      done
      ;;
    *)
      fail "No matching package manager provided."
      ;;
  esac
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
