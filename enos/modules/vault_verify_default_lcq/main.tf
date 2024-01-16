# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "vault_autopilot_default_max_leases" {
  type        = string
  description = "The autopilot upgrade expected max_leases"
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
  public_ips = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "smoke_verify_default_lcq" {
  for_each = local.public_ips

  environment = {
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
    VAULT_ADDR      = "http://localhost:8200"
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
