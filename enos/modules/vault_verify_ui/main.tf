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

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

resource "enos_remote_exec" "smoke-verify-ui" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR = var.vault_addr,
  }

  scripts = [abspath("${path.module}/scripts/smoke-verify-ui.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
