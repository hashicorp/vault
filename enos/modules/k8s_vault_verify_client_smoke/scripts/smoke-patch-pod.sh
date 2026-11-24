#!/bin/sh
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1
#
# smoke-patch-pod.sh — Smoke test: patch-pod single label (add operation).
#
# Adds /metadata/labels/smoke-test=passed to the target pod, then immediately
# re-fetches it with get-pod and asserts the label is present in the JSON.
# Verifies end-to-end: patch serialisation, Kubernetes API acceptance, and
# that get-pod reflects the updated state.
#
# Required env:
#   KUBECLIENT_BIN  — path to the kubeclient binary inside the container
#   SMOKE_NAMESPACE — kubernetes namespace of the target pod
#   SMOKE_POD       — name of the target pod

set -eu

LABEL_PATH="/metadata/labels/smoke-test"
LABEL_VALUE="passed"

fail() {
  echo "[smoke-patch-pod] FAIL: $1" >&2
  exit 1
}

echo "[smoke-patch-pod] Patching ${SMOKE_NAMESPACE}/${SMOKE_POD}: ${LABEL_PATH}=${LABEL_VALUE}" >&2

"${KUBECLIENT_BIN}" -call=patch-pod \
  -namespace="${SMOKE_NAMESPACE}" \
  -pod-name="${SMOKE_POD}" \
  -patches="${LABEL_PATH}:${LABEL_VALUE}"

echo "[smoke-patch-pod] Patch exited 0 — re-fetching pod to verify label persisted" >&2

OUTPUT="$("${KUBECLIENT_BIN}" -call=get-pod -namespace="${SMOKE_NAMESPACE}" -pod-name="${SMOKE_POD}")"

if ! echo "${OUTPUT}" | grep -q "${LABEL_VALUE}"; then
  fail "label value '${LABEL_VALUE}' not found in pod JSON after patch. got: ${OUTPUT}"
fi

echo "[smoke-patch-pod] PASS" >&2
