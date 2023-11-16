# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "id" {
  description = "Created VPC ID"
  value       = aws_vpc.vpc.id
}

output "cidr" {
  description = "CIDR for whole VPC"
  value       = var.cidr
}

output "cluster_id" {
  description = "A unique string associated with the VPC"
  value       = random_string.cluster_id.result
}
