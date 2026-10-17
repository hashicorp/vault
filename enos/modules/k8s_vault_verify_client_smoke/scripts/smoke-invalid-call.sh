#!/bin/sh
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1
#
# smoke-invalid-call.sh — Smoke test: unsupported -call value must exit non-zero.
#
# Passes -call=delete-pod (an operation that does not exist) and asserts the
# binary exits with a non-zero status code immediately, without making any
# Kubernetes API call. The binary panics with:
#   panic: unsupported call provided: "delete-pod"
#
# Required env:
#   KUBECLIENT_BIN — path to the kubeclient binary inside the container

set -eu

fail() {
  echo "[smoke-invalid-call] FAIL: $1" >&2
  exit 1
}

echo "[smoke-invalid-call] Calling ${KUBECLIENT_BIN} with unsupported -call=delete-pod" >&2

# Capture exit code without letting set -e abort the script.
set +e
"${KUBECLIENT_BIN}" -call=delete-pod -namespace=default -pod-name=some-pod 2> /tmp/invalid-call-stderr
EXIT_CODE=$?
set -e

echo "[smoke-invalid-call] exit code: ${EXIT_CODE}" >&2
cat /tmp/invalid-call-stderr >&2 || true

if [ "${EXIT_CODE}" -eq 0 ]; then
  fail "expected non-zero exit for unsupported -call, got exit 0"
fi

echo "[smoke-invalid-call] PASS (exit code ${EXIT_CODE} as expected)" >&2
