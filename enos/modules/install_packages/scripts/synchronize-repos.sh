#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "${PACKAGE_MANAGER}" ]] && fail "PACKAGE_MANAGER env variable has not been set"
[[ -z "${RETRY_INTERVAL}" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "${TIMEOUT_SECONDS}" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

# The SLES AMI's do not come configured with Zypper repositories by default. To get them you
# have to run SUSEConnect to register the instance with SUSE. On the AMI this is handled
# automatically by a oneshot systemd unit called guestregister.service. This oneshot service needs
# to complete before any other repo or package steps are completed. At the time of writing it's very
# unreliable so we have to ensure that it has correctly executed ourselves or restart it. We do this
# by checking if the guestregister.service has reached the correct "inactive" state that we need.
# If it hasn't reached that state it's usually in some sort of active state, i.e. running, or it has
# failed. If it's in one of the active states we need to let it continue and check the status when
# it completes. If it has completed but is failed we'll restart the service to re-run the script that
# executes SUSEConnect.
sles_check_guestregister_service_and_restart_if_failed() {
  local active_state
  local failed_state

  # systemctl returns non-zero exit codes. We rely on output here because all states don't have
  # their own exit code.
  set +e
  active_state=$(sudo systemctl is-active guestregister.service)
  failed_state=$(sudo systemctl is-failed guestregister.service)
  set -e

  case "$active_state" in
    active | activating | deactivating)
      # It's running so we'll return 1 and get retried by the caller
      echo "the guestregister.service is still in the ${active_state} state" 1>&2
      return 1
      ;;
    *)
      if [ "$active_state" == "inactive" ] && [ "$failed_state" == "inactive" ]; then
        # The oneshot has completed and hasn't "failed"
        echo "the guestregister.service is 'inactive' for both active and failed states"
        return 0
      fi

      # Our service is stopped and failed, restart it and hope it works the next time
      sudo systemctl restart --wait guestregister.service
      ;;
  esac
}

# Check or restart the guestregister service if it has failed. If it passes do another check to make
# sure that the zypper repositories list isn't empty.
sles_ensure_suseconnect() {
  local health_output
  if ! health_output=$(sles_check_guestregister_service_and_restart_if_failed); then
    echo "the guestregister.service failed to reach a healthy state: ${health_output}" 1>&2
    return 1
  fi

  # Make sure Zypper has repositories.
  if ! lr_output=$(zypper lr); then
    echo "The guestregister.service failed. Unable to SUSEConnect and thus have no Zypper repositories: ${lr_output}: ${health_output}." 1>&2
    return 1
  fi

  return 0
}

# Try each mirror in order until one succeeds for apt update. We support Ubuntu
# 22.04 (jammy), 24.04 (noble), and 26.04 (resolute) on four architectures:
#   - amd64        (AWS)  → archive.ubuntu.com/ubuntu only
#   - arm64        (AWS)  → ports.ubuntu.com/ubuntu-ports only
#   - s390x        (Fyre) → ports.ubuntu.com/ubuntu-ports only
#   - ppc64el      (Fyre) → ports.ubuntu.com/ubuntu-ports only
#
# Despite the Release file listing all architectures, actual package files are
# split: archive.ubuntu.com only hosts amd64/i386; ports.ubuntu.com hosts
# arm64, s390x, ppc64el, and other non-x86 architectures. We order mirrors by
# the running architecture so the correct one is tried first, with the other as
# a fallback in case of transient unavailability.
#
# Each fallback sources.list includes the base release plus the -updates and
# -security pockets so that a full apt-get upgrade succeeds, not just apt-get update.
apt_update_with_fallback() {
  # Select mirror order based on architecture: amd64 uses archive first;
  # arm64/s390x/ppc64el use ports first.
  local arch
  arch=$(dpkg --print-architecture)
  local mirrors
  if [[ "${arch}" == "amd64" ]]; then
    mirrors=(
      "http://archive.ubuntu.com/ubuntu"
      "http://ports.ubuntu.com/ubuntu-ports"
    )
  else
    mirrors=(
      "http://ports.ubuntu.com/ubuntu-ports"
      "http://archive.ubuntu.com/ubuntu"
    )
  fi

  # First, try a plain apt update using whatever sources the system already has.
  # This is the normal happy path and avoids touching sources.list at all.
  if sudo apt-get update 2>&1; then
    return 0
  fi

  echo "apt update failed with default sources, trying fallback mirrors..." 1>&2

  local codename
  codename=$(lsb_release -cs)

  for mirror in "${mirrors[@]}"; do
    echo "Trying mirror: ${mirror}" 1>&2
    # Include the base release pocket plus -updates and -security so that a
    # subsequent apt-get upgrade can resolve all packages for every supported
    # Ubuntu version (22.04/jammy, 24.04/noble, 26.04/resolute) on all four
    # architectures (amd64, arm64, s390x, ppc64el).
    if sudo apt-get update -o "Dir::Etc::SourceList=/dev/stdin" \
      -o "Dir::Etc::SourceParts=/dev/null" \
      -o "APT::Get::List-Cleanup=false" \
      <<< "deb ${mirror} ${codename} main restricted universe multiverse
deb ${mirror} ${codename}-updates main restricted universe multiverse
deb ${mirror} ${codename}-security main restricted universe multiverse" 2>&1; then
      echo "Successfully updated with mirror: ${mirror}" 1>&2
      return 0
    fi
    echo "Mirror ${mirror} failed, trying next..." 1>&2
  done

  echo "All mirrors failed" 1>&2
  return 1
}

# Synchronize our repositories so that further installation steps are working with updated cache
# and repo metadata.
synchronize_repos() {
  case $PACKAGE_MANAGER in
    apt)
      apt_update_with_fallback
      ;;
    dnf)
      sudo dnf makecache
      ;;
    yum)
      sudo yum makecache
      ;;
    zypper)
      if [ "$DISTRO" == "sles" ]; then
        if ! sles_ensure_suseconnect; then
          return 1
        fi
      fi
      sudo zypper --gpg-auto-import-keys --non-interactive ref
      sudo zypper --gpg-auto-import-keys --non-interactive refs
      ;;
    *)
      return 0
      ;;
  esac
}

# Function to check cloud-init status and retry on failure
# Before we start to modify repositories and install packages we'll wait for cloud-init to finish
# so it doesn't race with any of our package installations.
# We run as sudo because Amazon Linux 2 throws Python 2.7 errors when running `cloud-init status` as
# non-root user (known bug).
wait_for_cloud_init() {
  if output=$(sudo cloud-init status --wait); then
    return 0
  else
    res=$?
    case $res in
      2)
        {
          echo "WARNING: cloud-init did not complete successfully but recovered."
          echo "Exit code: $res"
          echo "Output: $output"
          echo "Here are the logs for the failure:"
          cat /var/log/cloud-init-*
        } 1>&2
        return 0
        ;;
      *)
        {
          echo "cloud-init did not complete successfully."
          echo "Exit code: $res"
          echo "Output: $output"
          echo "Here are the logs for the failure:"
          cat /var/log/cloud-init-*
        } 1>&2
        return 1
        ;;
    esac
  fi
}

# Wait for cloud-init if it exists
type cloud-init && wait_for_cloud_init

# Synchronizing repos
begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if synchronize_repos; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for distro repos to be set up"
