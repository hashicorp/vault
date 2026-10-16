#!/bin/sh
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

set -eu

SECRET_FILE="/vault/data/rendered-secret.txt"

for i in $(seq 1 30); do
  if [ -s "$SECRET_FILE" ]; then
    echo "Found rendered secret file:" >&2
    cat "$SECRET_FILE"
    exit 0
  fi
  echo "Waiting for rendered secret file ($i/30)..." >&2
  sleep 2
done

echo "Timed out waiting for rendered secret file at $SECRET_FILE" >&2
exit 1
