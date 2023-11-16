# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# A shim unseal key module for shamir seal types

variable "cluster_id" { default = null }
variable "cluster_meta" { default = null }
variable "common_tags" { default = null }
variable "names" {
  type    = list(string)
  default = []
}

output "alias" { value = null }
output "id" { value = null }
output "resource_name" { value = null }
output "resource_names" { value = var.names }
