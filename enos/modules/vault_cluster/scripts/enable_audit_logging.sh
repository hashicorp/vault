#!/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1
set -eux

# Run nc to listen to port 9090 for the socket audit log
nohup nc -l 9090 &>/dev/null </dev/null &

$VAULT_BIN_PATH audit enable file file_path="$LOG_FILE_PATH"
$VAULT_BIN_PATH audit enable syslog tag="vault" facility="AUTH"
$VAULT_BIN_PATH audit enable socket address="127.0.0.1:9090"
