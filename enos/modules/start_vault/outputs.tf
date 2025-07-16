# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "api_addr_localhost" {
  description = "The localhost API address"
  value       = local.api_addr_localhost
}

output "api_addrs" {
  description = "The external API addresses of all nodes the cluster"
  value       = local.api_addrs
}

output "cluster_name" {
  description = "The Vault cluster name"
  value       = var.cluster_name
}

output "cluster_port" {
  description = "The Vault cluster request forwarding listener port"
  value       = var.cluster_port
}

output "external_storage_port" {
  description = "The Vault cluster non-raft external storage port"
  value       = var.external_storage_port
}

output "followers" {
  description = "The follower enos_vault_start resources"
  value       = enos_vault_start.followers
}

output "leader" {
  description = "The leader enos_vault_start resource"
  value       = enos_vault_start.leader
}

output "ipv6s" {
  description = "Vault cluster target host ipv6s"
  value       = [for host in var.hosts : host.ipv6]
}

output "listener_port" {
  description = "The Vault cluster TCP listener port"
  value       = var.listener_port
}

output "private_ips" {
  description = "Vault cluster target host private_ips"
  value       = [for host in var.hosts : host.private_ip]
}

output "public_ips" {
  description = "Vault cluster target host public_ips"
  value       = [for host in var.hosts : host.public_ip]
}

output "hosts" {
  description = "The vault cluster instances that were created"

  value = var.hosts
}
