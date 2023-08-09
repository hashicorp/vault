# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_edition" {
  type        = string
  description = "The vault product edition"
}

variable "vault_pods" {
  type = list(object({
    name      = string
    namespace = string
  }))
  description = "The vault instances for the cluster to verify"
}

variable "kubeconfig_base64" {
  type        = string
  description = "The base64 encoded version of the Kubernetes configuration file"
}

variable "context_name" {
  type        = string
  description = "The name of the k8s context for Vault"
}
