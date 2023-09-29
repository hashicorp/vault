# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

// Shim module to handle the fact that Vault doesn't actually need a backend module when we use raft.
terraform {
  required_version = ">= 1.2.0"

  required_providers {
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.4.0"
    }
  }
}

variable "cluster_name" {
  default = null
}

variable "cluster_tag_key" {
  default = null
}

variable "config_dir" {
  default = null
}

variable "consul_log_level" {
  default = null
}

variable "data_dir" {
  default = null
}

variable "install_dir" {
  default = null
}

variable "license" {
  default = null
}

variable "log_dir" {
  default = null
}

variable "log_level" {
  default = null
}

variable "release" {
  default = null
}

variable "target_hosts" {
  default = null
}

output "private_ips" {
  value = [for host in var.target_hosts : host.private_ip]
}

output "public_ips" {
  value = [for host in var.target_hosts : host.public_ip]
}

output "target_hosts" {
  value = var.target_hosts
}
