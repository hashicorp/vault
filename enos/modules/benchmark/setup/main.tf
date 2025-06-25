# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "hashicorp-forge/enos"
    }
  }
}

module "set_up_k6" {
  source = "../set_up_k6"

  cluster_id             = var.vpc_id
  host                   = var.k6_host
  leader_addr            = var.leader_addr
  metrics_collector_host = var.metrics_host
  retry_interval         = var.retry_interval
  timeout                = var.timeout
  vault_token            = var.vault_token
  vault_hosts            = var.vault_hosts
}

module "set_up_telemetry_collector" {
  source = "../set_up_telemetry_collector"

  cluster_id         = var.vpc_id
  consul_hosts       = var.consul_hosts
  grafana_version    = var.grafana_version
  grafana_http_port  = var.grafana_http_port
  host               = var.metrics_host
  k6_host            = var.k6_host
  prometheus_version = var.prometheus_version
  retry_interval     = var.retry_interval
  timeout            = var.timeout
  vault_hosts        = var.vault_hosts
}

module "enable_telemetry_consul" {
  source = "../enable_telemetry_consul"

  hosts = var.consul_hosts
}

locals {
  vault_hosts = {
    for k, v in var.vault_hosts : "vault_${k}" => v
  }
  consul_hosts = {
    for k, v in var.consul_hosts : "consul_${k}" => v if length(var.consul_hosts) > 0
  }
  all_hosts = merge(local.vault_hosts, local.consul_hosts, module.set_up_k6.hosts)
}

module "enable_telemetry_node_exporter" {
  depends_on = [
    module.set_up_telemetry_collector
  ]
  source = "../enable_telemetry_node_exporter"

  hosts                            = local.all_hosts
  prometheus_node_exporter_version = var.prometheus_node_exporter_version
  retry_interval                   = var.retry_interval
  timeout                          = var.timeout
}

output "dashboard_url" {
  value = module.set_up_telemetry_collector.dashboard_url
}
