# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

/*

Given our expected hosts, determine which is currently the leader and verify that all expected
nodes are either the leader or a follower.

*/

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster hosts that are expected to be in the cluster"
}

variable "ip_version" {
  type        = number
  description = "The IP version used for the Vault TCP listener"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  follower_hosts_list = [
    for idx in range(length(var.hosts)) : var.hosts[idx] if var.ip_version == 6 ?
    contains(tolist(local.follower_ipv6s), var.hosts[idx].ipv6) :
    contains(tolist(local.follower_private_ips), var.hosts[idx].private_ip)
  ]
  follower_hosts = {
    for idx in range(local.host_count - 1) : idx => try(local.follower_hosts_list[idx], null)
  }
  single_follower_hosts = {
    0: try(local.follower_hosts_list[0], null)
  }
  follower_ipv6s       = jsondecode(enos_remote_exec.follower_ipv6s.stdout)
  follower_private_ips = jsondecode(enos_remote_exec.follower_private_ipv4s.stdout)
  follower_public_ips  = [for host in local.follower_hosts : host.public_ip]
  host_count           = length(var.hosts)
  ipv6s                = [for k, v in values(tomap(var.hosts)) : tostring(v["ipv6"])]
  leader_host_list = [
    for idx in range(length(var.hosts)) : var.hosts[idx] if var.ip_version == 6 ?
    var.hosts[idx].ipv6 == local.leader_ipv6 :
    var.hosts[idx].private_ip == local.leader_private_ip
  ]
  leader_host       = try(local.leader_host_list[0], null)
  leader_ipv6       = trimspace(enos_remote_exec.leader_ipv6.stdout)
  leader_private_ip = trimspace(enos_remote_exec.leader_private_ipv4.stdout)
  leader_public_ip  = try(local.leader_host.public_ip, null)
  private_ips       = [for k, v in values(tomap(var.hosts)) : tostring(v["private_ip"])]
}

resource "enos_remote_exec" "leader_private_ipv4" {
  environment = {
    IP_VERSION        = var.ip_version
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/get-leader-ipv4.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

resource "enos_remote_exec" "leader_ipv6" {
  environment = {
    IP_VERSION        = var.ip_version
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/get-leader-ipv6.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

resource "enos_remote_exec" "follower_private_ipv4s" {
  environment = {
    IP_VERSION              = var.ip_version
    VAULT_ADDR              = var.vault_addr
    VAULT_INSTALL_DIR       = var.vault_install_dir
    VAULT_LEADER_PRIVATE_IP = local.leader_private_ip
    VAULT_PRIVATE_IPS       = jsonencode(local.private_ips)
    VAULT_TOKEN             = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/get-follower-ipv4s.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

resource "enos_remote_exec" "follower_ipv6s" {
  environment = {
    IP_VERSION        = var.ip_version
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_IPV6S       = jsonencode(local.ipv6s)
    VAULT_LEADER_IPV6 = local.leader_ipv6
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/get-follower-ipv6s.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

output "follower_hosts" {
  value = local.follower_hosts
}

output "follower_ipv6s" {
  value = local.follower_ipv6s
}

output "follower_private_ips" {
  value = local.follower_private_ips
}

output "follower_public_ips" {
  value = local.follower_public_ips
}

output "leader_host" {
  value = local.leader_host
}

output "leader_hosts" {
  value = { 0 : local.leader_host }
}

output "leader_ipv6" {
  value = local.leader_ipv6
}

output "leader_private_ip" {
  value = local.leader_private_ip
}

output "leader_public_ip" {
  value = local.leader_public_ip
}

output "single_follower_hosts" {
  value = local.single_follower_hosts
}