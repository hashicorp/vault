#!/bin/bash
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

set -x

if [[ ! -f /etc/apt/sources.list.d/ubuntu.sources ]]; then
  echo "Cannot update sources: /etc/apt/sources.list.d/ubuntu.sources is not present" 1>&2
  exit 0
fi

# Enable our components so we can find more packages
sed -i 's|^Components: main$|Components: main universe restricted multiverse|g' /etc/apt/sources.list.d/ubuntu.sources

function fail() {
  echo "$1" 1>&2
  exit 1
}

count=0
retries=3
while :; do
  if apt-get update --fix-missing; then
    exit 0
  fi

  wait=$((2 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
    sleep "$wait"
  else
    fail "Timed out updating ubuntu sources"
  fi
done
