# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "container_image_archive" {
  description = "The path to the location of the container image archive to test"
  type        = string
  default     = null # If none is given we'll simply load a container from a repo
}

variable "log_level" {
  description = "The server log level for Vault logs. Supported values (in order of detail) are trace, debug, info, warn, and err."
  type        = string
  default     = "trace"
}

variable "instance_count" {
  description = "How many instances to create for the Vault cluster"
  type        = number
  default     = 3
}

variable "terraform_plugin_cache_dir" {
  description = "The directory to cache Terraform modules and providers"
  type        = string
  default     = null
}

variable "vault_build_date" {
  description = "The expected vault build date"
  type        = string
  default     = ""
}

variable "vault_revision" {
  type        = string
  description = "The expected vault revision"
}

variable "vault_version" {
  description = "The expected vault version"
  type        = string
}
