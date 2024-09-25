# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "cluster_id" {
  type = string
}

variable "cluster_meta" {
  type    = string
  default = null
}

variable "cluster_ssh_keypair" {
  type    = string
  default = null
}

variable "common_tags" {
  type    = map(string)
  default = null
}

variable "other_resources" {
  type    = list(string)
  default = []
}

locals {
  cluster_name = var.cluster_meta == null ? var.cluster_id : "${var.cluster_id}-${var.cluster_meta}"
}

resource "aws_kms_key" "key" {
  description             = "auto-unseal-key-${local.cluster_name}"
  deletion_window_in_days = 7 // 7 is the shortest allowed window
  tags                    = var.common_tags
}

resource "aws_kms_alias" "alias" {
  name          = "alias/auto-unseal-key-${local.cluster_name}"
  target_key_id = aws_kms_key.key.key_id
}

output "attributes" {
  description = "Seal device specific attributes"
  value = {
    kms_key_id = aws_kms_key.key.arn
  }
}

// We output our resource name and a collection of those passed in to create a full list of key
// resources that might be required for instance roles that are associated with some unseal types.
output "resource_name" {
  description = "The awskms key name"
  value       = aws_kms_key.key.arn
}

output "resource_names" {
  description = "The list of awskms key names to associate with a role"
  value       = compact(concat([aws_kms_key.key.arn], var.other_resources))
}
