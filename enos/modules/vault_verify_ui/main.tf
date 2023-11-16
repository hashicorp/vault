# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "smoke-verify-ui" {
  for_each = local.instances

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
