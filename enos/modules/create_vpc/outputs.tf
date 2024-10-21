# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "id" {
  description = "Created VPC ID"
  value       = aws_vpc.vpc.id
}

output "ipv4_cidr" {
  description = "The VPC subnet CIDR for ipv4 mode"
  value       = var.ipv4_cidr
}

output "ipv6_cidr" {
  description = "The VPC subnet CIDR for ipv6 mode"
  value       = aws_vpc.vpc.ipv6_cidr_block
}

output "cluster_id" {
  description = "A unique string associated with the VPC"
  value       = random_string.cluster_id.result
}
