#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1
#
# setup-rbac.sh — Grant the app pod's default service account GET + PATCH
# on pods so the kubeclient binary inside the sidecar container can reach the
# Kubernetes API server.
#
# Uses kubectl apply (idempotent) proxied through the kind control-plane
# container because kubectl is not installed on the Enos runner.
#
# Required env:
#   KIND_CLUSTER_NAME  — e.g. "kind"
#   TARGET_NAMESPACE   — namespace of the app pod (e.g. "default")

set -euo pipefail

CONTROL_PLANE="${KIND_CLUSTER_NAME}-control-plane"

echo "==> Applying RBAC for kubeclient smoke tests (namespace: ${TARGET_NAMESPACE})" >&2

docker exec -i "${CONTROL_PLANE}" kubectl apply -f - << EOF
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubeclient-smoke-role
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "patch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubeclient-smoke-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubeclient-smoke-role
subjects:
  # The app pod runs under the "default" service account in its namespace.
  - kind: ServiceAccount
    name: default
    namespace: ${TARGET_NAMESPACE}
EOF

echo "==> RBAC applied" >&2
