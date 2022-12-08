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

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "added_vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were added"
  default     = {}
}

locals {
  leftover_primary_instances = var.node_public_ip != "" ? {
    for k, v in var.vault_instances : k => v if contains(values(v), trimspace(var.node_public_ip))
  } : null
  all_instances          = var.node_public_ip != "" ? merge(var.added_vault_instances, local.leftover_primary_instances) : var.vault_instances
  updated_instance_count = length(local.all_instances)
  updated_instances = {
    for idx in range(local.updated_instance_count) : idx => {
      public_ip  = values(local.all_instances)[idx].public_ip
      private_ip = values(local.all_instances)[idx].private_ip
    }
  }
  node_ip = var.node_public_ip != "" ? var.node_public_ip : local.updated_instances[0].public_ip
  instance_private_ips = [
    for k, v in values(tomap(local.updated_instances)) :
    tostring(v["private_ip"])
  ]
  follower_public_ips = [
    for k, v in values(tomap(local.updated_instances)) :
    tostring(v["public_ip"]) if v["private_ip"] != trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  ]
  follower_private_ips = [
    for k, v in values(tomap(local.updated_instances)) :
    tostring(v["private_ip"]) if v["private_ip"] != trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  ]
}

resource "enos_remote_exec" "get_leader_private_ip" {
  environment = {
    VAULT_ADDR                 = "http://127.0.0.1:8200"
    VAULT_TOKEN                = var.vault_root_token
    vault_install_dir          = var.vault_install_dir
    vault_instance_private_ips = jsonencode(local.instance_private_ips)
  }

  scripts = ["${path.module}/scripts/get-leader-private-ip.sh"]

  transport = {
    ssh = {
      host = local.node_ip
    }
  }
}

output "leftover_primary_instances" {
  value = local.leftover_primary_instances
}

output "all_instances" {
  value = local.all_instances
}

output "updated_instance_count" {
  value = local.updated_instance_count
}

output "updated_instances" {
  value = local.updated_instances
}

output "leader_private_ip" {
  value = trimspace(enos_remote_exec.get_leader_private_ip.stdout)
}

output "leader_public_ip" {
  value = element([
    for k, v in values(tomap(local.all_instances)) :
    tostring(v["public_ip"]) if v["private_ip"] == trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  ], 0)
}

output "vault_instance_private_ips" {
  value = jsonencode(local.instance_private_ips)
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
