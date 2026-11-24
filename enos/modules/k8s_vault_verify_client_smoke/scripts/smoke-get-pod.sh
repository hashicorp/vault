#!/bin/sh
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1
#
# smoke-get-pod.sh — Smoke test: get-pod happy path.
#
# Runs the kubeclient binary (already present inside the sidecar container image)
# with -call=get-pod against the pod this script is executing inside.
# KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT are set automatically by
# Kubernetes for every container in the cluster.
#
# Asserts:
#   - exit code 0
#   - stdout is JSON containing the pod name
#
# Required env:
#   KUBECLIENT_BIN  — path to the kubeclient binary inside the container
#   SMOKE_NAMESPACE — kubernetes namespace of the target pod
#   SMOKE_POD       — name of the target pod

set -eu

fail() {
  echo "[smoke-get-pod] FAIL: $1" >&2
  exit 1
}

echo "[smoke-get-pod] ${KUBECLIENT_BIN} -call=get-pod -namespace=${SMOKE_NAMESPACE} -pod-name=${SMOKE_POD}" >&2

OUTPUT="$("${KUBECLIENT_BIN}" -call=get-pod -namespace="${SMOKE_NAMESPACE}" -pod-name="${SMOKE_POD}")"

echo "[smoke-get-pod] stdout: ${OUTPUT}" >&2

# Response must contain the pod name we asked for.
if ! echo "${OUTPUT}" | grep -q "${SMOKE_POD}"; then
  fail "expected pod name '${SMOKE_POD}' in output, got: ${OUTPUT}"
fi

# Response must look like JSON.
case "${OUTPUT}" in
  "{"*) ;;
  *) fail "expected JSON starting with '{', got: ${OUTPUT}" ;;
esac

echo "[smoke-get-pod] PASS" >&2
