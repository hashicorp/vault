# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "vault_product_version" {}

output "storage_addl_config" {
  value = {
    autopilot_upgrade_version = var.vault_product_version
  }
}
