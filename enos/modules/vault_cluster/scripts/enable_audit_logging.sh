#!/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -exo pipefail

# Run nc to listen on port 9090 for the socket auditor. We spawn nc
# with nohup to ensure that the listener doesn't expect a SIGHUP and
# thus block the SSH session from exiting or terminating on exit.
# We immediately write to STDIN from /dev/null to give nc an
# immediate EOF so as to not block on expecting STDIN.
nohup nc -kl 9090 &> /dev/null < /dev/null &

# Wait for nc to be listening before we attempt to enable the socket auditor.
retries=3
count=0
until nc -zv 127.0.0.1 9090 &> /dev/null < /dev/null; do
  wait=$((2 ** count))
  count=$((count + 1))

  if [ "$count" -lt "$retries" ]; then
    sleep "$wait"
  else

    echo "Timed out waiting for nc to listen on 127.0.0.1:9090" 1>&2
    exit 1
  fi
done

# Enable the auditors.
$VAULT_BIN_PATH audit enable file file_path="$LOG_FILE_PATH"
$VAULT_BIN_PATH audit enable syslog tag="vault" facility="AUTH"
$VAULT_BIN_PATH audit enable socket address="127.0.0.1:9090"
