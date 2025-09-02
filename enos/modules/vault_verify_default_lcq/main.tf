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

variable "vault_autopilot_default_max_leases" {
  type        = string
  description = "The autopilot upgrade expected max_leases"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

resource "enos_remote_exec" "smoke_verify_default_lcq" {
  for_each = var.hosts

  environment = {
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
    VAULT_ADDR      = var.vault_addr
    VAULT_TOKEN     = var.vault_root_token
    DEFAULT_LCQ     = var.vault_autopilot_default_max_leases
  }

  scripts = [abspath("${path.module}/scripts/smoke-verify-default-lcq.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
