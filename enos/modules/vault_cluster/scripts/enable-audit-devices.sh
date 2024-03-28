#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -exo pipefail

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$LOG_FILE_PATH" ]] && fail "LOG_FILE_PATH env variable has not been set"
[[ -z "$SOCKET_PORT" ]] && fail "SOCKET_PORT env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_BIN_PATH" ]] && fail "VAULT_BIN_PATH env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"

enable_file_audit_device() {
  $VAULT_BIN_PATH audit enable file file_path="$LOG_FILE_PATH"
}

enable_syslog_audit_device(){
  $VAULT_BIN_PATH audit enable syslog tag="vault" facility="AUTH"
}

enable_socket_audit_device() {
  "$VAULT_BIN_PATH" audit enable socket address="127.0.0.1:$SOCKET_PORT"
}

main() {
  if ! enable_file_audit_device; then
    fail "Failed to enable vault file audit device"
  fi

  if ! enable_syslog_audit_device; then
    fail "Failed to enable vault syslog audit device"
  fi

  if ! enable_socket_audit_device; then
    local log
    log=$(cat /tmp/vault-socket.log)
    fail "Failed to enable vault socket audit device: listener log: $log"
  fi

  return 0
}

main
