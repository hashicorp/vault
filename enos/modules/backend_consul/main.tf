# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_version = ">= 1.2.0"

  required_providers {
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.4.4"
    }
  }
}

locals {
  bin_path = "${var.install_dir}/consul"
}

resource "enos_bundle_install" "consul" {
  for_each = var.target_hosts

  destination = var.install_dir
  release     = merge(var.release, { product = "consul" })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_consul_start" "consul" {
  for_each = enos_bundle_install.consul

  bin_path   = local.bin_path
  data_dir   = var.data_dir
  config_dir = var.config_dir
  config = {
    data_dir         = var.data_dir
    datacenter       = "dc1"
    retry_join       = ["provider=aws tag_key=${var.cluster_tag_key} tag_value=${var.cluster_name}"]
    server           = true
    bootstrap_expect = length(var.target_hosts)
    log_level        = var.log_level
    log_file         = var.log_dir
  }
  license   = var.license
  unit_name = "consul"
  username  = "consul"

  transport = {
    ssh = {
      host = var.target_hosts[each.key].public_ip
    }
  }
}
