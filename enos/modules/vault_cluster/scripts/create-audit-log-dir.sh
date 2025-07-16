#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -eux

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$LOG_FILE_PATH" ]] && fail "LOG_FILE_PATH env variable has not been set"
[[ -z "$SERVICE_USER" ]] && fail "SERVICE_USER env variable has not been set"

LOG_DIR=$(dirname "$LOG_FILE_PATH")

function retry {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=10
    count=$((count + 1))

    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

retry 7 id -a "$SERVICE_USER"

sudo mkdir -p "$LOG_DIR"
sudo chown -R "$SERVICE_USER":"$SERVICE_USER" "$LOG_DIR"
