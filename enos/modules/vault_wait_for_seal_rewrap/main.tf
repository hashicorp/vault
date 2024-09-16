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

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  private_ips = [for k, v in values(tomap(var.hosts)) : tostring(v["private_ip"])]
  first_key   = element(keys(enos_remote_exec.wait_for_seal_rewrap_to_be_completed), 0)
}

resource "enos_remote_exec" "wait_for_seal_rewrap_to_be_completed" {
  for_each = var.hosts
  environment = {
    RETRY_INTERVAL    = var.retry_interval
    TIMEOUT_SECONDS   = var.timeout
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/wait-for-seal-rewrap.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

output "stdout" {
  value = enos_remote_exec.wait_for_seal_rewrap_to_be_completed[local.first_key].stdout
}

output "stderr" {
  value = enos_remote_exec.wait_for_seal_rewrap_to_be_completed[local.first_key].stdout
}
