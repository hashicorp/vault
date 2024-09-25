# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
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
  description = "The vault cluster hosts that can be expected as a leader"
}

variable "ip_version" {
  type        = number
  description = "The IP version used for the Vault TCP listener"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out"
  default     = 60
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

locals {
  ipv6s       = [for k, v in values(tomap(var.hosts)) : tostring(v["ipv6"])]
  private_ips = [for k, v in values(tomap(var.hosts)) : tostring(v["private_ip"])]
}

resource "enos_remote_exec" "wait_for_leader_in_hosts" {
  environment = {
    IP_VERSION                 = var.ip_version
    TIMEOUT_SECONDS            = var.timeout
    RETRY_INTERVAL             = var.retry_interval
    VAULT_ADDR                 = var.vault_addr
    VAULT_TOKEN                = var.vault_root_token
    VAULT_INSTANCE_IPV6S       = jsonencode(local.ipv6s)
    VAULT_INSTANCE_PRIVATE_IPS = jsonencode(local.private_ips)
    VAULT_INSTALL_DIR          = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/wait-for-leader.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}
