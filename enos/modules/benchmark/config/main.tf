# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "consul_node_instance_types" {
  description = "The instance types to use depending on architecture"
  type = object({
    amd64 = string
    arm64 = string
  })
  default = {
    amd64 = "i3.4xlarge"
    arm64 = "t4g.small"
  }
}

variable "grafana_version" {
  type    = string
  default = "11.6.0"
}

variable "k6_instance_types" {
  description = "The instance types to use depending on architecture"
  type = object({
    amd64 = string
    arm64 = string
  })
  default = {
    amd64 = "c5d.4xlarge"
    arm64 = "t4g.small" # not actually used right now
  }
}

variable "metrics_instance_types" {
  description = "The instance types to use depending on architecture"
  type = object({
    amd64 = string
    arm64 = string
  })
  default = {
    amd64 = "m4.large"
    arm64 = "t4g.small" # not actually used right now
  }
}

variable "ports_ingress" {
  description = "Additional port mappings to allow for ingress"
  type = list(object({
    description = string
    port        = number
    protocol    = string
  }))
}

variable "prometheus_node_exporter_version" {
  type    = string
  default = "1.9.1"
}

variable "prometheus_version" {
  type    = string
  default = "3.3.0"
}

variable "storage_disk_iops" {
  description = <<-EOF
    The IOPS to request for storage disk. AWS accounts have a 100,000 IOPS limit
    by default so our default limit is 16,000. Some scenarios (backend:raft) can
    be configured to use the 24,000 maximum without requests. If you wish to
    test the backend:consul scenarios at the maximum you'll need to request
    a limit increase for your account.
  EOF
  type        = number
  default     = 16000
}

variable "vault_node_instance_types" {
  description = "The instance types to use depending on architecture"
  type = object({
    amd64 = string
    arm64 = string
  })
  default = {
    amd64 = "i3.4xlarge"
    arm64 = "t4g.small"
  }
}

locals {
  consul_http_port                 = 8500
  grafana_http_port                = 3000
  grafana_version                  = var.grafana_version
  prometheus_exporter_port         = 9100
  prometheus_http_port             = 9090
  prometheus_node_exporter_version = var.prometheus_node_exporter_version
  prometheus_version               = var.prometheus_version
  benchmark_ports = [
    {
      description = "CONSUL_HTTP"
      port        = local.consul_http_port
      protocol    = "tcp"
    },
    {
      description = "GRAFANA_HTTP"
      port        = local.grafana_http_port
      protocol    = "tcp"
    },
    {
      description = "PROMETHEUS_EXPORTER"
      port        = local.prometheus_exporter_port
      protocol    = "tcp"
    },
    {
      description = "PROMETHEUS_HTTP"
      port        = local.prometheus_http_port
      protocol    = "tcp"
    },
  ]
  required_ports = concat(var.ports_ingress, local.benchmark_ports)
}

output "consul_http_port" {
  value = local.consul_http_port
}

output "consul_instance_types" {
  value = var.consul_node_instance_types
}

output "grafana_http_port" {
  value = local.grafana_http_port
}

output "grafana_version" {
  value = local.grafana_version
}

output "k6_instance_types" {
  value = var.k6_instance_types
}

output "metrics_instance_types" {
  value = var.metrics_instance_types
}

output "prometheus_exporter_port" {
  value = local.prometheus_exporter_port
}

output "prometheus_http_port" {
  value = local.prometheus_http_port
}

output "prometheus_node_exporter_version" {
  value = local.prometheus_node_exporter_version
}

output "prometheus_version" {
  value = local.prometheus_version
}

output "required_ports" {
  value = local.required_ports
}

output "storage_disk_iops" {
  value = var.storage_disk_iops
}

output "vault_node_instance_types" {
  value = var.vault_node_instance_types
}
