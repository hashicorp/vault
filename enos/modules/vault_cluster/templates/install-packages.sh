#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -ex -o pipefail

packages="${packages}"

if [ "$packages" == "" ]
then
  echo "No dependencies to install."
  exit 0
fi

function retry {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

echo "Installing Dependencies: $packages"
if [ -f /etc/debian_version ]; then
  # Do our best to make sure that we don't race with cloud-init. Wait a reasonable time until we
  # see ec2 in the sources list. Very rarely cloud-init will take longer than we wait. In that case
  # we'll just install our packages.
  retry 7 grep ec2 /etc/apt/sources.list || true

  cd /tmp
  retry 5 sudo apt update
  retry 5 sudo apt install -y "$${packages[@]}"
else
  cd /tmp
  retry 7 sudo yum -y install "$${packages[@]}"
fi
