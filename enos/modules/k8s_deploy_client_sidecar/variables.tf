# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

variable "context_name" {
  type        = string
  description = "The name of the k8s context for the kind cluster"
}

variable "kubeconfig_base64" {
  type        = string
  description = "The base64 encoded version of the Kubernetes configuration file"
}

variable "client_image_repository" {
  type        = string
  description = "The name of the Vault client image repository to deploy"
}

variable "client_image_tag" {
  type        = string
  description = "The tag of the Vault client image to deploy"
}

variable "vault_root_token" {
  type        = string
  description = "Root token for Vault server cluster to set up auth and policies"
}

variable "vault_pods" {
  type = list(object({
    name      = string
    namespace = string
  }))
  description = "List of Vault server pods in the cluster"
}

variable "vault_bin_path" {
  type        = string
  description = "The path to the Vault binary inside the server pod"
  default     = "/bin/vault"
}

variable "kind_cluster_name" {
  type        = string
  description = "The name of the kind cluster; used to resolve the control-plane container (docker exec)"
}
