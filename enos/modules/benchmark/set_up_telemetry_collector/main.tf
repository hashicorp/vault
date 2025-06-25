# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "hashicorp-forge/enos"
    }
  }
}

variable "cluster_id" {
  type = string
}

variable "consul_hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The consul hosts backing the vault cluster instances"
}

variable "grafana_version" {
  type = string
}

variable "grafana_http_port" {
  type = number
}

variable "host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
}

variable "k6_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
}

variable "prometheus_version" {
  type = string
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out. This is applied to each step so total timeout will be longer."
  default     = 120
}

variable "vault_hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances"
}

resource "random_string" "metrics" {
  length  = 8
  numeric = false
  special = false
  upper   = false
}

locals {
  metrics_id = "${var.cluster_id}-${random_string.metrics.result}"
  base_prometheus_environment = {
    PROMETHEUS_VERSION = var.prometheus_version
    RETRY_INTERVAL     = var.retry_interval
    TIMEOUT_SECONDS    = var.timeout
    K6_ADDR            = var.k6_host.private_ip
  }
  consul_prometheus_environment = {
    for k, v in var.consul_hosts : "CONSUL_${k}_ADDR" => var.consul_hosts[k].private_ip if length(var.consul_hosts) > 0
  }
  vault_prometheus_environment = {
    for k, v in var.vault_hosts : "VAULT_${k + 1}_ADDR" => var.vault_hosts[k].private_ip if length(var.vault_hosts) > 0
  }
  prometheus_environment = merge(
    local.base_prometheus_environment,
    local.consul_prometheus_environment,
    local.vault_prometheus_environment,
  )
}

resource "enos_remote_exec" "install_prometheus" {
  environment = local.prometheus_environment
  scripts     = [abspath("${path.module}/scripts/install-prometheus.sh")]

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

resource "enos_remote_exec" "install_grafana" {
  depends_on = [
    enos_remote_exec.install_prometheus,
  ]

  environment = {
    GRAFANA_VERSION = var.grafana_version
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/install-grafana.sh")]

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

resource "enos_file" "copy_grafana_dashboards" {
  depends_on = [
    enos_remote_exec.install_grafana,
  ]
  for_each = fileset(abspath("${path.module}/grafana-dashboards"), "*.json")

  source      = abspath("${path.module}/grafana-dashboards/${each.value}")
  destination = "/etc/grafana/dashboards/${basename(each.value)}"

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

resource "enos_remote_exec" "run_prometheus" {
  depends_on = [
    enos_remote_exec.install_grafana,
  ]

  scripts = [abspath("${path.module}/scripts/run-prometheus.sh")]

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

resource "enos_remote_exec" "run_grafana" {
  depends_on = [
    enos_file.copy_grafana_dashboards,
  ]

  environment = {
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/run-grafana.sh")]

  transport = {
    ssh = {
      host = var.host.public_ip
    }
  }
}

output "dashboard_url" {
  value = "http://${var.host.public_ip}:${var.grafana_http_port}"
}
