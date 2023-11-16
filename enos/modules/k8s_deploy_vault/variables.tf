# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "context_name" {
  type        = string
  description = "The name of the k8s context for Vault"
}

variable "ent_license" {
  type        = string
  description = "The value of a valid Vault Enterprise license"
}

variable "image_repository" {
  type        = string
  description = "The name of the Vault repository, ie hashicorp/vault or hashicorp/vault-enterprise for the image to deploy"
}

variable "image_tag" {
  type        = string
  description = "The tag of the vault image to deploy"
}

variable "kubeconfig_base64" {
  type        = string
  description = "The base64 encoded version of the Kubernetes configuration file"
}

variable "vault_edition" {
  type        = string
  description = "The Vault product edition"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_log_level" {
  description = "The server log level for Vault logs. Supported values (in order of detail) are trace, debug, info, warn, and err."
  type        = string
}
