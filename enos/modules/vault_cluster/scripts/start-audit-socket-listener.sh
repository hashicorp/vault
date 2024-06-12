#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -exo pipefail

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$NETCAT_COMMAND" ]] && fail "NETCAT_COMMAND env variable has not been set"
[[ -z "$SOCKET_PORT" ]] && fail "SOCKET_PORT env variable has not been set"

socket_listener_procs() {
  pgrep -x "${NETCAT_COMMAND}"
}

kill_socket_listener() {
  pkill  "${NETCAT_COMMAND}"
}

test_socket_listener() {
   "${NETCAT_COMMAND}" -zvw 2 127.0.0.1 "$SOCKET_PORT" < /dev/null
}

start_socket_listener() {
  if socket_listener_procs; then
    test_socket_listener
    return $?
  fi

  # Run nc to listen on port 9090 for the socket auditor. We spawn nc
  # with nohup to ensure that the listener doesn't expect a SIGHUP and
  # thus block the SSH session from exiting or terminating on exit.
  nohup nc -kl "$SOCKET_PORT" >> /tmp/vault-socket.log 2>&1 < /dev/null &
}

read_log() {
  local f
  f=/tmp/vault-socket.log
  [[ -f "$f" ]] && cat "$f"
}

main() {

  if socket_listener_procs; then
    # Clean up old nc's that might not be working
    kill_socket_listener
  fi

  if ! start_socket_listener; then
    fail "Failed to start audit socket listener: socket listener log: $(read_log)"
  fi

  # wait for nc to listen
  sleep 1

  if ! test_socket_listener; then
    fail "Error testing socket listener: socket listener log: $(read_log)"
  fi

  return 0
}

main
