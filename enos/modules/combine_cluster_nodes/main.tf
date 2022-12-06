terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "primary_vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "added_vault_instances" {}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "node_public_ip" {
  type        = string
  description = "The primary node public ip"
  default     = ""
}

locals {
  leftover_primary_instances = {
    for k, v in var.primary_vault_instances : k => v if contains(values(v), trimspace(var.node_public_ip))
  }
  all_primary_instances = merge(var.added_vault_instances, local.leftover_primary_instances)
  updated_primary_count = length(local.all_primary_instances)
  updated_primary_instances = {
    for idx in range(local.updated_primary_count) : idx => {
      public_ip  = values(local.all_primary_instances)[idx].public_ip
      private_ip = values(local.all_primary_instances)[idx].private_ip
    }
  }
  updated_primary_instance_private_ips = [
    for k, v in values((tomap(local.updated_primary_instances))) :
    tostring(v["private_ip"])
  ]
  follower_public_ips = [
    for k, v in values((tomap(local.updated_primary_instances))) :
    tostring(v["public_ip"]) if v["private_ip"] != trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  ]
  follower_private_ips = [
    for k, v in values((tomap(local.updated_primary_instances))) :
    tostring(v["private_ip"]) if v["private_ip"] != trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  ]
}

resource "enos_remote_exec" "get_leader_private_ip" {
  environment = {
    VAULT_ADDR                 = "http://127.0.0.1:8200"
    VAULT_TOKEN                = var.vault_root_token
    vault_install_dir          = var.vault_install_dir
    vault_instance_private_ips = jsonencode(local.updated_primary_instance_private_ips)
  }

  scripts = ["${path.module}/scripts/get-leader-private-ip.sh"]

  transport = {
    ssh = {
      host = var.node_public_ip
    }
  }
}

output "all_primary_instance" {
  value = local.all_primary_instances
}

output "vault_instances" {
  value = local.updated_primary_instances
}

output "new_instance_count" {
  value = local.updated_primary_count
}


output "leader_private_ip" {
  value = trimspace(enos_remote_exec.get_leader_private_ip.stdout)
}

output "leader_public_ip" {
  value = element([
    for k, v in values((tomap(local.updated_primary_instances))) :
    tostring(v["public_ip"]) if v["private_ip"] == trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  ], 0)
}

output "vault_instance_private_ips" {
  value = jsonencode(local.updated_primary_instance_private_ips)
}

output "follower_public_ips" {
  value = local.follower_public_ips
}

output "follower_public_ip_1" {
  value = element(local.follower_public_ips, 0)
}

output "follower_public_ip_2" {
  value = element(local.follower_public_ips, 1)
}

output "follower_private_ips" {
  value = local.follower_private_ips
}

output "follower_private_ip_1" {
  value = element(local.follower_private_ips, 0)
}

output "follower_private_ip_2" {
  value = element(local.follower_private_ips, 1)
}
