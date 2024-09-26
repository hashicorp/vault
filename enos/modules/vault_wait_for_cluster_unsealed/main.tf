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
  description = "The vault cluster instances that were created"
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

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

resource "enos_remote_exec" "verify_node_unsealed" {
  for_each = var.hosts

  scripts = [abspath("${path.module}/scripts/verify-vault-node-unsealed.sh")]

  environment = {
    HOST_IPV4         = each.value.public_ip
    HOST_IPV6         = each.value.ipv6
    RETRY_INTERVAL    = var.retry_interval
    TIMEOUT_SECONDS   = var.timeout
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
