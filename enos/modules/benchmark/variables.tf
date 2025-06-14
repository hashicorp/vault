# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "ami_id" {
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
  type    = string
  default = "11.6.0"
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances"
}

variable "k6_instance_type" {
  type    = string
  default = "c5d.4xlarge"
}

variable "leader_addr" {
  type = string
}

variable "metrics_instance_type" {
  type    = string
  default = "m4.large"
}

variable "metrics_security_group_ids" {
  type = object({
    grafana    = string
    prometheus = string
  })
  default = null
}

variable "project_name" {
  type    = string
  default = ""
}

variable "prometheus_node_exporter_version" {
  type    = string
  default = "1.9.1"
}

variable "prometheus_version" {
  type    = string
  default = "3.3.0"
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

variable "ssh_keypair" {
  type = string
}

variable "target_security_group_id" {
  type = string
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out. This is applied to each step so total timeout will be longer."
  default     = 120
}

variable "vault_token" {
  type = string
}

variable "vpc_id" {
  type = string
}
