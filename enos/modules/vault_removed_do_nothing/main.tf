# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster followers"
}


variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out"
  default     = 60
}

variable "listener_port" {
  type        = number
  description = "The listener port for vault"
}
variable "vault_leader_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The leader's host information"
}
variable "vault_local_addr" {
  type        = string
  description = "The local address to use to query vault"
}
variable "cluster_port" {
  type        = number
  description = "The cluster port for vault"
}

variable "ip_version" {
  type        = number
  description = "The IP version to use for the Vault TCP listeners"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}
variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}
variable "vault_seal_type" {
  type        = string
  description = "The Vault seal type"
}

variable "add_back_nodes" {
  type        = bool
  description = "whether to add the nodes back"
}

variable "vault_unseal_keys" {}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the vault binary is installed"
}