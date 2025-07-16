# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

output "private_ips" {
  description = "Consul cluster target host private_ips"
  value       = [for host in var.hosts : host.private_ip]
}

output "public_ips" {
  description = "Consul cluster target host public_ips"
  value       = [for host in var.hosts : host.public_ip]
}

output "hosts" {
  description = "The Consul cluster instances that were created"

  value = var.hosts
}
