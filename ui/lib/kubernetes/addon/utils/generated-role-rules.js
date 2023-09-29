/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

const example = `# The below is an example that you can use as a starting point.
#
# rules:
#   - apiGroups: [""]
#     resources: ["serviceaccounts", "serviceaccounts/token"]
#     verbs: ["create", "update", "delete"]
#   - apiGroups: ["rbac.authorization.k8s.io"]
#     resources: ["rolebindings", "clusterrolebindings"]
#     verbs: ["create", "update", "delete"]
#   - apiGroups: ["rbac.authorization.k8s.io"]
#     resources: ["roles", "clusterroles"]
#     verbs: ["bind", "escalate", "create", "update", "delete"]
`;

const readResources = `rules:
- apiGroups: [""]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["extensions"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["apps"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["batch"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["policy"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["networking.k8s.io"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["autoscaling"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
`;

const editResources = `rules:
- apiGroups: [""]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: [""]
  resources:
    ["pods", "pods/attach", "pods/exec", "pods/portforward", "pods/proxy"]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
- apiGroups: [""]
  resources:
    [
      "configmaps",
      "events",
      "persistentvolumeclaims",
      "replicationcontrollers",
      "replicationcontrollers/scale",
      "secrets",
      "serviceaccounts",
      "services",
      "services/proxy",
    ]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
- apiGroups: [""]
  resources: ["serviceaccounts/token"]
  verbs: ["create"]
- apiGroups: ["extensions"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["extensions"]
  resources:
    [
      "daemonsets",
      "deployments",
      "deployments/rollback",
      "deployments/scale",
      "ingresses",
      "networkpolicies",
      "replicasets",
      "replicasets/scale",
      "replicationcontrollers/scale",
    ]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
- apiGroups: ["apps"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["apps"]
  resources:
    [
      "daemonsets",
      "deployments",
      "deployments/rollback",
      "deployments/scale",
      "replicasets",
      "replicasets/scale",
      "statefulsets",
      "statefulsets/scale",
    ]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
- apiGroups: ["batch"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["batch"]
  resources: ["cronjobs", "jobs"]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
- apiGroups: ["policy"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["policy"]
  resources: ["poddisruptionbudgets"]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
- apiGroups: ["networking.k8s.io"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses", "networkpolicies"]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
- apiGroups: ["autoscaling"]
  resources: ["*"]
  verbs: ["get", "watch", "list"]
- apiGroups: ["autoscaling"]
  resources: ["horizontalpodautoscalers"]
  verbs: ["create", "delete", "deletecollection", "patch", "update"]
`;

const updatePods = `rules:
- apiGroups: [""]
  resources: ["secrets", "configmaps", "pods", "endpoints"]
  verbs: ["get", "watch", "list", "create", "delete", "deletecollection", "patch", "update"]
`;

const updateServices = `rules:
- apiGroups: [""]
  resources: ["secrets", "services"]
  verbs: ["get", "watch", "list", "create", "delete", "deletecollection", "patch", "update"]
`;

const usePolicies = `rules:
- apiGroups: ['policy']
  resources: ['podsecuritypolicies']
  verbs:     ['use']
  resourceNames:
  - <list of policies to authorize>
`;

export const getRules = () => [
  { id: '1', label: 'No template', rules: example },
  { id: '2', label: 'Read resources in a namespace', rules: readResources },
  { id: '3', label: 'Edit resources in a namespace', rules: editResources },
  { id: '4', label: 'Update pods, secrets, configmaps, and endpoints', rules: updatePods },
  { id: '5', label: 'Update services and secrets', rules: updateServices },
  { id: '6', label: 'Use pod security policies', rules: usePolicies },
];
