# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.3.24"
    }
  }
}

variable "ami_id" { default = null }
variable "cluster_name" { default = null }
variable "cluster_tag_key" { default = null }
variable "common_tags" { default = null }
variable "disable_selinux" { default = true }
variable "instance_count" { default = 3 }
variable "instance_cpu_max" { default = null }
variable "instance_cpu_min" { default = null }
variable "instance_mem_max" { default = null }
variable "instance_mem_min" { default = null }
variable "instance_types" { default = null }
variable "max_price" { default = null }
variable "project_name" { default = null }
variable "seal_key_names" { default = null }
variable "ssh_allow_ips" { default = null }
variable "ssh_keypair" { default = null }
variable "vpc_id" { default = null }

resource "random_string" "cluster_name" {
  length  = 8
  lower   = true
  upper   = false
  numeric = false
  special = false
}

output "cluster_name" {
  value = coalesce(var.cluster_name, random_string.cluster_name.result)
}

output "hosts" {
  value = { for idx in range(var.instance_count) : idx => {
    public_ip  = "null-public-${idx}"
    private_ip = "null-private-${idx}"
  } }
}
