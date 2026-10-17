# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

variable "context_name" {
  type        = string
  description = "The name of the k8s context for the kind cluster"
}

variable "kubeconfig_base64" {
  type        = string
  description = "The base64 encoded kubeconfig for the kind cluster"
}

variable "app_pod_name" {
  type        = string
  description = "Name of the application pod that contains the vault-client-sidecar container"
  default     = "vault-client-app"
}

variable "app_pod_namespace" {
  type        = string
  description = "Namespace of the application pod"
  default     = "default"
}

# ---------------------------------------------------------------------------
# kubeclient_bin_path — the only thing you change when the standalone client
# image ships.
#
# TODAY  : default "/bin/vault"
#   The vault-enterprise image is used as the client image. It contains the
#   vault binary at /bin/vault. The kubeclient CLI flags (-call, -namespace,
#   -pod-name, -patches) do not exist on vault, so the smoke tests are
#   skipped automatically via skip_smoke = true when this is set to the
#   default. RBAC is still applied so the infrastructure is ready.
#
# FUTURE : set to the path of the kubeclient binary inside the standalone
#   client image (e.g. "/bin/kubeclient"). Pass it in the scenario:
#     -var kubeclient_bin_path=/bin/kubeclient
#   The module detects it is not the vault fallback and runs all smoke tests.
# ---------------------------------------------------------------------------
variable "kubeclient_bin_path" {
  type        = string
  description = "Path to the kubeclient binary inside the vault-client-sidecar container. Set to /bin/vault (default) to skip binary smoke tests; set to the real kubeclient path when the standalone image is available."
  default     = "/bin/vault"
}

variable "kind_cluster_name" {
  type        = string
  description = "Name of the kind cluster; used to resolve the control-plane container for docker exec kubectl calls"
}
