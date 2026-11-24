# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.3.24"
    }
    fyre = {
      source = "registry.terraform.io/hashicorp-forge/fyre"
    }
  }
}

locals {
  cluster_name = coalesce(var.cluster_name, random_string.cluster_name.result)
  instances    = toset([for idx in range(var.instance_count) : tostring(idx)])
  name_prefix  = "${var.project_name}-${local.cluster_name}-${random_string.unique_id.result}"
  platform     = var.arch == "amd64" ? "x" : "z"

  hosts = {
    for idx in range(var.instance_count) : idx => {
      ipv6 = ""
      public_ip = try(
        fyre_vm.targets[tostring(idx)].ips[index(fyre_vm.targets[tostring(idx)].ips[*].type, "public")].ip,
        try(
          fyre_vm.targets[tostring(idx)].ips[index(fyre_vm.targets[tostring(idx)].ips[*].scope, "public")].ip,
          ""
        )
      )
      private_ip = try(
        fyre_vm.targets[tostring(idx)].ips[index(fyre_vm.targets[tostring(idx)].ips[*].type, "private")].ip,
        try(
          fyre_vm.targets[tostring(idx)].ips[index(fyre_vm.targets[tostring(idx)].ips[*].scope, "private")].ip,
          try(
            fyre_vm.targets[tostring(idx)].ips[index(fyre_vm.targets[tostring(idx)].ips[*].type, "public")].ip,
            fyre_vm.targets[tostring(idx)].ips[0].ip
          )
        )
      )
    }
  }
}

resource "random_string" "cluster_name" {
  length  = 8
  lower   = true
  upper   = false
  numeric = false
  special = false
}

resource "random_string" "unique_id" {
  length  = 4
  lower   = true
  upper   = false
  numeric = false
  special = false
}

resource "fyre_vm" "targets" {
  for_each = local.instances

  os               = var.os
  platform         = local.platform
  cpu              = var.cpu
  memory           = var.memory
  additional_disks = var.additional_disks
  hostname         = "${local.name_prefix}-${var.cluster_tag_key}-${each.key}"
  description      = var.description != null ? var.description : "${local.name_prefix}-${var.cluster_tag_key}"
  expiration       = var.expiration
  public_network   = var.public_network
  dns              = var.dns
  disable_delete   = var.disable_delete
  quota_type       = var.quota_type
  time_to_live     = var.quota_type == "quick_burn" ? var.quick_burn_ttl : null
  product_group_id = var.product_group_id
  ssh_keys         = [file(var.public_key_path)]
}

module "disable_selinux" {
  depends_on = [fyre_vm.targets]
  source     = "../disable_selinux"
  count      = var.disable_selinux ? 1 : 0

  hosts = local.hosts
}

resource "enos_remote_exec" "enable_apt_repos" {
  for_each = local.hosts

  scripts = [abspath("${path.module}/scripts/maybe-enable-apt-repos.sh")]
  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

output "cluster_name" {
  description = "The logical cluster name for the provisioned Fyre VMs"
  value       = local.cluster_name
}

output "hosts" {
  description = "The Fyre VM target hosts"
  value       = local.hosts
}

output "vm_ids" {
  description = "The created Fyre VM identifiers"
  value       = { for idx, vm in fyre_vm.targets : idx => vm.vm_id }
}

output "fqdns" {
  description = "The created Fyre VM FQDNs"
  value       = { for idx, vm in fyre_vm.targets : idx => vm.fqdn }
}
