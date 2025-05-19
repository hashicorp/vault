# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  vault_nodes = {
    for k, v in var.hosts : "vault_${k}" => v
  }
  consul_nodes = {
    for k, v in var.consul_hosts : "consul_${k}" => v if length(var.consul_hosts) > 0
  }
  k6_node = {
    "k6" = {
      ipv6       = try(aws_instance.k6.ipv6_addresses[0], "")
      private_ip = aws_instance.k6.private_ip
      public_ip  = aws_instance.k6.public_ip
    }
  }
  all_hosts = merge(local.vault_nodes, local.consul_nodes, local.k6_node)
  base_prometheus_environment = {
    PROMETHEUS_VERSION = var.prometheus_version
    RETRY_INTERVAL     = var.retry_interval
    TIMEOUT_SECONDS    = var.timeout
    VAULT_1_ADDR       = var.hosts[0].private_ip
    VAULT_2_ADDR       = var.hosts[1].private_ip
    VAULT_3_ADDR       = var.hosts[2].private_ip
    K6_ADDR            = aws_instance.k6.private_ip
  }
  consul_prometheus_environment = length(var.consul_hosts) > 0 ? {
    CONSUL_1_ADDR = var.consul_hosts[0].private_ip
    CONSUL_2_ADDR = var.consul_hosts[1].private_ip
    CONSUL_3_ADDR = var.consul_hosts[2].private_ip
  } : {}
  prometheus_environment = merge(
    local.base_prometheus_environment,
    local.consul_prometheus_environment
  )
}
