#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -ex -o pipefail

if [ "$PACKAGES" == "" ]
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
      exit "$exit"
    fi
  done

  return 0
}

echo "Installing Dependencies: $PACKAGES"
if [ -f /etc/debian_version ]; then
  # Do our best to make sure that we don't race with cloud-init. Wait a reasonable time until we
  # see ec2 in the sources list. Very rarely cloud-init will take longer than we wait. In that case
  # we'll just install our packages.
  retry 7 grep ec2 /etc/apt/sources.list || true

  cd /tmp
  retry 5 sudo apt update
  # shellcheck disable=2068
  retry 5 sudo apt install -y ${PACKAGES[@]}
else
  cd /tmp
  # shellcheck disable=2068
  retry 7 sudo yum -y install ${PACKAGES[@]}
fi
