# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "vault_product_version" {}
variable "addition" {}

locals {
  version = "${var.vault_product_version}${var.addition}" 
}

output "storage_addl_config" {
  value = {
    autopilot_upgrade_version = local.version 
  }
}

output "version" {
  value = local.version
}