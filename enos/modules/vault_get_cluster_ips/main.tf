# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "vault_instance_count" {
  type        = number
  description = "The number of instances in the vault cluster"
}

variable "vault_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster hosts. These are required to map private ip addresses to public addresses."
}

locals {
  follower_hosts_list = [for idx in range(var.vault_instance_count - 1) : {
    private_ip = local.follower_private_ips[idx]
    public_ip  = local.follower_public_ips[idx]
    }
  ]
  follower_hosts = {
    for idx in range(var.vault_instance_count - 1) : idx => try(local.follower_hosts_list[idx], null)
  }
  follower_private_ips = jsondecode(enos_remote_exec.get_follower_private_ips.stdout)
  follower_public_ips = [for idx in range(var.vault_instance_count) : var.vault_hosts[idx].public_ip if contains(
    local.follower_private_ips, var.vault_hosts[idx].private_ip)
  ]
  leader_host = {
    private_ip = local.leader_private_ip
    public_ip  = local.leader_public_ip
  }
  leader_private_ip = trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  leader_public_ip = element([
    for idx in range(var.vault_instance_count) : var.vault_hosts[idx].public_ip if var.vault_hosts[idx].private_ip == local.leader_private_ip
  ], 0)
  private_ips = [for k, v in values(tomap(var.vault_hosts)) : tostring(v["private_ip"])]
}

resource "enos_remote_exec" "get_leader_private_ip" {
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/get-leader-private-ip.sh")]

  transport = {
    ssh = {
      host = var.vault_hosts[0].public_ip
    }
  }
}

resource "enos_remote_exec" "get_follower_private_ips" {
  environment = {
    VAULT_ADDR                 = "http://127.0.0.1:8200"
    VAULT_TOKEN                = var.vault_root_token
    VAULT_LEADER_PRIVATE_IP    = local.leader_private_ip
    VAULT_INSTANCE_PRIVATE_IPS = jsonencode(local.private_ips)
    VAULT_INSTALL_DIR          = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/get-follower-private-ips.sh")]

  transport = {
    ssh = {
      host = var.vault_hosts[0].public_ip
    }
  }
}

output "follower_hosts" {
  value = local.follower_hosts
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

output "leader_private_ip" {
  value = local.leader_private_ip
}

output "leader_public_ip" {
  value = local.leader_public_ip
}
