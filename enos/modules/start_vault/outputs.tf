# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "cluster_name" {
  description = "The Vault cluster name"
  value       = var.cluster_name
}

output "followers" {
  description = "The follower enos_vault_start resources"
  value       = enos_vault_start.followers
}

output "leader" {
  description = "The leader enos_vault_start resource"
  value       = enos_vault_start.leader
}

output "private_ips" {
  description = "Vault cluster target host private_ips"
  value       = [for host in var.target_hosts : host.private_ip]
}

output "public_ips" {
  description = "Vault cluster target host public_ips"
  value       = [for host in var.target_hosts : host.public_ip]
}

output "target_hosts" {
  description = "The vault cluster instances that were created"

  value = var.target_hosts
}
