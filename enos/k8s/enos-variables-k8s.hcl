# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "vault_image_repository" {
  description = "The repository for the docker image to load, i.e. hashicorp/vault"
  type        = string
  default     = null
}

variable "vault_log_level" {
  description = "The server log level for Vault logs. Supported values (in order of detail) are trace, debug, info, warn, and err."
  type        = string
  default     = "info"
}

variable "vault_product_version" {
  description = "The vault product version to test"
  type        = string
  default     = null
}

variable "vault_product_revision" {
  type        = string
  description = "The vault product revision to test"
  default     = null
}

variable "vault_docker_image_archive" {
  description = "The path to the location of the docker image archive to test"
  type        = string
  default     = null
}

variable "vault_instance_count" {
  description = "How many instances to create for the Vault cluster"
  type        = number
  default     = 3
}

variable "terraform_plugin_cache_dir" {
  description = "The directory to cache Terraform modules and providers"
  type        = string
  default     = null
}

variable "tfc_api_token" {
  description = "The Terraform Cloud QTI Organization API token."
  type        = string
}

variable "vault_build_date" {
  description = "The build date for the vault docker image"
  type        = string
  default     = ""
}

variable "cluster_node_count" {
  description = "The number of nodes to create for the kubernetes cluster"
  type        = number
  default     = 1
}

variable "vault_license_path" {
  description = "The path to a valid Vault enterprise edition license. This is only required for non-oss editions"
  type        = string
  default     = null
}

variable "kubeconfig_path" {
  description = "The path to the kubeconfig file that will be created."
  type        = string
}

variable "vault_agent_image_name" {
  description = "The name of the docker image to run for the vault agent, e.g. hashicorp/vault"
  type        = string
  default     = "hashicorp/vault"
}

variable "vault_agent_image_tag" {
  description = "The tag of the docker image to run for the vault agent, e.g. 1.13.1"
  type        = string
  default     = "1.13.1"
}
