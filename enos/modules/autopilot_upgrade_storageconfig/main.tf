# Copyright IBM Corp. 2026, 2025
# SPDX-License-Identifier: BUSL-1.1

variable "vault_product_version" {}

output "storage_addl_config" {
  value = {
    autopilot_upgrade_version = var.vault_product_version
  }
}
