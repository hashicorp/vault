#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$ARTIFACT_NAME" ]] && fail "ARTIFACT_NAME env variable has not been set"

if [ "${ARTIFACT_NAME##*.}" == "zip" ]; then
  echo "Skipped removing unit file because new artifact is a zip bundle"
  exit 0
fi

# Get the unit file for the vault.service that is running. If it's not in /etc/systemd then it
# should be a package provided unit file so we don't need to delete anything.
#
# Note that we use -p instead of -P so that we support ancient amzn 2 systemctl.
if ! unit_path=$(systemctl show -p FragmentPath vault | cut -d = -f2 2>&1); then
  echo "Skipped removing unit file because and existing path could not be found: $unit_path"
  exit 0
fi

if [[ "$unit_path" == *"/etc/systemd"* ]]; then
  if [ -f "$unit_path" ]; then
    echo "Removing old systemd unit file: $unit_path"
    if ! out=$(sudo rm "$unit_path" 2>&1); then
      fail "Failed to remove old unit file: $unit_path: $out"
    fi
  else
    echo "Skipped removing old systemd unit file because it no longer exists: $unit_path"
  fi
else
  echo "Skipped removing old systemd unit file because it was not created in /etc/systemd/: $unit_path"
fi
