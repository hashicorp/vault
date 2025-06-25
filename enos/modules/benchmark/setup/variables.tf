# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

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

variable "k6_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The k6 target host"
}

variable "leader_addr" {
  type = string
}

variable "metrics_host" {
  type = object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  })
  description = "The metrics target host"
}

variable "prometheus_node_exporter_version" {
  type = string
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
  description = "The vault cluster hosts"
}

variable "vault_token" {
  type = string
}

variable "vpc_id" {
  type = string
}
