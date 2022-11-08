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

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_api_addr" {
  type        = string
  description = "The API address of the Vault cluster"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
  follower_public_ips = {
    follower_public_ips = [
      for k, v in values((tomap(local.instances))) :
      tostring(v["public_ip"]) if v["private_ip"] != trimspace(enos_remote_exec.get_leader_private_ip.stdout)
    ]
  }
  follower_private_ips = {
    follower_private_ips = [
      for k, v in values((tomap(local.instances))) :
      tostring(v["private_ip"]) if v["private_ip"] != trimspace(enos_remote_exec.get_leader_private_ip.stdout)
    ]
  }
  vault_bin_path = "${var.vault_install_dir}/vault"
}

resource "enos_remote_exec" "get_leader_private_ip" {
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_TOKEN       = var.vault_root_token
    vault_install_dir = var.vault_install_dir
    vault_instances   = jsonencode(local.instances)
  }

  scripts = ["${path.module}/scripts/get-leader-private-ip.sh"]

  transport = {
    ssh = {
      host = local.instances[0].public_ip
    }
  }
}

output "leader_private_ip" {
  value = trimspace(enos_remote_exec.get_leader_private_ip.stdout)
}

output "leader_public_ip" {
  value = element([
    for k, v in values((tomap(local.instances))) :
    tostring(v["public_ip"]) if v["private_ip"] == trimspace(enos_remote_exec.get_leader_private_ip.stdout)
  ], 0)
}

output "follower_public_ips" {
  value = tomap(local.follower_public_ips)
}

output "follower_public_ip_1" {
  value = element(flatten(values(tomap(local.follower_public_ips))), 0)
}

output "follower_public_ip_2" {
  value = element(flatten(values(tomap(local.follower_public_ips))), 1)
}

output "follower_private_ips" {
  value = tomap(local.follower_private_ips)
}

output "follower_private_ip_1" {
  value = element(flatten(values(tomap(local.follower_private_ips))), 0)
}

output "follower_private_ip_2" {
  value = element(flatten(values(tomap(local.follower_private_ips))), 1)
}
