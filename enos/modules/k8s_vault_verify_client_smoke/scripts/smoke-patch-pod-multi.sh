#!/bin/sh
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1
#
# smoke-patch-pod-multi.sh — Smoke test: patch-pod with multiple labels in a
# single call using comma-separated -patches.
#
# Sends two label patches at once, then re-fetches the pod and asserts both
# labels are present. This exercises the comma-splitting loop in main.go and
# verifies the JSON patch array is built correctly for multiple operations.
#
# Required env:
#   KUBECLIENT_BIN  — path to the kubeclient binary inside the container
#   SMOKE_NAMESPACE — kubernetes namespace of the target pod
#   SMOKE_POD       — name of the target pod

set -eu

PATCHES="/metadata/labels/smoke-env:staging,/metadata/labels/smoke-multi:true"

fail() {
  echo "[smoke-patch-pod-multi] FAIL: $1" >&2
  exit 1
}

echo "[smoke-patch-pod-multi] Patching with multiple labels: ${PATCHES}" >&2

"${KUBECLIENT_BIN}" -call=patch-pod \
  -namespace="${SMOKE_NAMESPACE}" \
  -pod-name="${SMOKE_POD}" \
  -patches="${PATCHES}"

echo "[smoke-patch-pod-multi] Patch exited 0 — re-fetching pod to verify both labels" >&2

OUTPUT="$("${KUBECLIENT_BIN}" -call=get-pod -namespace="${SMOKE_NAMESPACE}" -pod-name="${SMOKE_POD}")"

if ! echo "${OUTPUT}" | grep -q "staging"; then
  fail "label 'smoke-env=staging' not found after multi-patch. got: ${OUTPUT}"
fi

if ! echo "${OUTPUT}" | grep -q "true"; then
  fail "label 'smoke-multi=true' not found after multi-patch. got: ${OUTPUT}"
fi

echo "[smoke-patch-pod-multi] PASS" >&2
