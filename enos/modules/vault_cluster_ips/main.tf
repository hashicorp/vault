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
  vault_bin_path = "${var.vault_install_dir}/vault"
}

resource "enos_remote_exec" "get_leader_public_ip" {
  content = templatefile("${path.module}/templates/get-leader-public-ip.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_instances   = jsonencode(local.instances)
  })

  transport = {
    ssh = {
      host = local.instances[0].public_ip
    }
  }
}

output "leader_public_ip" {
  value = trimspace(enos_remote_exec.get_leader_public_ip.stdout)
}

resource "enos_remote_exec" "get_follower_public_ips" {
  content = templatefile("${path.module}/templates/get-follower-public-ips.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_instances   = jsonencode(local.instances)
  })

  transport = {
    ssh = {
      host = local.instances[0].public_ip
    }
  }
}

output "follower_public_ips" {
  value = trimspace(enos_remote_exec.get_follower_public_ips.stdout)
}

resource "enos_remote_exec" "get_leader_private_ip" {
  content = templatefile("${path.module}/templates/get-leader-private-ip.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_instances   = jsonencode(local.instances)
  })

  transport = {
    ssh = {
      host = local.instances[0].public_ip
    }
  }
}

output "leader_private_ip" {
  value = trimspace(enos_remote_exec.get_leader_private_ip.stdout)
}

resource "enos_remote_exec" "get_follower_private_ips" {
  content = templatefile("${path.module}/templates/get-follower-private-ips.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_instances   = jsonencode(local.instances)
  })

  transport = {
    ssh = {
      host = local.instances[0].public_ip
    }
  }
}

output "follower_private_ips" {
  value = trimspace(enos_remote_exec.get_follower_private_ips.stdout)
}
