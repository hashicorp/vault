# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "cluster_id" {
  type = string
}

variable "cluster_meta" {
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

output "alias" {
  description = "The key alias name"
  value       = aws_kms_alias.alias.name
}

output "id" {
  description = "The key ID"
  value       = aws_kms_key.key.key_id
}

output "resource_name" {
  description = "The ARN"
  value       = aws_kms_key.key.arn
}

output "resource_names" {
  description = "The list of names"
  value       = compact(concat([aws_kms_key.key.arn], var.other_resources))
}
