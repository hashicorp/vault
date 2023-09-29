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

variable "old_vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances to be shutdown"
}

locals {
  public_ips = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.old_vault_instances)[idx].public_ip
      private_ip = values(var.old_vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "shutdown_multiple_nodes" {
  for_each = local.public_ips
  inline   = ["sudo shutdown -H --no-wall; exit 0"]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
