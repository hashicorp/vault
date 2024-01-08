# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# A shim seal module for shamir seals. For Shamir seals the enos_vault_init resource will take care
# of creating our seal.

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "cluster_id" { default = null }
variable "cluster_meta" { default = null }
variable "cluster_ssh_keypair" { default = null }
variable "common_tags" { default = null }
variable "image_id" { default = null }
variable "other_resources" {
  type    = list(string)
  default = []
}

output "resource_name" { value = null }
output "resource_names" { value = var.other_resources }
output "attributes" { value = null }
