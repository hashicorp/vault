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

# Synchronize our repositories so that futher installation steps are working with updated cache
# and repo metadata.
synchronize_repos() {
  case $PACKAGE_MANAGER in
    apt)
      sudo apt update
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
check_cloud_init() {
  sudo cloud-init status --wait
  exit_code=$?
  if [ "$?" -ne 0 ] && [ "$?" -ne 2 ]; then
    echo "cloud-init did not complete successfully. Exit code: $exit_code" 1>&2
    echo "Here are the logs for the failure:"
    cat /var/log/cloud-init-* | grep "Failed"
    exit 1
  fi
}

# Checking cloud-init
check_cloud_init

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
