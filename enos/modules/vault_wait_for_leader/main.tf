# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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

variable "vault_instance_count" {
  type        = number
  description = "The number of instances in the vault cluster"
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
}

resource "enos_remote_exec" "wait_for_leader_in_vault_hosts" {
  environment = {
    RETRY_INTERVAL             = var.retry_interval
    TIMEOUT_SECONDS            = var.timeout
    VAULT_ADDR                 = "http://127.0.0.1:8200"
    VAULT_TOKEN                = var.vault_root_token
    VAULT_INSTANCE_PRIVATE_IPS = jsonencode(local.private_ips)
    VAULT_INSTALL_DIR          = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/wait-for-leader.sh")]

  transport = {
    ssh = {
      host = var.vault_hosts[0].public_ip
    }
  }
}
