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

variable "vault_autopilot_upgrade_version" {
  type        = string
  description = "The Vault upgraded version"
}

variable "vault_autopilot_upgrade_status" {
  type        = string
  description = "The autopilot upgrade expected status"
}

locals {
  public_ips = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "smoke-verify-autopilot" {
  for_each = local.public_ips

  content = templatefile("${path.module}/templates/smoke-verify-autopilot.sh", {
    vault_install_dir               = var.vault_install_dir
    vault_token                     = var.vault_root_token
    vault_autopilot_upgrade_status  = var.vault_autopilot_upgrade_status,
    vault_autopilot_upgrade_version = var.vault_autopilot_upgrade_version,
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
