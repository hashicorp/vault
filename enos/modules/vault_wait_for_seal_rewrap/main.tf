# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "vault_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster hosts that can be expected as a leader"
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
  private_ips = [for k, v in values(tomap(var.vault_hosts)) : tostring(v["private_ip"])]
  first_key   = element(keys(enos_remote_exec.wait_for_seal_rewrap_to_be_completed), 0)
}

resource "enos_remote_exec" "wait_for_seal_rewrap_to_be_completed" {
  for_each = var.vault_hosts
  environment = {
    RETRY_INTERVAL    = var.retry_interval
    TIMEOUT_SECONDS   = var.timeout
    VAULT_ADDR        = "http://127.0.0.1:8200"
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
