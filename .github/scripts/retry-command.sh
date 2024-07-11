#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1
set -euo pipefail

tries=5
count=0

until "$@"
do
  if [ $count -eq $tries ]; then
    echo "tried $count times, exiting"
    exit 1
  fi
  ((count++))
  echo "trying again, attempt $count"
  sleep 2
done
