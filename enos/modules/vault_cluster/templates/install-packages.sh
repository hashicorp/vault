#!/bin/bash

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
  # Make sure cloud-init is not modifying our sources list while we're trying
  # to install.
  retry 7 grep ec2 /etc/apt/sources.list

  cd /tmp
  retry 5 sudo apt update
  retry 5 sudo apt install -y "$${packages[@]}"
else
  cd /tmp
  retry 7 sudo yum -y install "$${packages[@]}"
fi
