#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -exo pipefail

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$IP_VERSION" ]] && fail "IP_VERSION env variable has not been set"
[[ -z "$NETCAT_COMMAND" ]] && fail "NETCAT_COMMAND env variable has not been set"
[[ -z "$SOCKET_PORT" ]] && fail "SOCKET_PORT env variable has not been set"

if [ "$IP_VERSION" = "4" ]; then
  export SOCKET_ADDR="127.0.0.1"
else
  export SOCKET_ADDR="::1"
fi

socket_listener_procs() {
  pgrep -x "${NETCAT_COMMAND}"
}

kill_socket_listener() {
  pkill  "${NETCAT_COMMAND}"
}

test_socket_listener() {
  case $IP_VERSION in
    4)
      "${NETCAT_COMMAND}" -zvw 2 "${SOCKET_ADDR}" "$SOCKET_PORT" < /dev/null
      ;;
    6)
      "${NETCAT_COMMAND}" -6 -zvw 2 "${SOCKET_ADDR}" "$SOCKET_PORT" < /dev/null
      ;;
    *)
      fail "unknown IP_VERSION: $IP_VERSION"
      ;;
  esac
}

start_socket_listener() {
  if socket_listener_procs; then
    test_socket_listener
    return $?
  fi

  # Run nc to listen on port 9090 for the socket auditor. We spawn nc
  # with nohup to ensure that the listener doesn't expect a SIGHUP and
  # thus block the SSH session from exiting or terminating on exit.
  case $IP_VERSION in
    4)
      nohup nc -kl "$SOCKET_PORT" >> /tmp/vault-socket.log 2>&1 < /dev/null &
      ;;
    6)
      nohup nc -6 -kl "$SOCKET_PORT" >> /tmp/vault-socket.log 2>&1 < /dev/null &
      ;;
    *)
      fail "unknown IP_VERSION: $IP_VERSION"
      ;;
  esac
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
