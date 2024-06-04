# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster target hosts to check"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "expected_state" {
  type        = number
  description = "The expected state to have in vault.core.replication.write_undo_logs telemetry. Must be either 1 for enabled or 0 for disabled."

  validation {
    condition     = contains([0, 1], var.expected_state)
    error_message = "The expected_state must be either 0 or 1"
  }
}

locals {
  public_ips = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_hosts)[idx].public_ip
      private_ip = values(var.vault_hosts)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "smoke-verify-undo-logs" {
  for_each = local.public_ips

  environment = {
    EXPECTED_STATE    = var.expected_state
    VAULT_ADDR        = "http://localhost:8200"
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/smoke-verify-undo-logs.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
