# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "expected_state" {
  type        = number
  description = "The expected state to have in vault.core.replication.write_undo_logs telemetry. Must be either 1 for enabled or 0 for disabled."

  validation {
    condition     = contains([0, 1], var.expected_state)
    error_message = "The expected_state must be either 0 or 1"
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster target hosts to check"
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

resource "enos_remote_exec" "smoke-verify-undo-logs" {
  for_each = var.hosts

  environment = {
    EXPECTED_STATE    = var.expected_state
    RETRY_INTERVAL    = var.retry_interval
    TIMEOUT_SECONDS   = var.timeout
    VAULT_ADDR        = var.vault_addr
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
