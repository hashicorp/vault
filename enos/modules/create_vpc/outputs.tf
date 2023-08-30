# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "vpc_id" {
  description = "Created VPC ID"
  value       = aws_vpc.vpc.id
}

output "vpc_cidr" {
  description = "CIDR for whole VPC"
  value       = var.cidr
}

output "kms_key_arn" {
  description = "ARN of the generated KMS key"
  value       = try(aws_kms_key.key[0].arn, null)
}

output "kms_key_alias" {
  description = "Alias of the generated KMS key"
  value       = try(aws_kms_alias.alias[0].name, null)
}
