# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

output "api_addr_localhost" {
  description = "The localhost API address"
  value       = module.start_vault.api_addr_localhost
}

output "api_addrs" {
  description = "The external API addresses of all nodes the cluster"
  value       = module.start_vault.api_addrs
}

output "audit_device_file_path" {
  description = "The file path for the audit device, if enabled"
  value       = var.enable_audit_devices ? local.audit_device_file_path : null
}

output "cluster_name" {
  description = "The Vault cluster name"
  value       = var.cluster_name
}

output "cluster_port" {
  description = "The Vault cluster request forwarding listener port"
  value       = module.start_vault.cluster_port
}

output "external_storage_port" {
  description = "The Vault cluster non-raft external storage port"
  value       = module.start_vault.external_storage_port
}

output "hosts" {
  description = "The vault cluster instances that were created"

  value = var.hosts
}

output "ipv6s" {
  description = "Vault cluster target host ipv6 addresses"
  value       = [for host in var.hosts : host.ipv6]
}

output "keys_base64" {
  value = try(module.start_vault.keys_base64, null)
}

output "keys_base64_secondary" {
  value = try(module.start_vault.keys_base64_secondary, null)
}

output "listener_port" {
  description = "The Vault cluster TCP listener port"
  value       = module.start_vault.listener_port
}

output "private_ips" {
  description = "Vault cluster target host private_ips"
  value       = [for host in var.hosts : host.private_ip]
}

output "public_ips" {
  description = "Vault cluster target host public_ips"
  value       = [for host in var.hosts : host.public_ip]
}

output "recovery_keys_b64" {
  value = try(enos_vault_init.leader[0].recovery_keys_b64, [])
}

output "recovery_keys_hex" {
  value = try(enos_vault_init.leader[0].recovery_keys_hex, [])
}

output "recovery_key_shares" {
  value = try(enos_vault_init.leader[0].recovery_keys_shares, -1)
}

output "recovery_threshold" {
  value = try(enos_vault_init.leader[0].recovery_keys_threshold, -1)
}

output "root_token" {
  value = coalesce(var.root_token, try(enos_vault_init.leader[0].root_token, null), "none")
}

output "unseal_keys_b64" {
  value = try(enos_vault_init.leader[0].unseal_keys_b64, [])
}

output "unseal_keys_hex" {
  value = try(enos_vault_init.leader[0].unseal_keys_hex, null)
}

output "unseal_shares" {
  value = try(enos_vault_init.leader[0].unseal_keys_shares, -1)
}

output "unseal_threshold" {
  value = try(enos_vault_init.leader[0].unseal_keys_threshold, -1)
}
