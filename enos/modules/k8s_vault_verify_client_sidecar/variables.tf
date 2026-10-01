# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

variable "context_name" {
  type        = string
  description = "The name of the k8s context for Vault"
}

variable "kubeconfig_base64" {
  type        = string
  description = "The base64 encoded version of the Kubernetes configuration file"
}

variable "app_pod_name" {
  type        = string
  description = "Name of the application pod running with the client sidecar"
  default     = "vault-client-app"
}

variable "app_pod_namespace" {
  type        = string
  description = "Namespace of the application pod running with the client sidecar"
  default     = "default"
}
