# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "cluster_name" {
  value = local.cluster_name
}

output "hosts" {
  description = "The ec2 instance target hosts"
  value       = local.hosts
}
