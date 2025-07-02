# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts on which to install and enable host metrics"
}

variable "prometheus_node_exporter_version" {
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

resource "enos_remote_exec" "install_prometheus_node_exporter" {
  for_each = var.hosts

  environment = {
    PROMETHEUS_NODE_EXPORTER_VERSION = var.prometheus_node_exporter_version
    RETRY_INTERVAL                   = var.retry_interval
    TIMEOUT_SECONDS                  = var.timeout
  }
  scripts = [abspath("${path.module}/scripts/install-prometheus-node-exporter.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "run_prometheus_node_exporter" {
  depends_on = [
    enos_remote_exec.install_prometheus_node_exporter,
  ]
  for_each = var.hosts

  environment = {
    PROMETHEUS_NODE_EXPORTER_VERSION = var.prometheus_node_exporter_version
  }
  scripts = [abspath("${path.module}/scripts/run-prometheus-node-exporter.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
