output "public_ips" {
  description = "Vault cluster target host public_ips"
  value       = [for host in var.target_hosts : host.public_ip]
}

output "private_ips" {
  description = "Vault cluster target host private_ips"
  value       = [for host in var.target_hosts : host.private_ip]
}

output "target_hosts" {
  description = "The vault cluster instances that were created"

  value = var.target_hosts
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

output "cluster_name" {
  description = "The Vault cluster name"
  value       = var.cluster_name
}
